package api

// Packet represents the packet sent for network testing.
type Packet struct {
	Sequence  int    // Sequence number.
	Timestamp int64  // Timestamp (sent).
	Data      []byte // Data transmitted.
}
