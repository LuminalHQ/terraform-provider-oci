// Copyright (c) 2017, 2019, Oracle and/or its affiliates. All rights reserved.

package oci

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	oci_marketplace "github.com/oracle/oci-go-sdk/marketplace"
)

func MarketplaceListingPackageDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readSingularMarketplaceListingPackage,
		Schema: map[string]*schema.Schema{
			"listing_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"package_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Computed
			"app_catalog_listing_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_catalog_listing_resource_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"package_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pricing": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required

						// Optional

						// Computed
						"currency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rate": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required

						// Optional

						// Computed
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"countries": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// Required

									// Optional

									// Computed
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"variables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required

						// Optional

						// Computed
						"data_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hint_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_mandatory": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func readSingularMarketplaceListingPackage(d *schema.ResourceData, m interface{}) error {
	sync := &MarketplaceListingPackageDataSourceCrud{}
	sync.D = d
	sync.Client = m.(*OracleClients).marketplaceClient

	return ReadResource(sync)
}

type MarketplaceListingPackageDataSourceCrud struct {
	D      *schema.ResourceData
	Client *oci_marketplace.MarketplaceClient
	Res    *oci_marketplace.ListingPackage
}

func (s *MarketplaceListingPackageDataSourceCrud) VoidState() {
	s.D.SetId("")
}

func (s *MarketplaceListingPackageDataSourceCrud) Get() error {
	request := oci_marketplace.GetPackageRequest{}

	if listingId, ok := s.D.GetOkExists("listing_id"); ok {
		tmp := listingId.(string)
		request.ListingId = &tmp
	}

	if packageVersion, ok := s.D.GetOkExists("package_version"); ok {
		tmp := packageVersion.(string)
		request.PackageVersion = &tmp
	}

	request.RequestMetadata.RetryPolicy = getRetryPolicy(false, "marketplace")

	response, err := s.Client.GetPackage(context.Background(), request)
	if err != nil {
		return err
	}

	s.Res = &response.ListingPackage
	return nil
}

func (s *MarketplaceListingPackageDataSourceCrud) SetData() error {
	if s.Res == nil {
		return nil
	}
	s.D.SetId(GenerateDataSourceID())
	switch v := (*s.Res).(type) {
	case oci_marketplace.ImageListingPackage:

		if v.AppCatalogListingId != nil {
			s.D.Set("app_catalog_listing_id", v.AppCatalogListingId)
		}

		if v.AppCatalogListingResourceVersion != nil {
			s.D.Set("app_catalog_listing_resource_version", v.AppCatalogListingResourceVersion)
		}

		if v.Description != nil {
			s.D.Set("description", v.Description)
		}

		if v.ListingId != nil {
			s.D.Set("Listing_id", v.ListingId)
		}

		s.D.Set("package_type", oci_marketplace.PackageTypeEnumImage)

		if v.Pricing != nil {
			s.D.Set("pricing", []interface{}{PricingModelToMap(v.Pricing)})
		}

		if v.Regions != nil {
			regions := []interface{}{}
			for _, item := range v.Regions {
				regions = append(regions, RegionToMap(item))
			}
			s.D.Set("regions", regions)
		}

		if v.ResourceId != nil {
			s.D.Set("resource_id", v.ResourceId)
		}

		if v.TimeCreated != nil {
			s.D.Set("time_created", v.TimeCreated.String())
		}

		if v.Version != nil {
			s.D.Set("version", v.Version)
		}
	case oci_marketplace.OrchestrationListingPackage:
		if v.Description != nil {
			s.D.Set("description", v.Description)
		}

		if v.ListingId != nil {
			s.D.Set("Listing_id", v.ListingId)
		}

		s.D.Set("package_type", oci_marketplace.PackageTypeEnumOrchestration)

		if v.Pricing != nil {
			s.D.Set("pricing", []interface{}{PricingModelToMap(v.Pricing)})
		}

		if v.ResourceId != nil {
			s.D.Set("resource_id", v.ResourceId)
		}

		if v.ResourceLink != nil {
			s.D.Set("resource_link", v.ResourceLink)
		}

		if v.TimeCreated != nil {
			s.D.Set("time_created", v.TimeCreated.String())
		}

		if v.Variables != nil {
			variables := []interface{}{}
			for _, item := range v.Variables {
				variables = append(variables, OrchestrationVariableToMap(item))
			}
			s.D.Set("variables", variables)
		}

		if v.Version != nil {
			s.D.Set("version", v.Version)
		}
	default:
		log.Printf("[WARN] Received 'ListingPackage' of unknown type %v", *s.Res)
		return nil
	}
	return nil
}

func OrchestrationVariableToMap(obj oci_marketplace.OrchestrationVariable) map[string]interface{} {
	result := map[string]interface{}{}

	if obj.DefaultValue != nil {
		result["default_value"] = string(*obj.DefaultValue)
	}

	if obj.Description != nil {
		result["description"] = string(*obj.Description)
	}

	if obj.HintMessage != nil {
		result["hint_message"] = string(*obj.HintMessage)
	}

	if obj.IsMandatory != nil {
		result["is_mandatory"] = bool(*obj.IsMandatory)
	}

	if obj.Name != nil {
		result["name"] = string(*obj.Name)
	}

	return result
}

func PricingModelToMap(obj *oci_marketplace.PricingModel) map[string]interface{} {
	result := map[string]interface{}{}

	if obj.Rate != nil {
		result["rate"] = float32(*obj.Rate)
	}
	return result
}