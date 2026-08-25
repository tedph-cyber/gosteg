package payload

import (
	"fmt"
)

// Payload represents the data structure for steganographic payloads
type Payload struct {
	Version uint8  // Protocol version
	Length  uint32 // Length of data in bytes
	Data    []byte // Actual payload data
}

// NewPayload creates a new payload with the given data
func NewPayload(data []byte) *Payload {
	return &Payload{
		Version: 1, // Current protocol version
		Length:  uint32(len(data)),
		Data:    data,
	}
}

// Marshal converts payload to bytes for embedding
// Format: [1 byte version][4 bytes length][data]
func (p *Payload) Marshal() []byte {
	result := make([]byte, 5+len(p.Data))

	// Write version (1 byte)
	result[0] = p.Version

	// Write length (4 bytes, big-endian)
	result[1] = byte((p.Length >> 24) & 0xFF)
	result[2] = byte((p.Length >> 16) & 0xFF)
	result[3] = byte((p.Length >> 8) & 0xFF)
	result[4] = byte(p.Length & 0xFF)

	// Write data
	copy(result[5:], p.Data)

	return result
}

// Unmarshal creates a payload from marshaled bytes
func Unmarshal(data []byte) (*Payload, error) {
	if len(data) < 5 {
		return nil, fmt.Errorf("payload data too short: expected at least 5 bytes, got %d", len(data))
	}

	version := data[0]

	// Read length (4 bytes, big-endian)
	length := uint32(data[1])<<24 | uint32(data[2])<<16 | uint32(data[3])<<8 | uint32(data[4])

	// Validate length
	expectedDataLen := len(data) - 5
	if int(length) != expectedDataLen {
		return nil, fmt.Errorf("length mismatch: header says %d bytes, but got %d bytes", length, expectedDataLen)
	}

	// Extract data
	payloadData := make([]byte, length)
	copy(payloadData, data[5:])

	payload := &Payload{
		Version: version,
		Length:  length,
		Data:    payloadData,
	}

	return payload, nil
}

// Validate checks if the payload is valid
func (p *Payload) Validate() error {
	if p.Version == 0 {
		return fmt.Errorf("invalid version: 0")
	}

	if uint32(len(p.Data)) != p.Length {
		return fmt.Errorf("length mismatch: header says %d bytes, but data has %d bytes", p.Length, len(p.Data))
	}

	if p.Length > 1000000 { // 1MB limit
		return fmt.Errorf("payload too large: %d bytes exceeds 1MB limit", p.Length)
	}

	return nil
}

// Size returns the size of just the payload data
func (p *Payload) Size() int {
	return len(p.Data)
}

// HeaderSize returns the size of the payload header (version + length)
func (p *Payload) HeaderSize() int {
	return 5 // 1 byte version + 4 bytes length
}

// TotalSize returns the total size when marshaled (header + data)
func (p *Payload) TotalSize() int {
	return p.HeaderSize() + p.Size()
}
