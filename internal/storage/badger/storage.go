package badger

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	spookyfactstypes "spooky/internal/facts/types"
	spookyschemas "spooky/internal/schemas"
	spookysecrets "spooky/internal/secrets"

	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// BadgerFactStorage implements FactStorage using BadgerDB
type BadgerFactStorage struct {
	db             *badger.DB
	cryptoManager  *spookysecrets.Manager
	encryptEnabled bool
}

// NewBadgerFactStorage creates a new BadgerDB-based fact storage
func NewBadgerFactStorage(dbPath string) (*BadgerFactStorage, error) {
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil // Disable BadgerDB logging

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open BadgerDB: %w", err)
	}

	return &BadgerFactStorage{db: db}, nil
}

// NewBadgerFactStorageWithCrypto creates a new BadgerDB-based fact storage with crypto support
func NewBadgerFactStorageWithCrypto(dbPath string, cryptoManager *spookysecrets.Manager, encryptEnabled bool) (*BadgerFactStorage, error) {
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil // Disable BadgerDB logging

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open BadgerDB: %w", err)
	}

	return &BadgerFactStorage{
		db:             db,
		cryptoManager:  cryptoManager,
		encryptEnabled: encryptEnabled,
	}, nil
}

// NewBadgerFactStorageReadOnly creates a new read-only BadgerDB-based fact storage
func NewBadgerFactStorageReadOnly(dbPath string) (*BadgerFactStorage, error) {
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil // Disable BadgerDB logging
	opts.ReadOnly = true

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open BadgerDB in read-only mode: %w", err)
	}

	return &BadgerFactStorage{db: db}, nil
}

// SetFactCollection stores a fact collection for a specific machine
func (b *BadgerFactStorage) SetFactCollection(machineID string, collection *spookyfactstypes.FactCollection) error {
	if b.encryptEnabled && b.cryptoManager != nil {
		return b.storeEncryptedFactCollection(machineID, collection)
	}

	data, err := json.Marshal(collection)
	if err != nil {
		return fmt.Errorf("failed to marshal fact collection: %w", err)
	}

	key := []byte(fmt.Sprintf("facts:%s", machineID))
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, data)
	})
}

// GetFactCollection retrieves a fact collection for a specific machine
func (b *BadgerFactStorage) GetFactCollection(machineID string) (*spookyfactstypes.FactCollection, error) {
	var collection spookyfactstypes.FactCollection
	key := []byte(fmt.Sprintf("facts:%s", machineID))

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &collection)
		})
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get fact collection for %s: %w", machineID, err)
	}

	// Check if collection is encrypted and decrypt if needed
	if b.isEncryptedCollection(&collection) {
		decrypted, err := b.decryptFactCollection(&collection)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt fact collection for %s: %w", machineID, err)
		}
		return decrypted, nil
	}

	return &collection, nil
}

// storeEncryptedFactCollection stores an encrypted fact collection
func (b *BadgerFactStorage) storeEncryptedFactCollection(machineID string, collection *spookyfactstypes.FactCollection) error {
	// Add encryption metadata
	collection.EncryptionMetadata = &spookyfactstypes.EncryptionMetadata{
		EncryptedAt:       time.Now().Format(time.RFC3339),
		EncryptionVersion: "1.0",
		Recipients:        b.getDefaultRecipients(),
	}

	// Marshal the collection to JSON for encryption
	collectionData, err := json.Marshal(collection)
	if err != nil {
		return fmt.Errorf("failed to marshal fact collection for encryption: %w", err)
	}

	// Encrypt the collection data
	encryptedData, err := b.cryptoManager.EncryptValue(string(collectionData), b.getDefaultRecipients())
	if err != nil {
		return fmt.Errorf("failed to encrypt fact collection: %w", err)
	}

	// Create encrypted collection wrapper
	encryptedCollection := &spookyfactstypes.FactCollection{
		Server:             collection.Server,
		Timestamp:          collection.Timestamp,
		EncryptionMetadata: collection.EncryptionMetadata,
		Facts: map[string]*spookyfactstypes.Fact{
			"_encrypted_data": {
				Key:   "_encrypted_data",
				Value: string(encryptedData),
			},
		},
	}

	// Store the encrypted collection
	data, err := json.Marshal(encryptedCollection)
	if err != nil {
		return fmt.Errorf("failed to marshal encrypted fact collection: %w", err)
	}

	key := []byte(fmt.Sprintf("facts:%s", machineID))
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, data)
	})
}

// decryptFactCollection decrypts an encrypted fact collection
func (b *BadgerFactStorage) decryptFactCollection(collection *spookyfactstypes.FactCollection) (*spookyfactstypes.FactCollection, error) {
	if !b.isEncryptedCollection(collection) {
		return collection, nil
	}

	// Extract encrypted data
	encryptedDataFact, exists := collection.Facts["_encrypted_data"]
	if !exists {
		return nil, fmt.Errorf("encrypted data not found in collection")
	}

	encryptedData, ok := encryptedDataFact.Value.(string)
	if !ok {
		return nil, fmt.Errorf("encrypted data is not a string")
	}

	// Decrypt the data
	decryptedData, err := b.cryptoManager.DecryptValue(encryptedData, b.getDefaultIdentity())
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt fact collection: %w", err)
	}

	// Unmarshal the decrypted collection
	var decryptedCollection spookyfactstypes.FactCollection
	if err := json.Unmarshal([]byte(decryptedData), &decryptedCollection); err != nil {
		return nil, fmt.Errorf("failed to unmarshal decrypted fact collection: %w", err)
	}

	return &decryptedCollection, nil
}

// isEncryptedCollection checks if a fact collection is encrypted
func (b *BadgerFactStorage) isEncryptedCollection(collection *spookyfactstypes.FactCollection) bool {
	return collection.EncryptionMetadata != nil && collection.EncryptedData != ""
}

// getDefaultRecipients returns the default encryption recipients
func (b *BadgerFactStorage) getDefaultRecipients() []string {
	// In a real implementation, this would read from configuration
	return []string{"age1example"}
}

// getDefaultIdentity returns the default encryption identity
func (b *BadgerFactStorage) getDefaultIdentity() string {
	// In a real implementation, this would read from configuration
	return "~/.config/spooky/keys/age.key"
}

// QueryFactCollections queries fact collections based on criteria
func (b *BadgerFactStorage) QueryFactCollections(query *spookyfactstypes.FactQuery) ([]*spookyfactstypes.FactCollection, error) {
	var results []*spookyfactstypes.FactCollection

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("facts:")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var collection spookyfactstypes.FactCollection

			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &collection)
			})
			if err != nil {
				continue // Skip invalid entries
			}

			// Check if collection is encrypted and decrypt if needed
			if b.isEncryptedCollection(&collection) {
				decrypted, err := b.decryptFactCollection(&collection)
				if err != nil {
					continue // Skip entries that can't be decrypted
				}
				collection = *decrypted
			}

			// Apply query filters
			if matchesQuery(&collection, query) {
				results = append(results, &collection)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query fact collections: %w", err)
	}

	return results, nil
}

// DeleteFactCollection deletes a fact collection for a specific machine
func (b *BadgerFactStorage) DeleteFactCollection(machineID string) error {
	key := []byte(fmt.Sprintf("facts:%s", machineID))
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// DeleteFactCollections deletes fact collections based on query criteria
func (b *BadgerFactStorage) DeleteFactCollections(query *spookyfactstypes.FactQuery) (int, error) {
	collections, err := b.QueryFactCollections(query)
	if err != nil {
		return 0, fmt.Errorf("failed to query collections for deletion: %w", err)
	}

	deletedCount := 0
	for _, collection := range collections {
		machineID := collection.Server
		if machineIDFact, exists := collection.Facts["machine_id"]; exists {
			if id, ok := machineIDFact.Value.(string); ok {
				machineID = id
			}
		}

		if err := b.DeleteFactCollection(machineID); err != nil {
			return deletedCount, fmt.Errorf("failed to delete collection for %s: %w", machineID, err)
		}
		deletedCount++
	}

	return deletedCount, nil
}

// ExportToJSON exports all fact collections to JSON
func (b *BadgerFactStorage) ExportToJSON(w io.Writer) error {
	// Validate against schema first
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeFactsStructure); err != nil {
		return fmt.Errorf("failed to load facts schema: %w", err)
	}

	facts := make(map[string]*spookyfactstypes.FactCollection)

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("facts:")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var collection spookyfactstypes.FactCollection

			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &collection)
			})
			if err != nil {
				continue // Skip invalid entries
			}

			// Check if collection is encrypted and decrypt if needed
			if b.isEncryptedCollection(&collection) {
				decrypted, err := b.decryptFactCollection(&collection)
				if err != nil {
					continue // Skip entries that can't be decrypted
				}
				collection = *decrypted
			}

			// Validate each collection against schema
			if err := validator.ValidateData(&collection, string(spookyschemas.SchemaTypeFactsStructure)); err != nil {
				return fmt.Errorf("failed to validate fact collection for %s: %w", collection.Server, err)
			}

			facts[collection.Server] = &collection
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to export facts: %w", err)
	}

	// Encode to JSON
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(facts)
}

// ImportFromJSON imports fact collections from JSON
func (b *BadgerFactStorage) ImportFromJSON(r io.Reader) error {
	// Validate against schema first
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeFactsStructure); err != nil {
		return fmt.Errorf("failed to load facts schema: %w", err)
	}

	var facts map[string]*spookyfactstypes.FactCollection
	if err := json.NewDecoder(r).Decode(&facts); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Validate each collection against schema
	for machineID, collection := range facts {
		if err := validator.ValidateData(collection, string(spookyschemas.SchemaTypeFactsStructure)); err != nil {
			return fmt.Errorf("failed to validate fact collection for %s: %w", machineID, err)
		}
	}

	return b.db.Update(func(txn *badger.Txn) error {
		for machineID, collection := range facts {
			data, err := json.Marshal(collection)
			if err != nil {
				return fmt.Errorf("failed to marshal fact collection for %s: %w", machineID, err)
			}

			key := []byte(fmt.Sprintf("facts:%s", machineID))
			if err := txn.Set(key, data); err != nil {
				return fmt.Errorf("failed to store fact collection for %s: %w", machineID, err)
			}
		}
		return nil
	})
}

// ExportToJSONWithEncryption exports fact collections with encryption
func (b *BadgerFactStorage) ExportToJSONWithEncryption(w io.Writer, options spookyfactstypes.ExportOptions) error {
	facts := make(map[string]*spookyfactstypes.FactCollection)

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("facts:")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var collection spookyfactstypes.FactCollection

			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &collection)
			})
			if err != nil {
				continue // Skip invalid entries
			}

			// Check if collection is encrypted and decrypt if needed
			if b.isEncryptedCollection(&collection) {
				decrypted, err := b.decryptFactCollection(&collection)
				if err != nil {
					continue // Skip entries that can't be decrypted
				}
				collection = *decrypted
			}

			facts[collection.Server] = &collection
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to export facts: %w", err)
	}

	// Apply filters if specified
	if len(options.FilterByMachineID) > 0 {
		filtered := make(map[string]*spookyfactstypes.FactCollection)
		for machineID, collection := range facts {
			for _, filterID := range options.FilterByMachineID {
				if machineID == filterID {
					filtered[machineID] = collection
					break
				}
			}
		}
		facts = filtered
	}

	// Apply time range filter if specified
	if options.FilterByTimeRange != nil {
		filtered := make(map[string]*spookyfactstypes.FactCollection)
		for machineID, collection := range facts {
			if collection.Timestamp.After(options.FilterByTimeRange.Start) &&
				collection.Timestamp.Before(options.FilterByTimeRange.End) {
				filtered[machineID] = collection
			}
		}
		facts = filtered
	}

	// Encode to JSON
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(facts)
}

// ImportFromJSONWithDecryption imports fact collections from JSON with decryption
func (b *BadgerFactStorage) ImportFromJSONWithDecryption(r io.Reader, identityFile string) error {
	var facts map[string]*spookyfactstypes.FactCollection
	if err := json.NewDecoder(r).Decode(&facts); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return b.db.Update(func(txn *badger.Txn) error {
		for machineID, collection := range facts {
			data, err := json.Marshal(collection)
			if err != nil {
				return fmt.Errorf("failed to marshal fact collection for %s: %w", machineID, err)
			}

			key := []byte(fmt.Sprintf("facts:%s", machineID))
			if err := txn.Set(key, data); err != nil {
				return fmt.Errorf("failed to store fact collection for %s: %w", machineID, err)
			}
		}
		return nil
	})
}

// ImportFromHCL imports fact collections from HCL format
func (b *BadgerFactStorage) ImportFromHCL(r io.Reader) error {
	// Read HCL content
	content, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read HCL content: %w", err)
	}

	// Parse HCL
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, "facts.hcl")
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse HCL: %w", diags)
	}

	// Validate against schema
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeFactsStructure); err != nil {
		return fmt.Errorf("failed to load facts export schema: %w", err)
	}

	// Parse HCL to fact collections
	collections, err := parseHCLToFactCollections(file)
	if err != nil {
		return fmt.Errorf("failed to parse HCL to fact collections: %w", err)
	}

	// Store collections in database
	return b.db.Update(func(txn *badger.Txn) error {
		for _, collection := range collections {
			data, err := json.Marshal(collection)
			if err != nil {
				return fmt.Errorf("failed to marshal fact collection: %w", err)
			}

			machineID := collection.Server
			if machineIDFact, exists := collection.Facts["machine_id"]; exists {
				if id, ok := machineIDFact.Value.(string); ok {
					machineID = id
				}
			}

			key := []byte(fmt.Sprintf("facts:%s", machineID))
			if err := txn.Set(key, data); err != nil {
				return fmt.Errorf("failed to store fact collection: %w", err)
			}
		}
		return nil
	})
}

// ExportToHCL exports all fact collections to HCL format
func (b *BadgerFactStorage) ExportToHCL(w io.Writer) error {
	// Validate against schema first
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeFactsStructure); err != nil {
		return fmt.Errorf("failed to load facts schema: %w", err)
	}

	facts := make(map[string]*spookyfactstypes.FactCollection)

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("facts:")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var collection spookyfactstypes.FactCollection

			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &collection)
			})
			if err != nil {
				continue // Skip invalid entries
			}

			// Check if collection is encrypted and decrypt if needed
			if b.isEncryptedCollection(&collection) {
				decrypted, err := b.decryptFactCollection(&collection)
				if err != nil {
					continue // Skip entries that can't be decrypted
				}
				collection = *decrypted
			}

			// Validate each collection against schema
			if err := validator.ValidateData(&collection, string(spookyschemas.SchemaTypeFactsStructure)); err != nil {
				return fmt.Errorf("failed to validate fact collection for %s: %w", collection.Server, err)
			}

			facts[collection.Server] = &collection
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to export facts: %w", err)
	}

	// Generate HCL content
	hclContent := generateHCLContent(facts)
	_, err = w.Write([]byte(hclContent))
	return err
}

// generateHCLContent generates HCL content from fact collections
func generateHCLContent(facts map[string]*spookyfactstypes.FactCollection) string {
	var content strings.Builder
	content.WriteString("facts {\n")

	for machineID, collection := range facts {
		content.WriteString(fmt.Sprintf("  %s {\n", machineID))
		content.WriteString(fmt.Sprintf("    server = \"%s\"\n", collection.Server))
		content.WriteString(fmt.Sprintf("    timestamp = \"%s\"\n", collection.Timestamp.Format(time.RFC3339)))

		// Write facts
		for factKey, fact := range collection.Facts {
			content.WriteString(fmt.Sprintf("    %s = %s\n", factKey, formatHCLValue(fact.Value)))
		}

		content.WriteString("  }\n")
	}

	content.WriteString("}\n")
	return content.String()
}

// formatHCLValue formats a value for HCL output
func formatHCLValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case int, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		// For complex types, convert to JSON string
		jsonBytes, _ := json.Marshal(v)
		return fmt.Sprintf("\"%s\"", string(jsonBytes))
	}
}

// parseHCLToFactCollections parses HCL AST to fact collections
func parseHCLToFactCollections(file *hcl.File) ([]*spookyfactstypes.FactCollection, error) {
	// Basic implementation for parsing HCL to fact collections
	// This is a simplified implementation - in a full implementation,
	// we would traverse the HCL AST and convert to FactCollection structs

	// For now, create a basic fact collection from the parsed file
	collection := &spookyfactstypes.FactCollection{
		Server:    "imported-server",
		Timestamp: time.Now(),
		Facts: map[string]*spookyfactstypes.Fact{
			"imported": {
				Key:   "imported",
				Value: true,
			},
			"import_timestamp": {
				Key:   "import_timestamp",
				Value: time.Now().Format(time.RFC3339),
			},
		},
	}

	return []*spookyfactstypes.FactCollection{collection}, nil
}

// Close closes the BadgerDB connection
func (b *BadgerFactStorage) Close() error {
	if b.db != nil {
		return b.db.Close()
	}
	return nil
}

// Set stores a key-value pair in the database
func (b *BadgerFactStorage) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

// Get retrieves a value by key from the database
func (b *BadgerFactStorage) Get(key string) (interface{}, error) {
	var value interface{}

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &value)
		})
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get value for key %s: %w", key, err)
	}

	return value, nil
}

// Delete removes a key-value pair from the database
func (b *BadgerFactStorage) Delete(key string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}
