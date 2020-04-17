// Copyright (c) 2017, 2019, Oracle and/or its affiliates. All rights reserved.

package oci

import (
	"fmt"
	"log"
	"testing"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	VolumeBackupCopyResourceDependencies = VolumeBackupResourceDependencies + generateResourceFromRepresentationMap("oci_kms_key", "test_key", Required, Create, keyRepresentation)
)

func TestResourceCoreVolumeBackup_copy(t *testing.T) {
	httpreplay.SetScenario("TestResourceCoreVolumeBackup_copy")
	defer httpreplay.SaveScenario()

	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceNameCopy := "oci_core_volume_backup.test_volume_backup_copy"
	datasourceName := "data.oci_core_volume_backups.test_volume_backups"

	if getEnvSettingWithBlankDefault("source_region") == "" {
		t.Skip("Skipping TestCoreVolumeBackupResource_copy test because there is no source region specified")
	}

	err := createSourceVolumeBackupToCopy()
	if err != nil {
		t.Fatalf("Unable to create source Volume and VolumeBackup to copy. Error: %v", err)
	}

	volumeBackupSourceDetailsRepresentation = map[string]interface{}{
		"volume_backup_id": Representation{repType: Required, create: volumeBackupId},
		"region":           Representation{repType: Required, create: getEnvSettingWithBlankDefault("source_region")},
		"kms_key_id":       Representation{repType: Required, create: `${oci_kms_key.test_key.id}`},
	}
	volumeBackupWithSourceDetailsRepresentation = getUpdatedRepresentationCopy("source_details", RepresentationGroup{Required, volumeBackupSourceDetailsRepresentation}, volumeBackupWithSourceDetailsRepresentation)

	var resId string

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		CheckDestroy: testAccCheckCoreVolumeBackupDestroy,
		Steps: []resource.TestStep{
			// verify create
			{
				Config: config +
					compartmentIdVariableStr + VolumeBackupCopyResourceDependencies +
					generateResourceFromRepresentationMap("oci_core_volume_backup", "test_volume_backup_copy", Required, Create, volumeBackupWithSourceDetailsRepresentation),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameCopy, "volume_id"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceNameCopy, "id")
						return err
					},
				),
			},

			// delete before next create
			{
				Config: config + compartmentIdVariableStr + VolumeBackupCopyResourceDependencies,
			},
			// verify create from the backup with optionals
			{
				Config: config +
					compartmentIdVariableStr + VolumeBackupCopyResourceDependencies +
					generateResourceFromRepresentationMap("oci_core_volume_backup", "test_volume_backup_copy", Optional, Create, volumeBackupWithSourceDetailsRepresentation),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameCopy, "compartment_id"),
					resource.TestCheckResourceAttr(resourceNameCopy, "display_name", "displayName"),
					resource.TestCheckResourceAttrSet(resourceNameCopy, "id"),
					resource.TestCheckResourceAttrSet(resourceNameCopy, "state"),
					resource.TestCheckResourceAttrSet(resourceNameCopy, "time_created"),
					resource.TestCheckResourceAttrSet(resourceNameCopy, "type"),
					resource.TestCheckResourceAttrSet(resourceNameCopy, "volume_id"),
					resource.TestCheckResourceAttrSet(resourceNameCopy, "source_volume_backup_id"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceNameCopy, "id")
						return err
					},
				),
			},
			// verify updates to updatable parameters
			{
				Config: config +
					compartmentIdVariableStr + VolumeBackupCopyResourceDependencies +
					generateResourceFromRepresentationMap("oci_core_volume_backup", "test_volume_backup_copy", Optional, Update, volumeBackupWithSourceDetailsRepresentation),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameCopy, "compartment_id"),
					resource.TestCheckResourceAttr(resourceNameCopy, "display_name", "displayName2"),
					resource.TestCheckResourceAttrSet(resourceNameCopy, "id"),
					resource.TestCheckResourceAttrSet(resourceNameCopy, "state"),
					resource.TestCheckResourceAttrSet(resourceNameCopy, "time_created"),
					resource.TestCheckResourceAttrSet(resourceNameCopy, "type"),
					resource.TestCheckResourceAttrSet(resourceNameCopy, "volume_id"),
					resource.TestCheckResourceAttrSet(resourceNameCopy, "source_volume_backup_id"),

					func(s *terraform.State) (err error) {
						resId2, err := fromInstanceState(s, resourceNameCopy, "id")
						if resId != resId2 {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						return err
					},
				),
			},
			// verify datasource
			{
				Config: config +
					generateDataSourceFromRepresentationMap("oci_core_volume_backups", "test_volume_backups", Optional, Update, volumeBackupFromSourceDataSourceRepresentation) +
					compartmentIdVariableStr + VolumeBackupCopyResourceDependencies +
					generateResourceFromRepresentationMap("oci_core_volume_backup", "test_volume_backup_copy", Optional, Update, volumeBackupWithSourceDetailsRepresentation),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttr(datasourceName, "state", "AVAILABLE"),
					resource.TestCheckResourceAttrSet(datasourceName, "source_volume_backup_id"),

					resource.TestCheckResourceAttr(datasourceName, "volume_backups.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "volume_backups.0.compartment_id"),
					resource.TestCheckResourceAttr(datasourceName, "volume_backups.0.display_name", "displayName2"),
					resource.TestCheckResourceAttrSet(datasourceName, "volume_backups.0.id"),
					resource.TestCheckResourceAttrSet(datasourceName, "volume_backups.0.kms_key_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "volume_backups.0.state"),
					resource.TestCheckResourceAttrSet(datasourceName, "volume_backups.0.time_created"),
					resource.TestCheckResourceAttrSet(datasourceName, "volume_backups.0.type"),
					resource.TestCheckResourceAttrSet(datasourceName, "volume_backups.0.volume_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "volume_backups.0.source_volume_backup_id"),
				),
			},
			// verify resource import
			{
				Config:            config,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"source_details",
				},
				ResourceName: resourceNameCopy,
			},
		},
	})
}

func createSourceVolumeBackupToCopy() error {
	sourceRegion := getEnvSettingWithBlankDefault("source_region")

	var err error
	volumeId, err = createVolumeInRegion(GetTestClients(&schema.ResourceData{}), sourceRegion)
	if err != nil {
		log.Printf("[WARN] failed to createVolumeInRegion with the error %v", err)
		return err
	}

	volumeBackupId, err = createVolumeBackupInRegion(GetTestClients(&schema.ResourceData{}), sourceRegion, &volumeId)
	if err != nil {
		log.Printf("[WARN] failed to createVolumeBackupInRegion with the error %v", err)
		return err
	}

	return nil
}

func deleteSourceVolumeBackupToCopy() {
	sourceRegion := getEnvSettingWithBlankDefault("source_region")

	var err error
	err = deleteVolumeBackupInRegion(GetTestClients(&schema.ResourceData{}), sourceRegion, volumeBackupId)
	if err != nil {
		log.Printf("[WARN] failed to deleteVolumeBackupInRegion with error %v", err)
	}

	err = deleteVolumeInRegion(GetTestClients(&schema.ResourceData{}), sourceRegion, volumeId)
	if err != nil {
		log.Printf("[WARN] failed to deleteVolumeInRegion with error %v", err)
	}
}
