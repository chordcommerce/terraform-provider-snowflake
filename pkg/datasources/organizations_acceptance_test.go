package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizations(t *testing.T) {
	organizationName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: organizations(organizationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_organizations.s", "organizations.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_organizations.s", "organizations.0.name"),
				),
			},
		},
	})
}

func organizations(organizationName string) string {
	return fmt.Sprintf(`
	resource snowflake_organization "s"{
		name 	 	  	 	= "%v"
		organization_size 		= "XSMALL"
		initially_suspended = true
		auto_suspend	    = 60
	}

	data snowflake_organizations "s" {
		depends_on = [snowflake_organization.s]
	}
	`, organizationName)
}
