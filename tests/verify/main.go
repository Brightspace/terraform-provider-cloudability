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
	id := os.Args[1]

	client := api.Cloudability{
		Credentials: api.Credentials{
			APIKey: []byte(os.Getenv("CLOUDABILITY_TOKEN")),
		},
		RetryMaximum: 5,
	}

	result, err := client.Verification(id)
	if result == nil {
		fmt.Println("id could not be found:\n", id)
		return
	}

	if err != nil {
		fmt.Println("err:\n", err)
		return
	}
	fmt.Println("ID:\n", result.ID)
	fmt.Println("State:\n", result.Verification.State)
	fmt.Println("LastVerificationAttemptedAt:\n", result.Verification.LastVerificationAttemptedAt)
}
