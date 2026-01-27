package nebel

import (
	_ "embed"
	"image/color"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed fonts/NotoSansCJKjp-Bold.otf
var fontData []byte

const (
	ogImageWidth  = 1200
	ogImageHeight = 630
)

func (p *Post) generateOGImage(outputDir string) error {
	dc := gg.NewContext(ogImageWidth, ogImageHeight)

	// Background color (#333)
	dc.SetColor(color.RGBA{51, 51, 51, 255})
	dc.Clear()

	// Load font
	f, err := opentype.Parse(fontData)
	if err != nil {
		return err
	}

	fontSize := calculateFontSize(p.Title)

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}

	dc.SetFontFace(face)

	// Text color (white)
	dc.SetColor(color.White)

	// Calculate text wrapping
	maxWidth := float64(ogImageWidth - 100) // 50px padding on each side
	lines := wrapText(dc, p.Title, maxWidth)

	// Calculate total text height
	lineHeight := fontSize * 1.4
	totalHeight := float64(len(lines)) * lineHeight

	// Starting Y position (center vertically)
	startY := (float64(ogImageHeight) - totalHeight) / 2

	// Draw each line centered
	for i, line := range lines {
		y := startY + float64(i)*lineHeight + fontSize
		dc.DrawStringAnchored(line, float64(ogImageWidth)/2, y, 0.5, 0.5)
	}

	// Save image
	outputPath := filepath.Join(outputDir, "ogp.png")
	return dc.SavePNG(outputPath)
}

func calculateFontSize(title string) float64 {
	length := utf8.RuneCountInString(title)

	switch {
	case length <= 15:
		return 72
	case length <= 25:
		return 60
	case length <= 40:
		return 48
	default:
		return 40
	}
}

func wrapText(dc *gg.Context, text string, maxWidth float64) []string {
	var lines []string
	var currentLine string

	// Split by spaces for English, but also handle Japanese characters
	words := splitTextForWrapping(text)

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		width, _ := dc.MeasureString(testLine)
		if width > maxWidth && currentLine != "" {
			lines = append(lines, strings.TrimSpace(currentLine))
			currentLine = word
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		lines = append(lines, strings.TrimSpace(currentLine))
	}

	return lines
}

func splitTextForWrapping(text string) []string {
	var result []string
	var current strings.Builder

	for _, r := range text {
		// Check if it's a CJK character (Japanese, Chinese, Korean)
		if isCJK(r) {
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
			result = append(result, string(r))
		} else if r == ' ' {
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

func isCJK(r rune) bool {
	// CJK Unified Ideographs
	if r >= 0x4E00 && r <= 0x9FFF {
		return true
	}
	// Hiragana
	if r >= 0x3040 && r <= 0x309F {
		return true
	}
	// Katakana
	if r >= 0x30A0 && r <= 0x30FF {
		return true
	}
	// CJK Unified Ideographs Extension A
	if r >= 0x3400 && r <= 0x4DBF {
		return true
	}
	// CJK Symbols and Punctuation
	if r >= 0x3000 && r <= 0x303F {
		return true
	}
	// Halfwidth and Fullwidth Forms
	if r >= 0xFF00 && r <= 0xFFEF {
		return true
	}
	return false
}
