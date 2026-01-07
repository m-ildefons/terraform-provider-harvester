# harvester_pci_device

Manages PCI device passthrough to VirtualMachines in Harvester using PCIDeviceClaim CRDs.

## Example Usage

```hcl
resource "harvester_virtualmachine" "example" {
  name        = "example-vm"
  namespace   = "default"
  cpu         = 2
  memory      = "4Gi"
  run_strategy = "RerunOnFailure"
  hostname     = "example-vm"

  network_interface {
    name         = "nic-1"
    network_name = "default/vlan1"
  }

  disk {
    name       = "disk-1"
    type       = "disk"
    size       = "20Gi"
    bus        = "virtio"
    boot_order = 1
    image      = "harvester-public/image-ubuntu20.04"
  }
}

resource "harvester_pci_device" "example" {
  name      = "example-pci-device"
  namespace = "default"

  # The VM to attach PCI devices to
  vm_name = "${harvester_virtualmachine.example.namespace}/${harvester_virtualmachine.example.name}"

  # REQUIRED: The node where the VM must be deployed
  # This ensures the VM runs on the correct node where PCI devices are available
  # This prevents scheduling issues when multiple nodes have the same PCI device type
  node_name = "harv1.home.lo"

  # List of PCI addresses to attach
  # Format: "0000:XX:YY.Z" (e.g., "0000:01:00.0")
  # The PCI devices must be enabled for passthrough in Harvester UI before
  # they can be attached via Terraform
  pci_addresses = [
    "0000:00:1f.3",  # Example: Audio device
  ]

  labels = {
    environment = "production"
    managed-by  = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the PCI device resource.
* `namespace` - (Required) The namespace where the resource will be created.
* `vm_name` - (Required) The name of the VirtualMachine to attach PCI devices to. Format: `namespace/name` or just `name` (defaults to the same namespace).
* `node_name` - (Required) The node where the VM must be deployed. This ensures the VM runs on the correct node where PCI devices are available. This is mandatory to prevent scheduling issues when multiple nodes have the same PCI device type.
* `pci_addresses` - (Required) List of PCI addresses to attach to the VM. Format: `0000:XX:YY.Z` (e.g., `0000:01:00.0`). The PCI devices must be enabled for passthrough in Harvester UI before they can be attached via Terraform.
* `labels` - (Optional) Map of labels to apply to the PCIDeviceClaim resource.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The unique identifier for the resource. Format: `namespace/vmname/claimname`.
* `message` - Status message from the PCIDeviceClaim (if available).

## Notes

### PCI Device Passthrough Requirements

1. **Enable Passthrough in Harvester UI**: Before using this resource, the PCI device must be enabled for passthrough in the Harvester UI. This is a prerequisite that cannot be automated via Terraform.

2. **Node Name is Mandatory**: The `node_name` field is required to ensure the VM runs on the correct node where PCI devices are available. This prevents scheduling issues when multiple nodes have the same PCI device type (see [Harvester issue #6648](https://github.com/harvester/harvester/issues/6648)).

3. **PCIDeviceClaim Naming**: The PCIDeviceClaim name is automatically generated based on the node name and PCI address: `{nodeName}-{address}` (e.g., `harv1.home.lo-0000001f3` for address `0000:00:1f.3` on node `harv1.home.lo`).

4. **One Claim per PCI Address**: Each PCI address requires a separate PCIDeviceClaim. The Terraform resource manages multiple addresses by creating multiple claims internally.

5. **Cluster-Scoped Resource**: PCIDeviceClaim is a cluster-scoped resource (not namespaced), but the Terraform resource uses the `namespace` field for consistency with other resources.

### Deletion Behavior

When a `harvester_pci_device` resource is destroyed:

1. The PCIDeviceClaim is marked for deletion with a `deletionTimestamp`.
2. Harvester processes the finalizer `wrangler.cattle.io/PCIDeviceClaimOnRemove` to detach the PCI device from the VM.
3. The claim is fully deleted after the finalizer is processed.

The deletion may take some time as Harvester needs to properly detach the PCI device from the running VM. If the VM is running, it may need to be stopped and restarted for the PCI device to be fully detached.

## Import

PCI device resources can be imported using the resource ID:

```bash
terraform import harvester_pci_device.example default/vm-name/claim-name
```

The ID format is: `namespace/vmname/claimname`

## Related Documentation

- [Harvester PCI Devices Documentation](https://docs.harvesterhci.io/v1.7/advanced/addons/pcidevices)

