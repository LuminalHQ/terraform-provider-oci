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
	listingSingularDataSourceRepresentation = map[string]interface{}{
		"listing_id": Representation{repType: Required, create: `${data.oci_marketplace_listings.test_listings.listings.0.id}`},
	}

	listingDataSourceRepresentation = map[string]interface{}{
		"category":     Representation{repType: Optional, create: []string{`category`}},
		"is_featured":  Representation{repType: Optional, create: `false`},
		"listing_id":   Representation{repType: Optional, create: `${oci_marketplace_listing.test_listing.id}`},
		"name":         Representation{repType: Optional, create: `name`},
		"package_type": Representation{repType: Optional, create: `packageType`},
		"pricing":      Representation{repType: Optional, create: []string{`pricing`}},
		"publisher_id": Representation{repType: Optional, create: `${oci_marketplace_publisher.test_publisher.id}`},
	}

	ListingResourceConfig = ``
)

func TestMarketplaceListingResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestMarketplaceListingResource_basic")
	defer httpreplay.SaveScenario()

	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_marketplace_listings.test_listings"
	singularDatasourceName := "data.oci_marketplace_listing.test_listing"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		Steps: []resource.TestStep{
			// verify datasource
			{
				Config: config +
					generateDataSourceFromRepresentationMap("oci_marketplace_listings", "test_listings", Required, Create, listingDataSourceRepresentation) +
					compartmentIdVariableStr + ListingResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet(datasourceName, "listings.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "listings.0.categories.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "listings.0.icon.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "listings.0.id"),
					resource.TestCheckResourceAttrSet(datasourceName, "listings.0.is_featured"),
					resource.TestCheckResourceAttrSet(datasourceName, "listings.0.name"),
					resource.TestCheckResourceAttrSet(datasourceName, "listings.0.package_type"),
					resource.TestCheckResourceAttrSet(datasourceName, "listings.0.publisher.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "listings.0.short_description"),
					resource.TestCheckResourceAttrSet(datasourceName, "listings.0.tagline"),
				),
			},
			// verify singular datasource
			{
				Config: config +
					generateDataSourceFromRepresentationMap("oci_marketplace_listings", "test_listings", Required, Create, listingDataSourceRepresentation) +
					generateDataSourceFromRepresentationMap("oci_marketplace_listing", "test_listing", Required, Create, listingSingularDataSourceRepresentation) +
					compartmentIdVariableStr + ListingResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(singularDatasourceName, "listing_id"),

					resource.TestCheckResourceAttrSet(singularDatasourceName, "categories.#"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "default_package_version"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "documentation_links.#"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "icon.#"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "is_featured"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "keywords"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "languages.#"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "links.#"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "long_description"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "name"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "package_type"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "publisher.#"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "regions.#"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "release_notes"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "screenshots.#"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "short_description"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "support_contacts.#"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "support_links.#"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "system_requirements"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "tagline"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "time_released"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "usage_information"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "version"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "videos.#"),
				),
			},
		},
	})
}
