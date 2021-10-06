package resources

import (
	"database/sql"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var organizationProperties = []string{
	"comment", "admin_name", "admin_password", "email",
	"edition", "region", "first_name", "last_name",
}

var organizationSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the virtual organization; must be unique for your account.",
	},
	"admin_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the login name of the initial administrative user of the account",
	},
	"admin_password": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the password for the initial administrative user of the account",
	},
	"email": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the email address of the initial administrative user of the account.",
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"first_name": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"last_name": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"region": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		Description: "Specifies the Snowflake Region where the account is created. If no value is provided, Snowflake creates the account in the same Snowflake Region as the current account",
	},
	"edition": {
		Type:     schema.TypeString,
		Required: true,
		ValidateFunc: validation.StringInSlice([]string{
			"STANDARD", "ENTERPRISE", "BUSINESS_CRITICAL",
		}, true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.Replace(s, "-", "", -1))
			}
			return normalize(old) == normalize(new)
		},
		Description: "Specifies the Snowflake Edition of the account.",
	},
}

// Organization returns a pointer to the resource representing a organization
func Organization() *schema.Resource {
	return &schema.Resource{
		Create: CreateOrganization,
		Read:   ReadOrganization,
		Update: UpdateOrganization,

		Schema: organizationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateOrganization implements schema.CreateFunc
func CreateOrganization(d *schema.ResourceData, meta interface{}) error {
	props := append(organizationProperties)
	return CreateResource("organization", props, organizationSchema, snowflake.Organization, ReadOrganization)(d, meta)
}

// ReadOrganization implements schema.ReadFunc
func ReadOrganization(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	stmt := snowflake.Organization(d.Id()).Show()

	row := snowflake.QueryRow(db, stmt)
	w, err := snowflake.ScanOrganization(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] organization (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	err = d.Set("name", w.Name)
	if err != nil {
		return err
	}
	err = d.Set("admin_name", w.AdminName)
	if err != nil {
		return err
	}
	err = d.Set("admin_password", w.AdminPassword)
	if err != nil {
		return err
	}
	err = d.Set("email", w.Email)
	if err != nil {
		return err
	}
	err = d.Set("comment", w.Comment)
	if err != nil {
		return err
	}
	err = d.Set("first_name", w.FirstName)
	if err != nil {
		return err
	}
	err = d.Set("last_name", w.LastName)
	if err != nil {
		return err
	}
	err = d.Set("region", w.Region)
	if err != nil {
		return err
	}
	err = d.Set("edition", w.Edition)

	return err
}

// UpdateOrganization implements schema.UpdateFunc
func UpdateOrganization(d *schema.ResourceData, meta interface{}) error {
	return UpdateResource("organization", organizationProperties, organizationSchema, snowflake.Organization, ReadOrganization)(d, meta)
}