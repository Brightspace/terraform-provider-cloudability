package main

import (
	"fmt"
	"github.com/Brightspace/terraform-provider-cloudability/cloudability/api"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing ID for AWS account")
	}
	arg := os.Args[1]

	client := api.Cloudability{
		Credentials: api.Credentials{
			APIKey: []byte(os.Getenv("CLOUDABILITY_TOKEN")),
		},
		RetryMaximum: 5,
	}

	result, err := client.Poll(arg)
	if result == nil {
		fmt.Println("id could not be found:\n", arg)
		return
	}

	if err != nil {
		fmt.Println("err:\n", err)
		return
	}
	fmt.Println("ID:\n", result.ID)
	fmt.Println("Name:\n", result.Name)
	fmt.Println("AccountID:\n", result.AccountID)
	fmt.Println("ParentAccountID:\n", result.ParentAccountID)
	fmt.Println("VendorKey:\n", result.VendorKey)
	fmt.Println("State:\n", result.Verification.State)
	fmt.Println("LastVerificationAttemptedAt:\n", result.Verification.LastVerificationAttemptedAt)
	fmt.Println("Type:\n", result.Authorization.Type)
	fmt.Println("RoleName:\n", result.Authorization.RoleName)
	fmt.Println("ExternalID:\n", result.Authorization.ExternalID)
}
