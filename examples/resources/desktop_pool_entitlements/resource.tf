resource "horizon_desktop_pool_entitlements" "example" {
  pool_id = horizon_desktop_pool_automated.example.id
  ad_user_or_group_ids = [
    data.horizon_active_directory_domain_user_or_group.domain_users.id,
  ]
}
