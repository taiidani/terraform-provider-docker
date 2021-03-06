package docker

import (
	"fmt"
	"testing"

	"context"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDockerConfig_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckDockerConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "docker_config" "foo" {
					name = "foo-config"
					data = "Ymxhc2RzYmxhYmxhMTI0ZHNkd2VzZA=="
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("docker_config.foo", "name", "foo-config"),
					resource.TestCheckResourceAttr("docker_config.foo", "data", "Ymxhc2RzYmxhYmxhMTI0ZHNkd2VzZA=="),
				),
			},
		},
	})
}
func TestAccDockerConfig_basicUpdatable(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckDockerConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "docker_config" "foo" {
					name 			 = "tftest-myconfig-${replace(timestamp(),":", ".")}"
					data 			 = "Ymxhc2RzYmxhYmxhMTI0ZHNkd2VzZA=="

					lifecycle {
						ignore_changes = ["name"]
						create_before_destroy = true
					}
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("docker_config.foo", "data", "Ymxhc2RzYmxhYmxhMTI0ZHNkd2VzZA=="),
				),
			},
			{
				Config: `
				resource "docker_config" "foo" {
					name 			 = "tftest-myconfig2-${replace(timestamp(),":", ".")}"
					data 			 = "U3VuIDI1IE1hciAyMDE4IDE0OjQ2OjE5IENFU1QK"

					lifecycle {
						ignore_changes = ["name"]
						create_before_destroy = true
					}
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("docker_config.foo", "data", "U3VuIDI1IE1hciAyMDE4IDE0OjQ2OjE5IENFU1QK"),
				),
			},
		},
	})
}

/////////////
// Helpers
/////////////
func testCheckDockerConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).DockerClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "configs" {
			continue
		}

		id := rs.Primary.Attributes["id"]
		_, _, err := client.ConfigInspectWithRaw(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Config with id '%s' still exists", id)
		}
		return nil
	}
	return nil
}
