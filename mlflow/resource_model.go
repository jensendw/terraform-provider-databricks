package mlflow

import (
	"context"

	"github.com/databrickslabs/terraform-provider-databricks/common"
	"github.com/databrickslabs/terraform-provider-databricks/mlflow/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MLFlowModel defines the parameters that can be set in the resource.
type Model struct {
	Name        string    `json:"name" tf:"force_new"`
	Tags        []api.Tag `json:"tags,omitempty" tf:"force_new"`
	Description string    `json:"description,omitempty"`
}

func (d *Model) fromAPIObject(ad *api.Model, schema map[string]*schema.Schema, data *schema.ResourceData) error {
	// Copy from API object.
	d.Name = ad.Name
	d.Tags = ad.Tags
	d.Description = ad.Description

	// Pass to ResourceData.
	if err := common.StructToData(*d, schema, data); err != nil {
		return err
	}

	return nil
}

// ResourceDashboard ...
func ResourceMLFlowModel() *schema.Resource {
	s := common.StructToSchema(
		Model{},
		func(m map[string]*schema.Schema) map[string]*schema.Schema {
			return m
		})

	return common.Resource{
		Create: func(ctx context.Context, data *schema.ResourceData, c *common.DatabricksClient) error {
			var ad api.Model
			if err := common.DataToStructPointer(data, s, &ad); err != nil {
				return nil
			}
			if err := api.NewModelAPI(ctx, c).Create(&ad); err != nil {
				return err
			}
			// No need to set anything because the resource is going to be
			// read immediately after being created.
			data.SetId(ad.Name)
			return nil
		},
		Read: func(ctx context.Context, data *schema.ResourceData, c *common.DatabricksClient) error {
			ad, err := api.NewModelAPI(ctx, c).Read(data.Id())
			if err != nil {
				return err
			}
			var d Model
			return d.fromAPIObject(ad, s, data)
		},
		Update: func(ctx context.Context, data *schema.ResourceData, c *common.DatabricksClient) error {
			var ad api.Model
			if err := common.DataToStructPointer(data, s, &ad); err != nil {
				return nil
			}

			return api.NewModelAPI(ctx, c).Update(&ad)
		},
		Delete: func(ctx context.Context, data *schema.ResourceData, c *common.DatabricksClient) error {
			var ad api.Model
			if err := common.DataToStructPointer(data, s, &ad); err != nil {
				return nil
			}
			return api.NewModelAPI(ctx, c).Delete(&ad)
		},
		Schema: s,
	}.ToResource()
}
