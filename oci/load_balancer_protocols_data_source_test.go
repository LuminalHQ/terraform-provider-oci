// Copyright (c) 2017, 2019, Oracle and/or its affiliates. All rights reserved.

package oci

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDatasourceLoadBalancerProtocols_basic(t *testing.T) {
	httpreplay.SetScenario("TestAccDatasourceLoadBalancerProtocols_basic")
	defer httpreplay.SaveScenario()
	providers := testAccProviders
	config := legacyTestProviderConfig() + `
	data "oci_load_balancer_protocols" "t" {
		compartment_id = "${var.compartment_id}"
	}`

	resourceName := "data.oci_load_balancer_protocols.t"

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		PreventPostDestroyRefresh: true,
		Providers:                 providers,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "protocols.#"),
					resource.TestCheckResourceAttrSet(resourceName, "protocols.0.name"),
				),
			},
		},
	})
}
