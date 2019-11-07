---
layout: "ovh"
page_title: "OVH: me_installation_templates"
sidebar_current: "docs-ovh-datasource-me-installation-templates"
description: |-
  Get the list of dedicated servers installation templates associated with your OVH Account.
---

# ovh_me_installation_templates

Use this data source to get the list of  dedicated servers installation templates associated with your OVH Account.

## Example Usage

```hcl
data "ovh_me_installation_templates" "templates" {}
```

## Argument Reference

This datasource takes no argument.

## Attributes Reference

The following attributes are exported:

* `result` - The list of dedicated servers installation templates IDs associated with your OVH Account.
