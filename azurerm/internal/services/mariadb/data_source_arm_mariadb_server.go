package mariadb

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mariadb/mgmt/2018-06-01/mariadb"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMariaDbServer() *schema.Resource {
	return &schema.Resource{
		Read:   resourceArmMariaDbServerRead,

		Timeouts: &schema.ResourceTimeout{
			Read:   schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{3,50}$"),
					"MariaDB server name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
				),
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"sku_name": {
				Type:          schema.TypeString,
				Computed:      true, // remove in 2.0
				ConflictsWith: []string{"sku"},
				ValidateFunc: validation.StringInSlice([]string{
					"B_Gen5_1",
					"B_Gen5_2",
					"GP_Gen5_2",
					"GP_Gen5_4",
					"GP_Gen5_8",
					"GP_Gen5_16",
					"GP_Gen5_32",
					"MO_Gen5_2",
					"MO_Gen5_4",
					"MO_Gen5_8",
					"MO_Gen5_16",
				}, false),
			},

			// remove in 2.0
			"sku": {
				Type:          schema.TypeList,
				Computed:      true,
				ConflictsWith: []string{"sku_name"},
				Deprecated:    "This property has been deprecated in favour of the 'sku_name' property and will be removed in version 2.0 of the provider",
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"B_Gen5_1",
								"B_Gen5_2",
								"GP_Gen5_2",
								"GP_Gen5_4",
								"GP_Gen5_8",
								"GP_Gen5_16",
								"GP_Gen5_32",
								"MO_Gen5_2",
								"MO_Gen5_4",
								"MO_Gen5_8",
								"MO_Gen5_16",
							}, false),
						},

						"capacity": {
							Type:     schema.TypeInt,
							Computed: true,
							ValidateFunc: validate.IntInSlice([]int{
								1,
								2,
								4,
								8,
								16,
								32,
							}),
						},

						"tier": {
							Type:     schema.TypeString,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(mariadb.Basic),
								string(mariadb.GeneralPurpose),
								string(mariadb.MemoryOptimized),
							}, false),
						},

						"family": {
							Type:     schema.TypeString,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Gen5",
							}, false),
						},
					},
				},
			},

			"administrator_login": {
				Type:         schema.TypeString,
				Computed:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"administrator_login_password": {
				Type:         schema.TypeString,
				Computed:     true,
				Sensitive:    true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"10.2",
					"10.3",
				}, false),
			},

			"storage_profile": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_mb": {
							Type:         schema.TypeInt,
							Computed:     true,
							ValidateFunc: validate.IntBetweenAndDivisibleBy(5120, 4096000, 1024),
						},

						"backup_retention_days": {
							Type:         schema.TypeInt,
							Computed:     true,
							ValidateFunc: validation.IntBetween(7, 35),
						},

						"geo_redundant_backup": {
							Type:     schema.TypeString,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(mariadb.Enabled),
								string(mariadb.Disabled),
							}, false),
						},

						"auto_grow": {
							Type:     schema.TypeString,
							Computed: true,
							Default:  string(mariadb.StorageAutogrowEnabled),
							ValidateFunc: validation.StringInSlice([]string{
								string(mariadb.StorageAutogrowEnabled),
								string(mariadb.StorageAutogrowDisabled),
							}, false),
						},
					},
				},
			},

			"ssl_enforcement": {
				Type:     schema.TypeString,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(mariadb.SslEnforcementEnumDisabled),
					string(mariadb.SslEnforcementEnumEnabled),
				}, false),
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}


func dataSourceArmMariaDbServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.ServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["servers"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] MariaDB Server %q was not found (Resource Group %q)", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure MariaDB Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if properties := resp.ServerProperties; properties != nil {
		d.Set("administrator_login", properties.AdministratorLogin)
		d.Set("version", string(properties.Version))
		d.Set("ssl_enforcement", string(properties.SslEnforcement))
		// Computed
		d.Set("fqdn", properties.FullyQualifiedDomainName)

		if err := d.Set("storage_profile", flattenMariaDbStorageProfile(properties.StorageProfile)); err != nil {
			return fmt.Errorf("Error setting `storage_profile`: %+v", err)
		}
	}

	if err := d.Set("sku", flattenMariaDbServerSku(resp.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func expandServerSkuName(skuName string) (*mariadb.Sku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 3 {
		return nil, fmt.Errorf("sku_name (%s) has the worng number of parts (%d) after splitting on _", skuName, len(parts))
	}

	var tier mariadb.SkuTier
	switch parts[0] {
	case "B":
		tier = mariadb.Basic
	case "GP":
		tier = mariadb.GeneralPurpose
	case "MO":
		tier = mariadb.MemoryOptimized
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, parts[0])
	}

	capacity, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("cannot convert skuname %s capcity %s to int", skuName, parts[2])
	}

	return &mariadb.Sku{
		Name:     utils.String(skuName),
		Tier:     tier,
		Capacity: utils.Int32(int32(capacity)),
		Family:   utils.String(parts[1]),
	}, nil
}

func expandAzureRmMariaDbServerSku(d *schema.ResourceData) *mariadb.Sku {
	skus := d.Get("sku").([]interface{})
	sku := skus[0].(map[string]interface{})

	name := sku["name"].(string)
	capacity := sku["capacity"].(int)
	tier := sku["tier"].(string)
	family := sku["family"].(string)

	return &mariadb.Sku{
		Name:     utils.String(name),
		Tier:     mariadb.SkuTier(tier),
		Capacity: utils.Int32(int32(capacity)),
		Family:   utils.String(family),
	}
}

func expandAzureRmMariaDbStorageProfile(d *schema.ResourceData) *mariadb.StorageProfile {
	storageprofiles := d.Get("storage_profile").([]interface{})
	storageprofile := storageprofiles[0].(map[string]interface{})

	backupRetentionDays := storageprofile["backup_retention_days"].(int)
	geoRedundantBackup := storageprofile["geo_redundant_backup"].(string)
	storageMB := storageprofile["storage_mb"].(int)
	autoGrow := storageprofile["auto_grow"].(string)

	return &mariadb.StorageProfile{
		BackupRetentionDays: utils.Int32(int32(backupRetentionDays)),
		GeoRedundantBackup:  mariadb.GeoRedundantBackup(geoRedundantBackup),
		StorageMB:           utils.Int32(int32(storageMB)),
		StorageAutogrow:     mariadb.StorageAutogrow(autoGrow),
	}
}

func flattenMariaDbServerSku(sku *mariadb.Sku) []interface{} {
	values := map[string]interface{}{}

	if sku == nil {
		return []interface{}{}
	}

	if name := sku.Name; name != nil {
		values["name"] = *name
	}

	if capacity := sku.Capacity; capacity != nil {
		values["capacity"] = *capacity
	}

	values["tier"] = string(sku.Tier)

	if family := sku.Family; family != nil {
		values["family"] = *family
	}

	return []interface{}{values}
}

func flattenMariaDbStorageProfile(storage *mariadb.StorageProfile) []interface{} {
	values := map[string]interface{}{}

	if storage == nil {
		return []interface{}{}
	}

	if storageMB := storage.StorageMB; storageMB != nil {
		values["storage_mb"] = *storageMB
	}

	if backupRetentionDays := storage.BackupRetentionDays; backupRetentionDays != nil {
		values["backup_retention_days"] = *backupRetentionDays
	}

	values["geo_redundant_backup"] = string(storage.GeoRedundantBackup)

	values["auto_grow"] = string(storage.StorageAutogrow)

	return []interface{}{values}
}