package img

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

// Channel represents a color channel
type Channel uint8

const (
	Red Channel = iota
	Green
	Blue
	Alpha
)

// LoadImage loads an image from a file path
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return img, nil
}

// SaveImage saves an image to a file path (PNG format)
func SaveImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}
	return nil
}

// Width returns the width of the image
func Width(img image.Image) int {
	bounds := img.Bounds()
	return bounds.Dx()
}

// Height returns the height of the image
func Height(img image.Image) int {
	bounds := img.Bounds()
	return bounds.Dy()
}

// Bounds returns the bounding rectangle of the image
func Bounds(img image.Image) image.Rectangle {
	return img.Bounds()
}

// GetPixel returns the color at coordinates (x, y)
func GetPixel(img image.Image, x, y int) color.Color {
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		panic("coordinates out of bounds")
	}
	return img.At(x, y)
}

// SetPixel sets the color at coordinates (x, y)
func SetPixel(img image.Image, x, y int, c color.Color) error {
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return fmt.Errorf("coordinates out of bounds")
	}
	drawImg, ok := img.(draw.Image)
	if !ok {
		return fmt.Errorf("image does not support pixel modification")
	}

	drawImg.Set(x, y, c)
	return nil
}

// GetChannel returns the value of a specific channel at (x, y)
func GetChannel(img image.Image, x, y int, channel Channel) uint8 {
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		panic("coordinates out of bounds")
	}

	r, g, b, a := img.At(x, y).RGBA()
	r, g, b, a = r>>8, g>>8, b>>8, a>>8

	switch channel {
	case Red:
		return uint8(r)
	case Green:
		return uint8(g)
	case Blue:
		return uint8(b)
	case Alpha:
		return uint8(a)
	default:
		panic("invalid channel")
	}
}

// SetChannel sets the value of a specific channel at (x, y)
func SetChannel(img image.Image, x, y int, channel Channel, value uint8) error {
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return fmt.Errorf("coordinates out of bounds")
	}
	drawImg, ok := img.(draw.Image)
	if !ok {
		return fmt.Errorf("image does not support pixel modification")
	}

	r, g, b, a := img.At(x, y).RGBA()
	r, g, b, a = r>>8, g>>8, b>>8, a>>8

	switch channel {
	case Red:
		r = uint32(value)
	case Green:
		g = uint32(value)
	case Blue:
		b = uint32(value)
	case Alpha:
		a = uint32(value)
	default:
		return fmt.Errorf("invalid channel")
	}

	newColor := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	drawImg.Set(x, y, newColor)
	return nil
}

// SetLSB clears the LSB and sets it to the new bit value
func SetLSB(value uint8, bit uint8) uint8 {
	return (value & 0xFE) | (bit & 0x01)
}

// GetLSB extracts the least significant bit
func GetLSB(value uint8) uint8 {
	return value & 0x01
}

// Capacity returns total pixel capacity (width * height * 3 channels)
func Capacity(img image.Image) int {
	bounds := img.Bounds()
	return bounds.Dx() * bounds.Dy() * 3 // 3 channels per pixel
}
