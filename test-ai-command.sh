#!/bin/bash

# OpenTofu AI Command Testing Script
# This script tests the AI command functionality with different configurations
# and captures screenshots of the results.

set -e

# Create a directory for test outputs and screenshots
mkdir -p test-results/screenshots

# Function to take a screenshot of terminal output
take_screenshot() {
  local test_name=$1
  echo "Taking screenshot for: $test_name"
  # Use the 'script' command to capture terminal output
  script -q -c "$2" test-results/$test_name.log
  # Copy the log to screenshots directory
  cp test-results/$test_name.log test-results/screenshots/
  echo "Screenshot saved to test-results/screenshots/$test_name.log"
  echo "------------------------------------------------------------"
}

# Function to run a test and capture its output
run_test() {
  local test_name=$1
  local command=$2
  
  echo "Running test: $test_name"
  echo "Command: $command"
  echo "------------------------------------------------------------"
  
  # Create a directory for this test
  mkdir -p test-results/$test_name
  
  # Run the command and capture its output
  take_screenshot "$test_name" "$command"
  
  echo "Test completed: $test_name"
  echo "------------------------------------------------------------"
}

# Check if we're running in GitHub Actions
if [ -n "$GITHUB_ACTIONS" ]; then
  echo "Running in GitHub Actions, using pre-built binary..."
  # In GitHub Actions, the binary should be pre-built by the workflow
  TOFU_BIN="./tofu"
else
  # For local testing, download the latest release
  echo "Running locally, downloading the latest release..."
  
  # Create a temporary directory for the download
  TEMP_DIR=$(mktemp -d)
  
  # Download the latest release for the current platform
  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  ARCH=$(uname -m)
  if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
  elif [ "$ARCH" = "aarch64" ]; then
    ARCH="arm64"
  fi
  
  # Download the latest release
  RELEASE_URL="https://github.com/opentofu/opentofu/releases/latest/download/tofu_${OS}_${ARCH}.zip"
  echo "Downloading from: $RELEASE_URL"
  curl -L "$RELEASE_URL" -o "$TEMP_DIR/tofu.zip"
  
  # Extract the binary
  unzip -q "$TEMP_DIR/tofu.zip" -d "$TEMP_DIR"
  cp "$TEMP_DIR/tofu" ./tofu
  chmod +x ./tofu
  
  # Clean up
  rm -rf "$TEMP_DIR"
  
  TOFU_BIN="./tofu"
fi

if [ ! -x "$TOFU_BIN" ]; then
  echo "Error: OpenTofu binary not found. Please check the build process."
  exit 1
fi

echo "Using OpenTofu binary: $TOFU_BIN"

# Create test directories for different scenarios
mkdir -p test-results/01_help_command
mkdir -p test-results/02_version_info
mkdir -p test-results/03_anthropic_basic
mkdir -p test-results/04_anthropic_with_registry
mkdir -p test-results/05_anthropic_custom_prompt
mkdir -p test-results/06_ollama_basic
mkdir -p test-results/07_ollama_with_registry
mkdir -p test-results/08_tf_format
mkdir -p test-results/09_json_format

# Test 1: Basic help command
run_test "01_help_command" "$TOFU_BIN ai --help"

# Test 2: Version information
run_test "02_version_info" "$TOFU_BIN version"

# Test 3: Generate with Anthropic (if API key is available)
if [ -n "$ANTHROPIC_API_KEY" ]; then
  # Basic test
  run_test "03_anthropic_basic" "$TOFU_BIN ai generate 'Create a simple AWS EC2 instance' --provider anthropic --model claude-3-sonnet-20240229 --output test-results/03_anthropic_basic"
  
  # Test with registry integration
  run_test "04_anthropic_with_registry" "$TOFU_BIN ai generate 'Create a simple AWS EC2 instance with a security group' --provider anthropic --model claude-3-sonnet-20240229 --use-registry --registry-db 'host=vultr-prod-860996d7-f3c4-4df8-b691-06ecc64db1c7-vultr-prod-c0b9.vultrdb.com port=16751 user=opentofu_user dbname=opentofu sslmode=require' --output test-results/04_anthropic_with_registry"
  
  # Test with system prompt override
  run_test "05_anthropic_custom_prompt" "$TOFU_BIN ai generate 'Create a simple AWS EC2 instance' --provider anthropic --model claude-3-sonnet-20240229 --system-prompt 'You are an expert in AWS infrastructure. Generate OpenTofu code for the user request.' --output test-results/05_anthropic_custom_prompt"
else
  echo "Skipping Anthropic tests - ANTHROPIC_API_KEY not set"
fi

# Test 4: Generate with Ollama (if available)
if command -v ollama &> /dev/null; then
  # Basic test
  run_test "06_ollama_basic" "$TOFU_BIN ai generate 'Create a simple AWS EC2 instance' --provider ollama --model llama3 --output test-results/06_ollama_basic"
  
  # Test with registry integration
  run_test "07_ollama_with_registry" "$TOFU_BIN ai generate 'Create a simple AWS EC2 instance with a security group' --provider ollama --model llama3 --use-registry --registry-db 'host=vultr-prod-860996d7-f3c4-4df8-b691-06ecc64db1c7-vultr-prod-c0b9.vultrdb.com port=16751 user=opentofu_user dbname=opentofu sslmode=require' --output test-results/07_ollama_with_registry"
else
  echo "Skipping Ollama tests - ollama command not found"
fi

# Test 5: Test with different output formats
if [ -n "$ANTHROPIC_API_KEY" ]; then
  # TF format
  run_test "08_tf_format" "$TOFU_BIN ai generate 'Create a simple AWS EC2 instance' --provider anthropic --model claude-3-sonnet-20240229 --output-format tf --output test-results/08_tf_format"
  
  # JSON format
  run_test "09_json_format" "$TOFU_BIN ai generate 'Create a simple AWS EC2 instance' --provider anthropic --model claude-3-sonnet-20240229 --output-format json --output test-results/09_json_format"
fi

echo "All tests completed. Results are in the test-results directory."
