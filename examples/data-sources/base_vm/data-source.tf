data "horizon_vcenter_base_vm" "example" {
  name                    = "vdi-win11"
  datacenter_id           = data.horizon_vcenter_datacenter.example.id
  vcenter_id              = data.horizon_vcenter_server.example.id
  filter_incompatible_vms = true
}
