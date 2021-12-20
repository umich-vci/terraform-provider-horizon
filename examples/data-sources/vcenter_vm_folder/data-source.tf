data "horizon_vcenter_vm_folder" "example" {
  datacenter_id = data.horizon_vcenter_datacenter.example.id
  path          = "/${data.horizon_vcenter_datacenter.example.name}/vm"
  vcenter_id    = data.horizon_vcenter_server.example.id
}
