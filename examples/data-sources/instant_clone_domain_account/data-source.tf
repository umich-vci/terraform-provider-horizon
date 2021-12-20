data "horizon_instant_clone_domain_account" "icuser" {
  username     = "icuser"
  ad_domain_id = data.horizon_active_directory_domain.contoso.id
}
