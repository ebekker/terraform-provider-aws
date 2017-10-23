package aws

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSCognitoUserPoolDomain_basic(t *testing.T) {
	name := fmt.Sprintf("%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	updatedName := fmt.Sprintf("%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSCognitoUserPoolDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSCognitoUserPoolDomainConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSCognitoUserPoolDomainExists("aws_cognito_identity_pool.main"),
					resource.TestCheckResourceAttr("aws_cognito_identity_pool.main", "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
					resource.TestCheckResourceAttr("aws_cognito_identity_pool.main", "allow_unauthenticated_identities", "false"),
				),
			},
			{
				Config: testAccAWSCognitoUserPoolDomainConfig_basic(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSCognitoUserPoolDomainExists("aws_cognito_identity_pool.main"),
					resource.TestCheckResourceAttr("aws_cognito_identity_pool.main", "identity_pool_name", fmt.Sprintf("identity pool %s", updatedName)),
				),
			},
		},
	})
}

func testAccCheckAWSCognitoUserPoolDomainExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Cognito Identity Pool Domain ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).cognitoidpconn

		_, err := conn.DescribeUserPoolDomain(&cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckAWSCognitoUserPoolDomainDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).cognitoidpconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_cognito_user_pool_domain" {
			continue
		}

		_, err := conn.DescribeUserPoolDomain(&cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: aws.String(rs.Primary.ID),
		})

		if err != nil {
			if isAWSErr(err, "ResourceNotFoundException", "") {
				return nil
			}
			return err
		}
	}

	return nil
}

func testAccAWSCognitoUserPoolDomainConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "aws_cognito_user_pool_domain" "main" {
  domain = "tf-acc-test-%s.com"
  user_pool_id = "${aws_cognito_identity_pool.main.id}"
}

resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = "identity pool %s"
  allow_unauthenticated_identities = false
  developer_provider_name          = "my.developer"
}
`, name, name)
}
