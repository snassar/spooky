#!/bin/bash
# Test script for spooky validation testing
# This script is harmless and returns predictable output for testing

echo "Test script executed successfully"
echo "Script name: $0"
echo "Arguments: $@"
echo "Current directory: $(pwd)"
echo "Current user: $(whoami)"
echo "Script timestamp: $(date)"

# Return a predictable exit code for testing
exit 0 