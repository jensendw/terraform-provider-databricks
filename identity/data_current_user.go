package identity

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var nonAlphanumeric = regexp.MustCompile(`\W`)

// DataSourceCurrentUser returns information about caller identity
func DataSourceCurrentUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"home": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repos": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alphanumeric": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			usersAPI := NewUsersAPI(ctx, m)
			me, err := usersAPI.Me()
			if err != nil {
				return diag.FromErr(err)
			}
			d.Set("user_name", me.UserName)
			d.Set("home", fmt.Sprintf("/Users/%s", me.UserName))
			d.Set("repos", fmt.Sprintf("/Repos/%s", me.UserName))
			d.Set("external_id", me.ExternalID)
			splits := strings.Split(me.UserName, "@")
			norm := nonAlphanumeric.ReplaceAllLiteralString(splits[0], "_")
			norm = strings.ToLower(norm)
			d.Set("alphanumeric", norm)
			d.SetId(me.ID)
			return nil
		},
	}
}
