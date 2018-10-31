package cloudability

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountCreate,
		Read:   resourceAccountRead,
		Delete: resourceAccountDelete,

		Schema: map[string]*schema.Schema{
			"role_arn": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Amazon Resource Name for the IAM role",
				ForceNew:    true,
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AWS Account ID",
				ForceNew:    true,
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

func resourceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.CloudabilityClient

	accountID := d.Get("account_id").(string)

	log.Printf("[DEBUG] account verify: (ID: %q)", accountID)
	account, err := client.verify(accountID)
	if err != nil {
		return err
	}

	if account.Verification.State != "verified" {
		return fmt.Errorf("The account failed to verify. %s", account.ID)
	}

	log.Printf("[DEBUG] account added: (ID: %q)", account.ID)
	d.SetId(account.ID)

	return resourceAccountRead(d, meta)
}

func resourceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.CloudabilityClient

	log.Printf("[DEBUG] account get: (ID: %q)", d.Id())
	account, err := client.get(d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	log.Printf("[DEBUG] account read: (ARN: %q, Name: %q, ExternalID: %q)", account.Authorization.RoleName, account.Authorization.ExternalID)
	d.Set("role_name", account.Authorization.RoleName)
	d.Set("external_id", account.Authorization.ExternalID)

	return nil
}

func resourceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.CloudabilityClient

	log.Printf("[DEBUG] account delete: (ID: %q)", d.Id())
	_, err := client.delete(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
