package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/opentofu/opentofu/internal/addrs"
	"github.com/opentofu/opentofu/internal/command"
	"github.com/posener/complete"
)

type RegistryProviderInstallCommand struct {
	Meta command.Meta
}

func (c *RegistryProviderInstallCommand) Run(args []string) int {
	// Create a new flag set for this command
	flags := flag.NewFlagSet("registry provider install", flag.ContinueOnError)
	flags.Usage = func() { c.Meta.Ui.Error(c.Help()) }
	
	var providerVersion string
	var autoApprove bool

	flags.StringVar(&providerVersion, "version", "", "provider version to install")
	flags.BoolVar(&autoApprove, "auto-approve", false, "skip interactive approval of provider installation")

	if err := flags.Parse(args); err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		c.Meta.Ui.Error("Please specify a provider to install")
		return 1
	}

	// Parse the provider source
	providerSource := args[0]
	
	// Parse the provider source using the addrs package
	provider, diags := addrs.ParseProviderSourceString(providerSource)
	if diags.HasErrors() {
		c.Meta.Ui.Error(fmt.Sprintf("Invalid provider format: %s", diags.Err()))
		c.Meta.Ui.Error("Expected format: namespace/name or hostname/namespace/name")
		return 1
	}
	
	// If no version is specified, use latest
	if providerVersion == "" {
		c.Meta.Ui.Output(fmt.Sprintf("No version specified, will use latest version constraint for %s", provider))
		providerVersion = ">= 0.1.0"
	}

	// Confirm with the user
	if !autoApprove {
		c.Meta.Ui.Output(fmt.Sprintf("Will add provider %s version %s to your configuration", provider, providerVersion))
		v, err := c.Meta.Ui.Ask("Do you want to proceed? (y/n)")
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error asking for confirmation: %s", err))
			return 1
		}
		if strings.ToLower(v) != "y" && strings.ToLower(v) != "yes" {
			c.Meta.Ui.Output("Installation cancelled")
			return 0
		}
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error getting current working directory: %s", err))
		return 1
	}

	// Check if there's an existing required_providers block
	found, err := c.updateExistingProviderBlock(cwd, provider.Namespace, provider.Type, providerVersion)
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error updating existing provider block: %s", err))
		return 1
	}

	// If no existing block was found, create a new providers.tf file
	if !found {
		err = c.createProvidersFile(cwd, provider.Namespace, provider.Type, providerVersion)
		if err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error creating providers.tf file: %s", err))
			return 1
		}
	}

	c.Meta.Ui.Output(fmt.Sprintf("Successfully added provider %s version %s to your configuration", provider, providerVersion))
	return 0
}

// updateExistingProviderBlock searches for an existing required_providers block and updates it
// Returns true if a block was found and updated, false otherwise
func (c *RegistryProviderInstallCommand) updateExistingProviderBlock(cwd, namespace, name, version string) (bool, error) {
	// Look for .tf files in the current directory
	files, err := filepath.Glob(filepath.Join(cwd, "*.tf"))
	if err != nil {
		return false, fmt.Errorf("error searching for .tf files: %s", err)
	}

	// Regular expressions for finding required_providers blocks
	requiredProvidersRegex := regexp.MustCompile(`(?m)^\s*required_providers\s*{`)
	
	for _, file := range files {
		// Read the file
		content, err := os.ReadFile(file)
		if err != nil {
			return false, fmt.Errorf("error reading file %s: %s", file, err)
		}

		contentStr := string(content)
		
		// Find all required_providers blocks
		requiredProvidersMatches := requiredProvidersRegex.FindAllStringIndex(contentStr, -1)
		if len(requiredProvidersMatches) == 0 {
			continue
		}

		// For each required_providers block
		for _, match := range requiredProvidersMatches {
			blockStart := match[0]
			
			// Find the closing brace for this block
			blockContent := contentStr[blockStart:]
			braceLevel := 1
			var blockEnd int
			
			for i, char := range blockContent {
				if char == '{' {
					braceLevel++
				} else if char == '}' {
					braceLevel--
					if braceLevel == 0 {
						blockEnd = blockStart + i + 1
						break
					}
				}
			}
			
			if blockEnd == 0 {
				continue // Couldn't find the end of the block
			}
			
			// Extract the block content
			block := contentStr[blockStart:blockEnd]
			
			// Check if the provider is already in the block
			providerRegex := regexp.MustCompile(fmt.Sprintf(`(?m)^\s*%s\s*=\s*{`, name))
			if providerRegex.MatchString(block) {
				// Provider already exists, update its version
				newBlock := c.updateProviderVersion(block, name, namespace, version)
				
				// Replace the old block with the new one
				newContent := contentStr[:blockStart] + newBlock + contentStr[blockEnd:]
				
				// Write the updated content back to the file
				err = os.WriteFile(file, []byte(newContent), 0644)
				if err != nil {
					return false, fmt.Errorf("error writing to file %s: %s", file, err)
				}
				
				return true, nil
			} else {
				// Provider doesn't exist in this block, add it
				newBlock := c.addProviderToBlock(block, name, namespace, version)
				
				// Replace the old block with the new one
				newContent := contentStr[:blockStart] + newBlock + contentStr[blockEnd:]
				
				// Write the updated content back to the file
				err = os.WriteFile(file, []byte(newContent), 0644)
				if err != nil {
					return false, fmt.Errorf("error writing to file %s: %s", file, err)
				}
				
				return true, nil
			}
		}
	}
	
	// No required_providers block found
	return false, nil
}

// updateProviderVersion updates the version of an existing provider in a required_providers block
func (c *RegistryProviderInstallCommand) updateProviderVersion(block, name, namespace, version string) string {
	// Find the provider block
	providerRegex := regexp.MustCompile(fmt.Sprintf(`(?m)(\s*%s\s*=\s*{[^}]*})`, name))
	providerMatch := providerRegex.FindStringSubmatch(block)
	
	if len(providerMatch) > 0 {
		providerBlock := providerMatch[1]
		
		// Create a new provider block with proper indentation
		newProviderBlock := fmt.Sprintf("    %s = {\n", name)
		newProviderBlock += fmt.Sprintf("      source  = \"%s/%s\"\n", namespace, name)
		newProviderBlock += fmt.Sprintf("      version = \"%s\"\n", version)
		newProviderBlock += "    }"
		
		// Replace the old provider block with the new one
		return strings.Replace(block, providerBlock, newProviderBlock, 1)
	}
	
	return block
}

// addProviderToBlock adds a new provider to an existing required_providers block
func (c *RegistryProviderInstallCommand) addProviderToBlock(block, name, namespace, version string) string {
	// Parse the block into lines for easier manipulation
	lines := strings.Split(block, "\n")
	
	// Find the required_providers block
	requiredProvidersRegex := regexp.MustCompile(`^\s*required_providers\s*{`)
	requiredProvidersClosingRegex := regexp.MustCompile(`^\s*}`)
	
	var requiredProvidersStart, requiredProvidersEnd int
	inRequiredProviders := false
	braceLevel := 0
	
	for i, line := range lines {
		if requiredProvidersRegex.MatchString(line) {
			requiredProvidersStart = i
			inRequiredProviders = true
			braceLevel = 1
			continue
		}
		
		if inRequiredProviders {
			if strings.Contains(line, "{") {
				braceLevel++
			}
			if strings.Contains(line, "}") {
				braceLevel--
			}
			
			if braceLevel == 0 && requiredProvidersClosingRegex.MatchString(line) {
				requiredProvidersEnd = i
				break
			}
		}
	}
	
	if requiredProvidersStart == 0 && requiredProvidersEnd == 0 {
		return block // No required_providers block found
	}
	
	// Create the provider entry
	providerLines := []string{
		fmt.Sprintf("    %s = {", name),
		fmt.Sprintf("      source  = \"%s/%s\"", namespace, name),
		fmt.Sprintf("      version = \"%s\"", version),
		"    }",
	}
	
	// Insert the provider entry before the closing brace
	result := append(
		lines[:requiredProvidersEnd],
		append(
			providerLines,
			lines[requiredProvidersEnd:]...,
		)...,
	)
	
	return strings.Join(result, "\n")
}

// createProvidersFile creates a new providers.tf file with the required_providers block
func (c *RegistryProviderInstallCommand) createProvidersFile(cwd, namespace, name, version string) error {
	// Check if providers.tf already exists
	providersPath := filepath.Join(cwd, "providers.tf")
	if _, err := os.Stat(providersPath); err == nil {
		// File exists, read it
		content, err := os.ReadFile(providersPath)
		if err != nil {
			return fmt.Errorf("error reading providers.tf: %s", err)
		}
		
		contentStr := string(content)
		
		// Check if terraform block exists
		terraformRegex := regexp.MustCompile(`(?m)^\s*terraform\s*{`)
		terraformMatch := terraformRegex.FindStringIndex(contentStr)
		
		if len(terraformMatch) > 0 {
			// Terraform block exists, check if required_providers block exists
			requiredProvidersRegex := regexp.MustCompile(`(?m)^\s*required_providers\s*{`)
			requiredProvidersMatch := requiredProvidersRegex.FindStringIndex(contentStr)
			
			if len(requiredProvidersMatch) > 0 {
				// required_providers block exists, update it
				found, err := c.updateExistingProviderBlock(cwd, namespace, name, version)
				if err != nil {
					return fmt.Errorf("error updating existing provider block: %s", err)
				}
				
				if !found {
					return fmt.Errorf("required_providers block found but couldn't update it")
				}
				
				return nil
			} else {
				// Terraform block exists but no required_providers block
				// Find the closing brace of the terraform block
				blockStart := terraformMatch[0]
				blockContent := contentStr[blockStart:]
				braceLevel := 0
				var blockEnd int
				
				for i, char := range blockContent {
					if char == '{' {
						braceLevel++
					} else if char == '}' {
						braceLevel--
						if braceLevel == 0 {
							blockEnd = blockStart + i
							break
						}
					}
				}
				
				if blockEnd == 0 {
					return fmt.Errorf("couldn't find the end of the terraform block")
				}
				
				// Insert the required_providers block before the closing brace
				requiredProvidersBlock := fmt.Sprintf("\n  required_providers {\n")
				requiredProvidersBlock += fmt.Sprintf("    %s = {\n", name)
				requiredProvidersBlock += fmt.Sprintf("      source  = \"%s/%s\"\n", namespace, name)
				requiredProvidersBlock += fmt.Sprintf("      version = \"%s\"\n", version)
				requiredProvidersBlock += "    }\n"
				requiredProvidersBlock += "  }\n"
				
				newContent := contentStr[:blockEnd] + requiredProvidersBlock + contentStr[blockEnd:]
				
				// Write the updated content back to the file
				err = os.WriteFile(providersPath, []byte(newContent), 0644)
				if err != nil {
					return fmt.Errorf("error writing to providers.tf: %s", err)
				}
				
				return nil
			}
		} else {
			// No terraform block, append a new one
			newContent := contentStr
			if len(contentStr) > 0 && !strings.HasSuffix(contentStr, "\n") {
				newContent += "\n"
			}
			
			newContent += "\nterraform {\n"
			newContent += "  required_providers {\n"
			newContent += fmt.Sprintf("    %s = {\n", name)
			newContent += fmt.Sprintf("      source  = \"%s/%s\"\n", namespace, name)
			newContent += fmt.Sprintf("      version = \"%s\"\n", version)
			newContent += "    }\n"
			newContent += "  }\n"
			newContent += "}\n"
			
			// Write the updated content back to the file
			err = os.WriteFile(providersPath, []byte(newContent), 0644)
			if err != nil {
				return fmt.Errorf("error writing to providers.tf: %s", err)
			}
			
			return nil
		}
	} else {
		// File doesn't exist, create it
		content := "terraform {\n"
		content += "  required_providers {\n"
		content += fmt.Sprintf("    %s = {\n", name)
		content += fmt.Sprintf("      source  = \"%s/%s\"\n", namespace, name)
		content += fmt.Sprintf("      version = \"%s\"\n", version)
		content += "    }\n"
		content += "  }\n"
		content += "}\n"
		
		// Write the content to the file
		err = os.WriteFile(providersPath, []byte(content), 0644)
		if err != nil {
			return fmt.Errorf("error creating providers.tf: %s", err)
		}
		
		return nil
	}
}

func (c *RegistryProviderInstallCommand) Help() string {
	helpText := `
Usage: tofu registry provider install [options] PROVIDER

  Installs a provider from a registry and adds it to your configuration.

Options:

  -version=<version>      Specific version constraint to install. If not specified, 
                          the latest version constraint (>= 0.1.0) will be used.
  -auto-approve           Skip interactive approval of provider installation.
`
	return strings.TrimSpace(helpText)
}

func (c *RegistryProviderInstallCommand) Synopsis() string {
	return "Install a provider from a registry"
}

func (c *RegistryProviderInstallCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RegistryProviderInstallCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{
		"-version":      complete.PredictAnything,
		"-auto-approve": complete.PredictNothing,
	}
}
