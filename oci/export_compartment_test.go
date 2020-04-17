package oci

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	resourceDiscoveryTestCompartmentOcid   = "ocid1.testcompartment.abc"
	resourceDiscoveryTestActiveLifecycle   = "ACTIVE"
	resourceDiscoveryTestInactiveLifecycle = "INACTIVE"
)

var exportParentDefinition = &TerraformResourceHints{
	resourceClass:               "oci_test_parent",
	datasourceClass:             "oci_test_parents",
	resourceAbbreviation:        "parent",
	datasourceItemsAttr:         "items",
	discoverableLifecycleStates: []string{resourceDiscoveryTestActiveLifecycle},
	alwaysExportable:            true,
}

var exportChildDefinition = &TerraformResourceHints{
	resourceClass:               "oci_test_child",
	datasourceClass:             "oci_test_children",
	resourceAbbreviation:        "child",
	datasourceItemsAttr:         "item_summaries",
	discoverableLifecycleStates: []string{resourceDiscoveryTestActiveLifecycle},
	requireResourceRefresh:      true,
}

var tenancyTestingResourceGraph = TerraformResourceGraph{
	"oci_identity_tenancy": {
		{
			TerraformResourceHints: exportParentDefinition,
		},
	},
	"oci_test_parent": {
		{
			TerraformResourceHints: exportChildDefinition,
			datasourceQueryParams:  map[string]string{"parent_id": "id"},
		},
	},
}

var compartmentTestingResourceGraph = TerraformResourceGraph{
	"oci_identity_compartment": {
		{
			TerraformResourceHints: exportParentDefinition,
		},
	},
	"oci_test_parent": {
		{
			TerraformResourceHints: exportChildDefinition,
			datasourceQueryParams:  map[string]string{"parent_id": "id"},
		},
	},
}

var childrenResources map[string]map[string]interface{}
var parentResources map[string]map[string]interface{}

// Test resources used by resource discovery tests
func testParentResource() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: createTestParent,
		Read:   readTestParent,
		Delete: deleteTestParent,
		Schema: map[string]*schema.Schema{
			// Required
			"compartment_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Optional
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"a_map": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     schema.TypeString,
			},
			"a_string": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"a_bool": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"a_int": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"a_float": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"a_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     schema.TypeString,
				ForceNew: true,
			},
			"a_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     schema.TypeString,
				ForceNew: true,
				Set:      literalTypeHashCodeForSets,
			},
			"a_nested": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nested_string": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"nested_bool": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"nested_int": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"nested_float": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},

			// Computed
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func testChildResource() *schema.Resource {
	// Reuse the parent schema and add a parent dependency attribute
	childResourceSchema := testParentResource().Schema
	childResourceSchema["parent_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}

	// Don't have a display_name attribute so a different name can be generated
	delete(childResourceSchema, "display_name")

	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: createTestChild,
		Read:   readTestChild,
		Delete: deleteTestChild,
		Schema: childResourceSchema,
	}
}

func testParentsDatasource() *schema.Resource {
	return &schema.Resource{
		Read: listTestParents,
		Schema: map[string]*schema.Schema{
			"compartment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(testParentResource()),
			},
		},
	}
}

func testChildrenDatasource() *schema.Resource {
	// Convert child resource schema to datasource schema
	childDatasourceSchema := GetDataSourceItemSchema(testChildResource())

	// Remove some attributes from datasource (i.e. treat the datasource results as incomplete representations of the resource)
	delete(childDatasourceSchema.Schema, "a_nested")

	return &schema.Resource{
		Read: listTestChildren,
		Schema: map[string]*schema.Schema{
			"compartment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"item_summaries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     childDatasourceSchema,
			},
		},
	}
}

func createTestParent(d *schema.ResourceData, m interface{}) error {
	return nil
}

func readTestParent(d *schema.ResourceData, m interface{}) error {
	initTestResources()
	if resource, exists := parentResources[d.Id()]; exists {
		for key, value := range resource {
			d.Set(key, value)
		}
	} else {
		return fmt.Errorf("could not find parent with id %s", d.Id())
	}
	return nil
}

func listTestParents(d *schema.ResourceData, m interface{}) error {
	initTestResources()
	results := make([]interface{}, len(parentResources))
	for i := 0; i < len(parentResources); i++ {
		id := getTestResourceId("parent", i)
		resource := parentResources[id]
		resource["id"] = id
		results[i] = resource
	}
	d.Set("items", results)
	return nil
}

func deleteTestParent(d *schema.ResourceData, m interface{}) error {
	return nil
}

func createTestChild(d *schema.ResourceData, m interface{}) error {
	return nil
}

func readTestChild(d *schema.ResourceData, m interface{}) error {
	initTestResources()
	if resource, exists := childrenResources[d.Id()]; exists {
		for key, value := range resource {
			d.Set(key, value)
		}
	} else {
		return fmt.Errorf("could not find child with id %s", d.Id())
	}
	return nil
}

func listTestChildren(d *schema.ResourceData, m interface{}) error {
	initTestResources()

	parentId, parentIdExists := d.GetOkExists("parent_id")
	results := []interface{}{}
	for i := 0; i < len(childrenResources); i++ {
		id := getTestResourceId("child", i)
		resource := childrenResources[id]
		resource["id"] = id

		if !parentIdExists || resource["parent_id"] == parentId {
			copyResource := map[string]interface{}{}
			for key, val := range resource {
				if key == "a_nested" {
					continue
				}
				copyResource[key] = val
			}
			results = append(results, copyResource)
		}
	}
	if err := d.Set("item_summaries", results); err != nil {
		return err
	}
	return nil
}

func deleteTestChild(d *schema.ResourceData, m interface{}) error {
	return nil
}

func initResourceDiscoveryTests() {
	resourceNameCount = map[string]int{}

	resourcesMap = ResourcesMap()
	datasourcesMap = DataSourcesMap()

	resourcesMap["oci_test_parent"] = testParentResource()
	resourcesMap["oci_test_child"] = testChildResource()
	datasourcesMap["oci_test_parents"] = testParentsDatasource()
	datasourcesMap["oci_test_children"] = testChildrenDatasource()

	tenancyResourceGraphs["tenancy_testing"] = tenancyTestingResourceGraph
	compartmentResourceGraphs["compartment_testing"] = compartmentTestingResourceGraph

	initTestResources()
}

func cleanupResourceDiscoveryTests() {
	delete(resourcesMap, "oci_test_parent")
	delete(resourcesMap, "oci_test_child")
	delete(datasourcesMap, "oci_test_parents")
	delete(datasourcesMap, "oci_test_children")
}

func initTestResources() {
	numParentResources := 4
	if parentResources == nil || len(parentResources) != numParentResources {
		parentResources = make(map[string]map[string]interface{}, numParentResources)
		for i := 0; i < numParentResources; i++ {
			parentResources[getTestResourceId("parent", i)] = generateTestResourceFromSchema(i, resourcesMap["oci_test_parent"].Schema)
		}
	}

	numChildrenResourcesPerParent := 2
	numChildrenResources := numParentResources * numChildrenResourcesPerParent
	if childrenResources == nil || len(childrenResources) != numChildrenResources {
		childrenResources = make(map[string]map[string]interface{}, numParentResources*numChildrenResourcesPerParent)
		childCount := 0
		for i := 0; i < len(parentResources); i++ {
			parentId := getTestResourceId("parent", i)
			for j := 0; j < numChildrenResourcesPerParent; j++ {
				childResource := generateTestResourceFromSchema(i, resourcesMap["oci_test_child"].Schema)
				childResource["parent_id"] = parentId

				childrenResources[getTestResourceId("child", childCount)] = childResource
				childCount++
			}
		}
	}
}

func getRootCompartmentResource() *OCIResource {
	return &OCIResource{
		compartmentId: resourceDiscoveryTestCompartmentOcid,
		TerraformResource: TerraformResource{
			id:             resourceDiscoveryTestCompartmentOcid,
			terraformClass: "oci_identity_compartment",
			terraformName:  "export",
		},
	}
}

func getTestResourceId(resourceType string, id int) string {
	return fmt.Sprintf("ocid1.%s.abcdefghiklmnop.%d", resourceType, id)
}

func generatePrimitiveValue(id int, valueType schema.ValueType) interface{} {
	switch valueType {
	case schema.TypeInt:
		return id
	case schema.TypeBool:
		return true
	case schema.TypeFloat:
		res, _ := strconv.ParseFloat(fmt.Sprintf("%d.%d", id, id), 64)
		return res
	case schema.TypeString:
		return fmt.Sprintf("string%d", id)
	}
	return nil
}

func generateTestResourceFromSchema(id int, resourceSchemaMap map[string]*schema.Schema) map[string]interface{} {
	result := map[string]interface{}{}
	for resourceAttribute, resourceSchema := range resourceSchemaMap {
		switch resourceAttribute {
		case "state":
			result[resourceAttribute] = resourceDiscoveryTestActiveLifecycle
			continue
		case "time_created":
			result[resourceAttribute] = time.Now().Format(time.RFC3339)
			continue
		}

		switch resourceSchema.Type {
		case schema.TypeInt, schema.TypeBool, schema.TypeFloat, schema.TypeString:
			result[resourceAttribute] = generatePrimitiveValue(id, resourceSchema.Type)
		case schema.TypeMap:
			mapResult := map[string]interface{}{}
			if elemType, ok := resourceSchema.Elem.(schema.ValueType); ok {
				for i := 0; i < id; i++ {
					mapKey := fmt.Sprintf("key%d", i)
					mapResult[mapKey] = generatePrimitiveValue(i, elemType)
				}
			}
			result[resourceAttribute] = mapResult
		case schema.TypeList, schema.TypeSet:
			listResult := make([]interface{}, id)
			switch elemType := resourceSchema.Elem.(type) {
			case schema.ValueType:
				for i := 0; i < id; i++ {
					listResult[i] = generatePrimitiveValue(i, elemType)
				}
			case *schema.Resource:
				for i := 0; i < id; i++ {
					listResult[i] = generateTestResourceFromSchema(i, elemType.Schema)
				}
			}
			result[resourceAttribute] = listResult
		}
	}
	return result
}

// Basic test to ensure that RunExportCommand generates TF artifacts
func TestUnitRunExportCommand_basic(t *testing.T) {
	initResourceDiscoveryTests()
	defer cleanupResourceDiscoveryTests()
	compartmentId := resourceDiscoveryTestCompartmentOcid
	outputDir, err := os.Getwd()
	outputDir = fmt.Sprintf("%s%sdiscoveryTest-%d", outputDir, string(os.PathSeparator), time.Now().Nanosecond())
	if err = os.Mkdir(outputDir, os.ModePerm); err != nil {
		t.Logf("unable to mkdir %s. err: %v", outputDir, err)
		t.Fail()
	}

	args := &ExportCommandArgs{
		CompartmentId: &compartmentId,
		Services:      []string{"compartment_testing", "tenancy_testing"},
		OutputDir:     &outputDir,
		GenerateState: false,
	}

	if err = RunExportCommand(args); err != nil {
		t.Logf("export command failed due to err: %v", err)
		t.Fail()
	}

	if _, err = os.Stat(fmt.Sprintf("%s%stenancy_testing.tf", outputDir, string(os.PathSeparator))); os.IsNotExist(err) {
		t.Logf("no tenancy_testing.tf file generated")
		t.Fail()
	}

	if _, err = os.Stat(fmt.Sprintf("%s%scompartment_testing.tf", outputDir, string(os.PathSeparator))); os.IsNotExist(err) {
		t.Logf("no compartment_testing.tf file generated")
		t.Fail()
	}

	if _, err = os.Stat(fmt.Sprintf("%s%sterraform.tfstate", outputDir, string(os.PathSeparator))); !os.IsNotExist(err) {
		t.Logf("found terraform.tfstate even though it wasn't expected")
	}

	os.RemoveAll(outputDir)
}

func TestUnitRunExportCommand_error(t *testing.T) {
	initResourceDiscoveryTests()
	defer cleanupResourceDiscoveryTests()
	compartmentId := resourceDiscoveryTestCompartmentOcid
	outputDir, err := os.Getwd()
	outputDir = fmt.Sprintf("%s%sdiscoveryTest-%d", outputDir, string(os.PathSeparator), time.Now().Nanosecond())
	if err = os.Mkdir(outputDir, os.ModePerm); err != nil {
		t.Logf("unable to mkdir %s. err: %v", outputDir, err)
		t.Fail()
	}

	nonexistentOutputDir := fmt.Sprintf("%s%s%s", outputDir, string(os.PathSeparator), "baddirectory")
	args := &ExportCommandArgs{
		CompartmentId: &compartmentId,
		Services:      []string{"compartment_testing", "tenancy_testing"},
		OutputDir:     &nonexistentOutputDir,
		GenerateState: false,
	}
	if err = RunExportCommand(args); err == nil {
		t.Logf("export command expected to fail due to non-existent path, but it succeeded")
		t.Fail()
	}

	os.RemoveAll(outputDir)
}

// Test that resources can be found using a resource dependency graph
func TestUnitFindResources_basic(t *testing.T) {
	initResourceDiscoveryTests()
	defer cleanupResourceDiscoveryTests()
	rootResource := getRootCompartmentResource()

	results, err := findResources(nil, rootResource, compartmentTestingResourceGraph, nil)
	if err != nil {
		t.Logf("got error from findResources: %v", err)
		t.Fail()
	}

	if len(results) != len(childrenResources)+len(parentResources) {
		t.Logf("got %d results but expected %d results", len(results), len(childrenResources)+len(parentResources))
		t.Fail()
	}

	for _, foundResource := range results {
		if foundResource.terraformClass == "oci_test_child" {
			if _, resourceRefreshAttributeExists := foundResource.sourceAttributes["a_nested"]; !resourceRefreshAttributeExists {
				t.Logf("child resource is missing an expected attribute that should have been filled by a resource refresh")
				t.Fail()
			}

			expectedTfNamePrefix := fmt.Sprintf("%s_%s", foundResource.parent.terraformName, exportChildDefinition.resourceAbbreviation)
			if !strings.HasPrefix(foundResource.terraformName, expectedTfNamePrefix) {
				t.Logf("child resource should have a name with prefix '%s' but name is '%s' instead", expectedTfNamePrefix, foundResource.terraformName)
				t.Fail()
			}
		}
	}
}

// Test that only targeted ocid resources are exportable
func TestUnitFindResources_restrictedOcids(t *testing.T) {
	initResourceDiscoveryTests()
	defer cleanupResourceDiscoveryTests()
	rootResource := getRootCompartmentResource()

	// Parent resources are defined as alwaysExportable. So even if it's not specified in the ocids, it should be exported.
	restrictedOcidTests := []map[string]interface{}{
		{
			"ocids":                map[string]bool{getTestResourceId("parent", 0): false, getTestResourceId("child", 0): false},
			"numExpectedResources": len(parentResources) + 1,
		},
		{
			"ocids":                map[string]bool{getTestResourceId("parent", 0): false, getTestResourceId("child", 3): false},
			"numExpectedResources": len(parentResources) + 1,
		},
		{
			"ocids":                map[string]bool{getTestResourceId("parent", 0): false, getTestResourceId("child", 0): false, "nonexistentID": false},
			"numExpectedResources": len(parentResources) + 1,
		},
		{
			"ocids":                map[string]bool{getTestResourceId("child", 0): false, getTestResourceId("child", 3): false, "nonexistentID": false},
			"numExpectedResources": len(parentResources) + 2,
		},
	}

	for idx, testCase := range restrictedOcidTests {
		t.Logf("running test #%d with following ocids: ", idx)
		restrictedOcids := testCase["ocids"].(map[string]bool)
		for ocid := range restrictedOcids {
			t.Logf(ocid)
		}

		results, err := findResources(nil, rootResource, compartmentTestingResourceGraph, restrictedOcids)
		if err != nil {
			t.Logf("got error from findResources: %v", err)
			t.Fail()
		}

		exportResourceCount := 0
		for _, resource := range results {
			if !resource.omitFromExport {
				exportResourceCount++
			}
		}
		if exportResourceCount != testCase["numExpectedResources"].(int) {
			t.Logf("expected %d resources to be exported, but got %d", testCase["numExpectedResources"].(int), len(results))
			t.Fail()
		}
	}
}

// Test that overriden find function is invoked if a resource has one
func TestUnitFindResources_overrideFn(t *testing.T) {
	initResourceDiscoveryTests()
	defer cleanupResourceDiscoveryTests()
	rootResource := getRootCompartmentResource()

	// Create an override function that returns nothing when discovering child test resources
	exportChildDefinition.findResourcesOverrideFn = func(*OracleClients, *TerraformResourceAssociation, *OCIResource) ([]*OCIResource, error) {
		return []*OCIResource{}, nil
	}
	defer func() { exportChildDefinition.findResourcesOverrideFn = nil }()

	results, err := findResources(nil, rootResource, compartmentTestingResourceGraph, nil)
	if err != nil {
		t.Logf("got error from findResources: %v", err)
		t.Fail()
	}

	// Check that we got only parent resources, because child resources have been overridden to return nothing in their discovery function
	if len(results) != len(parentResources) {
		t.Logf("got %d results but expected %d results", len(results), len(parentResources))
		t.Fail()
	}

	for _, foundResource := range results {
		if foundResource.terraformClass == "oci_test_child" {
			t.Logf("oci_test_child resource was returned when not expected")
			t.Fail()
		}
	}
}

// Test that process resource function is invoked if a resource has one
func TestUnitFindResources_processResourceFn(t *testing.T) {
	initResourceDiscoveryTests()
	defer cleanupResourceDiscoveryTests()
	rootResource := getRootCompartmentResource()

	// Create a processing function that adds a new attribute to every discovered child resource
	exportChildDefinition.processDiscoveredResourcesFn = func(clients *OracleClients, resources []*OCIResource) ([]*OCIResource, error) {
		for _, resource := range resources {
			resource.sourceAttributes["added_by_process_function"] = true
		}
		return resources, nil
	}
	defer func() { exportChildDefinition.processDiscoveredResourcesFn = nil }()

	results, err := findResources(nil, rootResource, compartmentTestingResourceGraph, nil)
	if err != nil {
		t.Logf("got error from findResources: %v", err)
		t.Fail()
	}

	// Check that we got only parent resources, because child resources have been overridden to return nothing in their discovery function
	if len(results) != len(parentResources)+len(childrenResources) {
		t.Logf("got %d results but expected %d results", len(results), len(childrenResources)+len(parentResources))
		t.Fail()
	}

	for _, foundResource := range results {
		if foundResource.terraformClass == "oci_test_child" {
			if _, ok := foundResource.sourceAttributes["added_by_process_function"]; !ok {
				t.Logf("oci_test_child resource was returned when not expected")
				t.Fail()
			}
		}
	}
}

// Test that Terraform names can be generated from discovered resources
func TestUnitGenerateTerraformNameFromResource_basic(t *testing.T) {
	type testCase struct {
		resource     map[string]interface{}
		schema       map[string]*schema.Schema
		expectError  bool
		expectedName string
	}

	testResourceSchema := testParentResource().Schema
	testCases := []testCase{
		{
			resource:     map[string]interface{}{"display_name": "abc"},
			schema:       testResourceSchema,
			expectedName: "export_abc",
		},
		{
			// Repeating it should result in a different name than the previous test
			resource:     map[string]interface{}{"display_name": "abc"},
			schema:       testResourceSchema,
			expectedName: "export_abc_1",
		},
		{
			// Non-alphanumeric or non-hyphen/underscore characters in resource name should be removed
			resource:     map[string]interface{}{"display_name": "?!@#$%^ABC:&*()def-+ghi123_"},
			schema:       testResourceSchema,
			expectedName: "export_-ABC-def--ghi123_",
		},
		{
			// Resources without display_name attribute should result in error
			resource:    map[string]interface{}{},
			schema:      testResourceSchema,
			expectError: true,
		},
		{
			// Resources with display_name attribute should result in error, because it's not part of the schema
			resource:    map[string]interface{}{"display_name": "abc"},
			schema:      map[string]*schema.Schema{},
			expectError: true,
		},
	}

	for idx, test := range testCases {
		t.Logf("Running test case %d", idx)
		result, err := generateTerraformNameFromResource(test.resource, test.schema)
		if (err != nil) != test.expectError {
			t.Logf("expect error was '%v' but got err '%v'", test.expectError, err)
			t.Fail()
		}

		if result != test.expectedName {
			t.Logf("expect generated TF name to be %s but got %s", test.expectedName, result)
			t.Fail()
		}
	}
}

// Test that correct HCL is generated from a discovered test resource
func TestUnitGetHCLString_basic(t *testing.T) {
	initResourceDiscoveryTests()
	defer cleanupResourceDiscoveryTests()
	rootResource := getRootCompartmentResource()

	results, err := findResources(nil, rootResource, compartmentTestingResourceGraph, nil)
	if err != nil {
		t.Logf("got error from findResources: %v", err)
		t.Fail()
	}

	targetResourceOcid := getTestResourceId("child", len(childrenResources)-1)
	testStringBuilder := &strings.Builder{}
	var targetResource *OCIResource
	for _, resource := range results {
		if resource.id == targetResourceOcid {
			targetResource = resource
			break
		}
	}

	if err := targetResource.getHCLString(testStringBuilder, nil); err != nil {
		t.Logf("got error '%v' when trying to get HCL string", err)
		t.Fail()
	}
	resultHcl := testStringBuilder.String()

	expectedHclResult := `resource oci_test_child export_string3_child_2 {
a_bool = "true"
a_float = "3.3"
a_int = "3"
a_list = [
"string0",
"string1",
"string2",
]
a_map = {
"key0" = "string0"
"key1" = "string1"
"key2" = "string2"
}
a_nested {
nested_bool = "true"
nested_float = "0"
nested_int = "0"
nested_string = "string0"
}
a_nested {
nested_bool = "true"
nested_float = "1.1"
nested_int = "1"
nested_string = "string1"
}
a_nested {
nested_bool = "true"
nested_float = "2.2"
nested_int = "2"
nested_string = "string2"
}
a_set = [
"string0",
"string2",
"string1",
]
a_string = "string3"
compartment_id = "string3"
parent_id = "ocid1.parent.abcdefghiklmnop.3"
}
`
	if expectedHclResult != resultHcl {
		t.Log("resulting Hcl does not match expected Hcl")
		t.Fail()
	}
}

// Test that HCL can be generated when optional or required fields are missing
func TestUnitGetHCLString_missingFields(t *testing.T) {
	initResourceDiscoveryTests()
	defer cleanupResourceDiscoveryTests()
	rootResource := getRootCompartmentResource()

	results, err := findResources(nil, rootResource, compartmentTestingResourceGraph, nil)
	if err != nil {
		t.Logf("got error from findResources: %v", err)
		t.Fail()
	}

	targetResourceOcid := getTestResourceId("child", len(childrenResources)-1)
	testStringBuilder := &strings.Builder{}
	var targetResource *OCIResource
	for _, resource := range results {
		if resource.id == targetResourceOcid {
			targetResource = resource
			break
		}
	}

	delete(targetResource.sourceAttributes, "compartment_id")
	delete(targetResource.sourceAttributes, "a_string")
	targetResource.sourceAttributes["a_map"] = nil
	if err := targetResource.getHCLString(testStringBuilder, nil); err != nil {
		t.Logf("got error '%v' when trying to get HCL string", err)
		t.Fail()
	}
	resultHcl := testStringBuilder.String()

	if !strings.Contains(resultHcl, "#compartment_id = <<Required") || !strings.Contains(resultHcl, "#a_string = <<Optional") {
		t.Logf("expected 'Required' compartment_id and 'Optional' a_string field to be commented out, but they weren't")
		t.Fail()
	}

	if strings.Contains(resultHcl, "a_map") {
		t.Logf("a_map was set to nil but it still shows up in result Hcl")
		t.Fail()
	}
}

// Test that HCL can be generated with values replaced by interpolation syntax
func TestUnitGetHCLString_interpolationMap(t *testing.T) {
	initResourceDiscoveryTests()
	defer cleanupResourceDiscoveryTests()
	rootResource := getRootCompartmentResource()

	results, err := findResources(nil, rootResource, compartmentTestingResourceGraph, nil)
	if err != nil {
		t.Logf("got error from findResources: %v", err)
		t.Fail()
	}

	targetResourceOcid := getTestResourceId("child", len(childrenResources)-1)
	testStringBuilder := &strings.Builder{}
	var targetResource *OCIResource
	for _, resource := range results {
		if resource.id == targetResourceOcid {
			targetResource = resource
			break
		}
	}

	// Test that ocids can be replaced with parent ID references
	interpolationMap := map[string]string{targetResource.parent.id: targetResource.parent.getHclReferenceIdString()}
	if err := targetResource.getHCLString(testStringBuilder, interpolationMap); err != nil {
		t.Logf("got error '%v' when trying to get HCL string", err)
		t.Fail()
	}
	resultHcl := testStringBuilder.String()

	if !strings.Contains(resultHcl, targetResource.parent.getHclReferenceIdString()) || strings.Contains(resultHcl, targetResource.parent.id) {
		t.Logf("expected hcl to replace parent ocid '%s' with '%s', but it wasn't", targetResource.parent.id, targetResource.parent.getHclReferenceIdString())
		t.Fail()
	}

	// Test that self-referencing IDs are ignored and do not show up in result hcl
	interpolationMap = map[string]string{targetResource.parent.id: targetResource.getHclReferenceIdString()}
	if err := targetResource.getHCLString(testStringBuilder, interpolationMap); err != nil {
		t.Logf("got error '%v' when trying to get HCL string", err)
		t.Fail()
	}
	resultHcl = testStringBuilder.String()

	if strings.Contains(resultHcl, targetResource.getHclReferenceIdString()) || !strings.Contains(resultHcl, targetResource.parent.id) {
		t.Logf("expected hcl to avoid cyclical reference '%s' but found one", targetResource.getHclReferenceIdString())
		t.Fail()
	}
}

func Test_getExportConfig(t *testing.T) {

	providerConfigTest(t, true, true, authAPIKeySetting, "", getExportConfig)              // ApiKey with required fields + disable auto-retries
	providerConfigTest(t, false, true, authAPIKeySetting, "", getExportConfig)             // ApiKey without required fields
	providerConfigTest(t, false, false, authInstancePrincipalSetting, "", getExportConfig) // InstancePrincipal
	providerConfigTest(t, true, false, "invalid-auth-setting", "", getExportConfig)        // Invalid auth + disable auto-retries
	configFile, keyFile, err := writeConfigFile()
	assert.Nil(t, err)
	providerConfigTest(t, true, true, authAPIKeySetting, "DEFAULT", getExportConfig)              // ApiKey with required fields + disable auto-retries
	providerConfigTest(t, false, true, authAPIKeySetting, "DEFAULT", getExportConfig)             // ApiKey without required fields
	providerConfigTest(t, false, false, authInstancePrincipalSetting, "DEFAULT", getExportConfig) // InstancePrincipal
	providerConfigTest(t, true, false, "invalid-auth-setting", "DEFAULT", getExportConfig)        // Invalid auth + disable auto-retries
	providerConfigTest(t, false, false, authAPIKeySetting, "PROFILE1", getExportConfig)           // correct profileName
	providerConfigTest(t, false, false, authAPIKeySetting, "wrongProfile", getExportConfig)       // Invalid profileName
	providerConfigTest(t, false, false, authAPIKeySetting, "PROFILE2", getExportConfig)           // correct profileName with mix and match
	providerConfigTest(t, false, false, authAPIKeySetting, "PROFILE3", getExportConfig)           // correct profileName with mix and match & env
	defer removeFile(configFile)
	defer removeFile(keyFile)
}
