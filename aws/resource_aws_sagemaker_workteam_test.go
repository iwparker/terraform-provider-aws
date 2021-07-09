package aws

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service/sagemaker/finder"
)

func init() {
	resource.AddTestSweepers("aws_sagemaker_workteam", &resource.Sweeper{
		Name: "aws_sagemaker_workteam",
		F:    testSweepSagemakerWorkteams,
	})
}

func testSweepSagemakerWorkteams(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*AWSClient).sagemakerconn
	var sweeperErrs *multierror.Error

	err = conn.ListWorkteamsPages(&sagemaker.ListWorkteamsInput{}, func(page *sagemaker.ListWorkteamsOutput, lastPage bool) bool {
		for _, workteam := range page.Workteams {

			r := resourceAwsSagemakerWorkteam()
			d := r.Data(nil)
			d.SetId(aws.StringValue(workteam.WorkteamName))
			err := r.Delete(d, client)
			if err != nil {
				log.Printf("[ERROR] %s", err)
				sweeperErrs = multierror.Append(sweeperErrs, err)
				continue
			}
		}

		return !lastPage
	})

	if testSweepSkipSweepError(err) {
		log.Printf("[WARN] Skipping SageMaker workteam sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil()
	}

	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error retrieving Sagemaker Workteams: %w", err))
	}

	return sweeperErrs.ErrorOrNil()
}

func TestAccAWSSagemakerWorkteam_cognitoConfig(t *testing.T) {
	var workteam sagemaker.Workteam
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_sagemaker_workteam.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, sagemaker.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSagemakerWorkteamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSagemakerWorkteamCognitoConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "workteam_name", rName),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexp.MustCompile(`workteam/.+`)),
					resource.TestCheckResourceAttr(resourceName, "description", rName),
					resource.TestCheckResourceAttr(resourceName, "member_definition.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "member_definition.0.cognito_member_definition.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.client_id", "aws_cognito_user_pool_client.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.user_pool", "aws_cognito_user_pool.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.user_group", "aws_cognito_user_group.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "subdomain"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"workforce_name"},
			},
			{
				Config: testAccAWSSagemakerWorkteamCognitoUpdatedConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "workteam_name", rName),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexp.MustCompile(`workteam/.+`)),
					resource.TestCheckResourceAttr(resourceName, "description", rName),
					resource.TestCheckResourceAttr(resourceName, "member_definition.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "member_definition.0.cognito_member_definition.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.client_id", "aws_cognito_user_pool_client.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.user_pool", "aws_cognito_user_pool.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.user_group", "aws_cognito_user_group.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "member_definition.1.cognito_member_definition.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.1.cognito_member_definition.0.client_id", "aws_cognito_user_pool_client.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.1.cognito_member_definition.0.user_pool", "aws_cognito_user_pool.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.1.cognito_member_definition.0.user_group", "aws_cognito_user_group.test2", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "subdomain"),
				),
			},
			{
				Config: testAccAWSSagemakerWorkteamCognitoConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "workteam_name", rName),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexp.MustCompile(`workteam/.+`)),
					resource.TestCheckResourceAttr(resourceName, "description", rName),
					resource.TestCheckResourceAttr(resourceName, "member_definition.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "member_definition.0.cognito_member_definition.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.client_id", "aws_cognito_user_pool_client.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.user_pool", "aws_cognito_user_pool.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.user_group", "aws_cognito_user_group.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "subdomain"),
				),
			},
		},
	})
}

// func TestAccAWSSagemakerWorkteam_oidcConfig(t *testing.T) {
// 	var workteam sagemaker.Workteam
// 	rName := acctest.RandomWithPrefix("tf-acc-test")
// 	resourceName := "aws_sagemaker_workteam.test"
// 	endpoint1 := "https://example.com"
// 	endpoint2 := "https://test.example.com"

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		ErrorCheck:   testAccErrorCheck(t, sagemaker.EndpointsID),
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckAWSSagemakerWorkteamDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccAWSSagemakerWorkteamOidcConfig(rName, endpoint1),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
// 					resource.TestCheckResourceAttr(resourceName, "workteam_name", rName),
// 					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexp.MustCompile(`workteam/.+`)),
// 					resource.TestCheckResourceAttr(resourceName, "cognito_config.#", "0"),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.#", "1"),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.authorization_endpoint", endpoint1),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.client_id", rName),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.client_secret", rName),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.issuer", endpoint1),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.jwks_uri", endpoint1),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.logout_endpoint", endpoint1),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.token_endpoint", endpoint1),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.user_info_endpoint", endpoint1),
// 					resource.TestCheckResourceAttr(resourceName, "source_ip_config.#", "1"),
// 					resource.TestCheckResourceAttr(resourceName, "source_ip_config.0.cidrs.#", "0"),
// 					resource.TestCheckResourceAttrSet(resourceName, "subdomain"),
// 				),
// 			},
// 			{
// 				ResourceName:            resourceName,
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"oidc_config.0.client_secret"},
// 			},
// 			{
// 				Config: testAccAWSSagemakerWorkteamOidcConfig(rName, endpoint2),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
// 					resource.TestCheckResourceAttr(resourceName, "workteam_name", rName),
// 					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexp.MustCompile(`workteam/.+`)),
// 					resource.TestCheckResourceAttr(resourceName, "cognito_config.#", "0"),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.#", "1"),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.authorization_endpoint", endpoint2),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.client_id", rName),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.client_secret", rName),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.issuer", endpoint2),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.jwks_uri", endpoint2),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.logout_endpoint", endpoint2),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.token_endpoint", endpoint2),
// 					resource.TestCheckResourceAttr(resourceName, "oidc_config.0.user_info_endpoint", endpoint2),
// 					resource.TestCheckResourceAttr(resourceName, "source_ip_config.#", "1"),
// 					resource.TestCheckResourceAttr(resourceName, "source_ip_config.0.cidrs.#", "0"),
// 					resource.TestCheckResourceAttrSet(resourceName, "subdomain"),
// 				),
// 			},
// 		},
// 	})
// }

func TestAccAWSSagemakerWorkteam_disappears(t *testing.T) {
	var workteam sagemaker.Workteam
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_sagemaker_workteam.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, sagemaker.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSagemakerWorkteamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSagemakerWorkteamCognitoConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					testAccCheckResourceDisappears(testAccProvider, resourceAwsSagemakerWorkteam(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAWSSagemakerWorkteamDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).sagemakerconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_sagemaker_workteam" {
			continue
		}

		workteam, err := finder.WorkteamByName(conn, rs.Primary.ID)
		if tfawserr.ErrMessageContains(err, "ValidationException", "The work team") {
			continue
		}

		if err != nil {
			return fmt.Errorf("error reading Sagemaker Workteam (%s): %w", rs.Primary.ID, err)
		}

		if aws.StringValue(workteam.WorkteamName) == rs.Primary.ID {
			return fmt.Errorf("SageMaker Workteam %q still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckAWSSagemakerWorkteamExists(n string, workteam *sagemaker.Workteam) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No sagmaker workteam ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).sagemakerconn
		resp, err := finder.WorkteamByName(conn, rs.Primary.ID)
		if err != nil {
			return err
		}

		*workteam = *resp

		return nil
	}
}

func testAccAWSSagemakerWorkteamBaseConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_cognito_user_pool" "test" {
  name = %[1]q
}

resource "aws_cognito_user_pool_client" "test" {
  name            = %[1]q
  generate_secret = true
  user_pool_id    = aws_cognito_user_pool.test.id
}

resource "aws_cognito_user_pool_domain" "test" {
  domain       = %[1]q
  user_pool_id = aws_cognito_user_pool.test.id
}

resource "aws_cognito_user_group" "test" {
  name         = %[1]q
  user_pool_id = aws_cognito_user_pool.test.id
}

resource "aws_sagemaker_workforce" "test" {
  workforce_name = %[1]q

  cognito_config {
    client_id = aws_cognito_user_pool_client.test.id
    user_pool = aws_cognito_user_pool_domain.test.user_pool_id
  }
}
`, rName)
}

func testAccAWSSagemakerWorkteamCognitoConfig(rName string) string {
	return testAccAWSSagemakerWorkteamBaseConfig(rName) + fmt.Sprintf(`
resource "aws_sagemaker_workteam" "test" {
  workteam_name  = %[1]q
  workforce_name = aws_sagemaker_workforce.test.id
  description    = %[1]q

  member_definition {
    cognito_member_definition {
      client_id  = aws_cognito_user_pool_client.test.id
      user_pool  = aws_cognito_user_pool_domain.test.user_pool_id
	  user_group = aws_cognito_user_group.test.id
	}
  }
}
`, rName)
}

func testAccAWSSagemakerWorkteamCognitoUpdatedConfig(rName string) string {
	return testAccAWSSagemakerWorkteamBaseConfig(rName) + fmt.Sprintf(`
resource "aws_cognito_user_group" "test2" {
  name         = "%[1]s-2"
  user_pool_id = aws_cognito_user_pool.test.id
}

resource "aws_sagemaker_workteam" "test" {
  workteam_name  = %[1]q
  workforce_name = aws_sagemaker_workforce.test.id
  description    = %[1]q

  member_definition {
    cognito_member_definition {
      client_id  = aws_cognito_user_pool_client.test.id
      user_pool  = aws_cognito_user_pool_domain.test.user_pool_id
	  user_group = aws_cognito_user_group.test.id
	}
  }

  member_definition {
	cognito_member_definition {
      client_id  = aws_cognito_user_pool_client.test.id
      user_pool  = aws_cognito_user_pool_domain.test.user_pool_id
	  user_group = aws_cognito_user_group.test2.id
	}
  }  
}
`, rName)
}

// func testAccAWSSagemakerWorkteamOidcConfig(rName, endpoint string) string {
// 	return testAccAWSSagemakerWorkteamBaseConfig(rName) + fmt.Sprintf(`
// resource "aws_sagemaker_workteam" "test" {
//   workteam_name = %[1]q

//   oidc_config {
//     authorization_endpoint = %[2]q
//     client_id              = %[1]q
//     client_secret          = %[1]q
//     issuer                 = %[2]q
//     jwks_uri               = %[2]q
//     logout_endpoint        = %[2]q
//     token_endpoint         = %[2]q
//     user_info_endpoint     = %[2]q
//   }
// }
// `, rName, endpoint)
// }
