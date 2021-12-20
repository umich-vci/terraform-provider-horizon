data "horizon_vcenter_datastore" "example" {
  name               = "Example Datastore"
  host_or_cluster_id = data.horizon_vcenter_host_or_cluster.example.id
  vcenter_id         = data.horizon_vcenter_server.example.id
}
