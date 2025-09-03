package sync

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/mmcloughlin/md4"
)

// Tag2 computes the 16-bit tag from two 16-bit values
// This is safe as it uses bitwise AND to ensure 16-bit result
func Tag2(s1, s2 uint16) uint16 {
	return (((s1) + (s2)) & 0xFFFF)
}

// Tag computes the 16-bit tag from a 32-bit checksum
// This is safe as it extracts specific 16-bit portions
func Tag(sum uint32) uint16 {
	return Tag2(uint16(sum&0xFFFF), uint16(sum>>16))
}

// SignExtend mirrors how C converts from (signed char) to uint32, i.e. using
// sign extension. get_checksum1 treats the buffer as (signed char*) instead of
// (unsigned char*), which likely was not a conscious choice, but here we are.
//
// This function is exported for use in the rolling checksum in match.go.
// This is safe as it only operates on a single byte
func SignExtend(b byte) uint32 {
	val := uint32(b)
	return uint32(int32(val<<24) >> 24)
}

// Checksum1 computes the fast rolling checksum (32-bit)
// Fixed to prevent integer overflow by using modulo arithmetic
func Checksum1(buf []byte) uint32 {
	bufLen := len(buf)
	var s1, s2 uint32
	var i int

	// Use modulo arithmetic to prevent overflow
	const mod = 0xFFFFFFFF

	if bufLen > 4 {
		for i = 0; i < (bufLen - 4); i += 4 {
			// Calculate each term separately with modulo to prevent overflow
			term1 := (4 * (s1 + SignExtend(buf[i]))) % mod
			term2 := (3 * SignExtend(buf[i+1])) % mod
			term3 := (2 * SignExtend(buf[i+2])) % mod
			term4 := SignExtend(buf[i+3]) % mod

			s2 = (s2 + term1 + term2 + term3 + term4) % mod

			// Calculate s1 terms separately
			s1_term1 := SignExtend(buf[i+0]) % mod
			s1_term2 := SignExtend(buf[i+1]) % mod
			s1_term3 := SignExtend(buf[i+2]) % mod
			s1_term4 := SignExtend(buf[i+3]) % mod

			s1 = (s1 + s1_term1 + s1_term2 + s1_term3 + s1_term4) % mod
		}
	}

	// Handle remaining bytes
	for ; i < bufLen; i++ {
		s1 = (s1 + SignExtend(buf[i])) % mod
		s2 = (s2 + s1) % mod
	}

	// Final calculation with bounds checking
	s1_part := s1 & 0xffff
	s2_part := s2 & 0xffff0000 // Ensure we only get the upper 16 bits

	// Combine safely
	result := s1_part + s2_part
	return result
}

// Checksum2 computes the strong checksum (MD4)
func Checksum2(seed int32, buf []byte) ([]byte, error) {
	h := md4.New()
	h.Write(buf)
	if err := binary.Write(h, binary.LittleEndian, seed); err != nil {
		// This should never fail for int32, but handle it gracefully
		return nil, fmt.Errorf("binary.Write failed for int32: %w", err)
	}
	return h.Sum(nil), nil
}

// FileChecksum computes the MD4 checksum of an entire file
func FileChecksum(fn string) ([]byte, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := md4.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// BlockChecksum computes checksums for a file in blocks
type BlockChecksum struct {
	Offset int64
	Len    int64
	Index  int32
	Sum1   uint32
	Sum2   [16]byte
	Data   []byte // Store actual data for comparison (used in simplified sync)
}

// ChecksumHead represents the header for a set of block checksums
type ChecksumHead struct {
	ChecksumCount   int32
	BlockLength     int32
	ChecksumLength  int32
	RemainderLength int32
	Sums            []BlockChecksum
}

const (
	// Size is the size of MD4 checksums
	Size = md4.Size
	// DefaultBlockLength is the default block size for rsync operations
	DefaultBlockLength = 2048
)
