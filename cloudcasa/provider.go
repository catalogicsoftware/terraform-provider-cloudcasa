package cloudcasa

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CC_EMAIL", os.Getenv("CC_EMAIL")),
				Description: "The email address of your CloudCasa user.",
			},
			"apikey": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CC_ACCESS_TOKEN", os.Getenv("CC_ACCESS_TOKEN")),
				Description: "The CloudCasa API key used to authenticate against the CloudCasa API.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cloudcasa_kubeclusters":	dataSourceKubeclusters(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// Perform validation here - this function will run when we 'terraform init'
func providerConfigure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	email := data.Get("email").(string)
	idToken := data.Get("apikey").(string)

	// TODO: check simple casa commands with supplied email/token
	// to validate login.
	// This can also be done in the Client library NewClient()

	// Warnings/errors are collected by this type
	var diags diag.Diagnostics

	// If email and idToken are supplied, create a new CC client
	if (email != "") && (idToken != "") {
		c, err := NewClient(&email, &idToken)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create CloudCasa client",
				Detail:   "Unable to authenticate user with supplied CloudCasa credentials",
			})

			return nil, diags
		}

		return c, diags
	}

	c, err := NewClient(nil, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create CloudCasa client",
			Detail:   "Unable to create anonymous CloudCasa client - credentials are missing.",
		})
		return nil, diags
	}

	return c, diags
}
