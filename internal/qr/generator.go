package qr

import (
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
)

// Generator generates QR codes
type Generator struct {
	size int
}

// NewGenerator creates a new QR code generator
func NewGenerator() *Generator {
	return &Generator{
		size: 256, // default size
	}
}

// Generate generates a QR code PNG for the given data
func (g *Generator) Generate(data string) ([]byte, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("create qr code: %w", err)
	}

	png, err := qr.PNG(g.size)
	if err != nil {
		return nil, fmt.Errorf("generate png: %w", err)
	}

	return png, nil
}

// GenerateBase64 generates a base64-encoded QR code
func (g *Generator) GenerateBase64(data string) (string, error) {
	png, err := g.Generate(data)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(png), nil
}

// GenerateDataURI generates a data URI for embedding in HTML
func (g *Generator) GenerateDataURI(data string) (string, error) {
	base64Data, err := g.GenerateBase64(data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("data:image/png;base64,%s", base64Data), nil
}

// GenerateWithSize generates a QR code with custom size
func (g *Generator) GenerateWithSize(data string, size int) ([]byte, error) {
	if size < 21 || size > 1024 {
		return nil, fmt.Errorf("size must be between 21 and 1024")
	}

	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("create qr code: %w", err)
	}

	return qr.PNG(size)
}
