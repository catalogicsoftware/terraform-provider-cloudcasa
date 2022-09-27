package cloudcasa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"context"
	"strconv"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKubeclusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubeclustersRead,
		Schema: map[string]*schema.Schema{
			"kubeclusters": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     	schema.TypeString,
							Computed:	true,
						},
						"name": &schema.Schema{
							Type:     	schema.TypeString,
							Computed:	true,
						},
						"cc_user_email": &schema.Schema{
							Type:		schema.TypeString,
							Optional:	true,
						},
						"updated": &schema.Schema{
							Type:		schema.TypeString,
							Computed:	true,
						},
						"created": &schema.Schema{
							Type:		schema.TypeString,
							Computed:	true,
						},
						"etag": &schema.Schema{
							Type:		schema.TypeString,
							Computed:	true,
						},
						"org_id": &schema.Schema{
							Type:		schema.TypeString,
							Optional: 	true,
						},
						"status": &schema.Schema{
							// Cant be typeset?
							Type:		schema.TypeMap,
							Optional:	true,
							Elem:		&schema.Schema{Type: schema.TypeString},
							// Elem:		&schema.Resource{
							// 	Schema:	map[string]*schema.Schema{
							// 		"agenturl": {
							// 			Type:		schema.TypeString,
							// 			Optional: 	true,
							// 		},
							// 		"update_required": {
							// 			Type:		schema.TypeBool,
							// 			Optional: 	true,
							// 		},
							// 		"dormant": {
							// 			Type:		schema.TypeBool,
							// 			Optional: 	true,
							// 		},
							// 		"state": {
							// 			Type:		schema.TypeString,
							// 			Optional: 	true,
							// 		},
							// 	},
							// },
						},
						"links": &schema.Schema{
							// Type cannot be TypeSet
							// But TypeMap cannot have schema.Resource underneath..
							Type:		schema.TypeMap,
							Optional:	true,
							// trying schema.Schema instead of schema.Resource
							Elem:	 	&schema.Schema{Type: schema.TypeMap},
							// Elem:		&schema.Schema{
							// 	Schema: map[string]*schema.Schema{
							// 		"self": {
							// 			Type:		schema.TypeMap,
							// 			Optional: 	true,
							// 		},
							// 		"parent": {
							// 			Type:		schema.TypeMap,
							// 			Optional: 	true,
							// 		},
							// 		"collection": {
							// 			Type:		schema.TypeMap,
							// 			Optional: 	true,
							// 		},
							// 	},
							// },
						},
						// "_links": &schema.Schema{
						// 	Type:		schema.TypeMap,
						// 	Optional:	true,
						// },
					},
				},
			},
		},
	}
}

// testing using Client in data source function
// func (c *Client) GetKubeclusters() ([]Kubeclusters, error) {
// 	req, err := http.NewRequest("GET", fmt.Sprintf("%s/kubeclusters", c.ApiURL), nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	body, err := c.doRequest(req, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	kubeclusters := []Kubeclusters{}
// 	err = json.Unmarshal(body, &ubeclusters)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return kubeclusters, nil
// }

type KubeclusterObjs struct {
	Items 				[]KubeclusterObj	`json:"_items"`
}

type KubeclusterObj struct {
	Id					string		`json:"_id"`
	Name				string		`json:"name"`
	Cc_user_email		string		`json:"cc_user_email"`
	Updated				string		`json:"_updated"`
	Created				string		`json:"_created"`
	Etag				string		`json:"_etag"`
	Org_id				string		`json:"org_id"`
	// Status 				struct {}	`json:"status"`
	// Links 				struct {}	`json:"_links"`
}


func dataSourceKubeclustersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: use cloudcasa go client for these requests

	client := &http.Client{Timeout: 10 * time.Second}

	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/kubeclusters", "https://api.staging.cloudcasa.io/api/v1"), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InM3NmtuNThRT2liTXRfZnNpVFlLMCJ9.eyJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9jb3VudHJ5IjoiVW5pdGVkIFN0YXRlcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL3RpbWV6b25lIjoiQW1lcmljYS9OZXdfWW9yayIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvdW50cnlfY29kZSI6IlVTIiwiaHR0cDovL3d3dy5jbG91ZGNhc2EuaW8vY291bnRyeV9jb2RlMyI6IlVTQSIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2ZpcnN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9sYXN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9qb2JUaXRsZSI6IkRldm9wcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvbXBhbnkiOiJDYXRhbG9naWMgU29mdHdhcmUiLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9hd3NfbWFya2V0cGxhY2VfdG9rZW4iOiItIiwibmlja25hbWUiOiJqZ2FybmVyIiwibmFtZSI6IkpvbmF0aGFuIEdhcm5lciIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8yOTlhNmJhNjhlNjEwOGFiYjY1MmY4ZTkwZTM0YjVhNj9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmpnLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDIyLTA5LTE0VDE2OjM4OjUzLjI4MVoiLCJlbWFpbCI6ImpnYXJuZXJAY2F0YWxvZ2ljc29mdHdhcmUuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vYXV0aC5jbG91ZGNhc2EuaW8vIiwic3ViIjoiYXV0aDB8NWZhYzQ4NDg0MWQ3MDgwMDY4YTA2ZGM5IiwiYXVkIjoiSkVKU3plblBGeE5FUFEwaDY0ZDIzZTZRMEdKNXpRanQiLCJpYXQiOjE2NjMxNzM1NTMsImV4cCI6MTY2MzE4MDc1MywiYWNyIjoiaHR0cDovL3NjaGVtYXMub3BlbmlkLm5ldC9wYXBlL3BvbGljaWVzLzIwMDcvMDYvbXVsdGktZmFjdG9yIiwiYW1yIjpbIm1mYSJdLCJzaWQiOiJpTmJ6elRDZG9RNVdYN3FXTEhQekRSTUFZNkJuaWVrSyIsIm5vbmNlIjoiWDBFelNta3hNVloyVlcxK2RuNUZPV0ZVTm5kU1RsOW5VbUpYVHpaUU1FTjNaVlpNYnpNeFVHZFRhQT09In0.w5KC0dq_WKBCXnYiqE8q0VRdIxJf3Jod9MLHKb7n3iKd7xrZpVVQ7fTdUaPMNcqLle2UPh4JfXFXcP5BI9T_gomG33SgbF-IaOGKyL9fzNe5naC179m7wYw7ntHJt9JQ5BQ-qdua3CaOl_PNzScyNRfabX4mdpvnopJxplcJ1_TDRvBUGFeymvBx1De4CFgB0eIAHcYtJdtV745OZGLKsNiFoj9TPzD5L5tHYs0IISJInkQnr2ppU1h8-urai-SAeOaYt-PyJud_EeJ8qWgYtEtuukf-DW_SAxf4I49opi4NBQMMXFApv2B-wyP2Q5-xCwo04mVYjVk8EUfH8_rpeg")

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	//var kubeclusters map[string]interface{}

	// Response will be a dict containing an array of KubeclusterObj
	kubeclusters := KubeclusterObjs{}
	err = json.NewDecoder(r.Body).Decode(&kubeclusters)
	if err != nil {
		return diag.FromErr(err)
	}

	// diags = append(diags, diag.Diagnostic{
	// 	Severity: 	diag.Error,
	// 	Summary:  	"kubeclusters.Items[0].Name",
	// 	Detail:  	kubeclusters.Items[0].Name,
	// })
	// return diags

	//if err := d.Set("kubeclusters", kubeclusters["_items"]); err != nil {
	if err := d.Set("kubeclusters", kubeclusters.Items); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags

}
