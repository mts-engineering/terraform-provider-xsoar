package xsoar

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccIntegrationInstanceDataSource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { testAccIntegrationInstanceDataSourcePreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"xsoar": func() (tfprotov6.ProviderServer, error) {
				return tfsdk.NewProtocol6Server(New()), nil
			},
		},
		CheckDestroy: testAccCheckIntegrationInstanceDataSourceDestroy(rName),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationInstanceDataSourceBasic(rName),
				Check:  resource.TestCheckResourceAttrPair("data.xsoar_integration_instance."+rName, "id", "xsoar_integration_instance."+rName, "id"),
			},
		},
	})
}

func testAccIntegrationInstanceDataSourcePreCheck(t *testing.T) {}

func testAccCheckIntegrationInstanceDataSourceDestroy(r string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources["xsoar_integration_instance."+r]
		if !ok {
			return fmt.Errorf("not found: %s in %s", r, state.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		resp, _, err := openapiClient.DefaultApi.GetIntegrationInstance(context.Background()).SetIdentifier(r).Execute()
		if err != nil {
			return fmt.Errorf("Error getting integration instance: " + err.Error())
		}
		if resp != nil {
			return fmt.Errorf("integration instance returned when it should be destroyed")
		}
		return nil
	}
}

func testAccIntegrationInstanceDataSourceBasic(name string) string {
	c := `
resource "xsoar_integration_instance" "{name}" {
  name               = "{name}"
  integration_name   = "threatcentral"
  propagation_labels = ["all"]
  config = {
    APIAddress : "https://threatcentral.io/tc/rest/summaries"
    APIKey : "123"
    useproxy : "true"
  }
}

data "xsoar_integration_instance" "{name}" {
  name = xsoar_integration_instance.{name}.name
}
`
	c = strings.Replace(c, "{name}", name, -1)
	return c
}
