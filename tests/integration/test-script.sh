#!/bin/bash

# Test script for spooky integration tests
# This script will be executed on the SSH container

set -e

echo "=== Spooky Integration Test Script ==="
echo "Script started at: $(date)"
echo "Running as user: $(whoami)"
echo "Current directory: $(pwd)"
echo "Hostname: $(hostname)"
echo ""

# Test basic commands
echo "Testing basic commands..."
echo "ls -la:"
ls -la

echo ""
echo "Testing environment..."
echo "PATH: $PATH"
echo "HOME: $HOME"
echo "USER: $USER"

echo ""
echo "Testing file operations..."
echo "Creating test file..."
echo "Hello from spooky test script!" > /tmp/spooky-test.txt
echo "File created successfully"
echo "File contents:"
cat /tmp/spooky-test.txt

echo ""
echo "Testing network connectivity..."
if command -v curl >/dev/null 2>&1; then
    echo "curl is available"
    echo "Testing localhost connectivity..."
    curl -s --connect-timeout 5 http://localhost || echo "No local web server running"
else
    echo "curl not available, skipping network test"
fi

echo ""
echo "Testing sudo access..."
if sudo -n true 2>/dev/null; then
    echo "Sudo access confirmed"
    echo "Testing sudo command..."
    sudo whoami
else
    echo "No sudo access available"
fi

echo ""
echo "=== Test script completed successfully ==="
echo "Script finished at: $(date)" 