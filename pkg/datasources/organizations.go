package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var organizationsSchema = map[string]*schema.Schema{
	"warehouses": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The organization in the database",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"admin_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"admin_password": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"email": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"edition": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"region": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"first_name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"last_name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"comment": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func Organizations() *schema.Resource {
	return &schema.Resource{
		Read:   ReadOrganizations,
		Schema: organizationsSchema,
	}
}

func ReadOrganizations(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	account, err := snowflake.ReadCurrentAccount(db)
	if err != nil {
		log.Print("[DEBUG] unable to retrieve current account")
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%s.%s", account.Account, account.Region))

	currentOrganizations, err := snowflake.ListOrganizations(db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] no organization found in account (%s)", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse organization in account (%s)", d.Id())
		d.SetId("")
		return nil
	}

	organizations := []map[string]interface{}{}

	for _, organization := range currentOrganizations {
		organizationMap := map[string]interface{}{}

		organizationMap["name"] = organization.Name
		organizationMap["admin_name"] = organization.AdminName
		organizationMap["admin_password"] = organization.AdminPassword
		organizationMap["email"] = organization.email
		organizationMap["edition"] = organization.Edition
		organizationMap["region"] = organization.Region
		organizationMap["first_name"] = organization.FirstName
		organizationMap["last_name"] = organization.LastName
		organizationMap["comment"] = organization.Comment

		organizations = append(organizations, organizationMap)
	}

	return d.Set("organizations", organizations)
}
