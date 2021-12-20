data "horizon_vcenter_resource_pool" "example" {
  host_or_cluster_id = data.horizon_vcenter_host_or_cluster.example.id
  name               = "Example"
  vcenter_id         = data.horizon_vcenter_server.example.id
}
