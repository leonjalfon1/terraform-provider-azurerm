package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/parse"
)

func DnsZoneV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Version: 0,
		Type:    dnsZoneSchemaForV0().CoreConfigSchema().ImpliedType(),
		Upgrade: dnsZoneUpgradeV0ToV1,
	}
}

func dnsZoneSchemaForV0() *schema.Resource {
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

			"number_of_record_sets": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_number_of_record_sets": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"name_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"soa_record": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Required: true,
						},

						"host_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"expire_time": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  2419200,
						},

						"minimum_ttl": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  300,
						},

						"refresh_time": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3600,
						},

						"retry_time": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  300,
						},

						"serial_number": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},

						"ttl": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3600,
						},

						"tags": {
							Type:     schema.TypeMap,
							Optional: true,
						},

						"fqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func dnsZoneUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	ctx := context.TODO()
	groupsClient := meta.(*clients.Client).Resource.GroupsClient
	oldId := rawState["id"].(string)
	id, err := parse.DnsZoneID(oldId)
	if err != nil {
		return rawState, err
	}
	resGroup, err := groupsClient.Get(ctx, id.ResourceGroup)
	if err != nil {
		return rawState, err
	}
	if resGroup.Name == nil {
		return rawState, fmt.Errorf("`name` was nil for Resource Group %q", id.ResourceGroup)
	}
	resourceGroup := *resGroup.Name
	name := rawState["name"].(string)
	newId := parse.NewDnsZoneID(id.SubscriptionId, resourceGroup, name).ID()
	log.Printf("Updating `id` from %q to %q", oldId, newId)
	rawState["id"] = newId
	return rawState, nil
}
