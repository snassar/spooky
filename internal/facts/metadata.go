package facts

import (
	"time"
)

// MetadataKey represents standard metadata keys
type MetadataKey string

const (
	// Source metadata
	MetadataKeySourceType MetadataKey = "source_type"
	MetadataKeySourcePath MetadataKey = "source_path"
	MetadataKeySourceURL  MetadataKey = "source_url"
	MetadataKeySourceFile MetadataKey = "source_file"

	// Collection metadata
	MetadataKeyCollectorType  MetadataKey = "collector_type"
	MetadataKeyFormat         MetadataKey = "format"
	MetadataKeyMergePolicy    MetadataKey = "merge_policy"
	MetadataKeyCollectionTime MetadataKey = "collection_time"

	// HTTP metadata
	MetadataKeyHTTPStatus  MetadataKey = "http_status"
	MetadataKeyHTTPHeaders MetadataKey = "http_headers"
	MetadataKeyHTTPTimeout MetadataKey = "http_timeout"

	// File metadata
	MetadataKeyFileSize     MetadataKey = "file_size"
	MetadataKeyFileModified MetadataKey = "file_modified"

	// SSH metadata
	MetadataKeySSHHost MetadataKey = "ssh_host"
	MetadataKeySSHUser MetadataKey = "ssh_user"
	MetadataKeySSHPort MetadataKey = "ssh_port"

	// Custom metadata
	MetadataKeyCustom MetadataKey = "custom"
)

// MetadataBuilder helps build standardized metadata
type MetadataBuilder struct {
	metadata map[string]interface{}
}

// NewMetadataBuilder creates a new metadata builder
func NewMetadataBuilder() *MetadataBuilder {
	return &MetadataBuilder{
		metadata: make(map[string]interface{}),
	}
}

// WithSourceType sets the source type
func (b *MetadataBuilder) WithSourceType(sourceType string) *MetadataBuilder {
	b.metadata[string(MetadataKeySourceType)] = sourceType
	return b
}

// WithSourcePath sets the source path
func (b *MetadataBuilder) WithSourcePath(path string) *MetadataBuilder {
	b.metadata[string(MetadataKeySourcePath)] = path
	return b
}

// WithSourceURL sets the source URL
func (b *MetadataBuilder) WithSourceURL(url string) *MetadataBuilder {
	b.metadata[string(MetadataKeySourceURL)] = url
	return b
}

// WithSourceFile sets the source file
func (b *MetadataBuilder) WithSourceFile(file string) *MetadataBuilder {
	b.metadata[string(MetadataKeySourceFile)] = file
	return b
}

// WithCollectorType sets the collector type
func (b *MetadataBuilder) WithCollectorType(collectorType string) *MetadataBuilder {
	b.metadata[string(MetadataKeyCollectorType)] = collectorType
	return b
}

// WithFormat sets the data format
func (b *MetadataBuilder) WithFormat(format string) *MetadataBuilder {
	b.metadata[string(MetadataKeyFormat)] = format
	return b
}

// WithMergePolicy sets the merge policy
func (b *MetadataBuilder) WithMergePolicy(policy MergePolicy) *MetadataBuilder {
	b.metadata[string(MetadataKeyMergePolicy)] = string(policy)
	return b
}

// WithCollectionTime sets the collection time
func (b *MetadataBuilder) WithCollectionTime(t time.Time) *MetadataBuilder {
	b.metadata[string(MetadataKeyCollectionTime)] = t.Format(time.RFC3339)
	return b
}

// WithHTTPStatus sets the HTTP status
func (b *MetadataBuilder) WithHTTPStatus(status int) *MetadataBuilder {
	b.metadata[string(MetadataKeyHTTPStatus)] = status
	return b
}

// WithHTTPHeaders sets the HTTP headers
func (b *MetadataBuilder) WithHTTPHeaders(headers map[string]string) *MetadataBuilder {
	b.metadata[string(MetadataKeyHTTPHeaders)] = headers
	return b
}

// WithHTTPTimeout sets the HTTP timeout
func (b *MetadataBuilder) WithHTTPTimeout(timeout time.Duration) *MetadataBuilder {
	b.metadata[string(MetadataKeyHTTPTimeout)] = timeout.String()
	return b
}

// WithFileSize sets the file size
func (b *MetadataBuilder) WithFileSize(size int64) *MetadataBuilder {
	b.metadata[string(MetadataKeyFileSize)] = size
	return b
}

// WithFileModified sets the file modification time
func (b *MetadataBuilder) WithFileModified(t time.Time) *MetadataBuilder {
	b.metadata[string(MetadataKeyFileModified)] = t.Format(time.RFC3339)
	return b
}

// WithSSHHost sets the SSH host
func (b *MetadataBuilder) WithSSHHost(host string) *MetadataBuilder {
	b.metadata[string(MetadataKeySSHHost)] = host
	return b
}

// WithSSHUser sets the SSH user
func (b *MetadataBuilder) WithSSHUser(user string) *MetadataBuilder {
	b.metadata[string(MetadataKeySSHUser)] = user
	return b
}

// WithSSHPort sets the SSH port
func (b *MetadataBuilder) WithSSHPort(port int) *MetadataBuilder {
	b.metadata[string(MetadataKeySSHPort)] = port
	return b
}

// WithCustom adds custom metadata
func (b *MetadataBuilder) WithCustom(key string, value interface{}) *MetadataBuilder {
	if b.metadata[string(MetadataKeyCustom)] == nil {
		b.metadata[string(MetadataKeyCustom)] = make(map[string]interface{})
	}
	custom := b.metadata[string(MetadataKeyCustom)].(map[string]interface{})
	custom[key] = value
	return b
}

// Build returns the built metadata map
func (b *MetadataBuilder) Build() map[string]interface{} {
	// Add collection time if not set
	if _, exists := b.metadata[string(MetadataKeyCollectionTime)]; !exists {
		b.metadata[string(MetadataKeyCollectionTime)] = time.Now().Format(time.RFC3339)
	}
	return b.metadata
}

// StandardMetadata creates standard metadata for a collector
func StandardMetadata(collectorType, source, format string) map[string]interface{} {
	return NewMetadataBuilder().
		WithCollectorType(collectorType).
		WithSourcePath(source).
		WithFormat(format).
		Build()
}
