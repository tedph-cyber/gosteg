package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"unicode"
	"unicode/utf8"

	img "github.com/theoo/go-pro/image"
	internalpkg "github.com/theoo/go-pro/internal"
	pl "github.com/theoo/go-pro/payload"
	"github.com/theoo/go-pro/qr"
)

func main() {
	if len(os.Args) < 2 {
		helpCommand()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "embed":
		embedCommand()
	case "extract":
		extractCommand()
	case "qr":
		qrCommand()
	case "inspect":
		inspectCommand()
	case "help":
		helpCommand()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		helpCommand()
		os.Exit(1)
	}
}

func embedCommand() {
	fs := flag.NewFlagSet("embed", flag.ExitOnError)
	inputPath := fs.String("input", "", "Path to cover image (input)")
	dataPath := fs.String("data", "", "Path to data file")
	text := fs.String("text", "", "Text to embed (alternative to -data)")
	outputPath := fs.String("output", "", "Path to output stego image")

	fs.Parse(os.Args[2:])

	if *inputPath == "" || *outputPath == "" || (*dataPath == "" && *text == "") {
		fmt.Println("Usage: steg embed -input <cover.png> (-data <file> | -text <string>) -output <stego.png>")
		os.Exit(1)
	}

	fmt.Printf("Embedding into %s → %s\n", *inputPath, *outputPath)
	// Load cover image
	coverImg, err := img.LoadImage(*inputPath)
	if err != nil {
		fmt.Printf("failed to load cover image: %v\n", err)
		os.Exit(1)
	}

	// Determine payload bytes: inline text or file
	var payloadBytes []byte
	sourceFlag := ""
	if *text != "" {
		payloadBytes = []byte(*text)
		sourceFlag = "-text"
	} else {
		var err error
		payloadBytes, err = internalpkg.ReadFile(*dataPath)
		if err != nil {
			fmt.Printf("failed to read data file: %v\n", err)
			os.Exit(1)
		}
		sourceFlag = "-data"
	}

	// Create payload
	payload := pl.NewPayload(payloadBytes)

	// Embed payload into image
	embedder := img.NewEmbedder(coverImg)
	if err := embedder.EmbedPayload(payload); err != nil {
		fmt.Printf("failed to embed payload: %v\n", err)
		os.Exit(1)
	}

	// Save stego image
	if err := img.SaveImage(coverImg, *outputPath); err != nil {
		fmt.Printf("failed to save stego image: %v\n", err)
		os.Exit(1)
	}

	// Print summary and echo inline text
	fmt.Printf("Embedded into %s (source: %s) → %s\n", *inputPath, sourceFlag, *outputPath)
	if sourceFlag == "-text" {
		fmt.Printf("Embedded text:\n%s\n", string(payloadBytes))
	} else {
		fmt.Printf("Embedded %d bytes from file %s\n", len(payloadBytes), *dataPath)
	}
	fmt.Println("Embedding complete")
}

func extractCommand() {
	fs := flag.NewFlagSet("extract", flag.ExitOnError)
	inputPath := fs.String("input", "", "Path to stego image (input)")
	outputPath := fs.String("output", "", "Path to output file")
	stdout := fs.Bool("stdout", false, "Print extracted payload to terminal")

	fs.Parse(os.Args[2:])

	if *inputPath == "" {
		fmt.Println("Usage: steg extract -input <stego.png> [-output <file> | -stdout]")
		os.Exit(1)
	}

	writeToStdout := *stdout || *outputPath == ""
	if writeToStdout {
		fmt.Printf("Extracting from %s → terminal\n", *inputPath)
	} else {
		fmt.Printf("Extracting from %s → %s\n", *inputPath, *outputPath)
	}

	// Load stego image
	stegoImg, err := img.LoadImage(*inputPath)
	if err != nil {
		fmt.Printf("failed to load stego image: %v\n", err)
		os.Exit(1)
	}

	// Extract payload
	payload, err := readEmbeddedPayload(stegoImg)
	if err != nil {
		fmt.Printf("failed to extract payload: %v\n", err)
		os.Exit(1)
	}

	if writeToStdout {
		fmt.Println(string(payload.Data))
		fmt.Printf("Extraction complete (%d bytes printed to terminal)\n", len(payload.Data))
		return
	}

	// Write payload data to output file
	if err := internalpkg.WriteFile(*outputPath, payload.Data); err != nil {
		fmt.Printf("failed to write output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Extraction complete (%d bytes written)\n", len(payload.Data))
}

func qrCommand() {
	fs := flag.NewFlagSet("qr", flag.ExitOnError)
	input := fs.String("input", "", "Data to encode in QR code (input)")
	output := fs.String("output", "", "Path to output QR image")

	fs.Parse(os.Args[2:])

	if *input == "" || *output == "" {
		fmt.Println("Usage: steg qr -input <text> -output <image>")
		os.Exit(1)
	}

	fmt.Printf("Generating QR code from '%s' → %s\n", *input, *output)
	img, err := qr.Generate(*input)
	if err != nil {
		fmt.Printf("failed to generate QR code: %v\n", err)
		os.Exit(1)
	}

	outFile, err := os.Create(*output)
	if err != nil {
		fmt.Printf("failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		fmt.Printf("failed to encode PNG: %v\n", err)
		os.Exit(1)
	}
}

func inspectCommand() {
	fs := flag.NewFlagSet("inspect", flag.ExitOnError)
	imagePath := fs.String("image", "", "Path to image to inspect")

	fs.Parse(os.Args[2:])

	if *imagePath == "" {
		fmt.Println("Usage: steg inspect -image <image>")
		os.Exit(1)
	}

	fmt.Printf("Inspecting %s\n", *imagePath)
	imgFile, err := img.LoadImage(*imagePath)
	if err != nil {
		fmt.Printf("failed to load image: %v\n", err)
		os.Exit(1)
	}

	bounds := img.Bounds(imgFile)
	capacityBits := img.Capacity(imgFile)
	capacityBytes := capacityBits / 8

	fmt.Printf("Dimensions: %d x %d\n", bounds.Dx(), bounds.Dy())
	fmt.Printf("Capacity: %d bits (%d bytes usable)\n", capacityBits, capacityBytes)
	fmt.Println("Note: capacity assumes 3 LSB channels per pixel.")

	payload, err := readEmbeddedPayload(imgFile)
	if err != nil {
		fmt.Println("Embedded payload: none detected")
		return
	}

	if isLikelyText(payload.Data) {
		fmt.Printf("Embedded payload: text (%d bytes)\n", len(payload.Data))
		fmt.Printf("Preview: %s\n", previewText(string(payload.Data), 200))
		return
	}

	fmt.Printf("Embedded payload: file/binary data (%d bytes)\n", len(payload.Data))
}

func readEmbeddedPayload(stegoImg image.Image) (*pl.Payload, error) {
	extractor := img.NewExtractor(stegoImg)
	headerBytes, err := extractor.ExtractBytes(5)
	if err != nil {
		return nil, fmt.Errorf("no embedded payload detected")
	}

	version := headerBytes[0]
	length := uint32(headerBytes[1])<<24 | uint32(headerBytes[2])<<16 | uint32(headerBytes[3])<<8 | uint32(headerBytes[4])

	if version == 0 {
		return nil, fmt.Errorf("no embedded payload detected")
	}

	capacityBytes := img.Capacity(stegoImg) / 8
	if int(length)+5 > capacityBytes {
		return nil, fmt.Errorf("no embedded payload detected")
	}

	data, err := extractor.ExtractBytes(int(length))
	if err != nil {
		return nil, fmt.Errorf("no embedded payload detected")
	}

	payload := &pl.Payload{
		Version: version,
		Length:  length,
		Data:    data,
	}

	if err := payload.Validate(); err != nil {
		return nil, fmt.Errorf("no embedded payload detected")
	}

	return payload, nil
}

func isLikelyText(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	if !utf8.Valid(data) {
		return false
	}

	for _, r := range string(data) {
		if unicode.IsPrint(r) || unicode.IsSpace(r) {
			continue
		}
		return false
	}

	return true
}

func previewText(text string, limit int) string {
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}

	return string(runes[:limit]) + "..."
}

func helpCommand() {
	fmt.Println(`
Steganography CLI - Hide data in images

Usage:
  steg <command> [flags]

Commands:
  embed      Embed data into a cover image
  extract    Extract hidden data from a stego image
  qr         Generate a QR code
  inspect    Inspect image properties and capacity
  help       Show this help message

Examples:
	steg embed -input cover.png -data secret.txt -output stego.png
	steg embed -input cover.png -text "Secret message" -output stego.png
	steg extract -input stego.png -output extracted.txt
	steg extract -input stego.png -stdout
	steg qr -input "Hello World" -output qr.png
	steg inspect -image cover.png
`)
}
