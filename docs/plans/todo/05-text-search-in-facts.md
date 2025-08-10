# Implementation Plan: Text Search in Facts

## Overview
Implement full-text search functionality for facts storage, enabling users to search through fact data using various search patterns, filters, and ranking algorithms.

## Task Details
- **Task ID**: 6.3
- **Priority**: Very Low
- **File**: `internal/facts/search/`
- **Functions**: Full-text search functionality

## Implementation Requirements

### Interface Compliance
The text search functionality must:
1. **Index fact data** for efficient full-text search
2. **Support multiple search patterns** (exact match, partial match, regex, fuzzy)
3. **Provide search ranking** and relevance scoring
4. **Support complex queries** with boolean operators
5. **Handle large datasets** efficiently
6. **Provide search metadata** and statistics
7. **Support search filters** by machine, fact type, timestamp
8. **Integrate with existing** fact storage system

## Detailed Implementation Plan

### Step 1: Define Search Interfaces

**File**: `internal/facts/search/interfaces.go`

```go
package search

import (
    "context"
    "time"
    "spooky/internal/facts/types"
)

// SearchEngine handles full-text search operations
type SearchEngine interface {
    IndexFacts(ctx context.Context, collection *types.FactCollection) error
    Search(ctx context.Context, query *SearchQuery) (*SearchResult, error)
    SearchByMachine(ctx context.Context, machine string, query *SearchQuery) (*SearchResult, error)
    SearchByTimeRange(ctx context.Context, start, end time.Time, query *SearchQuery) (*SearchResult, error)
    GetSearchStats(ctx context.Context) (*SearchStats, error)
    RebuildIndex(ctx context.Context) error
    Close() error
}

// SearchQuery represents a search query
type SearchQuery struct {
    Query         string                 `json:"query"`
    Pattern       SearchPattern          `json:"pattern"`
    Filters       *SearchFilters         `json:"filters"`
    Limit         int                    `json:"limit"`
    Offset        int                    `json:"offset"`
    SortBy        string                 `json:"sort_by"`
    SortOrder     SortOrder              `json:"sort_order"`
    IncludeScore  bool                   `json:"include_score"`
    Metadata      map[string]interface{} `json:"metadata"`
}

// SearchPattern represents the search pattern type
type SearchPattern string

const (
    SearchPatternExact    SearchPattern = "exact"
    SearchPatternPartial  SearchPattern = "partial"
    SearchPatternRegex    SearchPattern = "regex"
    SearchPatternFuzzy    SearchPattern = "fuzzy"
    SearchPatternWildcard SearchPattern = "wildcard"
)

// SortOrder represents the sort order
type SortOrder string

const (
    SortOrderAsc  SortOrder = "asc"
    SortOrderDesc SortOrder = "desc"
)

// SearchFilters represents search filters
type SearchFilters struct {
    Machines     []string               `json:"machines"`
    FactKeys     []string               `json:"fact_keys"`
    Sources      []string               `json:"sources"`
    StartTime    *time.Time             `json:"start_time"`
    EndTime      *time.Time             `json:"end_time"`
    MinScore     float64                `json:"min_score"`
    MaxScore     float64                `json:"max_score"`
    Metadata     map[string]interface{} `json:"metadata"`
}

// SearchResult represents search results
type SearchResult struct {
    Query       *SearchQuery    `json:"query"`
    Results     []*SearchMatch  `json:"results"`
    Total       int             `json:"total"`
    Limit       int             `json:"limit"`
    Offset      int             `json:"offset"`
    Duration    time.Duration   `json:"duration"`
    Stats       *SearchStats    `json:"stats"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// SearchMatch represents a search match
type SearchMatch struct {
    Fact        *types.Fact     `json:"fact"`
    Score       float64         `json:"score"`
    Highlights  []*Highlight    `json:"highlights"`
    Machine     string          `json:"machine"`
    Timestamp   time.Time       `json:"timestamp"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// Highlight represents a text highlight
type Highlight struct {
    Field       string          `json:"field"`
    Start       int             `json:"start"`
    End         int             `json:"end"`
    Text        string          `json:"text"`
    Score       float64         `json:"score"`
}

// SearchStats represents search statistics
type SearchStats struct {
    TotalFacts      int64                 `json:"total_facts"`
    TotalMachines   int64                 `json:"total_machines"`
    IndexSize       int64                 `json:"index_size"`
    LastIndexed     time.Time             `json:"last_indexed"`
    SearchCount     int64                 `json:"search_count"`
    AvgSearchTime   time.Duration         `json:"avg_search_time"`
    Metadata        map[string]interface{} `json:"metadata"`
}
```

### Step 2: Implement Bleve Search Engine

**File**: `internal/facts/search/bleve.go`

```go
package search

import (
    "context"
    "fmt"
    "sync"
    "time"
    "github.com/blevesearch/bleve/v2"
    "github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
    "github.com/blevesearch/bleve/v2/analysis/analyzer/standard"
    "github.com/blevesearch/bleve/v2/mapping"
    "spooky/internal/facts/types"
    "spooky/internal/logging"
)

// BleveSearchEngine implements search using Bleve
type BleveSearchEngine struct {
    index   bleve.Index
    mutex   sync.RWMutex
    logger  logging.Logger
}

// NewBleveSearchEngine creates a new Bleve search engine
func NewBleveSearchEngine(indexPath string, logger logging.Logger) (*BleveSearchEngine, error) {
    // Try to open existing index
    index, err := bleve.Open(indexPath)
    if err != nil {
        // Create new index
        index, err = createBleveIndex(indexPath)
        if err != nil {
            return nil, fmt.Errorf("failed to create Bleve index: %w", err)
        }
    }

    return &BleveSearchEngine{
        index:  index,
        logger: logger,
    }, nil
}

// createBleveIndex creates a new Bleve index
func createBleveIndex(indexPath string) (bleve.Index, error) {
    // Create document mapping
    docMapping := bleve.NewDocumentMapping()
    
    // Add field mappings
    docMapping.AddFieldMappingsAt("key", bleve.NewTextFieldMapping())
    docMapping.AddFieldMappingsAt("value", bleve.NewTextFieldMapping())
    docMapping.AddFieldMappingsAt("source", bleve.NewTextFieldMapping())
    docMapping.AddFieldMappingsAt("server", bleve.NewTextFieldMapping())
    docMapping.AddFieldMappingsAt("timestamp", bleve.NewDateTimeFieldMapping())
    
    // Add sub-fields for value
    valueMapping := bleve.NewTextFieldMapping()
    valueMapping.Analyzer = standard.Name
    docMapping.AddFieldMappingsAt("value_text", valueMapping)
    
    // Add keyword field for exact matching
    keywordMapping := bleve.NewTextFieldMapping()
    keywordMapping.Analyzer = keyword.Name
    docMapping.AddFieldMappingsAt("value_keyword", keywordMapping)

    // Create index mapping
    indexMapping := bleve.NewIndexMapping()
    indexMapping.DefaultMapping = docMapping
    indexMapping.DefaultAnalyzer = standard.Name

    // Create index
    return bleve.New(indexPath, indexMapping)
}

// IndexFacts indexes facts for search
func (e *BleveSearchEngine) IndexFacts(ctx context.Context, collection *types.FactCollection) error {
    e.mutex.Lock()
    defer e.mutex.Unlock()

    batch := e.index.NewBatch()
    
    for key, fact := range collection.Facts {
        // Create document for indexing
        doc := map[string]interface{}{
            "id":          fmt.Sprintf("%s:%s", collection.Server, key),
            "key":         fact.Key,
            "value":       fact.Value,
            "value_text":  fmt.Sprintf("%v", fact.Value),
            "value_keyword": fmt.Sprintf("%v", fact.Value),
            "source":      fact.Source,
            "server":      fact.Server,
            "timestamp":   fact.Timestamp,
            "machine":     collection.Server,
        }

        // Index document
        if err := batch.Index(fmt.Sprintf("%s:%s", collection.Server, key), doc); err != nil {
            return fmt.Errorf("failed to index fact: %w", err)
        }
    }

    // Commit batch
    if err := e.index.Batch(batch); err != nil {
        return fmt.Errorf("failed to commit batch: %w", err)
    }

    return nil
}

// Search performs a search query
func (e *BleveSearchEngine) Search(ctx context.Context, query *SearchQuery) (*SearchResult, error) {
    e.mutex.RLock()
    defer e.mutex.RUnlock()

    startTime := time.Now()

    // Build Bleve query
    bleveQuery, err := e.buildBleveQuery(query)
    if err != nil {
        return nil, fmt.Errorf("failed to build Bleve query: %w", err)
    }

    // Create search request
    searchRequest := bleve.NewSearchRequest(bleveQuery)
    searchRequest.Size = query.Limit
    searchRequest.From = query.Offset
    searchRequest.IncludeLocations = true

    // Add sorting
    if query.SortBy != "" {
        searchRequest.SortBy([]string{query.SortBy})
    }

    // Add highlighting
    if query.IncludeScore {
        searchRequest.Highlight = bleve.NewHighlight()
        searchRequest.Highlight.Fields = []string{"value_text", "key"}
    }

    // Act search
    searchResult, err := e.index.Search(searchRequest)
    if err != nil {
        return nil, fmt.Errorf("failed to act search: %w", err)
    }

    // Convert results
    results := e.convertSearchResults(searchResult, query)

    duration := time.Since(startTime)

    return &SearchResult{
        Query:    query,
        Results:  results,
        Total:    int(searchResult.Total),
        Limit:    query.Limit,
        Offset:   query.Offset,
        Duration: duration,
    }, nil
}

// buildBleveQuery builds a Bleve query from search query
func (e *BleveSearchEngine) buildBleveQuery(query *SearchQuery) (bleve.Query, error) {
    var queries []bleve.Query

    // Build main query based on pattern
    switch query.Pattern {
    case SearchPatternExact:
        queries = append(queries, bleve.NewTermQuery(query.Query))
    case SearchPatternPartial:
        queries = append(queries, bleve.NewWildcardQuery("*" + query.Query + "*"))
    case SearchPatternRegex:
        regexQuery := bleve.NewRegexpQuery(query.Query)
        queries = append(queries, regexQuery)
    case SearchPatternFuzzy:
        fuzzyQuery := bleve.NewFuzzyQuery(query.Query)
        queries = append(queries, fuzzyQuery)
    case SearchPatternWildcard:
        wildcardQuery := bleve.NewWildcardQuery(query.Query)
        queries = append(queries, wildcardQuery)
    default:
        // Default to text query
        textQuery := bleve.NewQueryStringQuery(query.Query)
        queries = append(queries, textQuery)
    }

    // Add filters
    if query.Filters != nil {
        if len(query.Filters.Machines) > 0 {
            machineQueries := make([]bleve.Query, len(query.Filters.Machines))
            for i, machine := range query.Filters.Machines {
                machineQueries[i] = bleve.NewTermQuery(machine)
            }
            queries = append(queries, bleve.NewDisjunctionQuery(machineQueries...))
        }

        if len(query.Filters.FactKeys) > 0 {
            keyQueries := make([]bleve.Query, len(query.Filters.FactKeys))
            for i, key := range query.Filters.FactKeys {
                keyQueries[i] = bleve.NewTermQuery(key)
            }
            queries = append(queries, bleve.NewDisjunctionQuery(keyQueries...))
        }

        if len(query.Filters.Sources) > 0 {
            sourceQueries := make([]bleve.Query, len(query.Filters.Sources))
            for i, source := range query.Filters.Sources {
                sourceQueries[i] = bleve.NewTermQuery(source)
            }
            queries = append(queries, bleve.NewDisjunctionQuery(sourceQueries...))
        }

        if query.Filters.StartTime != nil || query.Filters.EndTime != nil {
            var min, max *time.Time
            if query.Filters.StartTime != nil {
                min = query.Filters.StartTime
            }
            if query.Filters.EndTime != nil {
                max = query.Filters.EndTime
            }
            timeQuery := bleve.NewDateRangeQuery(min, max)
            timeQuery.Field = "timestamp"
            queries = append(queries, timeQuery)
        }
    }

    // Combine queries with AND
    if len(queries) == 1 {
        return queries[0], nil
    }
    return bleve.NewConjunctionQuery(queries...), nil
}

// convertSearchResults converts Bleve search results
func (e *BleveSearchEngine) convertSearchResults(bleveResult *bleve.SearchResult, query *SearchQuery) []*SearchMatch {
    results := make([]*SearchMatch, 0, len(bleveResult.Hits))

    for _, hit := range bleveResult.Hits {
        // Extract fact data from hit
        fact := &types.Fact{
            Key:       hit.Fields["key"].(string),
            Value:     hit.Fields["value"],
            Source:    hit.Fields["source"].(string),
            Server:    hit.Fields["server"].(string),
            Timestamp: hit.Fields["timestamp"].(time.Time),
        }

        // Create highlights
        var highlights []*Highlight
        if hit.Locations != nil {
            for field, locations := range hit.Locations {
                for _, location := range locations {
                    highlight := &Highlight{
                        Field: field,
                        Start: location.Start,
                        End:   location.End,
                        Text:  location.Term,
                        Score: float64(location.Pos),
                    }
                    highlights = append(highlights, highlight)
                }
            }
        }

        match := &SearchMatch{
            Fact:       fact,
            Score:      hit.Score,
            Highlights: highlights,
            Machine:    hit.Fields["machine"].(string),
            Timestamp:  hit.Fields["timestamp"].(time.Time),
            Metadata:   make(map[string]interface{}),
        }

        results = append(results, match)
    }

    return results
}

// SearchByMachine searches facts for a specific machine
func (e *BleveSearchEngine) SearchByMachine(ctx context.Context, machine string, query *SearchQuery) (*SearchResult, error) {
    // Add machine filter
    if query.Filters == nil {
        query.Filters = &SearchFilters{}
    }
    query.Filters.Machines = []string{machine}

    return e.Search(ctx, query)
}

// SearchByTimeRange searches facts within a time range
func (e *BleveSearchEngine) SearchByTimeRange(ctx context.Context, start, end time.Time, query *SearchQuery) (*SearchResult, error) {
    // Add time range filter
    if query.Filters == nil {
        query.Filters = &SearchFilters{}
    }
    query.Filters.StartTime = &start
    query.Filters.EndTime = &end

    return e.Search(ctx, query)
}

// GetSearchStats returns search statistics
func (e *BleveSearchEngine) GetSearchStats(ctx context.Context) (*SearchStats, error) {
    e.mutex.RLock()
    defer e.mutex.RUnlock()

    // Get index stats
    stats := e.index.Stats()

    return &SearchStats{
        TotalFacts:    stats.Total,
        IndexSize:     stats.Size,
        LastIndexed:   time.Now(),
        SearchCount:   0,
        AvgSearchTime: 0,
        Metadata:      make(map[string]interface{}),
    }, nil
}

// RebuildIndex rebuilds the search index
func (e *BleveSearchEngine) RebuildIndex(ctx context.Context) error {
    e.mutex.Lock()
    defer e.mutex.Unlock()

    // Close current index
    if err := e.index.Close(); err != nil {
        return fmt.Errorf("failed to close index: %w", err)
    }

    // Recreate index
    index, err := createBleveIndex(e.index.Name())
    if err != nil {
        return fmt.Errorf("failed to recreate index: %w", err)
    }

    e.index = index
    return nil
}

// Close closes the search engine
func (e *BleveSearchEngine) Close() error {
    e.mutex.Lock()
    defer e.mutex.Unlock()

    return e.index.Close()
}
```

### Step 3: Implement Search Manager

**File**: `internal/facts/search/manager.go`

```go
package search

import (
    "context"
    "fmt"
    "time"
    "spooky/internal/facts/storage"
    "spooky/internal/facts/types"
    "spooky/internal/logging"
)

// Manager manages search operations
type Manager struct {
    engine      SearchEngine
    storage     storage.Storage
    logger      logging.Logger
    indexPath   string
}

// NewManager creates a new search manager
func NewManager(storage storage.Storage, indexPath string, logger logging.Logger) (*Manager, error) {
    // Create search engine
    engine, err := NewBleveSearchEngine(indexPath, logger)
    if err != nil {
        return nil, fmt.Errorf("failed to create search engine: %w", err)
    }

    return &Manager{
        engine:    engine,
        storage:   storage,
        logger:    logger,
        indexPath: indexPath,
    }, nil
}

// IndexAllFacts indexes all facts from storage
func (m *Manager) IndexAllFacts(ctx context.Context) error {
    m.logger.Info("Starting full fact indexing")

    // List all servers
    servers, err := m.storage.List(ctx)
    if err != nil {
        return fmt.Errorf("failed to list servers: %w", err)
    }

    indexedCount := 0
    for _, server := range servers {
        // Load collection
        collection, err := m.storage.Load(ctx, server)
        if err != nil {
            m.logger.Warn("Failed to load collection for indexing",
                logging.String("server", server),
                logging.Error(err))
            continue
        }

        // Index collection
        if err := m.engine.IndexFacts(ctx, collection); err != nil {
            m.logger.Warn("Failed to index collection",
                logging.String("server", server),
                logging.Error(err))
            continue
        }

        indexedCount += len(collection.Facts)
    }

    m.logger.Info("Completed fact indexing",
        logging.Int("server_count", len(servers)),
        logging.Int("fact_count", indexedCount))

    return nil
}

// Search performs a search operation
func (m *Manager) Search(ctx context.Context, queryString string) (*SearchResult, error) {
    // Create search query
    query := &SearchQuery{
        Query:        queryString,
        Pattern:      SearchPatternPartial, // Default pattern
        Limit:        100,                  // Default limit
        IncludeScore: true,
        Metadata:     make(map[string]interface{}),
    }

    // Act search
    result, err := m.engine.Search(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to act search: %w", err)
    }

    m.logger.Debug("Search completed",
        logging.String("query", queryString),
        logging.Int("result_count", len(result.Results)),
        logging.Duration("duration", result.Duration))

    return result, nil
}

// SearchByMachine searches facts for a specific machine
func (m *Manager) SearchByMachine(ctx context.Context, machine, queryString string) (*SearchResult, error) {
    // Create search query
    query := &SearchQuery{
        Query:        queryString,
        Pattern:      SearchPatternPartial,
        Limit:        100,
        IncludeScore: true,
        Filters: &SearchFilters{
            Machines: []string{machine},
        },
        Metadata: make(map[string]interface{}),
    }

    // Act search
    result, err := m.engine.Search(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to act search: %w", err)
    }

    return result, nil
}

// SearchByTimeRange searches facts within a time range
func (m *Manager) SearchByTimeRange(ctx context.Context, start, end time.Time, queryString string) (*SearchResult, error) {
    // Create search query
    query := &SearchQuery{
        Query:        queryString,
        Pattern:      SearchPatternPartial,
        Limit:        100,
        IncludeScore: true,
        Filters: &SearchFilters{
            StartTime: &start,
            EndTime:   &end,
        },
        Metadata: make(map[string]interface{}),
    }

    // Act search
    result, err := m.engine.Search(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to act search: %w", err)
    }

    return result, nil
}

// GetSearchStats returns search statistics
func (m *Manager) GetSearchStats(ctx context.Context) (*SearchStats, error) {
    return m.engine.GetSearchStats(ctx)
}

// RebuildIndex rebuilds the search index
func (m *Manager) RebuildIndex(ctx context.Context) error {
    m.logger.Info("Rebuilding search index")

    // Rebuild index
    if err := m.engine.RebuildIndex(ctx); err != nil {
        return fmt.Errorf("failed to rebuild index: %w", err)
    }

    // Re-index all facts
    if err := m.IndexAllFacts(ctx); err != nil {
        return fmt.Errorf("failed to re-index facts: %w", err)
    }

    m.logger.Info("Search index rebuild completed")
    return nil
}

// Close closes the search manager
func (m *Manager) Close() error {
    return m.engine.Close()
}

// GetSearchInfo returns search information
func (m *Manager) GetSearchInfo() map[string]interface{} {
    return map[string]interface{}{
        "index_path": m.indexPath,
        "engine":     "bleve",
        "enabled":    true,
    }
}
```

### Step 4: Implement Search CLI Commands

**File**: `internal/cli/commands/search.go`

```go
package commands

import (
    "context"
    "fmt"
    "strings"
    "time"

    "github.com/spf13/cobra"
    "spooky/internal/facts/search"
    "spooky/internal/logging"
)

// SearchManager handles search commands
type SearchManager struct {
    searchManager *search.Manager
    logger        logging.Logger
}

// NewSearchManager creates a new search manager
func NewSearchManager(searchManager *search.Manager, logger logging.Logger) *SearchManager {
    return &SearchManager{
        searchManager: searchManager,
        logger:        logger,
    }
}

// NewSearchCommand creates the search command
func NewSearchCommand(manager *SearchManager) *cobra.Command {
    var (
        machine   string
        timeRange string
        limit     int
        output    string
    )

    cmd := &cobra.Command{
        Use:   "search <project> <query>",
        Short: "Search facts",
        Long: `Search facts using full-text search.

Examples:
  spooky search ./my-project "nginx"
  spooky search ./my-project "machine:server1 nginx"
  spooky search ./my-project "time:1d error"`,
        Args: cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            projectPath := args[0]
            queryString := args[1]

            return manager.SearchFacts(cmd.Context(), projectPath, queryString, &SearchOptions{
                Machine:   machine,
                TimeRange: timeRange,
                Limit:     limit,
                Output:    output,
            })
        },
    }

    // Add flags
    cmd.Flags().StringVarP(&machine, "machine", "m", "", "Search in specific machine")
    cmd.Flags().StringVarP(&timeRange, "time", "t", "", "Time range (e.g., 1h, 1d, 1w)")
    cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Limit number of results")
    cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format (text, json, table)")

    return cmd
}

// SearchOptions represents search options
type SearchOptions struct {
    Machine   string
    TimeRange string
    Limit     int
    Output    string
}

// SearchFacts searches facts
func (m *SearchManager) SearchFacts(ctx context.Context, projectPath, queryString string, options *SearchOptions) error {
    m.logger.Info("Searching facts",
        logging.String("project", projectPath),
        logging.String("query", queryString))

    var result *search.SearchResult
    var err error

    // Act search based on options
    if options.Machine != "" {
        result, err = m.searchManager.SearchByMachine(ctx, options.Machine, queryString)
    } else if options.TimeRange != "" {
        start, end, err := m.parseTimeRange(options.TimeRange)
        if err != nil {
            return fmt.Errorf("failed to parse time range: %w", err)
        }
        result, err = m.searchManager.SearchByTimeRange(ctx, start, end, queryString)
    } else {
        result, err = m.searchManager.Search(ctx, queryString)
    }

    if err != nil {
        return fmt.Errorf("search failed: %w", err)
    }

    // Display results
    return m.displaySearchResults(result, options.Output)
}

// parseTimeRange parses time range specification
func (m *SearchManager) parseTimeRange(timeSpec string) (time.Time, time.Time, error) {
    now := time.Now()
    var startTime time.Time

    switch {
    case strings.HasSuffix(timeSpec, "h"):
        hours := parseInt(strings.TrimSuffix(timeSpec, "h"))
        startTime = now.Add(-time.Duration(hours) * time.Hour)
    case strings.HasSuffix(timeSpec, "d"):
        days := parseInt(strings.TrimSuffix(timeSpec, "d"))
        startTime = now.AddDate(0, 0, -days)
    case strings.HasSuffix(timeSpec, "w"):
        weeks := parseInt(strings.TrimSuffix(timeSpec, "w"))
        startTime = now.AddDate(0, 0, -weeks*7)
    case strings.HasSuffix(timeSpec, "m"):
        months := parseInt(strings.TrimSuffix(timeSpec, "m"))
        startTime = now.AddDate(0, -months, 0)
    case strings.HasSuffix(timeSpec, "y"):
        years := parseInt(strings.TrimSuffix(timeSpec, "y"))
        startTime = now.AddDate(-years, 0, 0)
    default:
        return time.Time{}, time.Time{}, fmt.Errorf("invalid time specification: %s", timeSpec)
    }

    return startTime, now, nil
}

// parseInt parses an integer string
func parseInt(s string) int {
    if i, err := strconv.Atoi(s); err == nil {
        return i
    }
    return 0
}

// displaySearchResults displays search results
func (m *SearchManager) displaySearchResults(result *search.SearchResult, format string) error {
    switch format {
    case "json":
        return m.displayJSONResults(result)
    case "table":
        return m.displayTableResults(result)
    case "text":
    default:
        return m.displayTextResults(result)
    }
    return nil
}

// displayTextResults displays results in text format
func (m *SearchManager) displayTextResults(result *search.SearchResult) error {
    fmt.Printf("Search Results (%d total, %d shown)\n", result.Total, len(result.Results))
    fmt.Printf("Query: %s\n", result.Query.Query)
    fmt.Printf("Duration: %v\n\n", result.Duration)

    for i, match := range result.Results {
        fmt.Printf("%d. [%s] %s = %v (score: %.2f)\n", 
            i+1, match.Machine, match.Fact.Key, match.Fact.Value, match.Score)
        
        if len(match.Highlights) > 0 {
            fmt.Printf("   Highlights: ")
            for j, highlight := range match.Highlights {
                if j > 0 {
                    fmt.Printf(", ")
                }
                fmt.Printf("%s", highlight.Text)
            }
            fmt.Printf("\n")
        }
        fmt.Printf("\n")
    }

    return nil
}

// displayJSONResults displays results in JSON format
func (m *SearchManager) displayJSONResults(result *search.SearchResult) error {
    fmt.Printf("{\n")
    fmt.Printf("  \"query\": \"%s\",\n", result.Query.Query)
    fmt.Printf("  \"total\": %d,\n", result.Total)
    fmt.Printf("  \"duration\": \"%v\",\n", result.Duration)
    fmt.Printf("  \"results\": [\n")
    
    for i, match := range result.Results {
        fmt.Printf("    {\n")
        fmt.Printf("      \"machine\": \"%s\",\n", match.Machine)
        fmt.Printf("      \"key\": \"%s\",\n", match.Fact.Key)
        fmt.Printf("      \"value\": %v,\n", match.Fact.Value)
        fmt.Printf("      \"score\": %.2f\n", match.Score)
        fmt.Printf("    }")
        if i < len(result.Results)-1 {
            fmt.Printf(",")
        }
        fmt.Printf("\n")
    }
    
    fmt.Printf("  ]\n")
    fmt.Printf("}\n")
    
    return nil
}

// displayTableResults displays results in table format
func (m *SearchManager) displayTableResults(result *search.SearchResult) error {
    fmt.Printf("Search Results (%d total, %d shown)\n", result.Total, len(result.Results))
    fmt.Printf("Query: %s | Duration: %v\n\n", result.Query.Query, result.Duration)
    
    fmt.Printf("%-20s %-15s %-30s %-10s\n", "Machine", "Key", "Value", "Score")
    fmt.Println(strings.Repeat("-", 80))
    
    for _, match := range result.Results {
        valueStr := fmt.Sprintf("%v", match.Fact.Value)
        if len(valueStr) > 28 {
            valueStr = valueStr[:25] + "..."
        }
        
        fmt.Printf("%-20s %-15s %-30s %-10.2f\n", 
            match.Machine, match.Fact.Key, valueStr, match.Score)
    }
    
    return nil
}
```


2. **Large dataset** performance testing
3. **Complex queries** with filters
4. **Index rebuilding** scenarios


- Index rebuilding
- Query validation
- Performance monitoring



## Dependencies

### Internal Dependencies
- `spooky/internal/facts/storage`
- `spooky/internal/facts/types`
- `spooky/internal/logging`

### External Dependencies
- `github.com/blevesearch/bleve/v2`
- `regexp` (standard library)
- `strings` (standard library)
- `time` (standard library)



## Implementation Order

1. Define search interfaces
2. Implement search index with Bleve
3. Implement search manager
4. Add search CLI commands
5. Write comprehensive tests
6. Performance optimization
7. Documentation and cleanup


