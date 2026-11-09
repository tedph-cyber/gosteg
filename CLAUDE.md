# CLAUDE.md

Go library and CLI for hiding payloads in images with LSB steganography. Personal project, early stage.

## Layout

- `cmd/stegcli/` CLI: `embed`, `extract`, `inspect`, `qr`, `help`
- `image/` (package `img`): image load/save, pixel and channel helpers, `Embedder`, `Extractor`
- `payload/` (package `payload`): bit helpers and the payload header (`[1 byte version][4 bytes big-endian length][data]`)
- `qr/`: QR code generation and embedding helpers
- `internal/`: small shared utilities

## Facts about the repo

- The Go module path is `github.com/theoo/go-pro`, not the repository name `gosteg`. Do not rename it unless asked.
- There are currently no tests.
- `go vet ./...` currently reports an issue in `cmd/stegcli/main.go` (redundant newline in `fmt.Println`).
- There is no license file yet. Do not choose one without asking.
- Embedding is sequential LSB across the R, G and B channels. The README lists planned phases: randomized embedding, compression, encryption, integrity checks.

## Working rules

- Read the README and the code you will touch before changing anything.
- Run `go build ./...`, `go vet ./...` and `go test ./...` before saying work is done, and report the real results, including failures that existed before your change.
- Keep changes focused. Prefer small, idiomatic Go and the standard library. Add dependencies only when needed and say why.
- Keep existing CLI commands and flags working unless asked to change them.
- Say clearly what is verified, what is not, and what is a limitation. Do not overstate what steganography or any security feature protects against.
- Ask before: renaming the module, choosing a license, deleting files, or changing the payload format in a way that breaks existing stego images.
- Do not run `git init` or `git clone` in this project, and do not commit or push unless asked.
