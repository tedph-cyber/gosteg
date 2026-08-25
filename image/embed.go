package img

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/theoo/go-pro/payload"
)

type Embedder struct {
	img     image.Image
	x       int
	y       int
	minX    int
	minY    int
	maxX    int
	maxY    int
	channel uint8 // Current channel: 0=Red, 1=Green, 2=Blue
}

func NewEmbedder(img image.Image) *Embedder {
	bounds := img.Bounds()
	return &Embedder{
		img:     img,
		x:       bounds.Min.X,
		y:       bounds.Min.Y,
		minX:    bounds.Min.X,
		minY:    bounds.Min.Y,
		maxX:    bounds.Max.X,
		maxY:    bounds.Max.Y,
		channel: 0, // Start with Red
	}
}

func (e *Embedder) EmbedBits(bits []uint8) error {
	for i := 0; i < len(bits); i++ {
		if e.y >= e.maxY {
			return fmt.Errorf("insufficient capacity in image")
		}

		// Get current pixel and extract channel values
		pixel := e.img.At(e.x, e.y)
		r, g, b, a := pixel.RGBA()
		r, g, b, a = r>>8, g>>8, b>>8, a>>8

		// Modify LSB of current channel
		switch e.channel {
		case 0: // Red
			r = uint32(SetLSB(uint8(r), bits[i]))
		case 1: // Green
			g = uint32(SetLSB(uint8(g), bits[i]))
		case 2: // Blue
			b = uint32(SetLSB(uint8(b), bits[i]))
		}

		// Set the modified pixel
		newColor := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
		drawImg, ok := e.img.(draw.Image)
		if !ok {
			return fmt.Errorf("image does not support pixel modification")
		}

		drawImg.Set(e.x, e.y, newColor)

		e.NextPixel()
	}

	return nil
}

func (e *Embedder) EmbedBytes(data []byte) error {
	bits := payload.BytesToBits(data)
	return e.EmbedBits(bits)
}

func (e *Embedder) EmbedPayload(pl *payload.Payload) error {
	// Serialize payload to bytes
	payloadBytes := pl.Marshal()

	// Embed as bits
	return e.EmbedBytes(payloadBytes)
}

func (e *Embedder) NextPixel() {
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

func (e *Embedder) ResetCursor() {
	e.x = e.minX
	e.y = e.minY
	e.channel = 0
}
