package payload

import (
	"fmt"
)

// BytesToBits converts a byte slice to a bit slice
func BytesToBits(data []byte) []uint8 {
	bits := make([]uint8, len(data)*8)
	for i, b := range data {
		for j := 0; j < 8; j++ {
			bits[i*8+j] = (b >> (7 - j)) & 0x01
		}
	}
	return bits
}

// BitsToBytes converts a bit slice to a byte slice
func BitsToBytes(bits []uint8) []byte {
	if len(bits)%8 != 0 {
		panic("bits length must be multiple of 8")
	}

	bytes := make([]byte, len(bits)/8)
	for i := 0; i < len(bytes); i++ {
		for j := 0; j < 8; j++ {
			bytes[i] = (bytes[i] << 1) | (bits[i*8+j] & 0x01)
		}
	}
	return bytes
}

// ByteToBits converts a single byte to 8 bits
func ByteToBits(b byte) []uint8 {
	bits := make([]uint8, 8)
	for i := 0; i < 8; i++ {
		bits[i] = (b >> (7 - i)) & 0x01
	}
	return bits
}

// BitsToByte converts 8 bits to a single byte
func BitsToByte(bits []uint8) byte {
	if len(bits) != 8 {
		panic("bits slice must be exactly 8 elements")
	}

	var b byte
	for i := 0; i < 8; i++ {
		b = (b << 1) | (bits[i] & 0x01)
	}
	return b
}

// BitAt extracts a specific bit at index (0-7, where 0 is MSB)
func BitAt(value byte, index int) uint8 {
	if index < 0 || index > 7 {
		panic("bit index must be between 0 and 7")
	}
	return (value >> (7 - index)) & 0x01
}

// SetBit sets a specific bit at index (0-7, where 0 is MSB)
func SetBit(value byte, index int, bit uint8) byte {
	if index < 0 || index > 7 {
		panic("bit index must be between 0 and 7")
	}
	if bit&0x01 == 1 {
		return value | (1 << (7 - index))
	}
	return value & ^(1 << (7 - index))
}

// BitReader provides sequential bit reading from a byte slice
type BitReader struct {
	data []byte
	pos  int // Current bit position
}

// NewBitReader creates a new BitReader
func NewBitReader(data []byte) *BitReader {
	return &BitReader{
		data: data,
		pos:  0,
	}
}

// ReadBit reads a single bit
func (br *BitReader) ReadBit() (uint8, error) {
	if br.pos >= len(br.data)*8 {
		return 0, fmt.Errorf("end of data reached")
	}

	byteIndex := br.pos / 8
	bitIndex := br.pos % 8
	bit := (br.data[byteIndex] >> (7 - bitIndex)) & 0x01

	br.pos++
	return bit, nil
}

// ReadBits reads multiple bits
func (br *BitReader) ReadBits(count int) ([]uint8, error) {
	bits := make([]uint8, count)
	for i := 0; i < count; i++ {
		bit, err := br.ReadBit()
		if err != nil {
			return nil, err
		}
		bits[i] = bit
	}
	return bits, nil
}

// ReadBytes reads bits as bytes
func (br *BitReader) ReadBytes(count int) ([]byte, error) {
	bits, err := br.ReadBits(count * 8)
	if err != nil {
		return nil, err
	}
	return BitsToBytes(bits), nil
}

// Position returns current bit position
func (br *BitReader) Position() int {
	return br.pos
}

// Reset resets the reader to start
func (br *BitReader) Reset() {
	br.pos = 0
}

// BitWriter provides sequential bit writing to a byte slice
type BitWriter struct {
	data []byte
	pos  int // Current bit position
}

// NewBitWriter creates a new BitWriter
func NewBitWriter() *BitWriter {
	return &BitWriter{
		data: make([]byte, 0),
		pos:  0,
	}
}

// WriteBit writes a single bit
func (bw *BitWriter) WriteBit(bit uint8) {
	byteIndex := bw.pos / 8
	bitIndex := bw.pos % 8

	// Expand data if needed
	if byteIndex >= len(bw.data) {
		bw.data = append(bw.data, 0)
	}

	if bit&0x01 == 1 {
		bw.data[byteIndex] |= (1 << (7 - bitIndex))
	}

	bw.pos++
}

// WriteBits writes multiple bits
func (bw *BitWriter) WriteBits(bits []uint8) {
	for _, bit := range bits {
		bw.WriteBit(bit)
	}
}

// WriteBytes writes bytes as bits
func (bw *BitWriter) WriteBytes(data []byte) {
	bits := BytesToBits(data)
	bw.WriteBits(bits)
}

// Bytes returns the written data
func (bw *BitWriter) Bytes() []byte {
	return bw.data
}

// Position returns current bit position
func (bw *BitWriter) Position() int {
	return bw.pos
}

// Reset resets the writer
func (bw *BitWriter) Reset() {
	bw.data = make([]byte, 0)
	bw.pos = 0
}
