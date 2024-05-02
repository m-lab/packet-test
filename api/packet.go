package api

import "time"

// Packet represents the packet sent for network testing.
type Packet struct {
	Sequence int           // Sequence.
	Sent     int64         // Sent timestamp.
	Gap      time.Duration // Probe gap.
	Data     []byte        // Data transmitted.
}

// Pair1Result represents the result of a Packet Pair test.
type Pair1Result struct {
	Server       string
	Client       string
	Measurements []Measurement
}

// Train1Result represents the result of the train1 test as written
// to disk.
type Train1Result struct {
	Server       string
	Client       string
	Measurements []Measurement
}

// Received encapsulates the structure received over the network
type Received struct {
	Sequence int
	Sent     time.Time
	Received time.Time
	Gap      int64 // usecs.
	Size     int64 // Bytes.
}

// Measurement represents a measurement result.
type Measurement struct {
	Packets []*Received
	Metrics
}

type Metrics struct {
	Dispersion int64 // usecs.
	Bandwidth  int64 // Mbps.
}
