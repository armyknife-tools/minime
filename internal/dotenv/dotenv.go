// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

// Package dotenv provides functionality for loading environment variables from .env files
package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Load loads environment variables from a .env file
// If path is empty, it looks for .env in the current directory and binary directory
// Variables are loaded in the format KEY=VALUE
// Returns a map of the loaded variables
func Load(path string) (map[string]string, error) {
	loadedVars := make(map[string]string)
	
	// If the path is empty, look for .env in multiple locations
	if path == "" {
		// Try current directory first
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("error getting current working directory: %v", err)
		}
		cwdEnvPath := filepath.Join(cwd, ".env")
		
		// Try to load from current directory
		if vars, err := loadFile(cwdEnvPath); err == nil {
			for k, v := range vars {
				loadedVars[k] = v
			}
			return loadedVars, nil
		}
		
		// If not found in current directory, try binary directory
		if exePath, err := os.Executable(); err == nil {
			binaryDir := filepath.Dir(exePath)
			binaryEnvPath := filepath.Join(binaryDir, ".env")
			
			// Try to load from binary directory
			if vars, err := loadFile(binaryEnvPath); err == nil {
				for k, v := range vars {
					loadedVars[k] = v
				}
				return loadedVars, nil
			}
		}
		
		// Also try one directory up from binary (for development scenarios)
		if exePath, err := os.Executable(); err == nil {
			parentDir := filepath.Dir(filepath.Dir(exePath))
			parentEnvPath := filepath.Join(parentDir, ".env")
			
			// Try to load from parent directory
			if vars, err := loadFile(parentEnvPath); err == nil {
				for k, v := range vars {
					loadedVars[k] = v
				}
				return loadedVars, nil
			}
		}
		
		// No .env file found in any location
		return loadedVars, nil
	}

	// If path is specified, load from that path
	vars, err := loadFile(path)
	if err != nil {
		return loadedVars, err
	}
	for k, v := range vars {
		loadedVars[k] = v
	}
	
	return loadedVars, nil
}

// loadFile loads environment variables from a specific file path
func loadFile(path string) (map[string]string, error) {
	loadedVars := make(map[string]string)
	
	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", path)
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening .env file: %v", err)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) > 1 && (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		// Set the environment variable
		if err := os.Setenv(key, value); err != nil {
			return loadedVars, fmt.Errorf("error setting environment variable %s: %v", key, err)
		}
		
		loadedVars[key] = value
	}

	if err := scanner.Err(); err != nil {
		return loadedVars, fmt.Errorf("error reading .env file: %v", err)
	}

	return loadedVars, nil
}

// LoadAll loads environment variables from multiple .env files
// Files are processed in order, with later files overriding earlier ones
// Returns a map of all loaded variables
func LoadAll(paths ...string) (map[string]string, error) {
	allVars := make(map[string]string)
	
	// If no paths are provided, use default .env in current directory
	if len(paths) == 0 {
		vars, err := Load("")
		if err != nil {
			return allVars, err
		}
		for k, v := range vars {
			allVars[k] = v
		}
		return allVars, nil
	}
	
	// Process each file in order
	for _, path := range paths {
		vars, err := Load(path)
		if err != nil {
			return allVars, err
		}
		for k, v := range vars {
			allVars[k] = v
		}
	}
	
	return allVars, nil
}

// GetWithDefault gets an environment variable or returns a default value if not set
func GetWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
