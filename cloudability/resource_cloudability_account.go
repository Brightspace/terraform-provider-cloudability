package cloudability

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/matryer/try"
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
	attempts := 5

	var account CloudabilityAccount
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		account, err = client.verify(accountID)
		if err != nil || account.Verification.State != "verified" {
			time.Sleep(5 * time.Second)
		}
		return attempt < attempts, err
	})

	if err != nil {
		return err
	}

	if account.Verification.State != "verified" {
		return fmt.Errorf("The account failed to verify. %s", account.ID)
	}

	d.SetId(account.ID)

	return resourceAccountRead(d, meta)
}

func resourceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.CloudabilityClient

	account, err := client.get(d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	d.Set("role_name", account.Authorization.RoleName)
	d.Set("external_id", account.Authorization.ExternalID)

	return nil
}

func resourceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.CloudabilityClient

	_, err := client.delete(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
