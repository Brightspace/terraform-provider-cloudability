package cloudability

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the Cloudability personal access token. It must be provided.",
			},
			"payer_account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the Cloudability personal access token. It must be provided.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cloudability_account": dataSourceAccount(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"cloudability_account": resourceAccount(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := Cloudability{
		Credentials: Credentials{
			APIKey: []byte(d.Get("api_key").(string)),
		},
	}
	config := Config{
		CloudabilityClient: client,
		PayerAccountID:     d.Get("payer_account_id").(string),
	}

	return &config, nil
}