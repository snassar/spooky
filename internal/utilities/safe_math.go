package utilities

import (
	"errors"
	"math"
	"math/bits"
)

var (
	// ErrOverflow indicates an arithmetic overflow error.
	ErrOverflow = errors.New("integer overflow detected")
	// ErrOutOfBounds indicates a value is out of acceptable bounds.
	ErrOutOfBounds = errors.New("value out of bounds")
)

// SafeInt64 converts int to int64 with bounds checking
func SafeInt64(n int) (int64, error) {
	if n < math.MinInt32 || n > math.MaxInt32 {
		return 0, ErrOutOfBounds
	}
	return int64(n), nil
}

// SafeUint64 converts uint to uint64 with bounds checking
func SafeUint64(n uint) (uint64, error) {
	if n > math.MaxUint32 {
		return 0, ErrOutOfBounds
	}
	return uint64(n), nil
}

// SafeInt converts int64 to int with bounds checking
func SafeInt(n int64) (int, error) {
	if n < math.MinInt32 || n > math.MaxInt32 {
		return 0, ErrOutOfBounds
	}
	return int(n), nil
}

// SafeUint converts uint64 to uint with bounds checking
func SafeUint(n uint64) (uint, error) {
	if n > math.MaxUint32 {
		return 0, ErrOutOfBounds
	}
	return uint(n), nil
}

// SafeMultiplyInt64 multiplies two int64 values with overflow checking
func SafeMultiplyInt64(a, b int64) (int64, error) {
	// Handle zero case
	if a == 0 || b == 0 {
		return 0, nil
	}

	// For negative numbers, we need to handle them differently
	if a > 0 && b > 0 {
		// Both positive - check if a * b > MaxInt64
		if a > math.MaxInt64/b {
			return 0, ErrOverflow
		}
	} else if a < 0 && b < 0 {
		// Both negative - result will be positive
		if a == math.MinInt64 || b == math.MinInt64 {
			// Special case: MinInt64 * anything will overflow
			return 0, ErrOverflow
		}
		if -a > math.MaxInt64/-b {
			return 0, ErrOverflow
		}
	} else {
		// One positive, one negative - result will be negative
		// Check if |a| * |b| > MaxInt64 + 1 (since result can be MinInt64)
		if a < 0 {
			a, b = b, a // Swap so a is positive
		}
		if a > 0 && b < 0 {
			if a > math.MaxInt64/-b {
				return 0, ErrOverflow
			}
		}
	}

	result := a * b
	return result, nil
}

// SafeMultiplyUint64 multiplies two uint64 values with overflow checking
func SafeMultiplyUint64(a, b uint64) (uint64, error) {
	if a == 0 || b == 0 {
		return 0, nil
	}

	// Check for overflow using math/bits
	if bits.LeadingZeros64(a)+bits.LeadingZeros64(b) < 64 {
		return 0, ErrOverflow
	}

	result := a * b
	if a != 0 && result/a != b {
		return 0, ErrOverflow
	}

	return result, nil
}

// SafeAddInt64 adds two int64 values with overflow checking
func SafeAddInt64(a, b int64) (int64, error) {
	// Use the standard overflow checking for signed integers
	if (b > 0 && a > math.MaxInt64-b) || (b < 0 && a < math.MinInt64-b) {
		return 0, ErrOverflow
	}
	return a + b, nil
}

// SafeAddUint64 adds two uint64 values with overflow checking
func SafeAddUint64(a, b uint64) (uint64, error) {
	if a > math.MaxUint64-b {
		return 0, ErrOverflow
	}
	return a + b, nil
}

// SafeSubtractInt64 subtracts two int64 values with overflow checking
func SafeSubtractInt64(a, b int64) (int64, error) {
	// Use the standard overflow checking for signed integers
	if (b > 0 && a < math.MinInt64+b) || (b < 0 && a > math.MaxInt64+b) {
		return 0, ErrOverflow
	}
	return a - b, nil
}

// SafeSubtractUint64 subtracts two uint64 values with overflow checking
func SafeSubtractUint64(a, b uint64) (uint64, error) {
	if a < b {
		return 0, ErrOutOfBounds
	}
	return a - b, nil
}

// SafeShiftLeftUint64 performs left shift with overflow checking
func SafeShiftLeftUint64(n uint64, shift uint) (uint64, error) {
	if shift >= 64 {
		return 0, ErrOutOfBounds
	}

	result := n << shift
	if result < n {
		return 0, ErrOverflow
	}

	return result, nil
}

// SafeShiftRightUint64 performs right shift (always safe)
func SafeShiftRightUint64(n uint64, shift uint) uint64 {
	if shift >= 64 {
		return 0
	}
	return n >> shift
}

// SafeBitwiseAndUint64 performs bitwise AND (always safe)
func SafeBitwiseAndUint64(a, b uint64) uint64 {
	return a & b
}

// SafeBitwiseOrUint64 performs bitwise OR (always safe)
func SafeBitwiseOrUint64(a, b uint64) uint64 {
	return a | b
}

// SafeBitwiseXorUint64 performs bitwise XOR (always safe)
func SafeBitwiseXorUint64(a, b uint64) uint64 {
	return a ^ b
}

// SafeBitwiseNotUint64 performs bitwise NOT (always safe)
func SafeBitwiseNotUint64(n uint64) uint64 {
	return ^n
}

// SafeExtractBits safely extracts bits from a uint64 value
func SafeExtractBits(n uint64, start, length uint) (uint64, error) {
	if start+length > 64 {
		return 0, ErrOutOfBounds
	}

	mask := uint64(1)<<length - 1
	return (n >> start) & mask, nil
}

// SafeSetBits safely sets bits in a uint64 value
func SafeSetBits(n uint64, start, length uint, value uint64) (uint64, error) {
	if start+length > 64 {
		return 0, ErrOutOfBounds
	}

	if value >= (uint64(1) << length) {
		return 0, ErrOutOfBounds
	}

	mask := uint64(1)<<length - 1
	cleared := n &^ (mask << start)
	set := (value & mask) << start
	return cleared | set, nil
}
