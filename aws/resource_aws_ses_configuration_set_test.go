package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAWSSESConfigurationSet_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAWSSES(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSESConfigurationSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSESConfigurationSetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsSESConfigurationSetExists("aws_ses_configuration_set.test"),
				),
			},
			{
				ResourceName:      "aws_ses_configuration_set.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSESConfigurationSetDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).sesconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_ses_configuration_set" {
			continue
		}

		response, err := conn.ListConfigurationSets(&ses.ListConfigurationSetsInput{})
		if err != nil {
			return err
		}

		found := false
		for _, element := range response.ConfigurationSets {
			if *element.Name == fmt.Sprintf("some-configuration-set-%d", escRandomInteger) {
				found = true
			}
		}

		if found {
			return fmt.Errorf("The configuration set still exists")
		}

	}

	return nil

}

func testAccCheckAwsSESConfigurationSetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("SES configuration set not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("SES configuration set ID not set")
		}

		conn := testAccProvider.Meta().(*AWSClient).sesconn

		response, err := conn.ListConfigurationSets(&ses.ListConfigurationSetsInput{})
		if err != nil {
			return err
		}

		found := false
		for _, element := range response.ConfigurationSets {
			if *element.Name == fmt.Sprintf("some-configuration-set-%d", escRandomInteger) {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("The configuration set was not created")
		}

		return nil
	}
}

var escRandomInteger = acctest.RandInt()
var testAccAWSSESConfigurationSetConfig = fmt.Sprintf(`
resource "aws_ses_configuration_set" "test" {
    name = "some-configuration-set-%d"
}
`, escRandomInteger)
