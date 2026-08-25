package img

import (
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder

	"github.com/theoo/go-pro/payload"
)

type Extractor struct {
	img     image.Image
	x       int
	y       int
	minX    int
	minY    int
	maxX    int
	maxY    int
	channel uint8
	bitPos  uint8 // Current bit position within a byte (0-7)
}

func NewExtractor(img image.Image) *Extractor {
	bounds := img.Bounds()
	return &Extractor{
		img:     img,
		x:       bounds.Min.X,
		y:       bounds.Min.Y,
		minX:    bounds.Min.X,
		minY:    bounds.Min.Y,
		maxX:    bounds.Max.X,
		maxY:    bounds.Max.Y,
		channel: 0,
		bitPos:  0,
	}
}

func (e *Extractor) ExtractBits(length int) ([]uint8, error) {
	bits := make([]uint8, length)

	for i := 0; i < length; i++ {
		if e.y >= e.maxY {
			return nil, fmt.Errorf("insufficient data in image")
		}

		// Get LSB from current pixel's channels
		pixel := e.img.At(e.x, e.y)
		r, g, b, _ := pixel.RGBA()

		var value uint8
		switch e.channel {
		case 0: // Red
			value = uint8(r >> 8)
		case 1: // Green
			value = uint8(g >> 8)
		case 2: // Blue
			value = uint8(b >> 8)
		}

		bits[i] = GetLSB(value)
		e.NextPixel()
	}

	return bits, nil
}

func (e *Extractor) ExtractBytes(length int) ([]byte, error) {
	bits, err := e.ExtractBits(length * 8)
	if err != nil {
		return nil, err
	}

	return payload.BitsToBytes(bits), nil
}

func (e *Extractor) ExtractPayload() (*payload.Payload, error) {
	// Extract header (version + length)
	headerBytes, err := e.ExtractBytes(5) // 1 byte version + 4 bytes length
	if err != nil {
		return nil, err
	}

	version := headerBytes[0]
	length := uint32(headerBytes[1])<<24 | uint32(headerBytes[2])<<16 | uint32(headerBytes[3])<<8 | uint32(headerBytes[4])

	// Extract payload data
	data, err := e.ExtractBytes(int(length))
	if err != nil {
		return nil, err
	}

	pl := &payload.Payload{
		Version: version,
		Length:  length,
		Data:    data,
	}

	return pl, nil
}

func (e *Extractor) NextPixel() {
	e.channel++
	if e.channel >= 3 {
		e.channel = 0
		e.x++
		if e.x >= e.maxX {
			e.x = e.minX
			e.y++
		}
	}
}

func (e *Extractor) ResetCursor() {
	e.x = e.minX
	e.y = e.minY
	e.channel = 0
}
