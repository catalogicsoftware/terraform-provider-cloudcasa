package cloudcasa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"context"
	"strconv"
	
	"terraform-provider-cloudcasa/cloudcasa/handler"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKubecluster() *schema.Resource {
	return &schema.Resource{
		//ReadContext: 	dataSourceKubeclustersRead,
		CreateContext:	resourceCreateKubecluster,
		ReadContext:	resourceReadKubecluster,
		UpdateContext:	resourceUpdateKubecluster,
		DeleteContext:	resourceDeleteKubecluster,

		Schema: map[string]*schema.Schema{
			// "kubecluster": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     	schema.TypeString,
				Computed:	true,
			},
			"name": &schema.Schema{
				Type:     	schema.TypeString,
				Required:	true,
			},
			"cc_user_email": &schema.Schema{
				Type:		schema.TypeString,
				Computed:	true,
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
				Computed:	true,
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
				Computed:	true,
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
			// 		},
			// 	},
			// },
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

// TODO: we should try to follow the structure commvault uses
// eg. return nill for Read - check resource_aws_storage
func dataSourceKubeclustersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: use cloudcasa go client for these requests

	client := &http.Client{Timeout: 10 * time.Second}

	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/kubeclusters", "https://api.staging.cloudcasa.io/api/v1"), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InM3NmtuNThRT2liTXRfZnNpVFlLMCJ9.eyJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9jb3VudHJ5IjoiVW5pdGVkIFN0YXRlcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL3RpbWV6b25lIjoiQW1lcmljYS9OZXdfWW9yayIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvdW50cnlfY29kZSI6IlVTIiwiaHR0cDovL3d3dy5jbG91ZGNhc2EuaW8vY291bnRyeV9jb2RlMyI6IlVTQSIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2ZpcnN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9sYXN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9qb2JUaXRsZSI6IkRldm9wcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvbXBhbnkiOiJDYXRhbG9naWMgU29mdHdhcmUiLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9hd3NfbWFya2V0cGxhY2VfdG9rZW4iOiItIiwibmlja25hbWUiOiJqZ2FybmVyIiwibmFtZSI6IkpvbmF0aGFuIEdhcm5lciIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8yOTlhNmJhNjhlNjEwOGFiYjY1MmY4ZTkwZTM0YjVhNj9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmpnLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDIyLTEwLTAzVDE0OjQzOjMxLjUxNVoiLCJlbWFpbCI6ImpnYXJuZXJAY2F0YWxvZ2ljc29mdHdhcmUuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vYXV0aC5jbG91ZGNhc2EuaW8vIiwic3ViIjoiYXV0aDB8NWZhYzQ4NDg0MWQ3MDgwMDY4YTA2ZGM5IiwiYXVkIjoiSkVKU3plblBGeE5FUFEwaDY0ZDIzZTZRMEdKNXpRanQiLCJpYXQiOjE2NjQ4MDgyMTMsImV4cCI6MTY2NDgxNTQxMywic2lkIjoiZEVNdkJQd1NjS0pnRW02NmVMcXBUaHNfVnQ4eFRYS3MiLCJub25jZSI6ImFqSmZlRmhyTW5CS05reFFlazFpT1c0MGRTNVNhRnB1TUdWbWJtZDBjemxtYVU1YU0zSlFWWGxYV1E9PSJ9.XE1p0KI29twtr1vTmU09aQIhM_G-PmNi1skZNWIZToTy4meoIpMuIR4w5WScGEx7HdC_VE6IpIu3g7a2FOoSOPS82G-yF9y9cKIkyl5D0VydRoMJHVhRIJnVMtduWXle94I8gSsgEOaYBNxbjyRvsv-e2r7Z9wrEHLAGdmYp0OJwi-J_tEyxkuWXZrv9CLfhJAEeLSqZzRJDm4nXcZadQie1122kiCT2R2P5ZocqbC8sE3pcPGRDafru3VainT03qaCLAK9ae8qH3MIs9gW7PAJrjFZSDZI514pcoy6QON1TfloWVRwtDra3GXAiHWLMlL8KZ6dF49oPoahXA5VlEA")

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

func resourceCreateKubecluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Steps (?):
	// Populate struct with user input
	// Post struct to apiserver?
	// set id? terraform ID or CC ID?
	// Return ???

	var diags diag.Diagnostics
	var kubeclusterData handler.CreateKubeclusterReq

	kubeclusterData.Name = d.Get("name").(string)

	kubeclusterResp := handler.CreateKubecluster(kubeclusterData)

	// TODO: Better failure check
	// test by using wrong auth token
	status := kubeclusterResp.Status
	if status != "OK" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create Kubecluster",
			Detail:   "Received status NOT OK",
		})
		return diags
	}

	d.SetId(kubeclusterResp.Id)

	return diags
}

func resourceReadKubecluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceUpdateKubecluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceDeleteKubecluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}