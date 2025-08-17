# Query Language System Design

## Overview

This document outlines the design and implementation plan for a complex query language system in spooky that provides SQL-like querying capabilities for machines, facts, variables, and other spooky data sources. The system will support advanced filtering, joins, aggregations, and subqueries while maintaining integration with the existing HCL configuration system.

## Goals

1. **Provide powerful querying capabilities** for complex data filtering and analysis
2. **Maintain familiarity** through SQL-like syntax
3. **Integrate seamlessly** with existing HCL configuration system
4. **Support multiple data sources** (machines, facts, variables, actions)
5. **Enable query composition** and reusability
6. **Ensure performance** through optimization and caching
7. **Maintain type safety** through existing type system integration

## Architecture

### Core Components

```
Query Language System
├── Parser (SQL-like syntax → AST)
├── Executor (AST → Results)
├── Data Sources (Machines, Facts, Variables, Actions)
├── CLI Integration (query command)
├── HCL Integration (query configuration)
└── Caching & Optimization
```

### Data Flow

1. **Query Input**: SQL string or HCL query file
2. **Parsing**: Convert to Abstract Syntax Tree (AST)
3. **Validation**: Type checking and semantic validation
4. **Execution**: Execute against data sources
5. **Result Processing**: Format and return results
6. **Caching**: Store results for reuse

## Query Language Design

### Core Syntax

#### Basic Queries
```sql
-- Simple selection
SELECT * FROM machines WHERE tags.environment = 'production'

-- Field selection with filtering
SELECT hostname, host, user FROM machines 
WHERE groups CONTAINS 'webservers' AND resources.cpu_cores > 4

-- Complex conditions
SELECT hostname, tags.role, resources.memory_gb 
FROM machines 
WHERE (tags.environment = 'production' OR tags.environment = 'staging')
  AND resources.cpu_cores >= 4
  AND resources.memory_gb >= 8
```

#### Joins and Relationships
```sql
-- Join machines with facts
SELECT m.hostname, f.cpu_cores, f.memory_gb 
FROM machines m 
JOIN facts f ON m.hostname = f.machine_id 
WHERE m.tags.environment = 'production'

-- Multiple joins
SELECT m.hostname, f.cpu_cores, v.environment 
FROM machines m 
JOIN facts f ON m.hostname = f.machine_id 
JOIN variables v ON m.hostname = v.machine_id 
WHERE f.cpu_cores > 4
```

#### Aggregations
```sql
-- Group by with aggregations
SELECT tags.environment, COUNT(*) as machine_count, 
       AVG(resources.cpu_cores) as avg_cpu 
FROM machines 
GROUP BY tags.environment

-- Having clauses
SELECT hostname, COUNT(*) as fact_count 
FROM machines m 
JOIN facts f ON m.hostname = f.machine_id 
GROUP BY m.hostname 
HAVING fact_count > 10
```

#### Subqueries
```sql
-- IN subquery
SELECT * FROM machines 
WHERE hostname IN (
  SELECT machine_id FROM facts 
  WHERE cpu_cores > 4 AND memory_gb > 16
)

-- EXISTS subquery
SELECT hostname FROM machines m 
WHERE EXISTS (
  SELECT 1 FROM facts f 
  WHERE f.machine_id = m.hostname 
    AND f.cpu_cores > 8
)
```

### HCL Integration

#### Query Configuration
```hcl
# queries/production_servers.hcl
query "production_webservers" {
  sql = "SELECT * FROM machines WHERE tags.environment = 'production' AND groups CONTAINS 'webservers'"
  
  parameters = {
    environment = "production"
    min_cpu_cores = 4
  }
  
  options = {
    timeout = "30s"
    max_results = 100
    cache_results = true
    format = "json"
  }
}

query "high_cpu_machines" {
  sql = "SELECT hostname, resources.cpu_cores FROM machines WHERE resources.cpu_cores >= $min_cpu"
  
  parameters = {
    min_cpu = 8
  }
  
  options = {
    timeout = "15s"
    format = "table"
  }
}
```

#### Query Templates
```hcl
# queries/templates.hcl
template "environment_machines" {
  sql = "SELECT * FROM machines WHERE tags.environment = $env"
  parameters = {
    env = "production"
  }
}

template "role_machines" {
  sql = "SELECT * FROM machines WHERE tags.role = $role AND tags.environment = $env"
  parameters = {
    role = "web"
    env = "production"
  }
}
```

## Implementation Plan

### Phase 1: Core Infrastructure (Week 1-2)

#### 1.1 Query Parser
```go
// internal/query/parser/lexer.go
type Lexer struct {
    input string
    pos   int
    tokens []Token
}

type Token struct {
    Type    TokenType
    Value   string
    Line    int
    Column  int
}

// internal/query/parser/parser.go
type Parser struct {
    tokens []Token
    pos    int
}

type ASTNode interface {
    Execute(ctx context.Context, data interface{}) (interface{}, error)
    String() string
}

type SelectStatement struct {
    Fields    []Field
    From      string
    Where     Expression
    GroupBy   []string
    Having    Expression
    OrderBy   []OrderBy
    Limit     int
    Offset    int
}
```

#### 1.2 Expression System
```go
// internal/query/expressions/expression.go
type Expression interface {
    Evaluate(ctx context.Context, data interface{}) (interface{}, error)
    Type() reflect.Type
}

type BinaryExpression struct {
    Left     Expression
    Operator string
    Right    Expression
}

type FieldAccess struct {
    Object string
    Field  string
}

type FunctionCall struct {
    Name      string
    Arguments []Expression
}

type Literal struct {
    Value interface{}
}
```

#### 1.3 Data Source Interface
```go
// internal/query/sources/interface.go
type DataSource interface {
    GetData(ctx context.Context, filter interface{}) ([]map[string]interface{}, error)
    GetSchema() map[string]reflect.Type
    GetName() string
}
```

### Phase 2: Data Sources (Week 3-4)

#### 2.1 Machine Data Source
```go
// internal/query/sources/machines.go
type MachineDataSource struct {
    manager spookyinterfaces.MachinesIntegration
}

func (m *MachineDataSource) GetData(ctx context.Context, filter interface{}) ([]map[string]interface{}, error) {
    machines, err := m.manager.LoadMachines(ctx, projectPath)
    if err != nil {
        return nil, err
    }
    
    var data []map[string]interface{}
    for _, machine := range machines {
        data = append(data, machineToMap(machine))
    }
    
    return data, nil
}

func machineToMap(machine spookytypes.Machine) map[string]interface{} {
    return map[string]interface{}{
        "hostname": machine.Hostname,
        "host":     machine.Host,
        "port":     machine.Port,
        "user":     machine.User,
        "tags":     machine.Tags,
        "groups":   machine.Groups,
        "roles":    machine.Roles,
        "classes":  machine.Classes,
        "resources": map[string]interface{}{
            "cpu_cores":    machine.Resources.CPUCores,
            "memory_gb":    machine.Resources.MemoryGB,
            "disk_gb":      machine.Resources.DiskGB,
            "network_speed": machine.Resources.NetworkSpeed,
        },
        "metadata": map[string]interface{}{
            "environment": machine.MachineMetadata.Environment,
            "location":    machine.MachineMetadata.Location,
            "owner":       machine.MachineMetadata.Owner,
        },
    }
}
```

#### 2.2 Facts Data Source
```go
// internal/query/sources/facts.go
type FactsDataSource struct {
    manager spookyinterfaces.FactsIntegration
}

func (f *FactsDataSource) GetData(ctx context.Context, filter interface{}) ([]map[string]interface{}, error) {
    facts, err := f.manager.LoadFacts(ctx)
    if err != nil {
        return nil, err
    }
    
    var data []map[string]interface{}
    for _, fact := range facts {
        data = append(data, factToMap(fact))
    }
    
    return data, nil
}
```

#### 2.3 Variables Data Source
```go
// internal/query/sources/variables.go
type VariablesDataSource struct {
    manager spookyinterfaces.VariablesIntegration
}

func (v *VariablesDataSource) GetData(ctx context.Context, filter interface{}) ([]map[string]interface{}, error) {
    variables, err := v.manager.LoadVariables(ctx)
    if err != nil {
        return nil, err
    }
    
    var data []map[string]interface{}
    for _, variable := range variables {
        data = append(data, variableToMap(variable))
    }
    
    return data, nil
}
```

### Phase 3: Query Executor (Week 5-6)

#### 3.1 Core Executor
```go
// internal/query/executor.go
type QueryExecutor struct {
    sources map[string]DataSource
    cache   *QueryCache
    logger  spookytypeslogging.Logger
}

func (e *QueryExecutor) ExecuteQuery(ctx context.Context, query string) (*QueryResult, error) {
    // Parse query into AST
    ast, err := e.parser.Parse(query)
    if err != nil {
        return nil, fmt.Errorf("query parsing failed: %w", err)
    }
    
    // Execute AST
    result, err := ast.Execute(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("query execution failed: %w", err)
    }
    
    return &QueryResult{
        Data:    result,
        Query:   query,
        Runtime: time.Since(startTime),
    }, nil
}
```

#### 3.2 Join Engine
```go
// internal/query/joins/engine.go
type JoinEngine struct {
    sources map[string]DataSource
}

func (j *JoinEngine) ExecuteJoin(left, right string, condition Expression) ([]map[string]interface{}, error) {
    leftData, err := j.sources[left].GetData(ctx, nil)
    if err != nil {
        return nil, err
    }
    
    rightData, err := j.sources[right].GetData(ctx, nil)
    if err != nil {
        return nil, err
    }
    
    // Perform join based on condition
    return j.performJoin(leftData, rightData, condition)
}
```

#### 3.3 Aggregation Engine
```go
// internal/query/aggregations/engine.go
type AggregationEngine struct{}

func (a *AggregationEngine) ExecuteAggregation(data []map[string]interface{}, groupBy []string, aggregations []Aggregation) ([]map[string]interface{}, error) {
    // Group data
    groups := a.groupData(data, groupBy)
    
    // Apply aggregations
    var results []map[string]interface{}
    for _, group := range groups {
        result := a.applyAggregations(group, aggregations)
        results = append(results, result)
    }
    
    return results, nil
}
```

### Phase 4: CLI Integration (Week 7)

#### 4.1 Query Command
```go
// cmd/query.go
var queryCmd = &cobra.Command{
    Use:   "query [sql-query]",
    Short: "Execute complex queries against spooky data",
    Long: `Execute SQL-like queries against machines, facts, and other spooky data.

Supports complex filtering, joins, aggregations, and subqueries.

Examples:
  spooky query "SELECT * FROM machines WHERE tags.environment = 'production'"
  spooky query "SELECT hostname, cpu_cores FROM machines WHERE resources.cpu_cores > 4"
  spooky query "SELECT m.hostname, COUNT(f.*) as fact_count FROM machines m JOIN facts f ON m.hostname = f.machine_id GROUP BY m.hostname"
  
  # Using HCL query files
  spooky query --file queries/production_servers.hcl
  spooky query --file queries/high_cpu_servers.hcl --parameters cpu_min=8`,
    Args: cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        return handleQuery(cmd, args)
    },
}

func handleQuery(cmd *cobra.Command, args []string) error {
    ctx := context.Background()
    
    var query string
    if len(args) > 0 {
        query = args[0]
    } else {
        // Load from file
        queryFile, _ := cmd.Flags().GetString("file")
        if queryFile == "" {
            return fmt.Errorf("either provide a query string or use --file")
        }
        query = loadQueryFromFile(queryFile)
    }
    
    // Execute query
    executor := spookyquery.NewExecutor(machinesManager, factsManager, variablesManager)
    result, err := executor.ExecuteQuery(ctx, query)
    if err != nil {
        return fmt.Errorf("query execution failed: %w", err)
    }
    
    // Output results
    return outputQueryResults(cmd, result)
}
```

#### 4.2 Output Formats
```go
// internal/query/output/formatters.go
type ResultFormatter interface {
    Format(data interface{}) (string, error)
}

type JSONFormatter struct{}
type HCLFormatter struct{}
type TableFormatter struct{}
type CSVFormatter struct{}

func outputQueryResults(cmd *cobra.Command, result *QueryResult) error {
    format, _ := cmd.Flags().GetString("format")
    
    var formatter ResultFormatter
    switch format {
    case "json":
        formatter = &JSONFormatter{}
    case "hcl":
        formatter = &HCLFormatter{}
    case "table":
        formatter = &TableFormatter{}
    case "csv":
        formatter = &CSVFormatter{}
    default:
        formatter = &TableFormatter{}
    }
    
    output, err := formatter.Format(result.Data)
    if err != nil {
        return fmt.Errorf("failed to format results: %w", err)
    }
    
    fmt.Println(output)
    return nil
}
```

### Phase 5: Advanced Features (Week 8-10)

#### 5.1 Query Functions
```go
// internal/query/functions/functions.go
type FunctionRegistry struct {
    functions map[string]Function
}

type Function interface {
    Execute(args []interface{}) (interface{}, error)
    Name() string
    Arity() int
    ReturnType() reflect.Type
}

// Built-in functions
type CountFunction struct{}
type AvgFunction struct{}
type SumFunction struct{}
type MaxFunction struct{}
type MinFunction struct{}
type UpperFunction struct{}
type LowerFunction struct{}
type ConcatFunction struct{}
type LengthFunction struct{}
type NowFunction struct{}
type DateDiffFunction struct{}
```

#### 5.2 Query Caching
```go
// internal/query/cache/cache.go
type QueryCache struct {
    cache map[string]*CachedResult
    mu    sync.RWMutex
    ttl   time.Duration
}

type CachedResult struct {
    Data      interface{}
    Timestamp time.Time
    TTL       time.Duration
}

func (c *QueryCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if result, exists := c.cache[key]; exists {
        if time.Since(result.Timestamp) < result.TTL {
            return result.Data, true
        }
        delete(c.cache, key)
    }
    return nil, false
}

func (c *QueryCache) Set(key string, data interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.cache[key] = &CachedResult{
        Data:      data,
        Timestamp: time.Now(),
        TTL:       ttl,
    }
}
```

#### 5.3 Query Optimization
```go
// internal/query/optimizer/optimizer.go
type QueryOptimizer struct {
    rules []OptimizationRule
}

type OptimizationRule interface {
    Apply(ast ASTNode) (ASTNode, bool)
    Name() string
}

type PredicatePushdownRule struct{}
type JoinReorderingRule struct{}
type IndexSelectionRule struct{}
type SubqueryFlatteningRule struct{}
```

### Phase 6: Type System Integration (Week 11-12)

#### 6.1 Query Types
```go
// internal/types/query/query.go
type Query struct {
    SQL        string                 `hcl:"sql"`
    Parameters map[string]interface{} `hcl:"parameters,optional"`
    Options    *QueryOptions          `hcl:"options,block"`
}

type QueryOptions struct {
    Timeout     time.Duration `hcl:"timeout,optional"`
    MaxResults  int           `hcl:"max_results,optional"`
    CacheResults bool         `hcl:"cache_results,optional"`
    Format      string        `hcl:"format,optional"` // json, hcl, table, csv
}

type QueryResult struct {
    Data    interface{}   `json:"data"`
    Query   string        `json:"query"`
    Runtime time.Duration `json:"runtime"`
    Rows    int           `json:"rows"`
    Columns []string      `json:"columns"`
    Error   string        `json:"error,omitempty"`
}

type QueryTemplate struct {
    Name       string                 `hcl:"name,label"`
    SQL        string                 `hcl:"sql"`
    Parameters map[string]interface{} `hcl:"parameters,optional"`
    Options    *QueryOptions          `hcl:"options,block"`
}
```

#### 6.2 Schema Validation
```go
// internal/query/schema/validator.go
type SchemaValidator struct {
    sources map[string]DataSource
}

func (v *SchemaValidator) ValidateQuery(query string) (*ValidationResult, error) {
    ast, err := v.parser.Parse(query)
    if err != nil {
        return nil, err
    }
    
    return v.validateAST(ast)
}

func (v *SchemaValidator) validateAST(ast ASTNode) (*ValidationResult, error) {
    // Validate field references
    // Validate function calls
    // Validate join conditions
    // Validate aggregation usage
    return &ValidationResult{Valid: true}, nil
}
```

## Usage Examples

### Basic Queries
```bash
# Simple filtering
spooky query "SELECT * FROM machines WHERE tags.environment = 'production'"

# Field selection
spooky query "SELECT hostname, host, user FROM machines WHERE groups CONTAINS 'webservers'"

# Complex conditions
spooky query "SELECT hostname FROM machines WHERE resources.cpu_cores > 4 AND resources.memory_gb > 8"
```

### Advanced Queries
```bash
# Joins
spooky query "SELECT m.hostname, f.cpu_cores FROM machines m JOIN facts f ON m.hostname = f.machine_id"

# Aggregations
spooky query "SELECT tags.environment, COUNT(*) as count FROM machines GROUP BY tags.environment"

# Subqueries
spooky query "SELECT * FROM machines WHERE hostname IN (SELECT machine_id FROM facts WHERE cpu_cores > 4)"
```

### HCL Query Files
```bash
# Execute from file
spooky query --file queries/production_servers.hcl

# With parameters
spooky query --file queries/high_cpu_servers.hcl --parameters cpu_min=8

# Different output formats
spooky query "SELECT * FROM machines" --format json
spooky query "SELECT * FROM machines" --format table
spooky query "SELECT * FROM machines" --format csv
```

## Benefits

1. **Familiarity**: SQL-like syntax is widely understood
2. **Powerful**: Supports complex queries, joins, aggregations
3. **Extensible**: Easy to add new functions and data sources
4. **HCL Integration**: Works seamlessly with existing configuration
5. **Performance**: Query optimization and caching
6. **Type Safety**: Leverages existing type system
7. **Composability**: Queries can reference other queries
8. **Reusability**: Query templates and parameterization

## Success Metrics

1. **Query Performance**: Sub-second response times for typical queries
2. **Memory Usage**: Efficient memory usage for large datasets
3. **Cache Hit Rate**: >80% cache hit rate for repeated queries
4. **User Adoption**: >50% of users using query language within 3 months
5. **Query Complexity**: Support for queries with >5 joins and >10 conditions

## Future Enhancements

1. **Query Builder**: Visual query builder interface
2. **Query History**: Track and replay previous queries
3. **Query Scheduling**: Scheduled query execution
4. **Query Alerts**: Alert on query results
5. **Query Analytics**: Analyze query patterns and performance
6. **Distributed Queries**: Query across multiple spooky instances
7. **Real-time Queries**: Streaming query results
8. **Query Templates**: Library of common query patterns

## Conclusion

The query language system will provide spooky users with powerful, familiar, and extensible querying capabilities while maintaining integration with the existing architecture. The phased implementation approach ensures steady progress and allows for feedback and iteration throughout the development process.
