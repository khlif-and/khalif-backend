package utils

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"sort"

	"khalif-backend/internal/platform/logger"

	"go.uber.org/zap"
	_ "golang.org/x/image/webp"
)

const (
	DefaultColor     = "#1DB954"
	QuantizeFactor   = 32
	DarkThreshold    = 30
	LightThreshold   = 225
	TopColorsToCheck = 5
)

// ExtractDominantColor extracts the dominant vibrant color from an image file path.
// Returns hex color string (e.g., "#FF5733") or error.
func ExtractDominantColor(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		logError("Failed to open image", imagePath, err)
		return "", fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	return ExtractDominantColorFromReader(file)
}

// ExtractDominantColorFromReader extracts the dominant vibrant color from an io.Reader.
func ExtractDominantColorFromReader(r io.Reader) (string, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}
	return findDominantColor(img), nil
}

// GetThumbnailColor returns dominant color or default if extraction fails.
func GetThumbnailColor(path string) string {
	if path == "" {
		return DefaultColor
	}

	color, err := ExtractDominantColor(path)
	if err != nil {
		logError("Color extraction failed, using default", path, err)
		return DefaultColor
	}
	return color
}

// GetThumbnailColorFromReader returns dominant color from reader or default.
func GetThumbnailColorFromReader(r io.Reader) string {
	if r == nil {
		return DefaultColor
	}
	color, err := ExtractDominantColorFromReader(r)
	if err != nil {
		// We don't log path here as we don't have it
		return DefaultColor
	}
	return color
}

// GetThumbnailColorPath is alias for GetThumbnailColor (backward compatibility).
func GetThumbnailColorPath(path string) string {
	return GetThumbnailColor(path)
}

// ExtractDominantColorFromPath is alias for ExtractDominantColor.
func ExtractDominantColorFromPath(path string) (string, error) {
	return ExtractDominantColor(path)
}

func findDominantColor(img image.Image) string {
	colorCounts := countColors(img)

	if len(colorCounts) == 0 {
		return DefaultColor
	}

	sorted := sortByCount(colorCounts)
	return selectMostVibrant(sorted)
}

func countColors(img image.Image) map[string]int {
	bounds := img.Bounds()
	step := calculateStep(bounds)
	counts := make(map[string]int)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += step {
		for x := bounds.Min.X; x < bounds.Max.X; x += step {
			if hex := processPixel(img, x, y); hex != "" {
				counts[hex]++
			}
		}
	}

	return counts
}

func calculateStep(bounds image.Rectangle) int {
	pixels := (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)

	switch {
	case pixels > 500000:
		return 5
	case pixels > 100000:
		return 3
	default:
		return 1
	}
}

func processPixel(img image.Image, x, y int) string {
	r, g, b, a := img.At(x, y).RGBA()

	if a < 128<<8 {
		return ""
	}

	r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

	if isDark(r8, g8, b8) || isLight(r8, g8, b8) {
		return ""
	}

	return quantizeToHex(r8, g8, b8)
}

func isDark(r, g, b uint8) bool {
	return (int(r)+int(g)+int(b))/3 < DarkThreshold
}

func isLight(r, g, b uint8) bool {
	return (int(r)+int(g)+int(b))/3 > LightThreshold
}

func quantizeToHex(r, g, b uint8) string {
	r = (r / QuantizeFactor) * QuantizeFactor
	g = (g / QuantizeFactor) * QuantizeFactor
	b = (b / QuantizeFactor) * QuantizeFactor
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

type colorEntry struct {
	hex   string
	count int
}

func sortByCount(counts map[string]int) []colorEntry {
	entries := make([]colorEntry, 0, len(counts))
	for hex, count := range counts {
		entries = append(entries, colorEntry{hex, count})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].count > entries[j].count
	})

	return entries
}

func selectMostVibrant(sorted []colorEntry) string {
	limit := TopColorsToCheck
	if len(sorted) < limit {
		limit = len(sorted)
	}

	bestHex := sorted[0].hex
	bestSat := 0.0

	for i := 0; i < limit; i++ {
		if sat := saturation(sorted[i].hex); sat > bestSat {
			bestSat = sat
			bestHex = sorted[i].hex
		}
	}

	return bestHex
}

func saturation(hex string) float64 {
	var r, g, b uint8
	fmt.Sscanf(hex, "#%02X%02X%02X", &r, &g, &b)

	rf, gf, bf := float64(r)/255, float64(g)/255, float64(b)/255

	max := max3(rf, gf, bf)
	min := min3(rf, gf, bf)

	if max == min {
		return 0
	}

	l := (max + min) / 2
	if l > 0.5 {
		return (max - min) / (2 - max - min)
	}
	return (max - min) / (max + min)
}

func max3(a, b, c float64) float64 {
	if a >= b && a >= c {
		return a
	}
	if b >= c {
		return b
	}
	return c
}

func min3(a, b, c float64) float64 {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}

func logError(msg, path string, err error) {
	if logger.Log != nil {
		logger.Log.Error(msg, zap.String("path", path), zap.Error(err))
	}
}
