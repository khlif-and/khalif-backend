package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"khalif-backend/pkg/utils"
)

func main() {
	testDir := "./test_images"
	os.MkdirAll(testDir, 0755)

	// Test 1: Pure Red Image
	fmt.Println("=== Testing Color Extraction ===")
	fmt.Println()

	redPath := createTestImage(testDir, "red.png", color.RGBA{R: 220, G: 50, B: 50, A: 255})
	redColor, err := utils.ExtractDominantColorFromPath(redPath)
	fmt.Printf("Red Image (RGB: 220,50,50)\n")
	fmt.Printf("  Path: %s\n", redPath)
	fmt.Printf("  Extracted: %s\n", redColor)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Println()

	// Test 2: Blue Image
	bluePath := createTestImage(testDir, "blue.png", color.RGBA{R: 50, G: 100, B: 220, A: 255})
	blueColor, err := utils.ExtractDominantColorFromPath(bluePath)
	fmt.Printf("Blue Image (RGB: 50,100,220)\n")
	fmt.Printf("  Path: %s\n", bluePath)
	fmt.Printf("  Extracted: %s\n", blueColor)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Println()

	// Test 3: Green Image
	greenPath := createTestImage(testDir, "green.png", color.RGBA{R: 50, G: 200, B: 80, A: 255})
	greenColor, err := utils.ExtractDominantColorFromPath(greenPath)
	fmt.Printf("Green Image (RGB: 50,200,80)\n")
	fmt.Printf("  Path: %s\n", greenPath)
	fmt.Printf("  Extracted: %s\n", greenColor)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Println()

	// Test 4: Orange Image
	orangePath := createTestImage(testDir, "orange.png", color.RGBA{R: 255, G: 165, B: 0, A: 255})
	orangeColor, err := utils.ExtractDominantColorFromPath(orangePath)
	fmt.Printf("Orange Image (RGB: 255,165,0)\n")
	fmt.Printf("  Path: %s\n", orangePath)
	fmt.Printf("  Extracted: %s\n", orangeColor)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Println()

	// Test 5: Purple Image
	purplePath := createTestImage(testDir, "purple.png", color.RGBA{R: 150, G: 50, B: 200, A: 255})
	purpleColor, err := utils.ExtractDominantColorFromPath(purplePath)
	fmt.Printf("Purple Image (RGB: 150,50,200)\n")
	fmt.Printf("  Path: %s\n", purplePath)
	fmt.Printf("  Extracted: %s\n", purpleColor)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Println()

	// Test 6: Error case - file not found
	fmt.Println("Error Case - File Not Found:")
	_, err = utils.ExtractDominantColorFromPath("/nonexistent/path.png")
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Println()

	fmt.Println("=== Test Complete ===")

	// Cleanup
	os.RemoveAll(testDir)
}

func createTestImage(dir, name string, c color.RGBA) string {
	path := filepath.Join(dir, name)
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, c)
		}
	}

	file, _ := os.Create(path)
	png.Encode(file, img)
	file.Close()

	return path
}
