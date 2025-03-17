// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"path/filepath"

	"github.com/opentofu/opentofu/internal/command"
	"github.com/opentofu/opentofu/internal/dotenv"
	"github.com/opentofu/opentofu/internal/templates"
)

// SetupCommand is the command that handles initial setup
type SetupCommand struct {
	Meta command.Meta
}

// Help returns help text for the Setup command
func (c *SetupCommand) Help() string {
	helpText := `
Usage: tofu setup [options]

  This command sets up the necessary tools and configurations for OpenTofu development.

Options:
  -db-only         Only setup the database
  -tools-only      Only install development tools
  -skip-tools      Skip installing development tools
  -skip-db         Skip database setup
  -type=TYPE       Database type (sqlite or postgres). Default: sqlite
  -path=PATH       Path to SQLite database file. Default: ~/.opentofu/tofu.db
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short description of the Setup command
func (c *SetupCommand) Synopsis() string {
	return "Setup OpenTofu development environment"
}

// Run runs the Setup command
func (c *SetupCommand) Run(args []string) int {
	// Load environment variables from .env file
	vars, err := dotenv.Load("")
	if err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Error loading .env file: %s", err))
	} else if len(vars) > 0 {
		c.Meta.Ui.Info("Loaded environment variables from .env file")
	}

	var dbOnly bool
	var toolsOnly bool
	var skipTools bool
	var skipDB bool
	var dbType string
	var dbPath string

	// Parse command line arguments
	for i := 0; i < len(args); i++ {
		if args[i] == "-db-only" {
			dbOnly = true
		} else if args[i] == "-tools-only" {
			toolsOnly = true
		} else if args[i] == "-skip-tools" {
			skipTools = true
		} else if args[i] == "-skip-db" {
			skipDB = true
		} else if strings.HasPrefix(args[i], "-type=") {
			dbType = args[i][6:]
		} else if strings.HasPrefix(args[i], "-path=") {
			dbPath = args[i][6:]
		}
	}

	// Set default values
	if dbType == "" {
		dbType = dotenv.GetWithDefault("TOFU_DB_TYPE", "sqlite")
	}

	// Display welcome message
	c.Meta.Ui.Output("\n=== OpenTofu Development Environment Setup ===\n")
	c.Meta.Ui.Output("This setup will prepare your environment for OpenTofu development.")

	// Check if we're running the minime fork
	c.Meta.Ui.Output("\nDetected project: minime (OpenTofu fork)")

	// Install development tools if not skipped
	if !skipTools && !dbOnly {
		if err := c.installDevelopmentTools(); err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error installing development tools: %s", err))
			return 1
		}
	}

	// Setup database if not skipped
	if !skipDB && !toolsOnly {
		dbCmd := &DBSetupCommand{
			Meta: c.Meta,
		}
		dbArgs := []string{}
		if dbType != "" {
			dbArgs = append(dbArgs, fmt.Sprintf("-type=%s", dbType))
		}
		if dbPath != "" {
			dbArgs = append(dbArgs, fmt.Sprintf("-path=%s", dbPath))
		}
		
		exitCode := dbCmd.Run(dbArgs)
		if exitCode != 0 {
			c.Meta.Ui.Error("Database setup failed")
			return exitCode
		}
	}

	// Create default .env file if it doesn't exist
	if !templates.FileExists(".env") {
		c.Meta.Ui.Info("Creating default .env file...")
		envContent := `# OpenTofu Development Environment Configuration
TOFU_DB_TYPE=sqlite
# Uncomment and modify the following lines for PostgreSQL
# TOFU_REGISTRY_DB_HOST=localhost
# TOFU_REGISTRY_DB_PORT=5432
# TOFU_REGISTRY_DB_USER=postgres
# TOFU_REGISTRY_DB_PASSWORD=postgres
# TOFU_REGISTRY_DB_NAME=opentofu
# TOFU_REGISTRY_DB_SSLMODE=disable
`
		if err := os.WriteFile(".env", []byte(envContent), 0644); err != nil {
			c.Meta.Ui.Error(fmt.Sprintf("Error creating .env file: %s", err))
			return 1
		}
	}

	// Display completion message
	c.Meta.Ui.Output("\n=== Setup Complete ===")
	c.Meta.Ui.Output("\nYour OpenTofu development environment is now ready.")
	c.Meta.Ui.Output("You can use the following commands:")
	c.Meta.Ui.Output("  tofu db setup      - Setup the database")
	c.Meta.Ui.Output("  tofu db configure  - Configure database connection")
	c.Meta.Ui.Output("  tofu db test       - Test database functionality")
	c.Meta.Ui.Output("  tofu db migrate    - Migrate database schema")
	
	return 0
}

// installDevelopmentTools installs necessary development tools
func (c *SetupCommand) installDevelopmentTools() error {
	c.Meta.Ui.Output("\n=== Installing Development Tools ===")

	// Check Go installation
	c.Meta.Ui.Info("Checking Go installation...")
	goCmd := exec.Command("go", "version")
	goOutput, err := goCmd.CombinedOutput()
	if err != nil {
		c.Meta.Ui.Error("Go is not installed or not in PATH")
		c.Meta.Ui.Error("Please install Go from https://golang.org/dl/")
		return fmt.Errorf("go not installed")
	}
	c.Meta.Ui.Info(fmt.Sprintf("Found %s", strings.TrimSpace(string(goOutput))))

	// Install goenv
	c.Meta.Ui.Info("Checking for goenv installation...")
	_, err = exec.LookPath("goenv")
	if err != nil {
		c.Meta.Ui.Info("goenv not found, installing...")
		
		// Determine home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			c.Meta.Ui.Warn(fmt.Sprintf("Failed to determine home directory: %v", err))
			c.Meta.Ui.Warn("Skipping goenv installation")
		} else {
			goenvDir := filepath.Join(homeDir, ".goenv")
			
			// Check if curl or wget is available
			_, curlErr := exec.LookPath("curl")
			_, wgetErr := exec.LookPath("wget")
			
			if curlErr == nil || wgetErr == nil {
				// Create .goenv directory if it doesn't exist
				if err := os.MkdirAll(goenvDir, 0755); err != nil {
					c.Meta.Ui.Warn(fmt.Sprintf("Failed to create .goenv directory: %v", err))
				} else {
					// Clone goenv repository
					var cloneCmd *exec.Cmd
					
					if curlErr == nil {
						c.Meta.Ui.Info("Using curl to install goenv...")
						cloneCmd = exec.Command("sh", "-c", "curl -fsSL https://github.com/syndbg/goenv/archive/master.tar.gz | tar xz --strip=1 -C \""+goenvDir+"\"")
					} else {
						c.Meta.Ui.Info("Using wget to install goenv...")
						cloneCmd = exec.Command("sh", "-c", "wget -qO- https://github.com/syndbg/goenv/archive/master.tar.gz | tar xz --strip=1 -C \""+goenvDir+"\"")
					}
					
					if output, err := cloneCmd.CombinedOutput(); err != nil {
						c.Meta.Ui.Warn(fmt.Sprintf("Failed to install goenv: %s", output))
					} else {
						c.Meta.Ui.Info("goenv installed successfully")
						
						// Add goenv to PATH and initialize
						shellType := "bash"
						shellRcPath := filepath.Join(homeDir, ".bashrc")
						
						// Check for zsh
						if _, err := os.Stat(filepath.Join(homeDir, ".zshrc")); err == nil {
							shellType = "zsh"
							shellRcPath = filepath.Join(homeDir, ".zshrc")
						}
						
						// Add goenv to shell configuration
						c.Meta.Ui.Info(fmt.Sprintf("Adding goenv to %s configuration...", shellType))
						
						goenvConfig := `
# goenv setup
export GOENV_ROOT="$HOME/.goenv"
export PATH="$GOENV_ROOT/bin:$PATH"
eval "$(goenv init -)"
export PATH="$GOROOT/bin:$PATH"
export PATH="$PATH:$GOPATH/bin"
`
						
						// Append to shell configuration file
						f, err := os.OpenFile(shellRcPath, os.O_APPEND|os.O_WRONLY, 0644)
						if err != nil {
							c.Meta.Ui.Warn(fmt.Sprintf("Failed to open %s: %v", shellRcPath, err))
						} else {
							defer f.Close()
							if _, err := f.WriteString(goenvConfig); err != nil {
								c.Meta.Ui.Warn(fmt.Sprintf("Failed to update %s: %v", shellRcPath, err))
							} else {
								c.Meta.Ui.Info(fmt.Sprintf("Updated %s with goenv configuration", shellRcPath))
							}
						}
						
						// Install Go 1.24 using goenv
						c.Meta.Ui.Info("Installing Go 1.24 using goenv...")
						
						// Set PATH to include goenv
						newPath := filepath.Join(goenvDir, "bin") + string(os.PathListSeparator) + os.Getenv("PATH")
						
						// Create a command with the updated PATH
						installGoCmd := exec.Command("sh", "-c", "PATH=\""+newPath+"\" goenv install 1.24.0")
						
						if output, err := installGoCmd.CombinedOutput(); err != nil {
							c.Meta.Ui.Warn(fmt.Sprintf("Failed to install Go 1.24.0: %s", output))
							c.Meta.Ui.Warn("You can install it manually with: goenv install 1.24.0")
						} else {
							c.Meta.Ui.Info("Go 1.24.0 installed successfully")
							
							// Set global Go version
							globalCmd := exec.Command("sh", "-c", "PATH=\""+newPath+"\" goenv global 1.24.0")
							if output, err := globalCmd.CombinedOutput(); err != nil {
								c.Meta.Ui.Warn(fmt.Sprintf("Failed to set global Go version: %s", output))
							} else {
								c.Meta.Ui.Info("Global Go version set to 1.24.0")
							}
						}
						
						c.Meta.Ui.Info("To use goenv, restart your shell or run: source " + shellRcPath)
					}
				}
			} else {
				c.Meta.Ui.Warn("Neither curl nor wget found. Please install one of them to install goenv")
				c.Meta.Ui.Warn("You can install goenv manually from: https://github.com/syndbg/goenv")
			}
		}
	} else {
		c.Meta.Ui.Info("goenv is already installed")
		
		// Check if Go 1.24 is installed with goenv
		checkVersionCmd := exec.Command("sh", "-c", "goenv versions | grep 1.24")
		if output, err := checkVersionCmd.CombinedOutput(); err != nil {
			// Go 1.24 not installed, install it
			c.Meta.Ui.Info("Installing Go 1.24 using goenv...")
			installCmd := exec.Command("goenv", "install", "1.24.0")
			if output, err := installCmd.CombinedOutput(); err != nil {
				c.Meta.Ui.Warn(fmt.Sprintf("Failed to install Go 1.24.0: %s", output))
				c.Meta.Ui.Warn("You can install it manually with: goenv install 1.24.0")
			} else {
				c.Meta.Ui.Info("Go 1.24.0 installed successfully")
			}
		} else {
			c.Meta.Ui.Info(fmt.Sprintf("Go 1.24 is already installed: %s", strings.TrimSpace(string(output))))
		}
	}

	// Install golangci-lint
	c.Meta.Ui.Info("Installing golangci-lint...")
	lintCmd := exec.Command("go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest")
	if output, err := lintCmd.CombinedOutput(); err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Failed to install golangci-lint: %s", output))
		return err
	}
	c.Meta.Ui.Info("golangci-lint installed successfully")

	// Install goimports
	c.Meta.Ui.Info("Installing goimports...")
	importsCmd := exec.Command("go", "install", "golang.org/x/tools/cmd/goimports@latest")
	if output, err := importsCmd.CombinedOutput(); err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Failed to install goimports: %s", output))
		return err
	}
	c.Meta.Ui.Info("goimports installed successfully")

	// Install mockgen
	c.Meta.Ui.Info("Installing mockgen...")
	mockgenCmd := exec.Command("go", "install", "github.com/golang/mock/mockgen@latest")
	if output, err := mockgenCmd.CombinedOutput(); err != nil {
		c.Meta.Ui.Error(fmt.Sprintf("Failed to install mockgen: %s", output))
		return err
	}
	c.Meta.Ui.Info("mockgen installed successfully")

	// Check if we need to install SQLite development libraries based on OS
	switch runtime.GOOS {
	case "linux":
		c.Meta.Ui.Info("Checking for SQLite development libraries on Linux...")
		
		// Try to detect package manager
		var installCmd *exec.Cmd
		
		// Check for apt (Debian/Ubuntu)
		_, err := exec.LookPath("apt-get")
		if err == nil {
			c.Meta.Ui.Info("Detected apt package manager")
			installCmd = exec.Command("sudo", "apt-get", "install", "-y", "libsqlite3-dev")
		}
		
		// Check for yum/dnf (RHEL/Fedora)
		if installCmd == nil {
			_, err := exec.LookPath("dnf")
			if err == nil {
				c.Meta.Ui.Info("Detected dnf package manager")
				installCmd = exec.Command("sudo", "dnf", "install", "-y", "sqlite-devel")
			} else {
				_, err := exec.LookPath("yum")
				if err == nil {
					c.Meta.Ui.Info("Detected yum package manager")
					installCmd = exec.Command("sudo", "yum", "install", "-y", "sqlite-devel")
				}
			}
		}
		
		// Check for pacman (Arch)
		if installCmd == nil {
			_, err := exec.LookPath("pacman")
			if err == nil {
				c.Meta.Ui.Info("Detected pacman package manager")
				installCmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "sqlite")
			}
		}
		
		// Install SQLite development libraries if package manager was detected
		if installCmd != nil {
			c.Meta.Ui.Info("Installing SQLite development libraries...")
			if output, err := installCmd.CombinedOutput(); err != nil {
				c.Meta.Ui.Warn(fmt.Sprintf("Failed to install SQLite development libraries: %s", output))
				c.Meta.Ui.Warn("You may need to install them manually")
			} else {
				c.Meta.Ui.Info("SQLite development libraries installed successfully")
			}
		} else {
			c.Meta.Ui.Warn("Could not detect package manager, you may need to install SQLite development libraries manually")
		}
	
	case "darwin":
		c.Meta.Ui.Info("Checking for SQLite development libraries on macOS...")
		
		// Check for Homebrew
		_, err := exec.LookPath("brew")
		if err == nil {
			c.Meta.Ui.Info("Detected Homebrew package manager")
			
			// Check if SQLite is already installed
			brewListCmd := exec.Command("brew", "list", "--formula")
			brewOutput, err := brewListCmd.CombinedOutput()
			if err != nil {
				c.Meta.Ui.Warn(fmt.Sprintf("Failed to check installed Homebrew packages: %s", brewOutput))
			} else {
				installedPackages := string(brewOutput)
				if !strings.Contains(installedPackages, "sqlite") {
					// Install SQLite
					c.Meta.Ui.Info("Installing SQLite via Homebrew...")
					installCmd := exec.Command("brew", "install", "sqlite")
					if output, err := installCmd.CombinedOutput(); err != nil {
						c.Meta.Ui.Warn(fmt.Sprintf("Failed to install SQLite: %s", output))
						c.Meta.Ui.Warn("You may need to install it manually")
					} else {
						c.Meta.Ui.Info("SQLite installed successfully")
					}
				} else {
					c.Meta.Ui.Info("SQLite is already installed via Homebrew")
				}
			}
			
			// Check for PostgreSQL client
			if !strings.Contains(string(brewOutput), "postgresql") && !strings.Contains(string(brewOutput), "libpq") {
				c.Meta.Ui.Info("Installing PostgreSQL client libraries via Homebrew...")
				installCmd := exec.Command("brew", "install", "libpq")
				if output, err := installCmd.CombinedOutput(); err != nil {
					c.Meta.Ui.Warn(fmt.Sprintf("Failed to install PostgreSQL client libraries: %s", output))
					c.Meta.Ui.Warn("You may need to install them manually")
				} else {
					c.Meta.Ui.Info("PostgreSQL client libraries installed successfully")
					
					// Add libpq to PATH
					c.Meta.Ui.Info("To use the PostgreSQL client libraries, you may need to add them to your PATH:")
					c.Meta.Ui.Info("  echo 'export PATH=\"/usr/local/opt/libpq/bin:$PATH\"' >> ~/.zshrc")
				}
			} else {
				c.Meta.Ui.Info("PostgreSQL client libraries are already installed")
			}
		} else {
			c.Meta.Ui.Warn("Homebrew not found. Please install Homebrew from https://brew.sh/")
			c.Meta.Ui.Warn("Then install SQLite and PostgreSQL client libraries:")
			c.Meta.Ui.Warn("  brew install sqlite libpq")
		}
	
	case "windows":
		c.Meta.Ui.Info("Checking for package managers on Windows...")
		
		// Check for Chocolatey
		_, err := exec.LookPath("choco")
		if err == nil {
			c.Meta.Ui.Info("Detected Chocolatey package manager")
			
			// Install SQLite
			c.Meta.Ui.Info("Installing SQLite via Chocolatey...")
			installCmd := exec.Command("choco", "install", "sqlite", "-y")
			if output, err := installCmd.CombinedOutput(); err != nil {
				c.Meta.Ui.Warn(fmt.Sprintf("Failed to install SQLite: %s", output))
				c.Meta.Ui.Warn("You may need to install it manually")
			} else {
				c.Meta.Ui.Info("SQLite installed successfully")
			}
			
			// Install PostgreSQL client
			c.Meta.Ui.Info("Installing PostgreSQL client via Chocolatey...")
			installCmd = exec.Command("choco", "install", "postgresql", "-y")
			if output, err := installCmd.CombinedOutput(); err != nil {
				c.Meta.Ui.Warn(fmt.Sprintf("Failed to install PostgreSQL client: %s", output))
				c.Meta.Ui.Warn("You may need to install it manually")
			} else {
				c.Meta.Ui.Info("PostgreSQL client installed successfully")
			}
		} else {
			// Check for Scoop
			_, err := exec.LookPath("scoop")
			if err == nil {
				c.Meta.Ui.Info("Detected Scoop package manager")
				
				// Install SQLite
				c.Meta.Ui.Info("Installing SQLite via Scoop...")
				installCmd := exec.Command("scoop", "install", "sqlite")
				if output, err := installCmd.CombinedOutput(); err != nil {
					c.Meta.Ui.Warn(fmt.Sprintf("Failed to install SQLite: %s", output))
					c.Meta.Ui.Warn("You may need to install it manually")
				} else {
					c.Meta.Ui.Info("SQLite installed successfully")
				}
				
				// Install PostgreSQL client
				c.Meta.Ui.Info("Installing PostgreSQL client via Scoop...")
				installCmd = exec.Command("scoop", "install", "postgresql")
				if output, err := installCmd.CombinedOutput(); err != nil {
					c.Meta.Ui.Warn(fmt.Sprintf("Failed to install PostgreSQL client: %s", output))
					c.Meta.Ui.Warn("You may need to install it manually")
				} else {
					c.Meta.Ui.Info("PostgreSQL client installed successfully")
				}
			} else {
				c.Meta.Ui.Warn("No package manager found. Please install Chocolatey (https://chocolatey.org/) or Scoop (https://scoop.sh/)")
				c.Meta.Ui.Warn("Then install SQLite and PostgreSQL client")
				c.Meta.Ui.Warn("For Chocolatey: choco install sqlite postgresql -y")
				c.Meta.Ui.Warn("For Scoop: scoop install sqlite postgresql")
			}
		}
	
	default:
		c.Meta.Ui.Warn(fmt.Sprintf("Unsupported operating system: %s", runtime.GOOS))
		c.Meta.Ui.Warn("Please install SQLite and PostgreSQL client libraries manually")
	}

	c.Meta.Ui.Info("Development tools installation complete")
	return nil
}
