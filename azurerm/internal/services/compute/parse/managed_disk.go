package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ManagedDiskId struct {
	ResourceGroup string
	Name          string
}

func ManagedDiskID(input string) (*ManagedDiskId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Managed Disk ID %q: %+v", input, err)
	}

	server := ManagedDiskId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.Name, err = id.PopSegment("disks"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}
