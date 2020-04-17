// Copyright (c) 2017, 2019, Oracle and/or its affiliates. All rights reserved.

package oci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	PurgeCacheRequiredOnlyResource = PurgeCacheResourceDependencies +
		generateResourceFromRepresentationMap("oci_waas_purge_cache", "test_purge_cache", Required, Create, purgeCacheRepresentation)

	purgeCacheRepresentation = map[string]interface{}{
		"waas_policy_id": Representation{repType: Required, create: `${oci_waas_waas_policy.test_scenario_waas_policy.id}`},
		"resources":      Representation{repType: Optional, create: []string{`/about`, `/home`}},
	}

	PurgeCacheResourceDependencies = WaasPolicyResourceCachingOnlyConfig
)

func TestWaasPurgeCacheResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestWaasPurgeCacheResource_basic")
	defer httpreplay.SaveScenario()

	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_waas_purge_cache.test_purge_cache"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		Steps: []resource.TestStep{
			// verify purge select resources
			{
				Config: config + compartmentIdVariableStr + PurgeCacheResourceDependencies +
					generateResourceFromRepresentationMap("oci_waas_purge_cache", "test_purge_cache", Optional, Create, purgeCacheRepresentation),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "resources.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "waas_policy_id"),
				),
			},

			// delete before next create
			{
				Config: config + compartmentIdVariableStr + PurgeCacheResourceDependencies,
			},
			// verify purge all resources
			{
				Config: config + compartmentIdVariableStr + PurgeCacheResourceDependencies +
					generateResourceFromRepresentationMap("oci_waas_purge_cache", "test_purge_cache", Required, Create, purgeCacheRepresentation),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "waas_policy_id"),
				),
			},
		},
	})
}
