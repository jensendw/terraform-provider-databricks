package mlflow

import (
	"context"

	"github.com/databrickslabs/terraform-provider-databricks/common"
	"github.com/databrickslabs/terraform-provider-databricks/mlflow/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Experiment defines the parameters that can be set in the resource.
type Experiment struct {
	Name             string    `json:"name"`
	Tags             []api.Tag `json:"tags,omitempty" tf:"force_new"`
	ArtifactLocation string    `json:"artifact_location,omitempty" tf:"force_new"`
	Description      string    `json:"description,omitempty"`
}

func (d *Experiment) fromAPIObject(ad *api.Experiment, schema map[string]*schema.Schema, data *schema.ResourceData) error {
	// Copy from API object.
	d.Name = ad.Name
	d.Tags = ad.Tags

	// Pass to ResourceData.
	if err := common.StructToData(*d, schema, data); err != nil {
		return err
	}

	// Overwrite `tags` in case they're empty on the server side.
	// This would have been skipped by `common.StructToData` because of slice emptiness.
	// Ideally, the reflection code also sets empty values, but we'd risk
	// clobbering values we actually want to keep around in existing code.
	data.Set("tags", ad.Tags)
	data.Set("name", ad.Name)
	data.SetId(ad.ExperimentId)

	return nil
}

///func ResourceMLFlowExperiment() {}
func ResourceMLFlowExperiment() *schema.Resource {
	s := common.StructToSchema(
		Experiment{},
		func(m map[string]*schema.Schema) map[string]*schema.Schema {
			return m
		})

	return common.Resource{
		Create: func(ctx context.Context, data *schema.ResourceData, c *common.DatabricksClient) error {
			var ad api.Experiment
			if err := common.DataToStructPointer(data, s, &ad); err != nil {
				return err
			}
			if err := api.NewExperimentAPI(ctx, c).Create(&ad); err != nil {
				return err
			}
			data.SetId(ad.ExperimentId)
			return nil
		},
		Read: func(ctx context.Context, data *schema.ResourceData, c *common.DatabricksClient) error {
			ad, err := api.NewExperimentAPI(ctx, c).Read(data.Id())
			if err != nil {
				return err
			}
			var d Experiment
			return d.fromAPIObject(ad, s, data)
		},
		Update: func(ctx context.Context, data *schema.ResourceData, c *common.DatabricksClient) error {
			var ad api.Experiment
			if err := common.DataToStructPointer(data, s, &ad); err != nil {
				return err
			}
			updateDoc := api.ExperimentUpdate{ExperimentId: data.Id(), NewName: ad.Name}
			return api.NewExperimentAPI(ctx, c).Update(&updateDoc)
		},
		Delete: func(ctx context.Context, data *schema.ResourceData, c *common.DatabricksClient) error {
			ad := api.Experiment{ExperimentId: data.Id()}
			return api.NewExperimentAPI(ctx, c).Delete(&ad)
		},
		StateUpgraders: []schema.StateUpgrader{},
		Schema:         s,
		SchemaVersion:  0,
		Timeouts:       &schema.ResourceTimeout{},
	}.ToResource()
}
