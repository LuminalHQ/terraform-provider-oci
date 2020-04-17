// Copyright (c) 2017, 2019, Oracle and/or its affiliates. All rights reserved.

package oci

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	oci_apigateway "github.com/oracle/oci-go-sdk/apigateway"
	"github.com/oracle/oci-go-sdk/common"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	DeploymentRequiredOnlyResource = DeploymentResourceDependencies +
		generateResourceFromRepresentationMap("oci_apigateway_deployment", "test_deployment", Required, Create, deploymentRepresentation)

	DeploymentResourceConfig = DeploymentResourceDependencies +
		generateResourceFromRepresentationMap("oci_apigateway_deployment", "test_deployment", Optional, Update, deploymentRepresentation)

	deploymentSingularDataSourceRepresentation = map[string]interface{}{
		"deployment_id": Representation{repType: Required, create: `${oci_apigateway_deployment.test_deployment.id}`},
	}

	deploymentDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{repType: Required, create: `${var.compartment_id}`},
		"display_name":   Representation{repType: Optional, create: `displayName`, update: `displayName2`},
		"gateway_id":     Representation{repType: Optional, create: `${oci_apigateway_gateway.test_gateway.id}`},
		"state":          Representation{repType: Optional, create: `ACTIVE`},
		"filter":         RepresentationGroup{Required, deploymentDataSourceFilterRepresentation}}
	deploymentDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{repType: Required, create: `id`},
		"values": Representation{repType: Required, create: []string{`${oci_apigateway_deployment.test_deployment.id}`}},
	}

	deploymentRepresentation = map[string]interface{}{
		"compartment_id": Representation{repType: Required, create: `${var.compartment_id}`},
		"gateway_id":     Representation{repType: Required, create: `${oci_apigateway_gateway.test_gateway.id}`},
		"path_prefix":    Representation{repType: Required, create: `/v1`},
		"specification":  RepresentationGroup{Required, deploymentSpecificationRepresentation},
		"defined_tags":   Representation{repType: Optional, create: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "value")}`, update: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "updatedValue")}`},
		"display_name":   Representation{repType: Optional, create: `displayName`, update: `displayName2`},
		"freeform_tags":  Representation{repType: Optional, create: map[string]string{"Department": "Finance"}, update: map[string]string{"Department": "Accounting"}},
	}
	deploymentSpecificationRepresentation = map[string]interface{}{
		"logging_policies": RepresentationGroup{Optional, deploymentSpecificationLoggingPoliciesRepresentation},
		"request_policies": RepresentationGroup{Optional, deploymentSpecificationRequestPoliciesRepresentation},
		"routes":           RepresentationGroup{Required, deploymentSpecificationRoutesRepresentation},
	}
	deploymentSpecificationLoggingPoliciesRepresentation = map[string]interface{}{
		"access_log":    RepresentationGroup{Optional, deploymentSpecificationLoggingPoliciesAccessLogRepresentation},
		"execution_log": RepresentationGroup{Optional, deploymentSpecificationLoggingPoliciesExecutionLogRepresentation},
	}
	deploymentSpecificationRequestPoliciesRepresentation = map[string]interface{}{
		"authentication": RepresentationGroup{Optional, deploymentSpecificationRequestPoliciesAuthenticationRepresentation},
		"cors":           RepresentationGroup{Optional, deploymentSpecificationRequestPoliciesCorsRepresentation},
		"rate_limiting":  RepresentationGroup{Optional, deploymentSpecificationRequestPoliciesRateLimitingRepresentation},
	}
	deploymentSpecificationRoutesRepresentation = map[string]interface{}{
		"backend":          RepresentationGroup{Required, deploymentSpecificationRoutesBackendRepresentation},
		"path":             Representation{repType: Required, create: `/hello`, update: `/world`},
		"logging_policies": RepresentationGroup{Optional, deploymentSpecificationRoutesLoggingPoliciesRepresentation},
		"methods":          Representation{repType: Required, create: []string{`GET`}, update: []string{`GET`, `POST`}},
		"request_policies": RepresentationGroup{Optional, deploymentSpecificationRoutesRequestPoliciesRepresentation},
	}
	deploymentSpecificationLoggingPoliciesAccessLogRepresentation = map[string]interface{}{
		"is_enabled": Representation{repType: Optional, create: `false`, update: `true`},
	}
	deploymentSpecificationLoggingPoliciesExecutionLogRepresentation = map[string]interface{}{
		"is_enabled": Representation{repType: Optional, create: `false`, update: `true`},
		"log_level":  Representation{repType: Optional, create: `INFO`, update: `WARN`},
	}
	deploymentSpecificationRequestPoliciesAuthenticationRepresentation = map[string]interface{}{
		"function_id":                 Representation{repType: Required, create: `${oci_functions_function.test_function.id}`},
		"type":                        Representation{repType: Required, create: `CUSTOM_AUTHENTICATION`, update: `CUSTOM_AUTHENTICATION`},
		"token_header":                Representation{repType: Optional, create: `Authorization`, update: `Authorization`},
		"is_anonymous_access_allowed": Representation{repType: Optional, create: `false`, update: `true`},
	}
	deploymentSpecificationRequestPoliciesAuthorizeScopeRepresentation = map[string]interface{}{
		"allowed_scope": Representation{repType: Optional, create: []string{`cors`}},
	}
	deploymentSpecificationRequestPoliciesCorsRepresentation = map[string]interface{}{
		"allowed_origins":              Representation{repType: Required, create: []string{`https://www.oracle.org`}, update: []string{`*`}},
		"allowed_headers":              Representation{repType: Optional, create: []string{`*`}, update: []string{`*`, `Content-Type`}},
		"allowed_methods":              Representation{repType: Optional, create: []string{`GET`}, update: []string{`GET`, `POST`}},
		"exposed_headers":              Representation{repType: Optional, create: []string{`*`}, update: []string{`*`, `Content-Type`}},
		"is_allow_credentials_enabled": Representation{repType: Optional, create: `false`, update: `true`},
		"max_age_in_seconds":           Representation{repType: Optional, create: `600`, update: `500`},
	}
	deploymentSpecificationRequestPoliciesRateLimitingRepresentation = map[string]interface{}{
		"rate_in_requests_per_second": Representation{repType: Required, create: `10`, update: `11`},
		"rate_key":                    Representation{repType: Required, create: `CLIENT_IP`, update: `TOTAL`},
	}
	deploymentSpecificationRoutesBackendRepresentation = map[string]interface{}{
		"type": Representation{repType: Required, create: `HTTP_BACKEND`, update: `HTTP_BACKEND`},
		"url":  Representation{repType: Required, create: `https://api.weather.gov`, update: `https://www.oracle.com`},
	}
	deploymentSpecificationRoutesLoggingPoliciesRepresentation = map[string]interface{}{
		"access_log":    RepresentationGroup{Optional, deploymentSpecificationRoutesLoggingPoliciesAccessLogRepresentation},
		"execution_log": RepresentationGroup{Optional, deploymentSpecificationRoutesLoggingPoliciesExecutionLogRepresentation},
	}
	deploymentSpecificationRoutesRequestPoliciesRepresentation = map[string]interface{}{
		"authorization": RepresentationGroup{Optional, deploymentSpecificationRoutesRequestPoliciesAuthorizationRepresentation},
		"cors":          RepresentationGroup{Optional, deploymentSpecificationRoutesRequestPoliciesCorsRepresentation},
	}
	deploymentSpecificationRoutesBackendHeadersRepresentation = map[string]interface{}{
		"name":  Representation{repType: Optional, create: `name`, update: `name2`},
		"value": Representation{repType: Optional, create: `value`, update: `value2`},
	}
	deploymentSpecificationRoutesLoggingPoliciesAccessLogRepresentation = map[string]interface{}{
		"is_enabled": Representation{repType: Optional, create: `false`, update: `true`},
	}
	deploymentSpecificationRoutesLoggingPoliciesExecutionLogRepresentation = map[string]interface{}{
		"is_enabled": Representation{repType: Optional, create: `false`, update: `true`},
		"log_level":  Representation{repType: Optional, create: `INFO`, update: `WARN`},
	}

	deploymentSpecificationRoutesRequestPoliciesAuthorizationRepresentation = map[string]interface{}{
		"type": Representation{repType: Optional, create: `AUTHENTICATION_ONLY`, update: `ANONYMOUS`},
	}
	deploymentSpecificationRoutesRequestPoliciesCorsRepresentation = map[string]interface{}{
		"allowed_origins":              Representation{repType: Required, create: []string{`*`}, update: []string{`*`}},
		"allowed_headers":              Representation{repType: Optional, create: []string{`*`}, update: []string{`*`, `Content-Type`}},
		"allowed_methods":              Representation{repType: Optional, create: []string{`GET`}, update: []string{`GET`, `POST`}},
		"exposed_headers":              Representation{repType: Optional, create: []string{`*`}, update: []string{`*`, `Content-Type`}},
		"is_allow_credentials_enabled": Representation{repType: Optional, create: `false`, update: `true`},
		"max_age_in_seconds":           Representation{repType: Optional, create: `600`, update: `500`},
	}

	DeploymentResourceDependencies = generateResourceFromRepresentationMap("oci_apigateway_gateway", "test_gateway", Required, Create, gatewayRepresentation) +
		generateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Required, Create, subnetRegionalRepresentation) +
		generateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnRepresentation) +
		DefinedTagsDependencies
)

func TestApigatewayDeploymentResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestApigatewayDeploymentResource_basic")
	defer httpreplay.SaveScenario()

	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	compartmentIdU := getEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)

	image := getEnvSettingWithBlankDefault("image")
	imageVariableStr := fmt.Sprintf("variable \"image\" { default = \"%s\" }\n", image)

	resourceName := "oci_apigateway_deployment.test_deployment"
	datasourceName := "data.oci_apigateway_deployments.test_deployments"
	singularDatasourceName := "data.oci_apigateway_deployment.test_deployment"

	var resId, resId2 string

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		CheckDestroy: testAccCheckApigatewayDeploymentDestroy,
		Steps: []resource.TestStep{
			// verify create
			{
				Config: config + compartmentIdVariableStr + DeploymentResourceDependencies +
					generateResourceFromRepresentationMap("oci_apigateway_deployment", "test_deployment", Required, Create, deploymentRepresentation),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "path_prefix", "/v1"),
					resource.TestCheckResourceAttr(resourceName, "specification.#", "1"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// delete before next create
			{
				Config: config + compartmentIdVariableStr + DeploymentResourceDependencies,
			},
			// verify create with optionals
			{
				Config: config + compartmentIdVariableStr + imageVariableStr + DeploymentResourceDependencies +
					generateResourceFromRepresentationMap("oci_functions_application", "test_application", Required, Create, applicationRepresentation) +
					generateResourceFromRepresentationMap("oci_functions_function", "test_function", Required, Create, functionRepresentation) +
					generateResourceFromRepresentationMap("oci_apigateway_deployment", "test_deployment", Optional, Create, deploymentRepresentation),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "path_prefix", "/v1"),
					resource.TestCheckResourceAttr(resourceName, "specification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.access_log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.access_log.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.execution_log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.execution_log.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.execution_log.0.log_level", "INFO"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.authentication.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "specification.0.request_policies.0.authentication.0.function_id"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.authentication.0.token_header", "Authorization"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.authentication.0.type", "CUSTOM_AUTHENTICATION"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.allowed_methods.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.exposed_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.is_allow_credentials_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.max_age_in_seconds", "600"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.rate_limiting.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.rate_limiting.0.rate_in_requests_per_second", "10"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.rate_limiting.0.rate_key", "CLIENT_IP"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.backend.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "specification.0.routes.0.backend.0.connect_timeout_in_seconds"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.backend.0.is_ssl_verify_disabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "specification.0.routes.0.backend.0.read_timeout_in_seconds"),
					resource.TestCheckResourceAttrSet(resourceName, "specification.0.routes.0.backend.0.send_timeout_in_seconds"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.backend.0.type", "HTTP_BACKEND"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.backend.0.url", "https://api.weather.gov"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.access_log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.access_log.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.execution_log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.execution_log.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.execution_log.0.log_level", "INFO"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.methods.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.path", "/hello"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.authorization.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.authorization.0.type", "AUTHENTICATION_ONLY"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.allowed_methods.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.exposed_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.is_allow_credentials_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.max_age_in_seconds", "600"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						if isEnableExportCompartment, _ := strconv.ParseBool(getEnvSettingWithDefault("enable_export_compartment", "false")); isEnableExportCompartment {
							if errExport := testExportCompartmentWithResourceName(&resId, &compartmentId, resourceName); errExport != nil {
								return errExport
							}
						}
						return err
					},
				),
			},

			// verify update to the compartment (the compartment will be switched back in the next step)
			{
				Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + imageVariableStr + DeploymentResourceDependencies +
					generateResourceFromRepresentationMap("oci_functions_application", "test_application", Required, Create, applicationRepresentation) +
					generateResourceFromRepresentationMap("oci_functions_function", "test_function", Required, Create, functionRepresentation) +
					generateResourceFromRepresentationMap("oci_apigateway_deployment", "test_deployment", Optional, Create,
						representationCopyWithNewProperties(deploymentRepresentation, map[string]interface{}{
							"compartment_id": Representation{repType: Required, create: `${var.compartment_id_for_update}`},
						})),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
					resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "path_prefix", "/v1"),
					resource.TestCheckResourceAttr(resourceName, "specification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.access_log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.access_log.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.execution_log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.execution_log.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.execution_log.0.log_level", "INFO"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.authentication.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "specification.0.request_policies.0.authentication.0.function_id"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.authentication.0.token_header", "Authorization"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.authentication.0.type", "CUSTOM_AUTHENTICATION"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.allowed_methods.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.exposed_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.is_allow_credentials_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.max_age_in_seconds", "600"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.rate_limiting.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.rate_limiting.0.rate_in_requests_per_second", "10"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.rate_limiting.0.rate_key", "CLIENT_IP"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.backend.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "specification.0.routes.0.backend.0.connect_timeout_in_seconds"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.backend.0.is_ssl_verify_disabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "specification.0.routes.0.backend.0.read_timeout_in_seconds"),
					resource.TestCheckResourceAttrSet(resourceName, "specification.0.routes.0.backend.0.send_timeout_in_seconds"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.backend.0.type", "HTTP_BACKEND"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.backend.0.url", "https://api.weather.gov"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.access_log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.access_log.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.execution_log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.execution_log.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.execution_log.0.log_level", "INFO"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.methods.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.path", "/hello"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.authorization.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.authorization.0.type", "AUTHENTICATION_ONLY"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.allowed_methods.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.exposed_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.is_allow_credentials_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.max_age_in_seconds", "600"),

					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId != resId2 {
							return fmt.Errorf("resource recreated when it was supposed to be updated")
						}
						return err
					},
				),
			},

			// verify updates to updatable parameters
			{
				Config: config + compartmentIdVariableStr + imageVariableStr + DeploymentResourceDependencies +
					generateResourceFromRepresentationMap("oci_functions_application", "test_application", Required, Create, applicationRepresentation) +
					generateResourceFromRepresentationMap("oci_functions_function", "test_function", Required, Create, functionRepresentation) +
					generateResourceFromRepresentationMap("oci_apigateway_deployment", "test_deployment", Optional, Update, deploymentRepresentation),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "path_prefix", "/v1"),
					resource.TestCheckResourceAttr(resourceName, "specification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.access_log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.access_log.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.execution_log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.execution_log.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.logging_policies.0.execution_log.0.log_level", "WARN"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.authentication.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "specification.0.request_policies.0.authentication.0.function_id"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.authentication.0.token_header", "Authorization"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.authentication.0.type", "CUSTOM_AUTHENTICATION"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.allowed_headers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.allowed_methods.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.exposed_headers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.is_allow_credentials_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.cors.0.max_age_in_seconds", "500"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.rate_limiting.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.rate_limiting.0.rate_in_requests_per_second", "11"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.request_policies.0.rate_limiting.0.rate_key", "TOTAL"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.backend.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "specification.0.routes.0.backend.0.connect_timeout_in_seconds"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.backend.0.is_ssl_verify_disabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "specification.0.routes.0.backend.0.read_timeout_in_seconds"),
					resource.TestCheckResourceAttrSet(resourceName, "specification.0.routes.0.backend.0.send_timeout_in_seconds"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.backend.0.type", "HTTP_BACKEND"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.backend.0.url", "https://www.oracle.com"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.access_log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.access_log.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.execution_log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.execution_log.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.logging_policies.0.execution_log.0.log_level", "WARN"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.methods.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.path", "/world"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.authorization.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.authorization.0.type", "ANONYMOUS"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.allowed_headers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.allowed_methods.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.exposed_headers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.is_allow_credentials_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "specification.0.routes.0.request_policies.0.cors.0.max_age_in_seconds", "500"),

					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId != resId2 {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						return err
					},
				),
			},
			// verify datasource
			{
				Config: config + imageVariableStr +
					generateResourceFromRepresentationMap("oci_functions_application", "test_application", Required, Create, applicationRepresentation) +
					generateResourceFromRepresentationMap("oci_functions_function", "test_function", Required, Create, functionRepresentation) +
					generateDataSourceFromRepresentationMap("oci_apigateway_deployments", "test_deployments", Optional, Update, deploymentDataSourceRepresentation) +
					compartmentIdVariableStr + DeploymentResourceDependencies +
					generateResourceFromRepresentationMap("oci_apigateway_deployment", "test_deployment", Optional, Update, deploymentRepresentation),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttrSet(datasourceName, "gateway_id"),
					resource.TestCheckResourceAttr(datasourceName, "state", "ACTIVE"),

					resource.TestCheckResourceAttr(datasourceName, "deployment_collection.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "deployment_collection.0.endpoint"),
					resource.TestCheckResourceAttrSet(datasourceName, "deployment_collection.0.gateway_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "deployment_collection.0.id"),
				),
			},
			// verify singular datasource
			{
				Config: config + imageVariableStr +
					generateResourceFromRepresentationMap("oci_functions_application", "test_application", Required, Create, applicationRepresentation) +
					generateResourceFromRepresentationMap("oci_functions_function", "test_function", Required, Create, functionRepresentation) +
					generateDataSourceFromRepresentationMap("oci_apigateway_deployment", "test_deployment", Required, Create, deploymentSingularDataSourceRepresentation) +
					compartmentIdVariableStr + DeploymentResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(singularDatasourceName, "deployment_id"),

					resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(singularDatasourceName, "defined_tags.%", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "endpoint"),
					resource.TestCheckResourceAttr(singularDatasourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
					resource.TestCheckResourceAttr(singularDatasourceName, "path_prefix", "/v1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.logging_policies.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.logging_policies.0.access_log.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.logging_policies.0.access_log.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.logging_policies.0.execution_log.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.logging_policies.0.execution_log.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.logging_policies.0.execution_log.0.log_level", "WARN"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.request_policies.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.request_policies.0.authentication.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.request_policies.0.authentication.0.token_header", "Authorization"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.request_policies.0.authentication.0.type", "CUSTOM_AUTHENTICATION"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.request_policies.0.cors.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.request_policies.0.cors.0.allowed_headers.#", "2"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.request_policies.0.cors.0.allowed_methods.#", "2"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.request_policies.0.cors.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.request_policies.0.cors.0.exposed_headers.#", "2"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.request_policies.0.cors.0.is_allow_credentials_enabled", "true"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.request_policies.0.cors.0.max_age_in_seconds", "500"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.request_policies.0.rate_limiting.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.request_policies.0.rate_limiting.0.rate_in_requests_per_second", "11"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.request_policies.0.rate_limiting.0.rate_key", "TOTAL"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.backend.#", "1"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "specification.0.routes.0.backend.0.connect_timeout_in_seconds"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.backend.0.is_ssl_verify_disabled", "false"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "specification.0.routes.0.backend.0.read_timeout_in_seconds"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "specification.0.routes.0.backend.0.send_timeout_in_seconds"),

					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.backend.0.type", "HTTP_BACKEND"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.backend.0.url", "https://www.oracle.com"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.logging_policies.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.logging_policies.0.access_log.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.logging_policies.0.access_log.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.logging_policies.0.execution_log.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.logging_policies.0.execution_log.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.logging_policies.0.execution_log.0.log_level", "WARN"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.methods.#", "2"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.path", "/world"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.request_policies.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.request_policies.0.cors.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.request_policies.0.cors.0.allowed_headers.#", "2"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.request_policies.0.cors.0.allowed_methods.#", "2"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.request_policies.0.cors.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.request_policies.0.cors.0.exposed_headers.#", "2"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.request_policies.0.cors.0.is_allow_credentials_enabled", "true"),
					resource.TestCheckResourceAttr(singularDatasourceName, "specification.0.routes.0.request_policies.0.cors.0.max_age_in_seconds", "500"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "time_updated"),
				),
			},
			// remove singular datasource from previous step so that it doesn't conflict with import tests
			{
				Config: config + compartmentIdVariableStr + imageVariableStr + DeploymentResourceConfig +
					generateResourceFromRepresentationMap("oci_functions_application", "test_application", Required, Create, applicationRepresentation) +
					generateResourceFromRepresentationMap("oci_functions_function", "test_function", Required, Create, functionRepresentation),
			},
			// verify resource import
			{
				Config:            config,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"lifecycle_details",
				},
				ResourceName: resourceName,
			},
		},
	})
}

func testAccCheckApigatewayDeploymentDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).deploymentClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_apigateway_deployment" {
			noResourceFound = false
			request := oci_apigateway.GetDeploymentRequest{}

			tmp := rs.Primary.ID
			request.DeploymentId = &tmp

			request.RequestMetadata.RetryPolicy = getRetryPolicy(true, "apigateway")

			response, err := client.GetDeployment(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_apigateway.DeploymentLifecycleStateDeleted): true,
				}
				if _, ok := deletedLifecycleStates[string(response.LifecycleState)]; !ok {
					//resource lifecycle state is not in expected deleted lifecycle states.
					return fmt.Errorf("resource lifecycle state: %s is not in expected deleted lifecycle states", response.LifecycleState)
				}
				//resource lifecycle state is in expected deleted lifecycle states. continue with next one.
				continue
			}

			//Verify that exception is for '404 not found'.
			if failure, isServiceError := common.IsServiceError(err); !isServiceError || failure.GetHTTPStatusCode() != 404 {
				return err
			}
		}
	}
	if noResourceFound {
		return fmt.Errorf("at least one resource was expected from the state file, but could not be found")
	}

	return nil
}

func init() {
	if DependencyGraph == nil {
		initDependencyGraph()
	}
	if !inSweeperExcludeList("ApigatewayDeployment") {
		resource.AddTestSweepers("ApigatewayDeployment", &resource.Sweeper{
			Name:         "ApigatewayDeployment",
			Dependencies: DependencyGraph["deployment"],
			F:            sweepApigatewayDeploymentResource,
		})
	}
}

func sweepApigatewayDeploymentResource(compartment string) error {
	deploymentClient := GetTestClients(&schema.ResourceData{}).deploymentClient
	deploymentIds, err := getDeploymentIds(compartment)
	if err != nil {
		return err
	}
	for _, deploymentId := range deploymentIds {
		if ok := SweeperDefaultResourceId[deploymentId]; !ok {
			deleteDeploymentRequest := oci_apigateway.DeleteDeploymentRequest{}

			deleteDeploymentRequest.DeploymentId = &deploymentId

			deleteDeploymentRequest.RequestMetadata.RetryPolicy = getRetryPolicy(true, "apigateway")
			_, error := deploymentClient.DeleteDeployment(context.Background(), deleteDeploymentRequest)
			if error != nil {
				fmt.Printf("Error deleting Deployment %s %s, It is possible that the resource is already deleted. Please verify manually \n", deploymentId, error)
				continue
			}
			waitTillCondition(testAccProvider, &deploymentId, deploymentSweepWaitCondition, time.Duration(3*time.Minute),
				deploymentSweepResponseFetchOperation, "apigateway", true)
		}
	}
	return nil
}

func getDeploymentIds(compartment string) ([]string, error) {
	ids := getResourceIdsToSweep(compartment, "DeploymentId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	deploymentClient := GetTestClients(&schema.ResourceData{}).deploymentClient

	listDeploymentsRequest := oci_apigateway.ListDeploymentsRequest{}
	listDeploymentsRequest.CompartmentId = &compartmentId
	listDeploymentsRequest.LifecycleState = oci_apigateway.DeploymentLifecycleStateActive
	listDeploymentsResponse, err := deploymentClient.ListDeployments(context.Background(), listDeploymentsRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting Deployment list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, deployment := range listDeploymentsResponse.Items {
		id := *deployment.Id
		resourceIds = append(resourceIds, id)
		addResourceIdToSweeperResourceIdMap(compartmentId, "DeploymentId", id)
	}
	return resourceIds, nil
}

func deploymentSweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if deploymentResponse, ok := response.Response.(oci_apigateway.GetDeploymentResponse); ok {
		return deploymentResponse.LifecycleState != oci_apigateway.DeploymentLifecycleStateDeleted
	}
	return false
}

func deploymentSweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.deploymentClient.GetDeployment(context.Background(), oci_apigateway.GetDeploymentRequest{
		DeploymentId: resourceId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	})
	return err
}
