# JSON Fact Collector Example

This example demonstrates how to use the JSON fact collector in the Spooky system.

## Overview

The JSON fact collector can parse JSON files and convert them into fact collections. It supports:

- Single JSON file processing
- Directory processing (multiple JSON files)
- Nested object and array handling
- Selective fact collection
- Individual fact retrieval

## Features

### 1. Single File Processing
Collects facts from a single JSON file, converting nested objects and arrays into flat fact keys.

### 2. Directory Processing
Processes all `.json` files in a directory, ignoring non-JSON files.

### 3. Nested Structure Support
- Objects: `config.port` → `8080`
- Arrays: `tags[0]` → `"web"`
- Deep nesting: `config.ssl.cert` → `"/etc/ssl/certs/server.crt"`

### 4. Selective Collection
Collect only specific facts by key names.

### 5. Individual Fact Retrieval
Get a single fact by its key.

## Usage

```bash
go run main.go
```

## Example Output

The example demonstrates:

1. **Single JSON File**: Collects 10 facts from `server1.json`
2. **Directory Processing**: Collects 13 facts from both JSON files in the directory
3. **Selective Collection**: Collects only `name` and `config.port` facts
4. **Single Fact**: Retrieves the `name` fact value

## Sample JSON Structure

The example uses JSON files with nested structures:

```json
{
  "name": "web-server-1",
  "config": {
    "port": 8080,
    "enabled": true,
    "ssl": {
      "enabled": true,
      "cert": "/etc/ssl/certs/server.crt"
    }
  },
  "tags": ["web", "production", "load-balancer"],
  "metadata": {
    "environment": "production",
    "region": "us-west-2"
  }
}
```

## Generated Facts

The JSON collector converts the above structure into these facts:

- `name`: `"web-server-1"`
- `config.port`: `8080`
- `config.enabled`: `true`
- `config.ssl.enabled`: `true`
- `config.ssl.cert`: `"/etc/ssl/certs/server.crt"`
- `tags[0]`: `"web"`
- `tags[1]`: `"production"`
- `tags[2]`: `"load-balancer"`
- `metadata.environment`: `"production"`
- `metadata.region`: `"us-west-2"`

## Implementation Details

The JSON collector:

- Uses the `BaseCollector` pattern for consistency
- Implements the `FactCollector` interface
- Supports file size validation (10MB limit)
- Provides comprehensive error handling
- Includes metadata about source files and data types
- Sets appropriate TTL and source information for facts
