package cloudcasa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"context"
	"strconv"
	"os/exec"
	
	"terraform-provider-cloudcasa/cloudcasa/handler"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKubecluster() *schema.Resource {
	return &schema.Resource{
		CreateContext:	resourceKubeclusterCreate,
		ReadContext:	resourceKubeclusterRead,
		UpdateContext:	resourceKubeclusterUpdate,
		DeleteContext:	resourceKubeclusterDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     	schema.TypeString,
				Required:	true,
			},
			"auto_install": &schema.Schema{
				Type:		schema.TypeBool,
				Optional:	true,
				Default:	false,
			},
			"id": &schema.Schema{
				Type:     	schema.TypeString,
				Computed:	true,
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
				Computed: 	true,
			},
			// TODO: REQUIRE status GET each time we apply changes,
			// so that we can verify that the cluster is ACTIVE before anything else?
			// counterpoint: consider situations where user might want to apply 
			// changes to a cluster that is pending or unavailable.
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

// TODO: implement READ DESTROY UPDATE in that order is easiest
// https://www.hashicorp.com/blog/writing-custom-terraform-providers simple overview

func resourceKubeclusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//c := m.(*hc.Client)
	
	// TODO: use apikey from terraform config
	var diags diag.Diagnostics

	// Create kubecluster in cloudcasa
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

	// Set fields in resourceData 'd'
	// TODO: set Links and Status separately because they are maps
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

	// if auto_install is false return now. Otherwise proceed with agent installation
	if d.Get("auto_install") == false {
		return diags
	}

	var getKubeclusterResp *handler.GetKubeclusterResp
	var kubeclusterStatus handler.KubeclusterStatus

	// wait 1m for agent URL
	for i:=1; i<12; i++ {
		getKubeclusterResp = handler.GetKubecluster(createKubeclusterResp.Id)
		kubeclusterStatus = getKubeclusterResp.Status
		if len(kubeclusterStatus.Agent_url) > 0 {
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
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to retrieve CloudCasa Agent manifest",
			Detail:   fmt.Sprintf("%s", err),
		})
		return diags
	}

	// TODO: add tip to make sure kubeconfig env var is set?
	// OR we can accept kubeconfig as an input option?
	kubectlCmd := exec.Command("kubectl",  "apply",  "-f",  fmt.Sprintf("%s", kubeclusterStatus.Agent_url))
	_, err := kubectlCmd.Output()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to apply kubeagent manifest",
			Detail:   fmt.Sprintf("%s", err),
		})
		return diags
	}

	// Now wait for cluster to be ACTIVE
	// Wait 5min?
	for i:=1; i<60; i++ {
		getKubeclusterResp = handler.GetKubecluster(createKubeclusterResp.Id)
		kubeclusterStatus = getKubeclusterResp.Status
		if kubeclusterStatus.State == "ACTIVE"{
			break
		}
		time.Sleep(5 * time.Second)
	}

	// Check if state was set to ACTIVE
	if kubeclusterStatus.State != "ACTIVE" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "CloudCasa Agent Installation Failed",
			Detail:   fmt.Sprintf("Timed out waiting for cluster to reach ACTIVE state"),
		})
		return diags
	}

	// if err := d.Set("links", createKubeclusterResp.Links); err != nil {
	// 	return diag.FromErr(err)
	// }
	// if err := d.Set("status", createKubeclusterResp.Status); err != nil {
	// 	return diag.FromErr(err)
	// }

	return diags
}

func resourceKubeclusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceKubeclusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceKubeclusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}