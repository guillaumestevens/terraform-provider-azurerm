package tests

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMServiceBusNamespace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}
func TestAccAzureRMServiceBusNamespace_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMServiceBusNamespace_requiresImport),
		},
	})
}

func TestAccAzureRMServiceBusNamespace_readDefaultKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(data.ResourceName),
					resource.TestMatchResourceAttr(
						data.ResourceName, "default_primary_connection_string", regexp.MustCompile("Endpoint=.+")),
					resource.TestMatchResourceAttr(
						data.ResourceName, "default_secondary_connection_string", regexp.MustCompile("Endpoint=.+")),
					resource.TestMatchResourceAttr(
						data.ResourceName, "default_primary_key", regexp.MustCompile(".+")),
					resource.TestMatchResourceAttr(
						data.ResourceName, "default_secondary_key", regexp.MustCompile(".+")),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusNamespace_NonStandardCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespaceNonStandardCasing(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(data.ResourceName),
				),
			},
			{
				Config:             testAccAzureRMServiceBusNamespaceNonStandardCasing(data),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusNamespace_premium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespace_premium(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusNamespace_basicCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMServiceBusNamespace_basicCapacity(data),
				ExpectError: regexp.MustCompile("Service Bus SKU \"Basic\" only supports `capacity` of 0"),
			},
		},
	})
}

func TestAccAzureRMServiceBusNamespace_premiumCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMServiceBusNamespace_premiumCapacity(data),
				ExpectError: regexp.MustCompile("Service Bus SKU \"Premium\" only supports `capacity` of 1, 2, 4 or 8"),
			},
		},
	})
}

func TestAccAzureRMServiceBusNamespace_networkrule_iprule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespace_networkrule_iprule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusNamespace_networkrule_vnetrule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespace_networkrule_vnet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusNamespace_rule_set_basicSKU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMServiceBusNamespace_rule_set_basicSKU(data),
				ExpectError: regex.MustCompile("network_rulesets cannot be used when the SKU is Basic or \"Standard\" "),
			},
		},
	})
}

func TestAccAzureRMServiceBusNamespace_rule_set_standardSKU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMServiceBusNamespace_rule_set_basicSKU(data),
				ExpectError: regex.MustCompile("network_rulesets cannot be used when the SKU is Basic or \"Standard\" "),
			},
		},
	})
}

func TestAccAzureRMServiceBusNamespace_zoneRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespace_zoneRedundant(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "zone_redundant", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMServiceBusNamespaceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.NamespacesClientPreview
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_servicebus_namespace" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("ServiceBus Namespace still exists:\n%+v", resp)
		}
	}

	return nil
}

func testCheckAzureRMServiceBusNamespaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.NamespacesClientPreview
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		namespaceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Service Bus Namespace: %s", namespaceName)
		}

		resp, err := client.Get(ctx, resourceGroup, namespaceName)
		if err != nil {
			return fmt.Errorf("Bad: Get on serviceBusNamespacesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Service Bus Namespace %q (resource group: %q) does not exist", namespaceName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMServiceBusNamespace_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMServiceBusNamespace_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace" "import" {
  name                = "${azurerm_servicebus_namespace.test.name}"
  location            = "${azurerm_servicebus_namespace.test.location}"
  resource_group_name = "${azurerm_servicebus_namespace.test.resource_group_name}"
  sku                 = "${azurerm_servicebus_namespace.test.sku}"
}
`, testAccAzureRMServiceBusNamespace_basic(data))
}

func testAccAzureRMServiceBusNamespaceNonStandardCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMServiceBusNamespace_premium(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Premium"
  capacity            = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMServiceBusNamespace_basicCapacity(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"
  capacity            = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMServiceBusNamespace_premiumCapacity(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Premium"
  capacity            = 0
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMServiceBusNamespace_zoneRedundant(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Premium"
  capacity            = 1
  zone_redundant      = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMServiceBusNamespace_networkrule_iprule(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name = "acctestRG-%[1]d"
	location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
	name 		= "acctestservicebusnamespace-%[1]d"
	location	= "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	sku = "Premium"
	capacity = 1
	
	network_rulesets {
		default_action = "Deny"
		ip_rule {
			ip_mask = "10.0.0.0/16"
			action = "Allow"
		}
	}
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMServiceBusNamespace_networkrule_vnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name = "acctestRG-[1]%d"
	location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
	name                = "acctvn-%[1]d"
	address_space       = ["10.0.0.0/16"]
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
  }
  
  resource "azurerm_subnet" "test" {
	name                 = "acctsub-%[1]d"
	resource_group_name  = "${azurerm_resource_group.test.name}"
	virtual_network_name = "${azurerm_virtual_network.test.name}"
	address_prefix       = "10.0.2.0/24"
  }

resource "azurerm_servicebus_namespace" "test" {
	name 		= "acctestservicebusnamespace-%[1]d"
	location	= "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	sku = "Premium"
	capacity = 1
	
	network_rulesets {
		default_action = "Deny"
		virtual_network_rule {
			subnet_id = "${azurerm_subnet.test.id}"
	  
			ignore_missing_virtual_network_service_endpoint = true
		}
	}
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMServiceBusNamespace_rule_set_basicSKU(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name = "acctestRG-%[1]d"
	location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
	name 		= "acctestservicebusnamespace-%[1]d"
	location	= "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	sku = "Basic"
	capacity = 0
	
	network_rulesets {
		default_action = "Deny"
		ip_rule {
			ip_mask = "10.0.0.0/16"
			action = "Allow"
		}
	}
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMServiceBusNamespace_rule_set_standardSKU(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name = "acctestRG-%[1]d"
	location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
	name 		= "acctestservicebusnamespace-%[1]d"
	location	= "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	sku = "Standard"
	capacity = 0
	
	network_rulesets {
		default_action = "Deny"
		ip_rule {
			ip_mask = "10.0.0.0/16"
			action = "Allow"
		}
	}
}
`, data.RandomInteger, data.Locations.Primary)
}
