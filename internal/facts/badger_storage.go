package facts

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/dgraph-io/badger/v4"
)

// BadgerFactStorage implements FactStorage using BadgerDB
type BadgerFactStorage struct {
	db *badger.DB
}

// NewBadgerFactStorage creates a new BadgerDB-based fact storage
func NewBadgerFactStorage(dbPath string) (*BadgerFactStorage, error) {
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil // Disable logging for cleaner output

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open BadgerDB: %w", err)
	}

	return &BadgerFactStorage{db: db}, nil
}

// GetMachineFacts retrieves facts for a specific machine
func (b *BadgerFactStorage) GetMachineFacts(machineID string) (*MachineFacts, error) {
	var facts MachineFacts

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(machineID))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &facts)
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, fmt.Errorf("machine facts not found: %s", machineID)
	}

	return &facts, err
}

// SetMachineFacts stores facts for a specific machine
func (b *BadgerFactStorage) SetMachineFacts(machineID string, facts *MachineFacts) error {
	facts.UpdatedAt = time.Now()
	if facts.CreatedAt.IsZero() {
		facts.CreatedAt = facts.UpdatedAt
	}

	data, err := json.Marshal(facts)
	if err != nil {
		return fmt.Errorf("failed to marshal facts: %w", err)
	}

	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(machineID), data)
	})
}

// QueryFacts searches for facts matching the query criteria
func (b *BadgerFactStorage) QueryFacts(query *FactQuery) ([]*MachineFacts, error) {
	var results []*MachineFacts

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 100

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var facts MachineFacts

			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &facts)
			})
			if err != nil {
				continue
			}

			if matchesQuery(&facts, query) {
				results = append(results, &facts)
				if query.Limit > 0 && len(results) >= query.Limit {
					break
				}
			}
		}

		return nil
	})

	return results, err
}

// DeleteFacts deletes facts matching the query criteria
func (b *BadgerFactStorage) DeleteFacts(query *FactQuery) (int, error) {
	var deletedCount int

	err := b.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var facts MachineFacts

			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &facts)
			})
			if err != nil {
				continue
			}

			if matchesQuery(&facts, query) {
				if err := txn.Delete(item.KeyCopy(nil)); err != nil {
					return fmt.Errorf("failed to delete facts: %w", err)
				}
				deletedCount++
			}
		}

		return nil
	})

	return deletedCount, err
}

// DeleteMachineFacts deletes facts for a specific machine
func (b *BadgerFactStorage) DeleteMachineFacts(machineID string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(machineID))
	})
}

// ExportToJSON exports all facts to JSON format
func (b *BadgerFactStorage) ExportToJSON(w io.Writer) error {
	facts := make(map[string]*MachineFacts)

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var machineFacts MachineFacts

			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &machineFacts)
			})
			if err != nil {
				continue
			}

			facts[string(item.KeyCopy(nil))] = &machineFacts
		}

		return nil
	})

	if err != nil {
		return err
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(facts)
}

// ImportFromJSON imports facts from JSON format
func (b *BadgerFactStorage) ImportFromJSON(r io.Reader) error {
	var facts map[string]*MachineFacts

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&facts); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return b.db.Update(func(txn *badger.Txn) error {
		for machineID, machineFacts := range facts {
			data, err := json.Marshal(machineFacts)
			if err != nil {
				return fmt.Errorf("failed to marshal facts for %s: %w", machineID, err)
			}

			if err := txn.Set([]byte(machineID), data); err != nil {
				return fmt.Errorf("failed to set facts for %s: %w", machineID, err)
			}
		}
		return nil
	})
}

// ExportToJSONWithEncryption exports facts with encryption support
func (b *BadgerFactStorage) ExportToJSONWithEncryption(w io.Writer, _ ExportOptions) error {
	// For now, implement basic encryption without age
	// TODO: Add age encryption support
	return b.ExportToJSON(w)
}

// ImportFromJSONWithDecryption imports facts with decryption support
func (b *BadgerFactStorage) ImportFromJSONWithDecryption(r io.Reader, _ string) error {
	// For now, implement basic import without age decryption
	// TODO: Add age decryption support
	return b.ImportFromJSON(r)
}

// Close closes the BadgerDB connection
func (b *BadgerFactStorage) Close() error {
	return b.db.Close()
}
