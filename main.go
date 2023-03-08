package main

import (
	"context"
	"terraform-provider-cloudcasa/cloudcasa"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), cloudcasa.New, providerserver.ServeOpts{
		Address: "cloudcasa.io/cloudcasa/cloudcasa",
	})
}
