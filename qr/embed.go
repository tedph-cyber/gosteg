package qr

import (
	"fmt"
	"image"

	img "github.com/theoo/go-pro/image"
	"github.com/theoo/go-pro/payload"
)

// EmbedPayload embeds a payload into a QR code image
func EmbedPayload(qrImg image.Image, pl *payload.Payload) (image.Image, error) {
	if qrImg == nil {
		return nil, fmt.Errorf("QR image cannot be nil")
	}

	if pl == nil {
		return nil, fmt.Errorf("payload cannot be nil")
	}

	if err := pl.Validate(); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	// Create embedder for the QR image
	embedder := img.NewEmbedder(qrImg)

	// Marshal payload to bytes
	payloadBytes := pl.Marshal()

	// Embed payload bytes into QR image
	if err := embedder.EmbedBytes(payloadBytes); err != nil {
		return nil, fmt.Errorf("failed to embed payload: %w", err)
	}

	return qrImg, nil
}

// ExtractPayload extracts a payload from a QR code image
func ExtractPayload(qrImg image.Image) (*payload.Payload, error) {
	if qrImg == nil {
		return nil, fmt.Errorf("QR image cannot be nil")
	}

	// Create extractor for the QR image
	extractor := img.NewExtractor(qrImg)

	// Extract payload from QR image
	pl, err := extractor.ExtractPayload()
	if err != nil {
		return nil, fmt.Errorf("failed to extract payload: %w", err)
	}

	// Validate extracted payload
	if err := pl.Validate(); err != nil {
		return nil, fmt.Errorf("extracted payload is invalid: %w", err)
	}

	return pl, nil
}

// EmbedIntoQRCode embeds data into a QR code and returns the modified image
func EmbedIntoQRCode(qrImg image.Image, data []byte) (image.Image, error) {
	if qrImg == nil {
		return nil, fmt.Errorf("QR image cannot be nil")
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	// Create payload from data
	pl := payload.NewPayload(data)

	// Embed payload into QR image
	return EmbedPayload(qrImg, pl)
}

// ExtractFromQRCode extracts data from a QR code image
func ExtractFromQRCode(qrImg image.Image) ([]byte, error) {
	if qrImg == nil {
		return nil, fmt.Errorf("QR image cannot be nil")
	}

	// Extract payload from QR image
	pl, err := ExtractPayload(qrImg)
	if err != nil {
		return nil, fmt.Errorf("failed to extract from QR code: %w", err)
	}

	return pl.Data, nil
}

// GetCapacity returns the maximum data capacity of a QR code image (in bytes)
func GetCapacity(qrImg image.Image) int {
	if qrImg == nil {
		return 0
	}

	// Capacity is 3 bits per pixel (R, G, B channels)
	return img.Capacity(qrImg) / 8
}

// CheckCapacity validates if the QR image has enough capacity for the payload
func CheckCapacity(qrImg image.Image, pl *payload.Payload) bool {
	if qrImg == nil || pl == nil {
		return false
	}

	capacity := GetCapacity(qrImg)
	requiredSize := pl.TotalSize()

	return capacity >= requiredSize
}
