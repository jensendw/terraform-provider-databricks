package mlflow

import (
	"context"
	"fmt"

	"github.com/databrickslabs/terraform-provider-databricks/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MLFlowModel defines the parameters that can be set in the resource.
type MLFLowModel struct {
	Name        string   `json:"name"`
	Tags        []string `json:"tags,omitempty"`
	Description string   `json:"description,omitempty"`
}

// MLFlowModelAPI defines the response object from the API
type MLFLowModelAPI struct {
	Name                 string   `json:"name"`
	CreationTimestamp    int64    `json:"creation_timestamp,omitempty"`
	LastUpdatedTimestamp int64    `json:"last_updated_timestamp,omitempty"`
	UserID               string   `json:"user_id,omitempty"`
	LatestVersions       []string `json:"latest_versions,omitempty"`
	Description          string   `json:"description,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
}

type MLFlowRegisteredModel struct {
	RegisteredModel MLFLowModelAPI `json:"registered_model"`
}

func (d *MLFLowModel) toAPIObject(schema map[string]*schema.Schema, data *schema.ResourceData) (*MLFLowModelAPI, error) {
	// Extract from ResourceData.
	if err := common.DataToStructPointer(data, schema, d); err != nil {
		return nil, err
	}

	// Copy to API object.
	var ad MLFLowModelAPI
	ad.Name = d.Name
	ad.Tags = append([]string{}, d.Tags...)

	return &ad, nil
}

func (d *MLFLowModel) fromAPIObject(ad *MLFLowModelAPI, schema map[string]*schema.Schema, data *schema.ResourceData) error {
	// Copy from API object.
	d.Name = ad.Name
	d.Tags = append([]string{}, ad.Tags...)

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
	return nil
}

// NewMLFlowAPI ...
func NewMLFlowAPI(ctx context.Context, m interface{}) MLFlowModelAPI {
	return MLFlowModelAPI{m.(*common.DatabricksClient), ctx}
}

// MLFlowModelAPI ...
type MLFlowModelAPI struct {
	client  *common.DatabricksClient
	context context.Context
}

// Create ...
func (a MLFlowModelAPI) Create(d *MLFLowModelAPI) error {
	return a.client.Post(a.context, "/mlflow/registered-models/create", d, &d)
}

// Read ...
func (a MLFlowModelAPI) Read(modelName string) (*MLFLowModelAPI, error) {
	//var d MLFLowModelAPI
	var d MLFlowRegisteredModel
	// need to figure out how to send param
	err := a.client.Get(a.context, fmt.Sprintf("/mlflow/registered-models/get?name=%s", modelName), nil, &d)
	if err != nil {
		return nil, err
	}

	return &d.RegisteredModel, nil
}

// Update ...
func (a MLFlowModelAPI) Update(modelName string, d *MLFLowModelAPI) error {
	return a.client.Patch(a.context, fmt.Sprintf("/mlflow/registered-models/update/%s", modelName), d)
}

// Delete ...
//func (a MLFlowModelAPI) Delete(modelName string, d *MLFLowModelAPI) error {
func (a MLFlowModelAPI) Delete(d *MLFLowModelAPI) error {
	return a.client.Delete(a.context, fmt.Sprintf("/mlflow/registered-models/delete"), d)
}

// ResourceDashboard ...
func ResourceMLFlowModel() *schema.Resource {
	s := common.StructToSchema(
		MLFLowModel{},
		func(m map[string]*schema.Schema) map[string]*schema.Schema {
			return m
		})

	return common.Resource{
		Create: func(ctx context.Context, data *schema.ResourceData, c *common.DatabricksClient) error {
			var d MLFLowModel
			ad, err := d.toAPIObject(s, data)
			if err != nil {
				return err
			}

			err = NewMLFlowAPI(ctx, c).Create(ad)
			if err != nil {
				return err
			}

			// No need to set anything because the resource is going to be
			// read immediately after being created.
			data.SetId(ad.Name)
			data.Set("name", ad.Name)
			return nil
		},
		Read: func(ctx context.Context, data *schema.ResourceData, c *common.DatabricksClient) error {
			ad, err := NewMLFlowAPI(ctx, c).Read(data.Id())
			if err != nil {
				return err
			}

			var d MLFLowModel
			return d.fromAPIObject(ad, s, data)
		},
		Update: func(ctx context.Context, data *schema.ResourceData, c *common.DatabricksClient) error {
			var d MLFLowModel
			ad, err := d.toAPIObject(s, data)
			if err != nil {
				return err
			}

			return NewMLFlowAPI(ctx, c).Update(data.Id(), ad)
		},
		Delete: func(ctx context.Context, data *schema.ResourceData, c *common.DatabricksClient) error {
			//return NewMLFlowAPI(ctx, c).Delete(data.Id())
			var d MLFLowModel
			ad, err := d.toAPIObject(s, data)
			if err != nil {
				return err
			}
			return NewMLFlowAPI(ctx, c).Delete(ad)

		},
		Schema: s,
	}.ToResource()
}
