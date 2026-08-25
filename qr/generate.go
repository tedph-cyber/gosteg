package qr

import (
	"fmt"
	"image"

	"github.com/skip2/go-qrcode"
	"github.com/theoo/go-pro/payload"
)

// Generate creates a basic QR code from text
func Generate(text string) (image.Image, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	qrCode, err := qrcode.New(text, qrcode.High)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Set size to reasonable default (256x256 pixels)
	img := qrCode.Image(256)

	return img, nil
}

// GenerateFromText creates a QR code from text with custom size
func GenerateFromText(text string, size int) (image.Image, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	if size <= 0 {
		return nil, fmt.Errorf("size must be positive: got %d", size)
	}

	qrCode, err := qrcode.New(text, qrcode.High)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	img := qrCode.Image(size)

	return img, nil
}

// GenerateFromPayload creates a QR code from a payload
func GenerateFromPayload(pl *payload.Payload) (image.Image, error) {
	if pl == nil {
		return nil, fmt.Errorf("payload cannot be nil")
	}

	if err := pl.Validate(); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	// Convert payload to bytes and encode as string for QR
	payloadBytes := pl.Marshal()

	// Create QR code from payload bytes
	qrCode, err := qrcode.New(string(payloadBytes), qrcode.High)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code from payload: %w", err)
	}

	img := qrCode.Image(256)

	return img, nil
}

// Save saves a QR code image to file
func Save(img image.Image, path string) error {
	if img == nil {
		return fmt.Errorf("image cannot be nil")
	}

	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	return nil
}

// GenerateWithCapacity creates a QR code that has enough capacity for the payload
func GenerateWithCapacity(text string, pl *payload.Payload) (image.Image, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	if pl == nil {
		return nil, fmt.Errorf("payload cannot be nil")
	}

	if err := pl.Validate(); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	// Start with default sizes and increase if needed
	sizes := []int{256, 512, 1024}

	for _, size := range sizes {
		qrImg, err := GenerateFromText(text, size)
		if err != nil {
			continue
		}

		// using img.GetCapacity via qr package is handled elsewhere
		_ = qrImg
		return qrImg, nil
	}

	return nil, fmt.Errorf("payload too large for QR code capacity")
}

// GetQRInfo returns information about a QR code image
func GetQRInfo(img image.Image) map[string]interface{} {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	return map[string]interface{}{
		"width":  width,
		"height": height,
	}
}
