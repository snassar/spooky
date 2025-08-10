package schemas

// SchemaType represents the type of schema
type SchemaType string

// Schema represents a loaded schema
type Schema struct {
	Type     SchemaType
	Content  string
	Filename string
}

// Config holds schema configuration
type Config struct {
	CacheEnabled bool   `json:"cache_enabled"`
	CacheSize    int    `json:"cache_size"`
	LogLevel     string `json:"log_level"`
}
