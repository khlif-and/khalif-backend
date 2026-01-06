package utils

import (
	"crypto/md5"
	"fmt"
	"hash/fnv"
	"math"
	"strings"
	"unicode"
)

func GenerateInitials(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "AD"
	}

	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "AD"
	}

	if len(parts) >= 2 {
		return strings.ToUpper(string(rune(parts[0][0])) + string(rune(parts[1][0])))
	}

	var initials []rune
	runes := []rune(parts[0])
	for i, r := range runes {
		if i == 0 || unicode.IsUpper(r) {
			initials = append(initials, r)
		}
	}

	if len(initials) >= 2 {
		return strings.ToUpper(string(initials[:2]))
	}

	if len(runes) >= 2 {
		return strings.ToUpper(string(runes[:2]))
	}

	return strings.ToUpper(string(runes[0]))
}

func GenerateAmbientColor(input string) string {
	h := fnv.New32a()
	h.Write([]byte(input))
	hashVal := h.Sum32()

	hue := int(hashVal % 360)

	saturation := 70

	lightness := 20 + int(hashVal%11)

	return hslToHex(hue, saturation, lightness)
}

func hslToHex(h, s, l int) string {
	sFloat := float64(s) / 100
	lFloat := float64(l) / 100

	c := (1 - math.Abs(2*lFloat-1)) * sFloat
	x := c * (1 - math.Abs(math.Mod(float64(h)/60, 2)-1))
	m := lFloat - c/2

	var r, g, b float64

	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	case h < 360:
		r, g, b = c, 0, x
	}

	rInt := int((r + m) * 255)
	gInt := int((g + m) * 255)
	bInt := int((b + m) * 255)

	return fmt.Sprintf("#%02x%02x%02x", rInt, gInt, bInt)
}

func GenerateGravatarURL(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	hash := md5.Sum([]byte(email))
	return fmt.Sprintf("https://www.gravatar.com/avatar/%x?d=identicon", hash)
}
