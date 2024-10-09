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
	NDT7DownloadURLPath    = "/v0/ndt7/download"
	EarlyExitParameterName = "early_exit"
	// MaxCwndGainParameterName is the name of a client parameter whose value indicates a BBR
	// congestion window (cwnd) gain after which the test should exit.
	MaxCwndGainParameterName = "max_cwnd_gain"
	// MaxElapsedTime is the name of a client parameter whose values indicates the
	// number of seconds after which the test should exit.
	MaxElapsedTimeParameterName = "max_elapsed_time"
	// ImmediateExitParameterName is the name of the parameter to indicate that the termination
	// behavior of the test should be immediate, instead of waiting for TCPInfo and BBRInfo
	// snapshots to be emitted.
	ImmediateExitParameterName = "immediate_exit"

	// MinPoissonSamplingInterval is the min acceptable time that we want
	// the lambda distribution to return. Smaller values will be clamped
	// to be this value instead.
	MinPoissonSamplingInterval = 25 * time.Millisecond

	// AveragePoissonSamplingInterval is the average of a lambda distribution
	// used to decide when to perform next measurement.
	AveragePoissonSamplingInterval = 100 * time.Millisecond

	// MaxPoissonSamplingInterval is the max acceptable time that we want
	// the lambda distribution to return. Bigger values will be clamped
	// to be this value instead.
	MaxPoissonSamplingInterval = 250 * time.Millisecond
)
