package migration

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func AccountV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		// this should have been applied from pre-0.12 migration system; backporting just in-case
		Type:    accountSchemaForV0AndV1().CoreConfigSchema().ImpliedType(),
		Upgrade: accountUpgradeV0ToV1,
		Version: 0,
	}
}

func AccountV1ToV2() schema.StateUpgrader {
	return schema.StateUpgrader{
		// this should have been applied from pre-0.12 migration system; backporting just in-case
		Type:    accountSchemaForV0AndV1().CoreConfigSchema().ImpliedType(),
		Upgrade: accountUpgradeV1ToV2,
		Version: 1,
	}
}

func accountSchemaForV0AndV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"account_kind": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Storage",
			},

			"account_type": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "This field has been split into `account_tier` and `account_replication_type`",
			},

			"account_tier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"account_replication_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Only valid for BlobStorage accounts, defaults to "Hot" in create function
			"access_tier": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"custom_domain": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"use_subdomain": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"enable_blob_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"enable_file_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"enable_https_traffic_only": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"primary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_blob_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_blob_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_queue_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_queue_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_table_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_table_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// NOTE: The API does not appear to expose a secondary file endpoint
			"primary_file_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_blob_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_blob_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func accountUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	accountType := rawState["account_type"].(string)
	split := strings.Split(accountType, "_")
	rawState["account_tier"] = split[0]
	rawState["account_replication_type"] = split[1]
	return rawState, nil
}

func accountUpgradeV1ToV2(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	rawState["account_encryption_source"] = "Microsoft.Storage"
	return rawState, nil
}
