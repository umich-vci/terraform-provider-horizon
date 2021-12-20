data "horizon_vcenter_base_vm_snapshot" "example" {
  base_vm_id = data.horizon_vcenter_base_vm.example.id
  path       = "/my snapshot"
  vcenter_id = data.horizon_vcenter_server.example.id
}
