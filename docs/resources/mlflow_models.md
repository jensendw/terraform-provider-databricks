---
subcategory: "MLFlow"
---
# databricks_mlflow_model Resource

This resource allows you to create MLFlow models in Databricks.

## Example Usage

```hcl
resource "databricks_mlflow_model" "test" {
  name = "My MLFlow Model"

  description = "My MLFlow model description"

  tags {
    key   = "key1"
    value = "value1"
  }
  tags {
    key   = "key2"
    value = "value2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of MLFLow model.
* `description` - The description of the MLFlow model.
* `tags` - Tags for the MLFlow model.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of the MLFlow model.
* `description` - The description of the MLFlow model.