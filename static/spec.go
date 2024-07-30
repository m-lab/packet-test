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
	// MaxCwndGainParameterName is the name of a client parameter whose value indicates a BBR
	// congestion window (cwnd) gain after which the test should exit.
	MaxCwndGainParameterName = "max_cwnd_gain"
)
