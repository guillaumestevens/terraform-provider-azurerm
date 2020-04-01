package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMServiceBusNamespaceNetworkRule_iprule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespaceNetworkRule_iprule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusNamespaceNetworkRule_vnetrule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespaceNetworkRule_vnet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusNamespaceNetworkRule_basicSKU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMServiceBusNamespaceNetworkRule_basicSKU(data),
				ExpectError: regexp.MustCompile("network_rulesets cannot be used when the SKU is Basic or \"Standard\" "),
			},
		},
	})
}

func TestAccAzureRMServiceBusNamespaceNetworkRule_standardSKU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMServiceBusNamespaceNetworkRule_standardSKU(data),
				ExpectError: regexp.MustCompile("network_rulesets cannot be used when the SKU is Basic or \"Standard\" "),
			},
		},
	})
}

func testAccAzureRMServiceBusNamespaceNetworkRule_iprule(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Premium"
  capacity            = 1

}

resource "azurerm_servicebus_namespace_network_rule_set" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"

  properties {
    default_action = "Deny"
    ip_rule {
      ip_mask = "10.0.0.0/16"
      action  = "Allow"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMServiceBusNamespaceNetworkRule_vnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-[1]%d"
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
  name                = "acctestservicebusnamespace-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Premium"
  capacity            = 1
}

resource "azurerm_servicebus_namespace_network_rule_set" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"

  properties {
    default_action = "Deny"
    virtual_network_rule {
      subnet_id = "${azurerm_subnet.test.id}"

      ignore_missing_virtual_network_service_endpoint = true
    }
  }
}

`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMServiceBusNamespaceNetworkRule_basicSKU(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"
  capacity            = 0
}

resource "azurerm_servicebus_namespace_network_rule_set" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"

  properties {
    default_action = "Deny"
    ip_rule {
      ip_mask = "10.0.0.0/16"
      action  = "Allow"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMServiceBusNamespaceNetworkRule_standardSKU(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  capacity            = 0

  network_rulesets {
    default_action = "Deny"
    ip_rule {
      ip_mask = "10.0.0.0/16"
      action  = "Allow"
    }
  }
}

resource "azurerm_servicebus_namespace_network_rule_set" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"

  properties {
    default_action = "Deny"
    ip_rule {
      ip_mask = "10.0.0.0/16"
      action  = "Allow"
    }
  }
}

`, data.RandomInteger, data.Locations.Primary)
}
