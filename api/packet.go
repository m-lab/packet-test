package api

import (
	"time"

	"github.com/m-lab/ndt-server/metadata"
)

// Packet represents the packet sent for network testing.
type Packet struct {
	Sequence int           // Sequence.
	Sent     int64         // Sent timestamp.
	Gap      time.Duration // Probe gap.
	Data     []byte        // Data transmitted.
}

// Result represents the result of a packet test.
type Result struct {
	Server         string
	Client         string
	ClientMetadata []metadata.NameValue
	Measurements   []Measurement
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
