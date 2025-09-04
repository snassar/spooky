package sync

import (
	"io"
	"log/slog"
	"os"

	"github.com/mmcloughlin/md4"
)

// FileChecksum computes the MD4 checksum of an entire file
func FileChecksum(fn string) ([]byte, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			// Log the error but don't fail the function since we've already read the data
			// This is a best-effort cleanup
			slog.Warn("failed to close file during cleanup",
				slog.String("file", fn),
				slog.String("error", closeErr.Error()))
		}
	}()
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
