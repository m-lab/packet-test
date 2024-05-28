package static

import "time"

// Constants used for packet testing.
const (
	BufferBytes            = 508
	PacketBytes            = 720
	PairCount              = 30
	PairDelay              = 1 * time.Second
	PairGap                = 10 * time.Microsecond
	TrainCount             = 10
	TrainDelay             = 1 * time.Second
	TrainLength            = 30
	EarlyExitParameterName = "early_exit"
	BBRExitParameterName   = "bbr_exit"
)
