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
	counter := 100

	client := api.Cloudability{
		Credentials: api.Credentials{
			APIKey: []byte(os.Getenv("CLOUDABILITY_TOKEN")),
		},
		RetryMaximum: 5,
	}

	for i := 0; i < counter; i++ {
		result, err := client.Get(id)
		if err != nil {
			fmt.Println("Encountered an error:\n", err)
		} else if result == nil {
			fmt.Println("404 not found")
		} else if result.ID != id {
			fmt.Println("ID mismatch:\n", result.ID)
		} else {
			// Ignore sometimes
			// fmt.Println("ID:\n", result.ID)
		}
	}
}