package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/tapirs/terraform-provider-aws-elb-to-lb/aws-elb-to-lb"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: aws.Provider})
}
