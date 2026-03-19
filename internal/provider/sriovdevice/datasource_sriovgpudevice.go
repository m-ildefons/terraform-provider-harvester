package sriovdevice

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSRIOVGPUDevice() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSRIOVGPUDeviceRead,
		Schema:      DataSourceSRIOVGPUDeviceSchema(),
	}
}

func dataSourceSRIOVGPUDeviceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
}
