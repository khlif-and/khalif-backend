package utils

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractDominantColorFromReader_Success(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	red := color.RGBA{R: 200, G: 50, B: 50, A: 255}
	fillImage(img, red)

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	assert.NoError(t, err)

	color, err := ExtractDominantColorFromReader(&buf)
	assert.NoError(t, err)
	assert.Regexp(t, `^#[0-9A-F]{6}$`, color)
}

func TestExtractDominantColor_RedImage(t *testing.T) {
	testDir := t.TempDir()
	testImagePath := filepath.Join(testDir, "red_test.png")

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	red := color.RGBA{R: 200, G: 50, B: 50, A: 255}
	fillImage(img, red)
	saveImage(t, testImagePath, img)

	extractedColor, err := ExtractDominantColor(testImagePath)

	assert.NoError(t, err)
	assert.NotEmpty(t, extractedColor)
	assert.Regexp(t, `^#[0-9A-F]{6}$`, extractedColor)
	t.Logf("INPUT: RGB(200, 50, 50) -> OUTPUT: %s", extractedColor)
}

func TestExtractDominantColor_BlueImage(t *testing.T) {
	testDir := t.TempDir()
	testImagePath := filepath.Join(testDir, "blue_test.png")

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	blue := color.RGBA{R: 50, G: 50, B: 200, A: 255}
	fillImage(img, blue)
	saveImage(t, testImagePath, img)

	extractedColor, err := ExtractDominantColor(testImagePath)

	assert.NoError(t, err)
	assert.NotEmpty(t, extractedColor)
	t.Logf("INPUT: RGB(50, 50, 200) -> OUTPUT: %s", extractedColor)
}

func TestExtractDominantColor_GreenImage(t *testing.T) {
	testDir := t.TempDir()
	testImagePath := filepath.Join(testDir, "green_test.png")

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	green := color.RGBA{R: 50, G: 180, B: 50, A: 255}
	fillImage(img, green)
	saveImage(t, testImagePath, img)

	extractedColor, err := ExtractDominantColor(testImagePath)

	assert.NoError(t, err)
	assert.NotEmpty(t, extractedColor)
	t.Logf("INPUT: RGB(50, 180, 50) -> OUTPUT: %s", extractedColor)
}

func TestExtractDominantColor_MixedColors(t *testing.T) {
	testDir := t.TempDir()
	testImagePath := filepath.Join(testDir, "mixed_test.png")

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	orange := color.RGBA{R: 255, G: 165, B: 0, A: 255}
	blue := color.RGBA{R: 0, G: 100, B: 200, A: 255}

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if y < 70 {
				img.Set(x, y, orange)
			} else {
				img.Set(x, y, blue)
			}
		}
	}
	saveImage(t, testImagePath, img)

	extractedColor, err := ExtractDominantColor(testImagePath)

	assert.NoError(t, err)
	assert.NotEmpty(t, extractedColor)
	t.Logf("INPUT: 70%% orange + 30%% blue -> OUTPUT: %s", extractedColor)
}

func TestExtractDominantColor_FileNotFound(t *testing.T) {
	_, err := ExtractDominantColor("/nonexistent/path/image.png")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open image")
}

func TestExtractDominantColor_InvalidImage(t *testing.T) {
	testDir := t.TempDir()
	testFilePath := filepath.Join(testDir, "not_an_image.txt")
	os.WriteFile(testFilePath, []byte("not an image"), 0644)

	_, err := ExtractDominantColor(testFilePath)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode image")
}

func TestGetThumbnailColor_EmptyPath(t *testing.T) {
	result := GetThumbnailColor("")
	assert.Equal(t, DefaultColor, result)
}

func TestGetThumbnailColor_InvalidPath(t *testing.T) {
	result := GetThumbnailColor("/invalid/path.png")
	assert.Equal(t, DefaultColor, result)
}

func TestSaturation(t *testing.T) {
	assert.Greater(t, saturation("#FF0000"), 0.9)
	assert.Equal(t, 0.0, saturation("#808080"))
}

func TestIsDark(t *testing.T) {
	assert.True(t, isDark(10, 10, 10))
	assert.False(t, isDark(100, 100, 100))
}

func TestIsLight(t *testing.T) {
	assert.True(t, isLight(240, 240, 240))
	assert.False(t, isLight(100, 100, 100))
}

// Helper functions
func fillImage(img *image.RGBA, c color.RGBA) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, c)
		}
	}
}

func saveImage(t *testing.T, path string, img *image.RGBA) {
	file, err := os.Create(path)
	assert.NoError(t, err)
	defer file.Close()
	png.Encode(file, img)
}
