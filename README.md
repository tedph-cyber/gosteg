# go-pro — Steganography Framework

A lightweight Go library and CLI for embedding hidden payloads in images using sequential LSB steganography.

## Features

- Image load/save and pixel access helpers in the `image` package.
- Bit-level payload utilities and header serialization in the `payload` package.
- Sequential LSB embedding and extraction with payload validation.
- QR code generation helpers in the `qr` package.
- CLI in `cmd/stegcli` with `embed`, `extract`, `inspect`, and `qr` commands.

## Build

```bash
go mod tidy
go build ./...
```

## CLI usage

### Embed inline text

```bash
go run ./cmd/stegcli embed -input cover.png -text "hello from stego" -output stego.png
```

This prints the embedded text to the terminal as a confirmation and saves the modified image.

### Embed a file

```bash
go run ./cmd/stegcli embed -input cover.png -data secret.txt -output stego.png
```

### Extract to a file

```bash
go run ./cmd/stegcli extract -input stego.png -output extracted.txt
```

### Extract directly to the terminal

```bash
go run ./cmd/stegcli extract -input stego.png -stdout
```

This prints the hidden message in the terminal instead of writing to a file.

### Generate a QR image

```bash
go run ./cmd/stegcli qr -input "Hello World" -output qr.png
```

### Inspect an image

```bash
go run ./cmd/stegcli inspect -image stego.png
```

This shows:

- image dimensions
- estimated payload capacity
- whether an embedded payload is detected
- a text preview when the hidden payload looks like text

## Current implementation notes

- Embedding uses a simple sequential LSB approach across RGB channels.
- Payloads are wrapped in a small version+length header before writing to the image.
- The project is laid out to support future phases such as adaptive/randomized embedding, compression, encryption, hashing/integrity checks, and more advanced data encoding.

## Contributing

Pull requests and issues are welcome.

## License

This project is currently provided without a formal license file. Add a license before publishing or distributing it broadly.
