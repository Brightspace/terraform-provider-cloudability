package cloudability

import (
	"github.com/Brightspace/terraform-provider-cloudability/cloudability/api"
)

type Config struct {
	CloudabilityClient api.Cloudability
	PayerAccountID     string
}
