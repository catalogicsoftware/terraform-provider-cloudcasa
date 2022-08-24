package cloudcasa

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKubeclusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: GetKubeclusters,
		Schema: map[string]*schema.Schema{
			"kubeclusters": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespaces": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"namespace_name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (c *Client) GetKubeclusters() ([]kubeclusters, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/kubeclusters", c.ApiURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	kubeclusters := []kubeclusters{}
	err = json.Unmarshal(body, &kubeclusters)
	if err != nil {
		return nil, err
	}

	return kubeclusters, nil
}

// func dataSourceKubeclustersRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	// TODO: use cloudcasa go client for these requests

// 	client := httpClient{Timeout: 10 * time.Second}

// 	var diags diag.Diagnostics

// 	req, err := http.NewRequest("GET", fmt.Sprintf("%s/kubeclusters", "https://api.staging.cloudcasa.io"), nil)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	r, err := client.Do(req)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	defer r.Body.Close()

// 	kubeclusters := make([]map[string]interface{}, 0)
// 	err = json.NewDecoder(r.Body).Decode(&kubeclusters)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	if err := d.Set("kubeclusters", kubeclusters); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	// always run
// 	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

// }
