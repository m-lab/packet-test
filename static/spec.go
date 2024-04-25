package static

import "time"

// Constants used for packet testing.
const (
	BufferBytes = 508
	PacketBytes = 726
	PairCount   = 30
	PairDelay   = 10000 * time.Microsecond
	TrainCount  = 10
	TrainDelay  = 1 * time.Second
	TrainLength = 30
)
