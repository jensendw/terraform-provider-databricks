---
subcategory: "MLFlow"
---
# databricks_mlflow_experiment Resource

This resource allows you to create MLFlow experiments in Databricks.

## Example Usage

```hcl
resource "databricks_mlflow_experiment" "test" {
  name = "My MLFlow Experiment"

  description = "My MLFlow experiment description"

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

* `name` - (Required) Name of MLFLow experiment.
* `description` - The description of the MLFlow experiment.
* `tags` - Tags for the MLFlow experiment.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of the MLFlow experiment.
* `description` - The description of the MLFlow experiment.