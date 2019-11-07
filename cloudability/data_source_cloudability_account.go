package cloudability

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAccountRead,

		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AWS Account ID",
				ForceNew:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The AWS Account ID",
			},
			"parent_account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The AWS Account ID",
			},
			"role_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The AWS Account ID",
			},
			"external_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The AWS Account ID",
			},
		},
	}
}

func dataSourceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.CloudabilityClient

	accountID := d.Get("account_id").(string)

	log.Printf("[DEBUG] data_account : (ID: %q ROOT: %q)", accountID, config.PayerAccountID)
	account, err := client.Pull(config.PayerAccountID, accountID)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] verification state : (ID: %q STATE: %q)", accountID, account.Verification.State)
	if account.Verification.State == "" {
		account, err = client.Add(accountID)
		if err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] account parameters: (Parent: %q, Name: %q, ExternalID: %q)", account.ParentAccountID, account.Authorization.RoleName, account.Authorization.ExternalID)
	d.SetId(account.ID)
	d.Set("name", account.Name)
	d.Set("parent_account_id", account.ParentAccountID)
	d.Set("role_name", account.Authorization.RoleName)
	d.Set("external_id", account.Authorization.ExternalID)

	return nil
}
