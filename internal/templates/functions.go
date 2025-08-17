// Package templates provides template function registry and security functionality.
package templates

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	spookytypeslogging "spooky/internal/types/logging"
)

// BuiltInFunctions provides a comprehensive set of secure template functions
type BuiltInFunctions struct {
	logger spookytypeslogging.Logger
}

// NewBuiltInFunctions creates a new built-in functions registry
func NewBuiltInFunctions(logger spookytypeslogging.Logger) *BuiltInFunctions {
	return &BuiltInFunctions{
		logger: logger,
	}
}

// GetFunctions returns all built-in template functions
func (b *BuiltInFunctions) GetFunctions() map[string]interface{} {
	return map[string]interface{}{
		// String manipulation functions
		"upper":      strings.ToUpper,
		"lower":      strings.ToLower,
		"title":      b.title,
		"trim":       strings.TrimSpace,
		"trimLeft":   strings.TrimLeft,
		"trimRight":  strings.TrimRight,
		"replace":    strings.Replace,
		"replaceAll": strings.ReplaceAll,
		"split":      strings.Split,
		"join":       strings.Join,
		"contains":   strings.Contains,
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
		"repeat":     strings.Repeat,
		"substr":     b.substring,
		"len":        b.length,

		// Mathematical functions
		"add":   b.add,
		"sub":   b.sub,
		"mul":   b.mul,
		"div":   b.div,
		"mod":   b.mod,
		"abs":   math.Abs,
		"ceil":  math.Ceil,
		"floor": math.Floor,
		"round": math.Round,
		"min":   b.min,
		"max":   b.max,
		"pow":   math.Pow,
		"sqrt":  math.Sqrt,

		// Array manipulation functions
		"first":        b.first,
		"last":         b.last,
		"index":        b.index,
		"slice":        b.slice,
		"append":       b.append,
		"prepend":      b.prepend,
		"reverse":      b.reverse,
		"sort":         b.sort,
		"uniq":         b.uniq,
		"containsItem": b.containsItem,

		// Hash and encoding functions
		"sha256":       b.sha256,
		"base64":       b.base64,
		"base64Decode": b.base64Decode,
		"hex":          b.hex,
		"hexDecode":    b.hexDecode,

		// Type conversion functions
		"toString": b.toString,
		"toInt":    b.toInt,
		"toFloat":  b.toFloat,
		"toBool":   b.toBool,

		// JSON functions
		"toJSON":     b.toJSON,
		"fromJSON":   b.fromJSON,
		"prettyJSON": b.prettyJSON,

		// Regular expression functions
		"regexMatch":   b.regexMatch,
		"regexReplace": b.regexReplace,

		// Random functions
		"random": b.random,
		"uuid":   b.uuid,
	}
}

// String manipulation functions

func (b *BuiltInFunctions) title(s string) string {
	caser := cases.Title(language.English)
	return caser.String(s)
}

func (b *BuiltInFunctions) substring(s string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	if start >= end {
		return ""
	}
	return s[start:end]
}

func (b *BuiltInFunctions) length(v interface{}) int {
	switch val := v.(type) {
	case string:
		return len(val)
	case []interface{}:
		return len(val)
	case map[string]interface{}:
		return len(val)
	default:
		return len(b.toString(val))
	}
}

// Mathematical functions

func (b *BuiltInFunctions) add(values ...interface{}) float64 {
	result := 0.0
	for _, v := range values {
		result += b.toFloat(v)
	}
	return result
}

func (b *BuiltInFunctions) sub(a, bVal interface{}) float64 {
	return b.toFloat(a) - b.toFloat(bVal)
}

func (b *BuiltInFunctions) mul(values ...interface{}) float64 {
	result := 1.0
	for _, v := range values {
		result *= b.toFloat(v)
	}
	return result
}

func (b *BuiltInFunctions) div(a, bVal interface{}) float64 {
	divisor := b.toFloat(bVal)
	if divisor == 0 {
		b.logger.Warn("Division by zero", map[string]interface{}{
			"dividend": a,
			"divisor":  divisor,
		})
		return 0
	}
	return b.toFloat(a) / divisor
}

func (b *BuiltInFunctions) mod(a, bVal interface{}) int {
	divisor := b.toInt(bVal)
	if divisor == 0 {
		b.logger.Warn("Modulo by zero", map[string]interface{}{
			"dividend": a,
			"divisor":  divisor,
		})
		return 0
	}
	return b.toInt(a) % divisor
}

func (b *BuiltInFunctions) min(values ...interface{}) float64 {
	if len(values) == 0 {
		return 0
	}
	minVal := b.toFloat(values[0])
	for _, v := range values[1:] {
		if val := b.toFloat(v); val < minVal {
			minVal = val
		}
	}
	return minVal
}

func (b *BuiltInFunctions) max(values ...interface{}) float64 {
	if len(values) == 0 {
		return 0
	}
	maxVal := b.toFloat(values[0])
	for _, v := range values[1:] {
		if val := b.toFloat(v); val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}

// Array manipulation functions

func (b *BuiltInFunctions) first(slice []interface{}) interface{} {
	if len(slice) == 0 {
		return nil
	}
	return slice[0]
}

func (b *BuiltInFunctions) last(slice []interface{}) interface{} {
	if len(slice) == 0 {
		return nil
	}
	return slice[len(slice)-1]
}

func (b *BuiltInFunctions) index(slice []interface{}, i int) interface{} {
	if i < 0 || i >= len(slice) {
		return nil
	}
	return slice[i]
}

func (b *BuiltInFunctions) slice(slice []interface{}, start, end int) []interface{} {
	if start < 0 {
		start = 0
	}
	if end > len(slice) {
		end = len(slice)
	}
	if start >= end {
		return []interface{}{}
	}
	return slice[start:end]
}

func (b *BuiltInFunctions) append(slice []interface{}, items ...interface{}) []interface{} {
	return append(slice, items...)
}

func (b *BuiltInFunctions) prepend(slice []interface{}, items ...interface{}) []interface{} {
	return append(items, slice...)
}

func (b *BuiltInFunctions) reverse(slice []interface{}) []interface{} {
	result := make([]interface{}, len(slice))
	for i, j := 0, len(slice)-1; i < len(slice); i, j = i+1, j-1 {
		result[i] = slice[j]
	}
	return result
}

func (b *BuiltInFunctions) sort(slice []interface{}) []interface{} {
	result := make([]interface{}, len(slice))
	copy(result, slice)
	sort.Slice(result, func(i, j int) bool {
		return b.toString(result[i]) < b.toString(result[j])
	})
	return result
}

func (b *BuiltInFunctions) uniq(slice []interface{}) []interface{} {
	seen := make(map[string]bool)
	result := []interface{}{}
	for _, item := range slice {
		key := b.toString(item)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}
	return result
}

func (b *BuiltInFunctions) containsItem(slice []interface{}, item interface{}) bool {
	for _, v := range slice {
		if b.toString(v) == b.toString(item) {
			return true
		}
	}
	return false
}

// Hash and encoding functions

func (b *BuiltInFunctions) sha256(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

func (b *BuiltInFunctions) base64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func (b *BuiltInFunctions) base64Decode(s string) string {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		b.logger.Warn("Failed to decode base64", map[string]interface{}{
			"input": s,
			"error": err.Error(),
		})
		return ""
	}
	return string(decoded)
}

func (b *BuiltInFunctions) hex(s string) string {
	return hex.EncodeToString([]byte(s))
}

func (b *BuiltInFunctions) hexDecode(s string) string {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		b.logger.Warn("Failed to decode hex", map[string]interface{}{
			"input": s,
			"error": err.Error(),
		})
		return ""
	}
	return string(decoded)
}

// Type conversion functions

func (b *BuiltInFunctions) toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func (b *BuiltInFunctions) toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return int(f)
		}
	}
	return 0
}

func (b *BuiltInFunctions) toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return 0
}

func (b *BuiltInFunctions) toBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return strings.EqualFold(val, "true") || val == "1"
	case int:
		return val != 0
	case float64:
		return val != 0
	}
	return false
}

// JSON functions

func (b *BuiltInFunctions) toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		b.logger.Warn("Failed to marshal JSON", map[string]interface{}{
			"value": v,
			"error": err.Error(),
		})
		return ""
	}
	return string(data)
}

func (b *BuiltInFunctions) fromJSON(s string) interface{} {
	var result interface{}
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		b.logger.Warn("Failed to unmarshal JSON", map[string]interface{}{
			"input": s,
			"error": err.Error(),
		})
		return nil
	}
	return result
}

func (b *BuiltInFunctions) prettyJSON(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		b.logger.Warn("Failed to marshal pretty JSON", map[string]interface{}{
			"value": v,
			"error": err.Error(),
		})
		return ""
	}
	return string(data)
}

// Regular expression functions

func (b *BuiltInFunctions) regexMatch(pattern, s string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		b.logger.Warn("Invalid regex pattern", map[string]interface{}{
			"pattern": pattern,
			"error":   err.Error(),
		})
		return false
	}
	return re.MatchString(s)
}

func (b *BuiltInFunctions) regexReplace(pattern, s, replacement string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		b.logger.Warn("Invalid regex pattern", map[string]interface{}{
			"pattern": pattern,
			"error":   err.Error(),
		})
		return s
	}
	return re.ReplaceAllString(s, replacement)
}

func (b *BuiltInFunctions) random(minVal, maxVal int) int {
	if minVal >= maxVal {
		return minVal
	}

	// Use crypto/rand for secure random number generation
	rangeSize := maxVal - minVal + 1
	if rangeSize <= 0 {
		return minVal
	}

	// Use a simpler approach to avoid integer overflow
	buf := make([]byte, 4)
	_, err := rand.Read(buf)
	if err != nil {
		b.logger.Warn("Failed to generate random number", map[string]interface{}{
			"error": err.Error(),
		})
		return minVal
	}

	// Convert to uint32 and scale to range
	randomUint := uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])

	// Use modulo to get a value in range [0, rangeSize)
	// Ensure rangeSize is positive and within uint32 bounds
	if rangeSize > 0 && rangeSize <= 1<<31-1 {
		scaledValue := randomUint % uint32(rangeSize)
		return minVal + int(scaledValue)
	}

	// Fallback for very large ranges
	fallbackRange := maxVal - minVal + 1
	if fallbackRange > 0 && fallbackRange <= 1<<31-1 {
		return minVal + int(randomUint%uint32(fallbackRange))
	}
	// If range is still too large, just return minVal
	return minVal
}

func (b *BuiltInFunctions) uuid() string {
	// Secure UUID v4 implementation using crypto/rand
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		b.logger.Warn("Failed to generate UUID", map[string]interface{}{
			"error": err.Error(),
		})
		return ""
	}

	// Set version (4) and variant bits
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}

// FunctionSecurityManager provides security restrictions for template functions
type FunctionSecurityManager struct {
	restrictedMode     bool
	allowedFunctions   map[string]bool
	forbiddenFunctions map[string]bool
	mu                 sync.RWMutex
}

// NewFunctionSecurityManager creates a new function security manager
func NewFunctionSecurityManager() *FunctionSecurityManager {
	return &FunctionSecurityManager{
		restrictedMode:     false,
		allowedFunctions:   make(map[string]bool),
		forbiddenFunctions: make(map[string]bool),
	}
}

// SetRestrictedMode sets the restricted mode
func (f *FunctionSecurityManager) SetRestrictedMode(restricted bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.restrictedMode = restricted
}

// AllowFunction allows a specific function
func (f *FunctionSecurityManager) AllowFunction(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.allowedFunctions[name] = true
	delete(f.forbiddenFunctions, name)
}

// ForbidFunction forbids a specific function
func (f *FunctionSecurityManager) ForbidFunction(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.forbiddenFunctions[name] = true
	delete(f.allowedFunctions, name)
}

// IsFunctionAllowed checks if a function is allowed
func (f *FunctionSecurityManager) IsFunctionAllowed(name string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Check if function is explicitly forbidden
	if f.forbiddenFunctions[name] {
		return false
	}

	// In restricted mode, only explicitly allowed functions are permitted
	if f.restrictedMode {
		return f.allowedFunctions[name]
	}

	// In non-restricted mode, all functions are allowed unless explicitly forbidden
	return true
}

// GetAllowedFunctions returns all allowed functions
func (f *FunctionSecurityManager) GetAllowedFunctions() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var functions []string
	for name := range f.allowedFunctions {
		functions = append(functions, name)
	}
	sort.Strings(functions)
	return functions
}

// GetForbiddenFunctions returns all forbidden functions
func (f *FunctionSecurityManager) GetForbiddenFunctions() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var functions []string
	for name := range f.forbiddenFunctions {
		functions = append(functions, name)
	}
	sort.Strings(functions)
	return functions
}
