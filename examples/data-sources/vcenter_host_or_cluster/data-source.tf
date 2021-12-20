data "horizon_vcenter_host_or_cluster" "example" {
  name          = "Example Cluster"
  datacenter_id = data.horizon_vcenter_datacenter.example.id
  vcenter_id    = data.horizon_vcenter_server.example.id
}
