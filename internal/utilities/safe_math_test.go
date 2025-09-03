package utilities

import (
	"math"
	"testing"
)

func TestSafeInt64(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		want    int64
		wantErr bool
	}{
		{"normal case", 42, 42, false},
		{"zero", 0, 0, false},
		{"negative", -42, -42, false},
		{"max int32", math.MaxInt32, math.MaxInt32, false},
		{"min int32", math.MinInt32, math.MinInt32, false},
		{"overflow positive", math.MaxInt32 + 1, 0, true},
		{"overflow negative", math.MinInt32 - 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeInt64(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SafeInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeUint64(t *testing.T) {
	tests := []struct {
		name    string
		input   uint
		want    uint64
		wantErr bool
	}{
		{"normal case", 42, 42, false},
		{"zero", 0, 0, false},
		{"max uint32", math.MaxUint32, math.MaxUint32, false},
		{"overflow", math.MaxUint32 + 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeUint64(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeUint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SafeUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeInt(t *testing.T) {
	tests := []struct {
		name    string
		input   int64
		want    int
		wantErr bool
	}{
		{"normal case", 42, 42, false},
		{"zero", 0, 0, false},
		{"negative", -42, -42, false},
		{"max int32", math.MaxInt32, math.MaxInt32, false},
		{"min int32", math.MinInt32, math.MinInt32, false},
		{"overflow positive", math.MaxInt32 + 1, 0, true},
		{"overflow negative", math.MinInt32 - 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SafeInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeUint(t *testing.T) {
	tests := []struct {
		name    string
		input   uint64
		want    uint
		wantErr bool
	}{
		{"normal case", 42, 42, false},
		{"zero", 0, 0, false},
		{"max uint32", math.MaxUint32, math.MaxUint32, false},
		{"overflow", math.MaxUint32 + 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeUint(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeUint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SafeUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeMultiplyInt64(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int64
		want    int64
		wantErr bool
	}{
		{"normal case", 2, 3, 6, false},
		{"zero", 0, 5, 0, false},
		{"negative", -2, 3, -6, false},
		{"large numbers", 1000000, 1000000, 1000000000000, false},
		{"overflow", math.MaxInt64, 2, 0, true},
		{"overflow negative", math.MinInt64, 2, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeMultiplyInt64(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeMultiplyInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SafeMultiplyInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeMultiplyUint64(t *testing.T) {
	tests := []struct {
		name    string
		a, b    uint64
		want    uint64
		wantErr bool
	}{
		{"normal case", 2, 3, 6, false},
		{"zero", 0, 5, 0, false},
		{"large numbers", 1000000, 1000000, 1000000000000, false},
		{"overflow", math.MaxUint64, 2, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeMultiplyUint64(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeMultiplyUint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SafeMultiplyUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeAddInt64(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int64
		want    int64
		wantErr bool
	}{
		{"normal case", 2, 3, 5, false},
		{"zero", 0, 5, 5, false},
		{"negative", -2, 3, 1, false},
		{"overflow positive", math.MaxInt64, 1, 0, true},
		{"overflow negative", math.MinInt64, -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeAddInt64(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeAddInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SafeAddInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeAddUint64(t *testing.T) {
	tests := []struct {
		name    string
		a, b    uint64
		want    uint64
		wantErr bool
	}{
		{"normal case", 2, 3, 5, false},
		{"zero", 0, 5, 5, false},
		{"overflow", math.MaxUint64, 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeAddUint64(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeAddUint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SafeAddUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeSubtractInt64(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int64
		want    int64
		wantErr bool
	}{
		{"normal case", 5, 3, 2, false},
		{"zero", 5, 0, 5, false},
		{"negative result", 2, 5, -3, false},
		{"overflow positive", math.MaxInt64, -1, 0, true},
		{"overflow negative", math.MinInt64, 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeSubtractInt64(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeSubtractInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SafeSubtractInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeSubtractUint64(t *testing.T) {
	tests := []struct {
		name    string
		a, b    uint64
		want    uint64
		wantErr bool
	}{
		{"normal case", 5, 3, 2, false},
		{"zero", 5, 0, 5, false},
		{"underflow", 2, 5, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeSubtractUint64(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeSubtractUint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SafeSubtractUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeShiftLeftUint64(t *testing.T) {
	tests := []struct {
		name    string
		n       uint64
		shift   uint
		want    uint64
		wantErr bool
	}{
		{"normal case", 1, 2, 4, false},
		{"zero", 0, 5, 0, false},
		{"large shift", 1, 63, 1 << 63, false},
		{"overflow shift", 1, 64, 0, true},
		{"overflow value", 1 << 63, 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeShiftLeftUint64(tt.n, tt.shift)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeShiftLeftUint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SafeShiftLeftUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeShiftRightUint64(t *testing.T) {
	tests := []struct {
		name  string
		n     uint64
		shift uint
		want  uint64
	}{
		{"normal case", 8, 2, 2},
		{"zero", 0, 5, 0},
		{"large shift", 1 << 63, 63, 1},
		{"excessive shift", 1, 64, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SafeShiftRightUint64(tt.n, tt.shift)
			if got != tt.want {
				t.Errorf("SafeShiftRightUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeBitwiseOperations(t *testing.T) {
	a := uint64(0b1010)
	b := uint64(0b1100)

	// Test AND
	if got := SafeBitwiseAndUint64(a, b); got != 0b1000 {
		t.Errorf("SafeBitwiseAndUint64() = %b, want %b", got, 0b1000)
	}

	// Test OR
	if got := SafeBitwiseOrUint64(a, b); got != 0b1110 {
		t.Errorf("SafeBitwiseOrUint64() = %b, want %b", got, 0b1110)
	}

	// Test XOR
	if got := SafeBitwiseXorUint64(a, b); got != 0b0110 {
		t.Errorf("SafeBitwiseXorUint64() = %b, want %b", got, 0b0110)
	}

	// Test NOT
	if got := SafeBitwiseNotUint64(a); got != ^a {
		t.Errorf("SafeBitwiseNotUint64() = %b, want %b", got, ^a)
	}
}

func TestSafeExtractBits(t *testing.T) {
	tests := []struct {
		name          string
		n             uint64
		start, length uint
		want          uint64
		wantErr       bool
	}{
		{"normal case", 0b11110000, 2, 4, 0b1100, false},
		{"zero length", 0b11110000, 2, 0, 0, false},
		{"full length", 0b11110000, 0, 8, 0b11110000, false},
		{"out of bounds", 0b11110000, 60, 8, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeExtractBits(tt.n, tt.start, tt.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeExtractBits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SafeExtractBits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeSetBits(t *testing.T) {
	tests := []struct {
		name          string
		n             uint64
		start, length uint
		value         uint64
		want          uint64
		wantErr       bool
	}{
		{"normal case", 0b11110000, 2, 4, 0b1010, 0b11101000, false},
		{"zero length", 0b11110000, 2, 0, 0, 0b11110000, false},
		{"out of bounds", 0b11110000, 60, 8, 0, 0, true},
		{"value too large", 0b11110000, 2, 2, 0b1010, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeSetBits(tt.n, tt.start, tt.length, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeSetBits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SafeSetBits() = %v, want %v", got, tt.want)
			}
		})
	}
}
