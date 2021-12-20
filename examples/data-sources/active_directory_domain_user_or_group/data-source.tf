locals {
  my_ad_filter = {
    type = "And"
    filters = [
      {
        type  = "Equals"
        name  = "name"
        value = "Domain Users"
      },
      {
        type  = "Equals"
        name  = "domain"
        value = "ad.contoso.com"
      }
    ]
  }
}

data "horizon_active_directory_domain_user_or_group" "example" {
  filter = jsonencode(local.my_ad_filter)
}
