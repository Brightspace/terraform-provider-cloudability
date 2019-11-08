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

	resp, _ := client.Delete(id)
	fmt.Println(resp)

	result, _ := client.Get(id)
	if result.Verification.State == "" {
		fmt.Println("success:\n", id)
	} else {
		fmt.Println("fail:\n", result.Verification.State)
	}
}
