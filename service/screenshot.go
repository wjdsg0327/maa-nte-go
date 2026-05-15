package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
)

type ScreenshotCapture struct {
	DataURL   string `json:"dataUrl"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Format    string `json:"format"`
	RawWidth  int32  `json:"rawWidth,omitempty"`
	RawHeight int32  `json:"rawHeight,omitempty"`
}

func encodeScreenshotImage(img image.Image) (*ScreenshotCapture, error) {
	if img == nil {
		return nil, errors.New("screenshot image is nil")
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid screenshot size: %dx%d", width, height)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encode screenshot png: %w", err)
	}

	return &ScreenshotCapture{
		DataURL: "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
		Width:   width,
		Height:  height,
		Format:  "png",
	}, nil
}

func (s *MaaService) CaptureScreenshot() (*ScreenshotCapture, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.inited || s.controller == nil {
		return nil, errors.New("请先连接窗口")
	}
	if !s.controller.Connected() {
		return nil, errors.New("请先连接窗口")
	}

	if err := waitMaaJob("截图", s.controller.PostScreencap()); err != nil {
		return nil, err
	}

	img, err := s.controller.CacheImage()
	if err != nil {
		return nil, fmt.Errorf("读取截图失败: %w", err)
	}

	capture, err := encodeScreenshotImage(img)
	if err != nil {
		return nil, err
	}

	if width, height, err := s.controller.GetResolution(); err == nil {
		capture.RawWidth = width
		capture.RawHeight = height
	}
	return capture, nil
}
