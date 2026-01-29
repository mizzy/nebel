package nebel

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed fonts/NotoSansCJKjp-Bold.otf
var fontData []byte

//go:embed assets/avatar.png
var avatarData []byte

const (
	ogImageWidth  = 1200
	ogImageHeight = 630
)

func (p *Post) generateOGImage(outputDir string) error {
	dc := gg.NewContext(ogImageWidth, ogImageHeight)

	// Draw background
	drawBackground(dc)

	// Load font
	f, err := opentype.Parse(fontData)
	if err != nil {
		return err
	}

	maxWidth := float64(ogImageWidth - 160) // 80px padding on each side

	// Calculate font size and wrap text, reducing font size if too many lines
	fontSize, lines, err := calculateFontSizeAndWrap(dc, f, p.Title, maxWidth)
	if err != nil {
		return err
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}

	dc.SetFontFace(face)

	// Text color (dark gray for white background)
	dc.SetColor(color.RGBA{45, 45, 45, 255})

	// Calculate total text height
	lineHeight := fontSize * 1.4
	totalHeight := float64(len(lines)) * lineHeight

	// Starting Y position (center vertically, with room for footer)
	startY := (float64(ogImageHeight)-totalHeight)/2 - 40

	// Draw each line centered
	for i, line := range lines {
		y := startY + float64(i)*lineHeight + fontSize
		dc.DrawStringAnchored(line, float64(ogImageWidth)/2, y, 0.5, 0.5)
	}

	// Draw footer with avatar, site name, and date
	drawFooter(dc, f, p.Date.Format("2006-01-02"))

	// Save image
	outputPath := filepath.Join(outputDir, "ogp.png")
	return dc.SavePNG(outputPath)
}

func drawBackground(dc *gg.Context) {
	// White background
	dc.SetColor(color.RGBA{255, 255, 255, 255})
	dc.Clear()
}

func drawFooter(dc *gg.Context, f *opentype.Font, date string) {
	// Load avatar image
	avatarImg, _, err := image.Decode(bytes.NewReader(avatarData))
	if err != nil {
		return
	}

	// Footer layout (right-aligned): 2026-01-27 | mizzy.org [avatar]
	avatarSize := 44.0
	rightMargin := 60.0
	avatarX := float64(ogImageWidth) - rightMargin - avatarSize/2
	footerY := float64(ogImageHeight) - 55

	// Scale avatar image to fit the circle
	avatarDC := gg.NewContext(int(avatarSize), int(avatarSize))
	avatarDC.DrawCircle(avatarSize/2, avatarSize/2, avatarSize/2)
	avatarDC.Clip()
	avatarDC.Scale(avatarSize/float64(avatarImg.Bounds().Dx()), avatarSize/float64(avatarImg.Bounds().Dy()))
	avatarDC.DrawImage(avatarImg, 0, 0)

	// Draw the circular avatar onto main context
	dc.DrawImageAnchored(avatarDC.Image(), int(avatarX), int(footerY), 0.5, 0.5)

	// Create smaller font for footer text
	smallFace, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return
	}

	dc.SetFontFace(smallFace)

	// Adjust text Y position to align with avatar center (font baseline offset)
	textY := footerY - 8

	// Draw site name to the left of avatar
	dc.SetColor(color.RGBA{45, 45, 45, 220})
	siteNameX := avatarX - avatarSize/2 - 20
	dc.DrawStringAnchored("mizzy.org", siteNameX, textY, 1, 0.5)

	// Draw separator between date and site name
	dc.SetColor(color.RGBA{45, 45, 45, 100})
	separatorX := siteNameX - 115
	dc.DrawStringAnchored("|", separatorX, textY, 0.5, 0.5)

	// Draw date to the left of separator
	dc.SetColor(color.RGBA{45, 45, 45, 180})
	dc.DrawStringAnchored(date, separatorX-15, textY, 1, 0.5)
}

func calculateFontSizeAndWrap(dc *gg.Context, f *opentype.Font, title string, maxWidth float64) (float64, []string, error) {
	fontSizes := []float64{72, 60, 48, 40, 32, 28, 24, 20, 18}
	maxLines := 1

	for _, fontSize := range fontSizes {
		face, err := opentype.NewFace(f, &opentype.FaceOptions{
			Size:    fontSize,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if err != nil {
			return 0, nil, err
		}

		dc.SetFontFace(face)
		lines := wrapText(dc, title, maxWidth)

		if len(lines) <= maxLines {
			return fontSize, lines, nil
		}
	}

	// Return smallest font size if still too many lines
	smallestSize := fontSizes[len(fontSizes)-1]
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    smallestSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return 0, nil, err
	}
	dc.SetFontFace(face)
	lines := wrapText(dc, title, maxWidth)
	return smallestSize, lines, nil
}

func wrapText(dc *gg.Context, text string, maxWidth float64) []string {
	totalWidth, _ := dc.MeasureString(text)
	if totalWidth <= maxWidth {
		return []string{text}
	}

	words, needsSpace := splitTextForWrappingWithSpace(text)
	bestSplitIdx := findBestSplitPoint(dc, words, needsSpace, totalWidth, maxWidth)
	lines := buildLines(words, needsSpace, bestSplitIdx)

	// If second line is still too long, recursively wrap it
	if len(lines) == 2 {
		secondWidth, _ := dc.MeasureString(lines[1])
		if secondWidth > maxWidth {
			subLines := wrapText(dc, lines[1], maxWidth)
			lines = append([]string{lines[0]}, subLines...)
		}
	}

	return lines
}

func findBestSplitPoint(dc *gg.Context, words []string, needsSpace []bool, totalWidth, maxWidth float64) int {
	targetWidth := totalWidth / 2
	bestSplitIdx := 0
	bestScore := totalWidth * 10

	var accumulated string
	for i, word := range words {
		if i > 0 && needsSpace[i] {
			accumulated += " "
		}
		accumulated += word

		accWidth, _ := dc.MeasureString(accumulated)
		if accWidth > maxWidth {
			continue
		}

		score := calculateSplitScore(word, words, i, accWidth, targetWidth, totalWidth)
		if score < bestScore {
			bestScore = score
			bestSplitIdx = i + 1
		}
	}

	return bestSplitIdx
}

func calculateSplitScore(word string, words []string, i int, accWidth, targetWidth, totalWidth float64) float64 {
	diff := abs(accWidth - targetWidth)
	score := diff

	// Bonus for splitting after particles (better word boundaries)
	ratio := accWidth / totalWidth
	if isJapaneseParticle(word) && ratio >= 0.3 && ratio <= 0.7 {
		score -= targetWidth * 0.5
	}

	// Penalty for splitting between consecutive kanji
	if i+1 < len(words) && len(word) > 0 && len(words[i+1]) > 0 {
		currentRunes := []rune(word)
		nextRunes := []rune(words[i+1])
		if isKanji(currentRunes[len(currentRunes)-1]) && isKanji(nextRunes[0]) {
			score += targetWidth * 1.0
		}
	}

	return score
}

func buildLines(words []string, needsSpace []bool, splitIdx int) []string {
	var lines []string
	var line1, line2 string

	for i, word := range words {
		if i < splitIdx {
			if i > 0 && needsSpace[i] {
				line1 += " "
			}
			line1 += word
		} else {
			if line2 != "" && needsSpace[i] {
				line2 += " "
			}
			line2 += word
		}
	}

	if line1 != "" {
		lines = append(lines, line1)
	}
	if line2 != "" {
		lines = append(lines, line2)
	}

	return lines
}

func isJapaneseParticle(s string) bool {
	particles := []string{"に", "を", "は", "が", "で", "と", "へ", "の", "も", "や", "から", "まで", "より"}
	for _, p := range particles {
		if s == p {
			return true
		}
	}
	return false
}

func isKanji(r rune) bool {
	// CJK Unified Ideographs
	if r >= 0x4E00 && r <= 0x9FFF {
		return true
	}
	// CJK Unified Ideographs Extension A
	if r >= 0x3400 && r <= 0x4DBF {
		return true
	}
	return false
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func splitTextForWrappingWithSpace(text string) ([]string, []bool) {
	var result []string
	var needsSpace []bool
	var current strings.Builder
	lastWasCJK := false

	for _, r := range text {
		// Check if it's a CJK character (Japanese, Chinese, Korean)
		if isCJK(r) {
			if current.Len() > 0 {
				result = append(result, current.String())
				needsSpace = append(needsSpace, !lastWasCJK)
				current.Reset()
			}
			result = append(result, string(r))
			needsSpace = append(needsSpace, false) // No space before CJK
			lastWasCJK = true
		} else if r == ' ' {
			if current.Len() > 0 {
				result = append(result, current.String())
				needsSpace = append(needsSpace, true)
				current.Reset()
			}
			lastWasCJK = false
		} else {
			current.WriteRune(r)
			lastWasCJK = false
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
		needsSpace = append(needsSpace, !lastWasCJK)
	}

	// First element never needs leading space
	if len(needsSpace) > 0 {
		needsSpace[0] = false
	}

	return result, needsSpace
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
