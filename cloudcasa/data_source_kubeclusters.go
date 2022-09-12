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
						"_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						// "namespaces": &schema.Schema{
						// 	Type:     schema.TypeList,
						// 	Computed: true,
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"namespace_name": &schema.Schema{
						// 				Type:     schema.TypeString,
						// 				Computed: true,
						// 			},
						// 		},
						// 	},
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

func dataSourceKubeclustersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: use cloudcasa go client for these requests

	client := &http.Client{Timeout: 10 * time.Second}

	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/kubeclusters", "https://api.staging.cloudcasa.io/v1"), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	kubeclusters := make([]map[string]string, 0)
	//var kubeclusters interface{}
	err = json.NewDecoder(r.Body).Decode(&kubeclusters)
	if err != nil {
		return diag.FromErr(err)
	}

	fmt.Println("Map: ", kubeclusters)
	fmt.Println("len: ", len(kubeclusters))

	if err := d.Set("kubeclusters", kubeclusters); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags

}
