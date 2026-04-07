package utils

import (
	"sync"
	"time"
)

// Snowflake ID generator
type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	dataCenterID  int64
	sequence      int64
}

const (
	workerIDBits     = 5
	dataCenterIDBits = 5
	sequenceBits     = 12

	workerIDShift     = sequenceBits
	dataCenterIDShift = sequenceBits + workerIDBits
	timestampShift    = sequenceBits + workerIDBits + dataCenterIDBits

	sequenceMask = -1 ^ (-1 << sequenceBits)
	maxWorkerID  = -1 ^ (-1 << workerIDBits)

	epoch = int64(1609459200000) // 2021-01-01 00:00:00 UTC
)

var (
	snowflake     *Snowflake
	snowflakeOnce sync.Once
)

// GetSnowflake returns singleton Snowflake instance
func GetSnowflake() *Snowflake {
	snowflakeOnce.Do(func() {
		snowflake = NewSnowflake(1, 1)
	})
	return snowflake
}

// NewSnowflake creates a new Snowflake instance
func NewSnowflake(workerID, dataCenterID int64) *Snowflake {
	if workerID < 0 || workerID > maxWorkerID {
		panic("worker ID out of range")
	}
	return &Snowflake{
		workerID:     workerID,
		dataCenterID: dataCenterID,
	}
}

// NextID generates next unique ID
func (s *Snowflake) NextID() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixMilli()

	if timestamp < s.lastTimestamp {
		panic("clock moved backwards")
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	id := ((timestamp - epoch) << timestampShift) |
		(s.dataCenterID << dataCenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return uint64(id)
}

// GenerateID generates a new Snowflake ID as uint64
func GenerateID() uint64 {
	return GetSnowflake().NextID()
}
