---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "horizon_vcenter_host_or_cluster Data Source - terraform-provider-horizon"
subcategory: ""
description: |-
  Data source to find the ID of a vCenter Datacenter.
---

# horizon_vcenter_host_or_cluster (Data Source)

Data source to find the ID of a vCenter Datacenter.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `datacenter_id` (String) Datacenter ID
- `name` (String) Host or cluster display name.
- `vcenter_id` (String) Virtual Center ID

### Optional

- `cluster` (Boolean) Whether or not this is a cluster or a host. Defaults to `true`.

### Read-Only

- `id` (String) The ID of this resource.
- `incompatible_reasons` (Set of String) Reasons that may preclude this Host Or Cluster from being used in desktop pool creation.
- `path` (String) Host or cluster path.
- `vgpu_types` (String) Types of NVIDIA GRID vGPUs supported by this host or at least one host on this cluster. If unset, this host or cluster does not support NVIDIA GRID vGPUs and cannot be used for desktop creation with NVIDIA GRID vGPU support enabled.


