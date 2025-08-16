package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypeslogging "spooky/internal/types/logging"
)

// HCLProcessor handles encryption and decryption of HCL values
type HCLProcessor struct {
	logger spookytypeslogging.Logger
}

// NewHCLProcessor creates a new HCL processor
func NewHCLProcessor(logger spookytypeslogging.Logger) *HCLProcessor {
	return &HCLProcessor{
		logger: logger,
	}
}

// EncryptHCLValues encrypts string values and entire objects in HCL-compatible structures
func (d *HCLProcessor) EncryptHCLValues(ctx context.Context, data interface{}, secretsIntegration spookyinterfaces.SecretsIntegration, recipients []string, dryRun bool) error {
	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}

	if secretsIntegration == nil {
		return fmt.Errorf("secrets integration cannot be nil")
	}

	if len(recipients) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	d.logger.Info("Starting HCL values encryption", map[string]interface{}{
		"recipients": recipients,
		"data_type":  fmt.Sprintf("%T", data),
		"dry_run":    dryRun,
	})

	var encryptedCount int
	var errors []string

	// Process the data structure
	if err := d.encryptValue(ctx, data, secretsIntegration, recipients, dryRun, &encryptedCount, &errors); err != nil {
		return fmt.Errorf("failed to encrypt HCL values: %w", err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("HCL encryption completed with errors: %s", strings.Join(errors, "; "))
	}

	d.logger.Info("HCL values encryption completed", map[string]interface{}{
		"encrypted_count": encryptedCount,
		"dry_run":         dryRun,
	})

	return nil
}

// DecryptHCLValues decrypts string values and entire objects in HCL-compatible structures
func (d *HCLProcessor) DecryptHCLValues(ctx context.Context, data interface{}, secretsIntegration spookyinterfaces.SecretsIntegration, identityPath string) error {
	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}

	if secretsIntegration == nil {
		return fmt.Errorf("secrets integration cannot be nil")
	}

	if identityPath == "" {
		return fmt.Errorf("identity path cannot be empty")
	}

	d.logger.Info("Starting HCL values decryption", map[string]interface{}{
		"identity_path": identityPath,
		"data_type":     fmt.Sprintf("%T", data),
	})

	var decryptedCount int
	var errors []string

	// Process the data structure
	if err := d.decryptValue(ctx, data, secretsIntegration, identityPath, &decryptedCount, &errors); err != nil {
		return fmt.Errorf("failed to decrypt HCL values: %w", err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("HCL decryption completed with errors: %s", strings.Join(errors, "; "))
	}

	d.logger.Info("HCL values decryption completed", map[string]interface{}{
		"decrypted_count": decryptedCount,
	})

	return nil
}

// encryptValue processes values for encryption - encrypts strings and entire objects
func (d *HCLProcessor) encryptValue(ctx context.Context, value interface{}, secretsIntegration spookyinterfaces.SecretsIntegration, recipients []string, dryRun bool, encryptedCount *int, errors *[]string) error {
	if value == nil {
		return nil
	}

	// Use reflection to handle different types
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		// Handle string values - check if they need encryption
		if strValue := v.String(); strValue != "" && !strings.HasPrefix(strValue, "age1") {
			if dryRun {
				d.logger.Info("Would encrypt string value", map[string]interface{}{
					"value_length": len(strValue),
				})
				*encryptedCount++
			} else {
				encryptedBytes, err := secretsIntegration.EncryptWithAge(ctx, []byte(strValue), recipients)
				if err != nil {
					errorMsg := fmt.Sprintf("failed to encrypt string value: %v", err)
					*errors = append(*errors, errorMsg)
					d.logger.Error("Failed to encrypt string value", err, map[string]interface{}{
						"value_length": len(strValue),
					})
					return nil // Continue processing other values
				}

				// Update the string value
				if reflect.ValueOf(value).Kind() == reflect.Ptr {
					reflect.ValueOf(value).Elem().SetString(string(encryptedBytes))
				} else {
					d.logger.Warn("Cannot modify non-pointer string value", map[string]interface{}{
						"value_type": fmt.Sprintf("%T", value),
					})
				}

				*encryptedCount++
				d.logger.Debug("Encrypted string value", map[string]interface{}{
					"value_length": len(strValue),
				})
			}
		}

	case reflect.Map, reflect.Struct, reflect.Slice, reflect.Array:
		// Encrypt entire objects as JSON
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			errorMsg := fmt.Sprintf("failed to serialize object for encryption: %v", err)
			*errors = append(*errors, errorMsg)
			d.logger.Error("Failed to serialize object for encryption", err, map[string]interface{}{
				"value_type": fmt.Sprintf("%T", value),
			})
			return nil // Continue processing other values
		}

		if dryRun {
			d.logger.Info("Would encrypt entire object", map[string]interface{}{
				"value_type":  fmt.Sprintf("%T", value),
				"json_length": len(jsonBytes),
			})
			*encryptedCount++
		} else {
			// Encrypt the entire object
			encryptedBytes, err := secretsIntegration.EncryptWithAge(ctx, jsonBytes, recipients)
			if err != nil {
				errorMsg := fmt.Sprintf("failed to encrypt object: %v", err)
				*errors = append(*errors, errorMsg)
				d.logger.Error("Failed to encrypt object", err, map[string]interface{}{
					"value_type":  fmt.Sprintf("%T", value),
					"json_length": len(jsonBytes),
				})
				return nil // Continue processing other values
			}

			// Replace the entire object with the encrypted string
			if reflect.ValueOf(value).Kind() == reflect.Ptr {
				reflect.ValueOf(value).Elem().SetString(string(encryptedBytes))
			} else {
				d.logger.Warn("Cannot modify non-pointer object value", map[string]interface{}{
					"value_type": fmt.Sprintf("%T", value),
				})
			}

			*encryptedCount++
			d.logger.Debug("Encrypted entire object", map[string]interface{}{
				"value_type":  fmt.Sprintf("%T", value),
				"json_length": len(jsonBytes),
			})
		}

	case reflect.Interface:
		// Handle interface values - get the underlying value and process it
		if v.IsNil() {
			return nil
		}
		underlyingValue := v.Elem()
		if underlyingValue.IsValid() {
			if err := d.encryptValue(ctx, underlyingValue.Interface(), secretsIntegration, recipients, dryRun, encryptedCount, errors); err != nil {
				return fmt.Errorf("failed to encrypt interface value: %w", err)
			}
		}

	default:
		// For other types (int, float, bool, etc.), no encryption needed
		d.logger.Debug("Skipping encryption for non-string/object type", map[string]interface{}{
			"type": v.Kind().String(),
		})
	}

	return nil
}

// decryptValue processes values for decryption - decrypts strings and entire objects
func (d *HCLProcessor) decryptValue(ctx context.Context, value interface{}, secretsIntegration spookyinterfaces.SecretsIntegration, identityPath string, decryptedCount *int, errors *[]string) error {
	if value == nil {
		return nil
	}

	// Use reflection to handle different types
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		// Handle string values - check if they need decryption
		if strValue := v.String(); strValue != "" && strings.HasPrefix(strValue, "age1") {
			// Try to decrypt as a simple string first
			decryptedBytes, err := secretsIntegration.DecryptWithAge(ctx, []byte(strValue), identityPath)
			if err != nil {
				// If simple decryption fails, try as JSON object
				var jsonData interface{}
				if jsonErr := json.Unmarshal(decryptedBytes, &jsonData); jsonErr != nil {
					errorMsg := fmt.Sprintf("failed to decrypt string value: %v", err)
					*errors = append(*errors, errorMsg)
					d.logger.Error("Failed to decrypt string value", err, map[string]interface{}{
						"value_length": len(strValue),
					})
					return nil // Continue processing other values
				}

				// Update the string value with the decrypted JSON
				if reflect.ValueOf(value).Kind() == reflect.Ptr {
					reflect.ValueOf(value).Elem().Set(reflect.ValueOf(jsonData))
				} else {
					d.logger.Warn("Cannot modify non-pointer string value", map[string]interface{}{
						"value_type": fmt.Sprintf("%T", value),
					})
				}

				*decryptedCount++
				d.logger.Debug("Decrypted string value to object", map[string]interface{}{
					"value_length": len(strValue),
				})
			} else {
				// Simple string decryption succeeded
				if reflect.ValueOf(value).Kind() == reflect.Ptr {
					reflect.ValueOf(value).Elem().SetString(string(decryptedBytes))
				} else {
					d.logger.Warn("Cannot modify non-pointer string value", map[string]interface{}{
						"value_type": fmt.Sprintf("%T", value),
					})
				}

				*decryptedCount++
				d.logger.Debug("Decrypted string value", map[string]interface{}{
					"value_length": len(strValue),
				})
			}
		}

	case reflect.Interface:
		// Handle interface values - get the underlying value and process it
		if v.IsNil() {
			return nil
		}
		underlyingValue := v.Elem()
		if underlyingValue.IsValid() {
			if err := d.decryptValue(ctx, underlyingValue.Interface(), secretsIntegration, identityPath, decryptedCount, errors); err != nil {
				return fmt.Errorf("failed to decrypt interface value: %w", err)
			}
		}

	default:
		// For other types (int, float, bool, etc.), no decryption needed
		d.logger.Debug("Skipping decryption for non-string type", map[string]interface{}{
			"type": v.Kind().String(),
		})
	}

	return nil
}
