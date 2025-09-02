package hcl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// SharedHCLUtils provides common HCL generation utilities used by multiple generators
type SharedHCLUtils struct{}

// NewSharedHCLUtils creates a new shared HCL utilities instance
func NewSharedHCLUtils() *SharedHCLUtils {
	return &SharedHCLUtils{}
}

// CtyValueToHCL converts a cty.Value to HCL and writes it to the body
// This is a shared implementation used by both ConfigGenerator and HCLGenerator
func (su *SharedHCLUtils) CtyValueToHCL(value cty.Value, body *hclwrite.Body, blockName string) error {
	if !value.IsKnown() {
		return fmt.Errorf("cannot convert unknown value to HCL")
	}

	if value.IsNull() {
		return nil
	}

	valueType := value.Type()
	if valueType.IsObjectType() {
		// Create a block for objects
		block := body.AppendNewBlock(blockName, nil)
		blockBody := block.Body()

		// Add each field to the block
		for key, val := range value.AsValueMap() {
			if !val.IsNull() {
				err := su.CtyValueToHCL(val, blockBody, key)
				if err != nil {
					return fmt.Errorf("failed to convert field %s: %v", key, err)
				}
			}
		}

	} else if valueType.IsListType() {
		// Handle lists by creating blocks for each element
		values := value.AsValueSlice()
		for _, val := range values {
			if !val.IsNull() {
				err := su.CtyValueToHCL(val, body, blockName)
				if err != nil {
					return fmt.Errorf("failed to convert list element: %v", err)
				}
			}
		}

	} else if valueType.IsMapType() {
		// Handle maps by creating blocks for each key-value pair
		valueMap := value.AsValueMap()
		for key, val := range valueMap {
			if !val.IsNull() {
				// For maps, we create a block with the key as the label
				block := body.AppendNewBlock(blockName, []string{key})
				blockBody := block.Body()

				err := su.CtyValueToHCL(val, blockBody, "")
				if err != nil {
					return fmt.Errorf("failed to convert map value for key %s: %v", key, err)
				}
			}
		}

	} else {
		// For primitive types, set as attribute
		body.SetAttributeValue(blockName, value)
	}

	return nil
}
