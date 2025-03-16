// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/opentofu/opentofu/internal/registry"
	svchost "github.com/hashicorp/terraform-svchost"
)

func main() {
	// Create a new registry client
	client := registry.NewClient(nil, nil)
	
	// Parse the hostname
	host, err := svchost.ForComparison("registry.terraform.io")
	if err != nil {
		log.Fatalf("Error parsing hostname: %s", err)
	}
	
	// Set up logging
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	
	// Fetch providers
	ctx := context.Background()
	providers, err := client.BulkFetchProviders(ctx, host)
	if err != nil {
		log.Fatalf("Error fetching providers: %s", err)
	}
	
	fmt.Printf("Successfully fetched %d providers\n", len(providers))
	
	// Print the first 10 providers
	for i, provider := range providers {
		if i >= 10 {
			break
		}
		fmt.Printf("Provider %d: %s\n", i+1, provider.Name)
	}
}
