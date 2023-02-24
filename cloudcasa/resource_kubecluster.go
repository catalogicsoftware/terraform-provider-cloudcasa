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
		CreateContext:	resourceCreateKubecluster,
		ReadContext:	resourceReadKubecluster,
		UpdateContext:	resourceUpdateKubecluster,
		DeleteContext:	resourceDeleteKubecluster,

		Schema: map[string]*schema.Schema{
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
			"agent_url": &schema.Schema{
				Type:		schema.TypeString,
				Computed:	true,
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

	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InM3NmtuNThRT2liTXRfZnNpVFlLMCJ9.eyJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9jb3VudHJ5IjoiVW5pdGVkIFN0YXRlcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL3RpbWV6b25lIjoiQW1lcmljYS9OZXdfWW9yayIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvdW50cnlfY29kZSI6IlVTIiwiaHR0cDovL3d3dy5jbG91ZGNhc2EuaW8vY291bnRyeV9jb2RlMyI6IlVTQSIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2ZpcnN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9sYXN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9qb2JUaXRsZSI6IkRldm9wcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvbXBhbnkiOiJDYXRhbG9naWMgU29mdHdhcmUiLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9hd3NfbWFya2V0cGxhY2VfdG9rZW4iOiItIiwibmlja25hbWUiOiJqZ2FybmVyIiwibmFtZSI6IkpvbmF0aGFuIEdhcm5lciIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8yOTlhNmJhNjhlNjEwOGFiYjY1MmY4ZTkwZTM0YjVhNj9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmpnLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDIzLTAyLTI0VDE2OjUxOjU0LjQzMVoiLCJlbWFpbCI6ImpnYXJuZXJAY2F0YWxvZ2ljc29mdHdhcmUuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vYXV0aC5jbG91ZGNhc2EuaW8vIiwiYXVkIjoiSkVKU3plblBGeE5FUFEwaDY0ZDIzZTZRMEdKNXpRanQiLCJpYXQiOjE2NzcyNTc1MTYsImV4cCI6MTY3NzI2NDcxNiwic3ViIjoiYXV0aDB8NWZhYzQ4NDg0MWQ3MDgwMDY4YTA2ZGM5Iiwic2lkIjoiZzFabm5TcUJnSVgzWmMya1ZJMEU3TkRIWXhHcnYtSHYiLCJub25jZSI6ImVtOU9NR3RpYTBaRGJHMU9SMkpDWTFBd1RWRjFibXRVYkdWWFJuSmpibmt4VUV0cVp5MXplbVJGZEE9PSJ9.mFH80HRANlJGo3Ti1eQ9KIfiEBjy_QCW7T-NTtRpG0nwd7o1BqWEu4TXIaMt515fsU6TylMpsOa5WRmS6AryyNVOn9rWhOCPy4SjjJQ3LIm5Ewl-imlo87eq7uM4fXLxVHi0quiHhGBIh3jTgSHnGXiCVXLHt1qgP-sfYUAzu9nkvjOv2bjrLHLBuiKpF146tw7O0kVsHfsjoABSQeOJwuHLGhaRdqq8wX5msh7rx6yurJhX2-7Oy4y90HkAmHy52SN9yCqFE0m0Kula8t5CG4SqLfkd-eX9OiCXpa4elTrXTITD0uQHSUjDfrYpRR3_Hku-cLj8plbeTv_d2r8E1g")

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
	var diags diag.Diagnostics

	var kubeclusterData handler.CreateKubeclusterReq
	kubeclusterData.Name = d.Get("name").(string)
	createKubeclusterResp := handler.CreateKubecluster(kubeclusterData)

	// TODO: Better failure check
	status := createKubeclusterResp.Status
	if status != "OK" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create Kubecluster",
			Detail:   "Received status NOT OK",
		})
		return diags
	}

	d.SetId(createKubeclusterResp.Id)

	// At this point the cluster resource is created, but resp does not
	// have the agent installation URL. So GET the kubecluster..
	// It takes a while for the URL to be available, so loop until
	// URL is ready (1 min?)

	var getKubeclusterResp *handler.GetKubeclusterResp
	var kubeclusterStatus handler.KubeclusterStatus

	for i:=1; i<12; i++ {
		getKubeclusterResp = handler.GetKubecluster(createKubeclusterResp.Id)
		kubeclusterStatus = getKubeclusterResp.Status
		if len(kubeclusterStatus.Agent_url) > 0{
			break
		}
		time.Sleep(5 * time.Second)
	}

	// Check that Agent URL was fetched successfully
	if len(kubeclusterStatus.Agent_url) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to get Agent URL for Kubecluster",
			Detail:   "Timed out waiting for Agent URL",
		})
		return diags
	}

	// try to set agent url from GET response
	if err := d.Set("agent_url", kubeclusterStatus.Agent_url); err != nil {
		return diag.FromErr(err)
	}

	// Set fields in resourceData 'd'
	// all string fields at once
	// TODO: set Links and Status below
	var kubeclusterFieldsMap = map[string]string{
		"cc_user_email": createKubeclusterResp.Cc_user_email,
		"created": createKubeclusterResp.Created,
		"etag": createKubeclusterResp.Etag,
		"org_id": createKubeclusterResp.Org_id,
		"updated": createKubeclusterResp.Updated,
	} 

	for k,v := range kubeclusterFieldsMap{
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	// if err := d.Set("links", createKubeclusterResp.Links); err != nil {
	// 	return diag.FromErr(err)
	// }
	// if err := d.Set("status", createKubeclusterResp.Status); err != nil {
	// 	return diag.FromErr(err)
	// }

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