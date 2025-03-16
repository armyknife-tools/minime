// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/opentofu/opentofu/internal/templates"
)

func main() {
	// Parse command line flags
	dbType := flag.String("db", "sqlite", "Database type: sqlite or postgres")
	dbPath := flag.String("path", "", "Database path (for SQLite)")
	flag.Parse()

	var path string
	if *dbType == "sqlite" {
		// If no path is provided, use the default location
		if *dbPath == "" {
			configDir := filepath.Join(os.Getenv("HOME"), ".opentofu")
			if err := os.MkdirAll(configDir, 0755); err != nil {
				log.Fatalf("Error creating config directory: %v", err)
			}
			path = filepath.Join(configDir, "templates.db")
		} else {
			path = *dbPath
		}
	} else if *dbType == "postgres" {
		// For PostgreSQL, we use environment variables for connection details
		path = ""
	} else {
		log.Fatalf("Unsupported database type: %s", *dbType)
	}

	// Load templates into the database
	fmt.Printf("Loading templates into %s database...\n", *dbType)
	err := templates.LoadTemplates(*dbType, path)
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	fmt.Println("Templates loaded successfully!")
}
