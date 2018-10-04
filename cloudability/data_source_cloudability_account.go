package cloudability

import (
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

	account, err := client.pull(config.PayerAccountID, accountID)
	if err != nil {
		return err
	}

	if account.Verification.State == "" {
		account, err = client.add(accountID)
		if err != nil {
			return err
		}
	}

	d.SetId(account.ID)
	d.Set("name", account.Name)
	d.Set("parent_account_id", account.ParentAccountID)
	d.Set("role_name", account.Authorization.RoleName)
	d.Set("external_id", account.Authorization.ExternalID)

	return nil
}
