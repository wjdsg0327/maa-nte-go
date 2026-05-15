package service

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestEncodeScreenshotImageReturnsPNGDataURLAndSize(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 0, color.RGBA{G: 255, A: 255})

	got, err := encodeScreenshotImage(img)
	if err != nil {
		t.Fatalf("encode screenshot: %v", err)
	}

	if got.Width != 2 || got.Height != 1 {
		t.Fatalf("expected 2x1 screenshot, got %dx%d", got.Width, got.Height)
	}
	if got.Format != "png" {
		t.Fatalf("expected png format, got %q", got.Format)
	}
	if !strings.HasPrefix(got.DataURL, "data:image/png;base64,") {
		t.Fatalf("expected PNG data URL, got %q", got.DataURL)
	}
}

func TestEncodeScreenshotImageRejectsNilImage(t *testing.T) {
	_, err := encodeScreenshotImage(nil)
	if err == nil {
		t.Fatal("expected nil image to be rejected")
	}
}
