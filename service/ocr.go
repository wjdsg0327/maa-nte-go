package service

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MaaXYZ/maa-framework-go/v4"
)

// OCRResult 表示 OCR 识别结果
type OCRResult struct {
	Text  string  `json:"text"`
	Score float64 `json:"score"`
	X     int     `json:"x"`
	Y     int     `json:"y"`
	W     int     `json:"w"`
	H     int     `json:"h"`
}

// PrintOCRResults 打印 OCR 识别结果
func PrintOCRResults(results []*OCRResult, nodeName string) {
	if len(results) == 0 {
		log.Printf("[%s] 未识别到任何文字", nodeName)
		return
	}

	log.Printf("[%s] OCR 识别结果 (共 %d 个):", nodeName, len(results))
	log.Println(strings.Repeat("-", 60))

	for i, r := range results {
		log.Printf("  [%d] 文字: \"%s\"", i+1, r.Text)
		log.Printf("      置信度: %.4f", r.Score)
		log.Printf("      位置: (%d, %d) 大小: %dx%d", r.X, r.Y, r.W, r.H)
		if i < len(results)-1 {
			log.Println()
		}
	}

	log.Println(strings.Repeat("-", 60))
}

// ExtractOCRResults 从识别结果中提取 OCR 结果
func ExtractOCRResults(detail *maa.RecognitionDetail) []*OCRResult {
	if detail == nil || detail.Results == nil {
		return nil
	}

	var ocrResults []*OCRResult

	for _, r := range detail.Results.All {
		if ocr, ok := r.AsOCR(); ok {
			ocrResults = append(ocrResults, &OCRResult{
				Text:  ocr.Text,
				Score: ocr.Score,
				X:     ocr.Box.X(),
				Y:     ocr.Box.Y(),
				W:     ocr.Box.Width(),
				H:     ocr.Box.Height(),
			})
		}
	}

	return ocrResults
}

// FormatOCRResults 格式化 OCR 结果为字符串
func FormatOCRResults(results []*OCRResult) string {
	if len(results) == 0 {
		return "无识别结果"
	}

	var sb strings.Builder
	for i, r := range results {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("\"%s\"(%.2f)", r.Text, r.Score))
	}
	return sb.String()
}

// FindTextByText 根据文字内容查找匹配的 OCR 结果
func FindTextByText(results []*OCRResult, text string) *OCRResult {
	for _, r := range results {
		if r.Text == text {
			return r
		}
	}
	return nil
}

// FindTextByContains 查找包含指定文字的 OCR 结果
func FindTextByContains(results []*OCRResult, substr string) []*OCRResult {
	var matched []*OCRResult
	for _, r := range results {
		if strings.Contains(r.Text, substr) {
			matched = append(matched, r)
		}
	}
	return matched
}

// GetHighestScoreResult 获取置信度最高的 OCR 结果
func GetHighestScoreResult(results []*OCRResult) *OCRResult {
	if len(results) == 0 {
		return nil
	}

	best := results[0]
	for _, r := range results[1:] {
		if r.Score > best.Score {
			best = r
		}
	}
	return best
}

// SaveDrawImage 保存带标记的识别结果图片
func SaveDrawImage(detail *maa.RecognitionDetail, nodeName string) error {
	if detail == nil {
		return fmt.Errorf("识别详情为空")
	}

	// 获取原始截图
	rawImg := detail.Raw
	if rawImg == nil {
		return fmt.Errorf("无原始截图")
	}

	// 创建可绘制的图像
	bounds := rawImg.Bounds()
	drawImg := image.NewRGBA(bounds)
	draw.Draw(drawImg, bounds, rawImg, bounds.Min, draw.Src)

	// 定义标记颜色
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	green := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	// 绘制识别框
	if detail.Results != nil {
		for i, r := range detail.Results.All {
			var box maa.Rect
			var clr color.RGBA

			if ocr, ok := r.AsOCR(); ok {
				box = ocr.Box
				clr = red // OCR 用红色
			} else if tm, ok := r.AsTemplateMatch(); ok {
				box = tm.Box
				clr = green // 模板匹配用绿色
				_ = i
			} else {
				continue
			}

			// 绘制矩形框
			drawRect(drawImg, box.X(), box.Y(), box.Width(), box.Height(), clr, 2)
		}
	}

	// 保存到文件
	visionDir := filepath.Join(".", "log", "vision")
	os.MkdirAll(visionDir, 0755)

	timestamp := time.Now().Format("2006.01.02-15.04.05.000")
	filename := fmt.Sprintf("%s_%s.png", timestamp, nodeName)
	filePath := filepath.Join(visionDir, filename)

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, drawImg); err != nil {
		return fmt.Errorf("编码图片失败: %v", err)
	}

	log.Printf("标记图片已保存: %s", filePath)
	return nil
}

// drawRect 在图像上绘制矩形框
func drawRect(img *image.RGBA, x, y, w, h int, clr color.RGBA, thickness int) {
	if w <= 0 || h <= 0 {
		return
	}

	bounds := img.Bounds()

	for t := 0; t < thickness; t++ {
		// 上边
		for i := x; i < x+w; i++ {
			if i >= bounds.Min.X && i < bounds.Max.X && y-t >= bounds.Min.Y && y-t < bounds.Max.Y {
				img.SetRGBA(i, y-t, clr)
			}
			if i >= bounds.Min.X && i < bounds.Max.X && y+h-1+t >= bounds.Min.Y && y+h-1+t < bounds.Max.Y {
				img.SetRGBA(i, y+h-1+t, clr)
			}
		}
		// 左右边
		for j := y; j < y+h; j++ {
			if x-t >= bounds.Min.X && x-t < bounds.Max.X && j >= bounds.Min.Y && j < bounds.Max.Y {
				img.SetRGBA(x-t, j, clr)
			}
			if x+w-1+t >= bounds.Min.X && x+w-1+t < bounds.Max.X && j >= bounds.Min.Y && j < bounds.Max.Y {
				img.SetRGBA(x+w-1+t, j, clr)
			}
		}
	}
}
