package main

// TODO: Add comments to top-level constructs

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
)

type RGB struct {
	R, G, B uint8
}

func (c RGB) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	return r, g, b, math.MaxUint16
}

func (c RGB) String() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

const Size = 16

var (
	// TODO: Add more colorschemes
	SchemeVGA = color.Palette{
		RGB{0x00, 0x00, 0x00},
		RGB{0xAA, 0x00, 0x00},
		RGB{0x00, 0xAA, 0x00},
		RGB{0xAA, 0x55, 0x00},
		RGB{0x00, 0x00, 0xAA},
		RGB{0xAA, 0x00, 0xAA},
		RGB{0xAA, 0xAA, 0xAA},
		RGB{0xAA, 0xAA, 0xAA},
		RGB{0x55, 0x55, 0x55},
		RGB{0xFF, 0x55, 0x55},
		RGB{0x55, 0xFF, 0x55},
		RGB{0xFF, 0xFF, 0x55},
		RGB{0x55, 0x55, 0xFF},
		RGB{0xFF, 0x55, 0xFF},
		RGB{0x55, 0xFF, 0xFF},
		RGB{0xFF, 0xFF, 0xFF},
	}
)

type ColorCount struct {
	Color color.Color
	Count int
}

func ToRGB(c color.Color) RGB {
	r, g, b, _ := c.RGBA()
	return RGB{
		uint8(r >> 8),
		uint8(g >> 8),
		uint8(b >> 8),
	}
}

func NewScheme(m image.Image, base color.Palette) color.Palette {
	bounds := m.Bounds()
	counts := make(map[color.Color]int)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := m.At(x, y)
			count := counts[c]
			counts[c] = count + 1
		}
	}
	var ccs [Size]ColorCount
	for c, count := range counts {
		// TODO: Try organizing colors by hue instead of Euclidean distance
		i := base.Index(c)
		if count > ccs[i].Count {
			ccs[i].Color = c
			ccs[i].Count = count
		}
	}
	s := make(color.Palette, Size)
	for i, cc := range ccs {
		if cc.Count == 0 {
			continue
		}
		s[i] = ToRGB(cc.Color)
	}
	return s
}

func main() {
	// TODO: Add a flag to select the base colorscheme
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [PATH]\n", os.Args[0])
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	file, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	m, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range NewScheme(m, SchemeVGA) {
		fmt.Println(c)
	}
}
