package api

import "time"

// Packet represents the packet sent for network testing.
type Packet struct {
	Sequence int    // Sequence.
	Sent     int64  // Sent timestamp.
	Data     []byte // Data transmitted.
}

// Pair1Result represents the result of a Packet Pair test.
type Pair1Result struct {
	Capacity float64 // Mbps
}

// Train1Result represents the result of the train1 test as written
// to disk.
type Train1Result struct {
	Server       string
	Client       string
	StartTime    time.Time
	EndTime      time.Time
	Measurements []Measurement
}

// Received encapsulates the structure received over the network
type Received struct {
	Sequence int
	Sent     time.Time
	Received time.Time
	Size     int64
}

// Measurement represents a measurement result.
type Measurement struct {
	Packets    []*Received
	Dispersion int64
	Bandwidth  int64
}
