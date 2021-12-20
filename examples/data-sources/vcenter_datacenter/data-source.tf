data "horizon_vcenter_datacenter" "example" {
  name       = "Example Datacenter"
  vcenter_id = data.horizon_vcenter_server.example.id
}
