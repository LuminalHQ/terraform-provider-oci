// Copyright (c) 2017, 2019, Oracle and/or its affiliates. All rights reserved.

package oci

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/oracle/oci-go-sdk/common"
	oci_containerengine "github.com/oracle/oci-go-sdk/containerengine"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	NodePoolRequiredOnlyResource = NodePoolResourceDependencies +
		generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool", Required, Create, nodePoolRepresentation)

	NodePoolResourceConfig = NodePoolResourceDependencies +
		generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool", Optional, Update, nodePoolRepresentation)

	nodePoolSingularDataSourceRepresentation = map[string]interface{}{
		"node_pool_id": Representation{repType: Required, create: `${oci_containerengine_node_pool.test_node_pool.id}`},
	}

	nodePoolSingularDataSourceRepresentationForImageId = map[string]interface{}{
		"node_pool_id": Representation{repType: Required, create: `${oci_containerengine_node_pool.test_node_pool_imageId.id}`},
	}

	nodePoolSingularDataSourceRepresentationForNodeSourceDetails = map[string]interface{}{
		"node_pool_id": Representation{repType: Required, create: `${oci_containerengine_node_pool.test_node_pool_node_source_details.id}`},
	}

	nodePoolDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{repType: Required, create: `${var.compartment_id}`},
		"cluster_id":     Representation{repType: Optional, create: `${oci_containerengine_cluster.test_cluster.id}`},
		"name":           Representation{repType: Optional, create: `name`, update: `name2`},
		"filter":         RepresentationGroup{Required, nodePoolDataSourceFilterRepresentation}}
	nodePoolDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{repType: Required, create: `id`},
		"values": Representation{repType: Required, create: []string{`${oci_containerengine_node_pool.test_node_pool.id}`}},
	}

	nodePoolDataSourceRepresentationForImageId = map[string]interface{}{
		"compartment_id": Representation{repType: Required, create: `${var.compartment_id}`},
		"cluster_id":     Representation{repType: Optional, create: `${oci_containerengine_cluster.test_cluster.id}`},
		"name":           Representation{repType: Optional, create: `name`, update: `name2`},
		"filter":         RepresentationGroup{Required, nodePoolDataSourceFilterRepresentationForImageId}}
	nodePoolDataSourceFilterRepresentationForImageId = map[string]interface{}{
		"name":   Representation{repType: Required, create: `id`},
		"values": Representation{repType: Required, create: []string{`${oci_containerengine_node_pool.test_node_pool_imageId.id}`}},
	}

	nodePoolDataSourceRepresentationForNodeSourceDetails = map[string]interface{}{
		"compartment_id": Representation{repType: Required, create: `${var.compartment_id}`},
		"cluster_id":     Representation{repType: Optional, create: `${oci_containerengine_cluster.test_cluster.id}`},
		"name":           Representation{repType: Optional, create: `name`, update: `name2`},
		"filter":         RepresentationGroup{Required, nodePoolDataSourceFilterRepresentationForNodeSourceDetails}}
	nodePoolDataSourceFilterRepresentationForNodeSourceDetails = map[string]interface{}{
		"name":   Representation{repType: Required, create: `id`},
		"values": Representation{repType: Required, create: []string{`${oci_containerengine_node_pool.test_node_pool_node_source_details.id}`}},
	}

	nodePoolRepresentation = map[string]interface{}{
		"cluster_id":          Representation{repType: Required, create: `${oci_containerengine_cluster.test_cluster.id}`},
		"compartment_id":      Representation{repType: Required, create: `${var.compartment_id}`},
		"kubernetes_version":  Representation{repType: Required, create: `${data.oci_containerengine_node_pool_option.test_node_pool_option.kubernetes_versions.0}`},
		"name":                Representation{repType: Required, create: `name`, update: `name2`},
		"node_shape":          Representation{repType: Required, create: `VM.Standard2.1`},
		"subnet_ids":          Representation{repType: Required, create: []string{`${oci_core_subnet.nodePool_Subnet_1.id}`, `${oci_core_subnet.nodePool_Subnet_2.id}`}},
		"initial_node_labels": RepresentationGroup{Optional, nodePoolInitialNodeLabelsRepresentation},
		"node_image_name":     Representation{repType: Required, create: `${data.oci_containerengine_node_pool_option.test_node_pool_option.images.0}`},
		"node_metadata":       Representation{repType: Optional, create: map[string]string{"nodeMetadata": "nodeMetadata"}, update: map[string]string{"nodeMetadata2": "nodeMetadata2"}},
		"quantity_per_subnet": Representation{repType: Optional, create: `1`, update: `2`},
		"ssh_public_key":      Representation{repType: Optional, create: `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample`},
	}
	nodePoolRepresentationForImageId = map[string]interface{}{
		"cluster_id":          Representation{repType: Required, create: `${oci_containerengine_cluster.test_cluster.id}`},
		"compartment_id":      Representation{repType: Required, create: `${var.compartment_id}`},
		"kubernetes_version":  Representation{repType: Required, create: `${data.oci_containerengine_node_pool_option.test_node_pool_option.kubernetes_versions.0}`},
		"name":                Representation{repType: Required, create: `name`, update: `name2`},
		"node_image_id":       Representation{repType: Required, create: `${data.oci_containerengine_node_pool_option.test_node_pool_option.sources.0.image_id}`},
		"node_shape":          Representation{repType: Required, create: `VM.Standard2.1`},
		"subnet_ids":          Representation{repType: Required, create: []string{`${oci_core_subnet.nodePool_Subnet_1.id}`, `${oci_core_subnet.nodePool_Subnet_2.id}`}},
		"initial_node_labels": RepresentationGroup{Optional, nodePoolInitialNodeLabelsRepresentation},
		"quantity_per_subnet": Representation{repType: Optional, create: `1`, update: `2`},
		"ssh_public_key":      Representation{repType: Optional, create: `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample`},
	}

	nodePoolRepresentationForNodeSourceDetails = map[string]interface{}{
		"cluster_id":          Representation{repType: Required, create: `${oci_containerengine_cluster.test_cluster.id}`},
		"compartment_id":      Representation{repType: Required, create: `${var.compartment_id}`},
		"kubernetes_version":  Representation{repType: Required, create: `${data.oci_containerengine_node_pool_option.test_node_pool_option.kubernetes_versions.0}`},
		"name":                Representation{repType: Required, create: `name`, update: `name2`},
		"node_shape":          Representation{repType: Required, create: `VM.Standard2.1`},
		"subnet_ids":          Representation{repType: Required, create: []string{`${oci_core_subnet.nodePool_Subnet_1.id}`, `${oci_core_subnet.nodePool_Subnet_2.id}`}},
		"initial_node_labels": RepresentationGroup{Optional, nodePoolInitialNodeLabelsRepresentation},
		"node_source_details": RepresentationGroup{Required, nodePoolNodeSourceDetailsRepresentation},
		"quantity_per_subnet": Representation{repType: Optional, create: `1`, update: `2`},
		"ssh_public_key":      Representation{repType: Optional, create: `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample`},
	}
	nodePoolInitialNodeLabelsRepresentation = map[string]interface{}{
		"key":   Representation{repType: Optional, create: `key`, update: `key2`},
		"value": Representation{repType: Optional, create: `value`, update: `value2`},
	}
	nodePoolNodeSourceDetailsRepresentation = map[string]interface{}{
		"image_id":    Representation{repType: Required, create: `${data.oci_containerengine_node_pool_option.test_node_pool_option.sources.0.image_id}`},
		"source_type": Representation{repType: Required, create: `${data.oci_containerengine_node_pool_option.test_node_pool_option.sources.0.source_type}`},
	}

	NodePoolResourceDependencies = generateDataSourceFromRepresentationMap("oci_containerengine_node_pool_option", "test_node_pool_option", Required, Create, nodePoolOptionSingularDataSourceRepresentation) +
		generateResourceFromRepresentationMap("oci_core_subnet", "nodePool_Subnet_1", Required, Create, representationCopyWithNewProperties(subnetRepresentation, map[string]interface{}{"availability_domain": Representation{repType: Required, create: `${lower("${data.oci_identity_availability_domains.test_availability_domains.availability_domains.0.name}")}`}, "cidr_block": Representation{repType: Required, create: `10.0.22.0/24`}, "dns_label": Representation{repType: Required, create: `nodepool1`}})) +
		generateResourceFromRepresentationMap("oci_core_subnet", "nodePool_Subnet_2", Required, Create, representationCopyWithNewProperties(subnetRepresentation, map[string]interface{}{"availability_domain": Representation{repType: Required, create: `${lower("${data.oci_identity_availability_domains.test_availability_domains.availability_domains.0.name}")}`}, "cidr_block": Representation{repType: Required, create: `10.0.23.0/24`}, "dns_label": Representation{repType: Required, create: `nodepool2`}})) +
		generateResourceFromRepresentationMap("oci_containerengine_cluster", "test_cluster", Required, Create, clusterRepresentation) +
		generateResourceFromRepresentationMap("oci_core_subnet", "clusterSubnet_1", Required, Create, representationCopyWithNewProperties(subnetRepresentation, map[string]interface{}{"availability_domain": Representation{repType: Required, create: `${lower("${data.oci_identity_availability_domains.test_availability_domains.availability_domains.0.name}")}`}, "cidr_block": Representation{repType: Required, create: `10.0.20.0/24`}, "dns_label": Representation{repType: Required, create: `cluster1`}})) +
		generateResourceFromRepresentationMap("oci_core_subnet", "clusterSubnet_2", Required, Create, representationCopyWithNewProperties(subnetRepresentation, map[string]interface{}{"availability_domain": Representation{repType: Required, create: `${lower("${data.oci_identity_availability_domains.test_availability_domains.availability_domains.0.name}")}`}, "cidr_block": Representation{repType: Required, create: `10.0.21.0/24`}, "dns_label": Representation{repType: Required, create: `cluster2`}})) +
		AvailabilityDomainConfig +
		generateDataSourceFromRepresentationMap("oci_containerengine_cluster_option", "test_cluster_option", Required, Create, clusterOptionSingularDataSourceRepresentation) +
		generateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, representationCopyWithNewProperties(vcnRepresentation, map[string]interface{}{
			"dns_label": Representation{repType: Required, create: `dnslabel`},
		}))
)

func TestContainerengineNodePoolResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestContainerengineNodePoolResource_basic")
	defer httpreplay.SaveScenario()

	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_containerengine_node_pool.test_node_pool"
	resourceNameForImageId := "oci_containerengine_node_pool.test_node_pool_imageId"
	resourceNameForNodeSourceDetails := "oci_containerengine_node_pool.test_node_pool_node_source_details"

	datasourceName := "data.oci_containerengine_node_pools.test_node_pools"
	datasourceNameForImageId := "data.oci_containerengine_node_pools.test_node_pools_imageId"
	datasourceNameForNodeSourceDetails := "data.oci_containerengine_node_pools.test_node_pools_node_source_details"

	singularDatasourceName := "data.oci_containerengine_node_pool.test_node_pool"
	singularDatasourceNameForImageId := "data.oci_containerengine_node_pool.test_node_pool_imageId"
	singularDatasourceNameForNodeSourceDetails := "data.oci_containerengine_node_pool.test_node_pool_node_source_details"

	var resId, resId2, resIdCreatedWithImageId, resId2CreatedWithImageId, resIdCreatedWithNodeSourceDetails, resId2CreatedWithNodeSourceDetails string

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		CheckDestroy: testAccCheckContainerengineNodePoolDestroy,
		Steps: []resource.TestStep{
			// verify create
			{
				Config: config + compartmentIdVariableStr + NodePoolResourceDependencies +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool", Required, Create, nodePoolRepresentation) +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool_imageId", Required, Create, nodePoolRepresentationForImageId) +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool_node_source_details", Required, Create, nodePoolRepresentationForNodeSourceDetails),

				Check: resource.ComposeAggregateTestCheckFunc(
					//Asserting Resource created with Image Name
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(resourceName, "kubernetes_version"),
					resource.TestCheckResourceAttr(resourceName, "name", "name"),
					resource.TestCheckResourceAttr(resourceName, "node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "2"),
					//Asserting Resource created with Image Id
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "cluster_id"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "kubernetes_version"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "name", "name"),
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "node_image_name"),
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "node_image_id"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "subnet_ids.#", "2"),
					//Asserting Resource created with Node Source Details
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "cluster_id"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "kubernetes_version"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "name", "name"),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "node_image_name"),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "node_image_id"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "subnet_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "node_source_details.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "node_source_details.0.image_id"),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "node_source_details.0.source_type"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						resIdCreatedWithImageId, err = fromInstanceState(s, resourceNameForImageId, "id")
						resIdCreatedWithNodeSourceDetails, err = fromInstanceState(s, resourceNameForNodeSourceDetails, "id")
						return err
					},
				),
			},

			// delete before next create
			{
				Config: config + compartmentIdVariableStr + NodePoolResourceDependencies,
			},
			// verify create with optionals
			{
				Config: config + compartmentIdVariableStr + NodePoolResourceDependencies +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool", Optional, Create, nodePoolRepresentation) +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool_imageId", Optional, Create, nodePoolRepresentationForImageId) +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool_node_source_details", Optional, Create, nodePoolRepresentationForNodeSourceDetails),
				Check: resource.ComposeAggregateTestCheckFunc(
					//Asserting Resource created with Image Name
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.0.key", "key"),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.0.value", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "kubernetes_version"),
					resource.TestCheckResourceAttr(resourceName, "name", "name"),
					resource.TestCheckResourceAttrSet(resourceName, "node_image_id"),
					resource.TestCheckResourceAttrSet(resourceName, "node_image_name"),
					resource.TestCheckResourceAttr(resourceName, "node_metadata.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(resourceName, "quantity_per_subnet", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "2"),
					//Asserting Resource created with Image Id
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "cluster_id"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceNameForImageId, "initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "initial_node_labels.0.key", "key"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "initial_node_labels.0.value", "value"),
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "kubernetes_version"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "name", "name"),
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "node_image_name"),
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "node_image_id"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "quantity_per_subnet", "1"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "subnet_ids.#", "2"),
					//Asserting Resource created with Node Source Details
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "cluster_id"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "initial_node_labels.0.key", "key"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "initial_node_labels.0.value", "value"),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "kubernetes_version"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "name", "name"),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "node_image_id"),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "node_image_name"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "node_source_details.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "node_source_details.0.image_id"),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "node_source_details.0.source_type"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "quantity_per_subnet", "1"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "subnet_ids.#", "2"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						if isEnableExportCompartment, _ := strconv.ParseBool(getEnvSettingWithDefault("enable_export_compartment", "false")); isEnableExportCompartment {
							if errExport := testExportCompartmentWithResourceName(&resId, &compartmentId, resourceName); errExport != nil {
								return errExport
							}
						}
						resIdCreatedWithImageId, err = fromInstanceState(s, resourceNameForImageId, "id")
						resIdCreatedWithNodeSourceDetails, err = fromInstanceState(s, resourceNameForNodeSourceDetails, "id")
						return err
					},
				),
			},

			// verify updates to updatable parameters
			{
				Config: config + compartmentIdVariableStr + NodePoolResourceDependencies +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool", Optional, Update, getUpdatedRepresentationCopy("node_metadata", Representation{repType: Optional, update: map[string]string{"nodeMetadata": "nodeMetadata"}}, nodePoolRepresentation)) +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool_imageId", Optional, Update, nodePoolRepresentationForImageId) +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool_node_source_details", Optional, Update, nodePoolRepresentationForNodeSourceDetails),
				Check: resource.ComposeAggregateTestCheckFunc(
					//Asserting Resource created with Image Name
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.0.key", "key2"),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.0.value", "value2"),
					resource.TestCheckResourceAttrSet(resourceName, "kubernetes_version"),
					resource.TestCheckResourceAttr(resourceName, "name", "name2"),
					resource.TestCheckResourceAttrSet(resourceName, "node_image_id"),
					resource.TestCheckResourceAttrSet(resourceName, "node_image_name"),
					resource.TestCheckResourceAttr(resourceName, "node_metadata.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(resourceName, "quantity_per_subnet", "2"),
					resource.TestCheckResourceAttr(resourceName, "ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "2"),
					//Asserting Resource created with Image Id
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "cluster_id"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceNameForImageId, "initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "initial_node_labels.0.key", "key2"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "initial_node_labels.0.value", "value2"),
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "kubernetes_version"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "name", "name2"),
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "node_image_name"),
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "node_image_id"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "quantity_per_subnet", "2"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "subnet_ids.#", "2"),
					//Asserting Resource created with Node Source Details
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "cluster_id"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "initial_node_labels.0.key", "key2"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "initial_node_labels.0.value", "value2"),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "kubernetes_version"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "name", "name2"),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "node_image_id"),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "node_image_name"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "node_source_details.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "node_source_details.0.image_id"),
					resource.TestCheckResourceAttrSet(resourceNameForNodeSourceDetails, "node_source_details.0.source_type"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "quantity_per_subnet", "2"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(resourceNameForNodeSourceDetails, "subnet_ids.#", "2"),
					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId != resId2 {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						resId2CreatedWithImageId, err = fromInstanceState(s, resourceNameForImageId, "id")
						if resIdCreatedWithImageId != resId2CreatedWithImageId {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						resId2CreatedWithNodeSourceDetails, err = fromInstanceState(s, resourceNameForNodeSourceDetails, "id")
						if resIdCreatedWithNodeSourceDetails != resId2CreatedWithNodeSourceDetails {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						return err
					},
				),
			},
			// verify updates to node_metadata should create new resource
			{
				Config: config + compartmentIdVariableStr + NodePoolResourceDependencies +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool", Optional, Update, nodePoolRepresentation) +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool_imageId", Optional, Update, nodePoolRepresentationForImageId) +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool_node_source_details", Optional, Update, nodePoolRepresentationForNodeSourceDetails),
				Check: resource.ComposeAggregateTestCheckFunc(
					//Asserting Resource created with Image Name
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.0.key", "key2"),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.0.value", "value2"),
					resource.TestCheckResourceAttrSet(resourceName, "kubernetes_version"),
					resource.TestCheckResourceAttr(resourceName, "name", "name2"),
					resource.TestCheckResourceAttrSet(resourceName, "node_image_name"),
					resource.TestCheckResourceAttr(resourceName, "node_metadata.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(resourceName, "quantity_per_subnet", "2"),
					resource.TestCheckResourceAttr(resourceName, "ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "2"),
					//Asserting Resource created with Image Id
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "cluster_id"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceNameForImageId, "initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "initial_node_labels.0.key", "key2"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "initial_node_labels.0.value", "value2"),
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "kubernetes_version"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "name", "name2"),
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "node_image_name"),
					resource.TestCheckResourceAttrSet(resourceNameForImageId, "node_image_id"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "quantity_per_subnet", "2"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(resourceNameForImageId, "subnet_ids.#", "2"),

					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId == resId2 {
							return fmt.Errorf("Resource updated when it was supposed to be recreated.")
						}
						resId2CreatedWithImageId, err = fromInstanceState(s, resourceNameForImageId, "id")
						if resIdCreatedWithImageId != resId2CreatedWithImageId {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						resId2CreatedWithNodeSourceDetails, err = fromInstanceState(s, resourceNameForNodeSourceDetails, "id")
						if resIdCreatedWithNodeSourceDetails != resId2CreatedWithNodeSourceDetails {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						return err
					},
				),
			},
			// verify datasource
			{
				Config: config +
					generateDataSourceFromRepresentationMap("oci_containerengine_node_pools", "test_node_pools", Optional, Update, nodePoolDataSourceRepresentation) +
					generateDataSourceFromRepresentationMap("oci_containerengine_node_pools", "test_node_pools_imageId", Optional, Update, nodePoolDataSourceRepresentationForImageId) +
					generateDataSourceFromRepresentationMap("oci_containerengine_node_pools", "test_node_pools_node_source_details", Optional, Update, nodePoolDataSourceRepresentationForNodeSourceDetails) +
					compartmentIdVariableStr + NodePoolResourceDependencies +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool", Optional, Update, nodePoolRepresentation) +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool_imageId", Optional, Update, nodePoolRepresentationForImageId) +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool_node_source_details", Optional, Update, nodePoolRepresentationForNodeSourceDetails),
				Check: resource.ComposeAggregateTestCheckFunc(
					//Asserting Datasource for NodePool created with Image Name
					resource.TestCheckResourceAttrSet(datasourceName, "cluster_id"),
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "name", "name2"),

					resource.TestCheckResourceAttr(datasourceName, "node_pools.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "node_pools.0.cluster_id"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(datasourceName, "node_pools.0.id"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.initial_node_labels.0.key", "key2"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.initial_node_labels.0.value", "value2"),
					resource.TestCheckResourceAttrSet(datasourceName, "node_pools.0.kubernetes_version"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.name", "name2"),
					resource.TestCheckResourceAttrSet(datasourceName, "node_pools.0.node_image_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "node_pools.0.node_image_name"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.quantity_per_subnet", "2"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.subnet_ids.#", "2"),
					//Asserting Datasource for NodePool created with Image Id
					resource.TestCheckResourceAttrSet(datasourceNameForImageId, "cluster_id"),
					resource.TestCheckResourceAttr(datasourceNameForImageId, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceNameForImageId, "name", "name2"),

					resource.TestCheckResourceAttr(datasourceNameForImageId, "node_pools.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameForImageId, "node_pools.0.cluster_id"),
					resource.TestCheckResourceAttr(datasourceNameForImageId, "node_pools.0.compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceNameForImageId, "node_pools.0.initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameForImageId, "node_pools.0.initial_node_labels.0.key", "key2"),
					resource.TestCheckResourceAttr(datasourceNameForImageId, "node_pools.0.initial_node_labels.0.value", "value2"),
					resource.TestCheckResourceAttrSet(datasourceNameForImageId, "node_pools.0.kubernetes_version"),
					resource.TestCheckResourceAttr(datasourceNameForImageId, "node_pools.0.name", "name2"),
					resource.TestCheckResourceAttrSet(datasourceNameForImageId, "node_pools.0.node_image_name"),
					resource.TestCheckResourceAttrSet(datasourceNameForImageId, "node_pools.0.node_image_id"),
					resource.TestCheckResourceAttr(datasourceNameForImageId, "node_pools.0.node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(datasourceNameForImageId, "node_pools.0.quantity_per_subnet", "2"),
					resource.TestCheckResourceAttr(datasourceNameForImageId, "node_pools.0.ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(datasourceNameForImageId, "node_pools.0.subnet_ids.#", "2"),
					//Asserting Datasource for NodePool created with Node Source Details
					resource.TestCheckResourceAttrSet(datasourceNameForNodeSourceDetails, "cluster_id"),
					resource.TestCheckResourceAttr(datasourceNameForNodeSourceDetails, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceNameForNodeSourceDetails, "name", "name2"),

					resource.TestCheckResourceAttr(datasourceNameForNodeSourceDetails, "node_pools.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameForNodeSourceDetails, "node_pools.0.cluster_id"),
					resource.TestCheckResourceAttr(datasourceNameForNodeSourceDetails, "node_pools.0.compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(datasourceNameForNodeSourceDetails, "node_pools.0.id"),
					resource.TestCheckResourceAttr(datasourceNameForNodeSourceDetails, "node_pools.0.initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameForNodeSourceDetails, "node_pools.0.initial_node_labels.0.key", "key2"),
					resource.TestCheckResourceAttr(datasourceNameForNodeSourceDetails, "node_pools.0.initial_node_labels.0.value", "value2"),
					resource.TestCheckResourceAttrSet(datasourceNameForNodeSourceDetails, "node_pools.0.kubernetes_version"),
					resource.TestCheckResourceAttr(datasourceNameForNodeSourceDetails, "node_pools.0.name", "name2"),
					resource.TestCheckResourceAttrSet(datasourceNameForNodeSourceDetails, "node_pools.0.node_image_id"),
					resource.TestCheckResourceAttrSet(datasourceNameForNodeSourceDetails, "node_pools.0.node_image_name"),
					resource.TestCheckResourceAttr(datasourceNameForNodeSourceDetails, "node_pools.0.node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(datasourceNameForNodeSourceDetails, "node_pools.0.node_source.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameForNodeSourceDetails, "node_pools.0.quantity_per_subnet", "2"),
					resource.TestCheckResourceAttr(datasourceNameForNodeSourceDetails, "node_pools.0.ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(datasourceNameForNodeSourceDetails, "node_pools.0.subnet_ids.#", "2"),
				),
			},
			// verify singular datasource
			{
				Config: config +
					generateDataSourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool", Required, Create, nodePoolSingularDataSourceRepresentation) +
					generateDataSourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool_imageId", Required, Create, nodePoolSingularDataSourceRepresentationForImageId) +
					generateDataSourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool_node_source_details", Required, Create, nodePoolSingularDataSourceRepresentationForNodeSourceDetails) +
					compartmentIdVariableStr + NodePoolResourceConfig +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool_imageId", Optional, Update, nodePoolRepresentationForImageId) +
					generateResourceFromRepresentationMap("oci_containerengine_node_pool", "test_node_pool_node_source_details", Optional, Update, nodePoolRepresentationForNodeSourceDetails),
				Check: resource.ComposeAggregateTestCheckFunc(
					//Asserting Singular Datasource for NodePool created with Image Name
					resource.TestCheckResourceAttrSet(singularDatasourceName, "cluster_id"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "node_pool_id"),

					resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
					resource.TestCheckResourceAttr(singularDatasourceName, "initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "initial_node_labels.0.key", "key2"),
					resource.TestCheckResourceAttr(singularDatasourceName, "initial_node_labels.0.value", "value2"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "kubernetes_version"),
					resource.TestCheckResourceAttr(singularDatasourceName, "name", "name2"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "node_image_id"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "node_image_name"),
					resource.TestCheckResourceAttr(singularDatasourceName, "node_metadata.%", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "quantity_per_subnet", "2"),
					resource.TestCheckResourceAttr(singularDatasourceName, "ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(singularDatasourceName, "subnet_ids.#", "2"),
					//Asserting Singular Datasource for NodePool created with Image Id
					resource.TestCheckResourceAttrSet(singularDatasourceNameForImageId, "cluster_id"),
					resource.TestCheckResourceAttrSet(singularDatasourceNameForImageId, "node_pool_id"),

					resource.TestCheckResourceAttr(singularDatasourceNameForImageId, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(singularDatasourceNameForImageId, "initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceNameForImageId, "initial_node_labels.0.key", "key2"),
					resource.TestCheckResourceAttr(singularDatasourceNameForImageId, "initial_node_labels.0.value", "value2"),
					resource.TestCheckResourceAttrSet(singularDatasourceNameForImageId, "kubernetes_version"),
					resource.TestCheckResourceAttr(singularDatasourceNameForImageId, "name", "name2"),
					resource.TestCheckResourceAttrSet(singularDatasourceNameForImageId, "node_image_name"),
					resource.TestCheckResourceAttrSet(singularDatasourceNameForImageId, "node_image_id"),
					resource.TestCheckResourceAttr(singularDatasourceNameForImageId, "node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(singularDatasourceNameForImageId, "quantity_per_subnet", "2"),
					resource.TestCheckResourceAttr(singularDatasourceNameForImageId, "ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(singularDatasourceNameForImageId, "subnet_ids.#", "2"),
					//Asserting Singular Datasource for NodePool created with Image Name
					resource.TestCheckResourceAttrSet(singularDatasourceNameForNodeSourceDetails, "cluster_id"),
					resource.TestCheckResourceAttrSet(singularDatasourceNameForNodeSourceDetails, "node_pool_id"),

					resource.TestCheckResourceAttr(singularDatasourceNameForNodeSourceDetails, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(singularDatasourceNameForNodeSourceDetails, "id"),
					resource.TestCheckResourceAttr(singularDatasourceNameForNodeSourceDetails, "initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceNameForNodeSourceDetails, "initial_node_labels.0.key", "key2"),
					resource.TestCheckResourceAttr(singularDatasourceNameForNodeSourceDetails, "initial_node_labels.0.value", "value2"),
					resource.TestCheckResourceAttrSet(singularDatasourceNameForNodeSourceDetails, "kubernetes_version"),
					resource.TestCheckResourceAttr(singularDatasourceNameForNodeSourceDetails, "name", "name2"),
					resource.TestCheckResourceAttrSet(singularDatasourceNameForNodeSourceDetails, "node_image_id"),
					resource.TestCheckResourceAttrSet(singularDatasourceNameForNodeSourceDetails, "node_image_name"),
					resource.TestCheckResourceAttr(singularDatasourceNameForNodeSourceDetails, "node_shape", "VM.Standard2.1"),
					resource.TestCheckResourceAttr(singularDatasourceNameForNodeSourceDetails, "node_source.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceNameForNodeSourceDetails, "quantity_per_subnet", "2"),
					resource.TestCheckResourceAttr(singularDatasourceNameForNodeSourceDetails, "ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(singularDatasourceNameForNodeSourceDetails, "subnet_ids.#", "2"),
				),
			},
			// remove singular datasource from previous step so that it doesn't conflict with import tests
			{
				Config: config + compartmentIdVariableStr + NodePoolResourceConfig,
			},
			// verify resource import
			{
				Config:                  config + compartmentIdVariableStr + NodePoolResourceConfig,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
				ResourceName:            resourceName,
			},
		},
	})
}

func testAccCheckContainerengineNodePoolDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).containerEngineClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_containerengine_node_pool" {
			noResourceFound = false
			request := oci_containerengine.GetNodePoolRequest{}

			tmp := rs.Primary.ID
			request.NodePoolId = &tmp

			request.RequestMetadata.RetryPolicy = getRetryPolicy(true, "containerengine")

			_, err := client.GetNodePool(context.Background(), request)

			if err == nil {
				return fmt.Errorf("resource still exists")
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
	if !inSweeperExcludeList("ContainerengineNodePool") {
		resource.AddTestSweepers("ContainerengineNodePool", &resource.Sweeper{
			Name:         "ContainerengineNodePool",
			Dependencies: DependencyGraph["nodePool"],
			F:            sweepContainerengineNodePoolResource,
		})
	}
}

func sweepContainerengineNodePoolResource(compartment string) error {
	containerEngineClient := GetTestClients(&schema.ResourceData{}).containerEngineClient
	nodePoolIds, err := getNodePoolIds(compartment)
	if err != nil {
		return err
	}
	for _, nodePoolId := range nodePoolIds {
		if ok := SweeperDefaultResourceId[nodePoolId]; !ok {
			deleteNodePoolRequest := oci_containerengine.DeleteNodePoolRequest{}

			deleteNodePoolRequest.NodePoolId = &nodePoolId

			deleteNodePoolRequest.RequestMetadata.RetryPolicy = getRetryPolicy(true, "containerengine")
			_, error := containerEngineClient.DeleteNodePool(context.Background(), deleteNodePoolRequest)
			if error != nil {
				fmt.Printf("Error deleting NodePool %s %s, It is possible that the resource is already deleted. Please verify manually \n", nodePoolId, error)
				continue
			}
		}
	}
	return nil
}

func getNodePoolIds(compartment string) ([]string, error) {
	ids := getResourceIdsToSweep(compartment, "NodePoolId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	containerEngineClient := GetTestClients(&schema.ResourceData{}).containerEngineClient

	listNodePoolsRequest := oci_containerengine.ListNodePoolsRequest{}
	listNodePoolsRequest.CompartmentId = &compartmentId
	listNodePoolsResponse, err := containerEngineClient.ListNodePools(context.Background(), listNodePoolsRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting NodePool list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, nodePool := range listNodePoolsResponse.Items {
		id := *nodePool.Id
		resourceIds = append(resourceIds, id)
		addResourceIdToSweeperResourceIdMap(compartmentId, "NodePoolId", id)
	}
	return resourceIds, nil
}
