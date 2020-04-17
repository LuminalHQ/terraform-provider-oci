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
	listingPackageAgreementManagementRepresentation = map[string]interface{}{
		"agreement_id":    Representation{repType: Required, create: `${data.oci_marketplace_listing_package_agreements.test_listing_package_agreements.agreements.0.id}`},
		"listing_id":      Representation{repType: Required, create: `${data.oci_marketplace_listing.test_listing.id}`},
		"package_version": Representation{repType: Required, create: `${data.oci_marketplace_listing.test_listing.default_package_version}`},
	}

	listingPackageAgreementDataSourceRepresentation = map[string]interface{}{
		"listing_id":      Representation{repType: Required, create: `${data.oci_marketplace_listing.test_listing.id}`},
		"package_version": Representation{repType: Required, create: `${data.oci_marketplace_listing.test_listing.default_package_version}`},
	}

	ListingPackageAgreementResourceConfig = generateDataSourceFromRepresentationMap("oci_marketplace_listing", "test_listing", Required, Create, listingSingularDataSourceRepresentation) +
		generateDataSourceFromRepresentationMap("oci_marketplace_listings", "test_listings", Required, Create, listingDataSourceRepresentation)
)

func TestMarketplaceListingPackageAgreementResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestMarketplaceListingPackageAgreementResource_basic")
	defer httpreplay.SaveScenario()

	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_marketplace_listing_package_agreements.test_listing_package_agreements"
	resourceName := "oci_marketplace_listing_package_agreement.test_listing_package_agreement"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		Steps: []resource.TestStep{
			// verify resource
			{
				Config: config +
					generateResourceFromRepresentationMap("oci_marketplace_listing_package_agreement", "test_listing_package_agreement", Required, Create, listingPackageAgreementManagementRepresentation) +
					generateDataSourceFromRepresentationMap("oci_marketplace_listing_package_agreements", "test_listing_package_agreements", Required, Create, listingPackageAgreementDataSourceRepresentation) +
					compartmentIdVariableStr + ListingPackageAgreementResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "agreement_id"),
					resource.TestCheckResourceAttrSet(resourceName, "listing_id"),
					resource.TestCheckResourceAttrSet(resourceName, "package_version"),

					resource.TestCheckResourceAttrSet(resourceName, "content_url"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "prompt"),
					resource.TestCheckResourceAttrSet(resourceName, "signature"),
				),
			},
			// verify datasource
			{
				Config: config +
					generateDataSourceFromRepresentationMap("oci_marketplace_listing_package_agreements", "test_listing_package_agreements", Required, Create, listingPackageAgreementDataSourceRepresentation) +
					compartmentIdVariableStr + ListingPackageAgreementResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "listing_id"),

					resource.TestCheckResourceAttrSet(datasourceName, "agreements.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "agreements.0.author"),
					resource.TestCheckResourceAttrSet(datasourceName, "agreements.0.content_url"),
					resource.TestCheckResourceAttrSet(datasourceName, "agreements.0.id"),
					resource.TestCheckResourceAttrSet(datasourceName, "agreements.0.prompt"),
				),
			},
		},
	})
}
