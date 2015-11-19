package main

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

// RGB is a 24-bit RGB color.
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

// Bases is a map of strings and their corresponding base color schemes.
var Bases = map[string]color.Palette{
	"vga": {
		RGB{0x00, 0x00, 0x00},
		RGB{0xaa, 0x00, 0x00},
		RGB{0x00, 0xaa, 0x00},
		RGB{0xaa, 0x55, 0x00},
		RGB{0x00, 0x00, 0xaa},
		RGB{0xaa, 0x00, 0xaa},
		RGB{0xaa, 0xaa, 0xaa},
		RGB{0xaa, 0xaa, 0xaa},
		RGB{0x55, 0x55, 0x55},
		RGB{0xff, 0x55, 0x55},
		RGB{0x55, 0xff, 0x55},
		RGB{0xff, 0xff, 0x55},
		RGB{0x55, 0x55, 0xff},
		RGB{0xff, 0x55, 0xff},
		RGB{0x55, 0xff, 0xff},
		RGB{0xff, 0xff, 0xff},
	},
	"xterm": {
		RGB{0x00, 0x00, 0x00},
		RGB{0xcd, 0x00, 0x00},
		RGB{0x00, 0xcd, 0x00},
		RGB{0xcd, 0xcd, 0x00},
		RGB{0x00, 0x00, 0xee},
		RGB{0xcd, 0x00, 0xcd},
		RGB{0x00, 0xcd, 0xcd},
		RGB{0xe5, 0xe5, 0xe5},
		RGB{0x7f, 0x7f, 0x7f},
		RGB{0xff, 0x00, 0x00},
		RGB{0x00, 0xff, 0x00},
		RGB{0xff, 0xff, 0x00},
		RGB{0x5c, 0x5c, 0xff},
		RGB{0xff, 0x00, 0xff},
		RGB{0x00, 0xff, 0xff},
		RGB{0xff, 0xff, 0xff},
	},
}

// Counts returns a map of the colors in m and the frequency of each color.
func Counts(m image.Image) map[color.Color]int {
	bounds := m.Bounds()
	counts := make(map[color.Color]int)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := m.At(x, y)
			count := counts[c]
			counts[c] = count + 1
		}
	}
	return counts
}

// A ColorCount contains a color and the number of times it appears in an image.
type ColorCount struct {
	Color color.Color
	Count int
}

// BestMatch returns the most frequent color in m for each place.
func BestMatch(m image.Image, base color.Palette) color.Palette {
	counts := Counts(m)
	ccs := make([]ColorCount, len(base))
	for c, count := range counts {
		// TODO: Try organizing colors by hue instead of Euclidean
		// distance - maybe use luminance for black/white
		i := base.Index(c)
		if count > ccs[i].Count {
			ccs[i].Color = c
			ccs[i].Count = count
		}
	}
	s := make(color.Palette, len(ccs))
	for i, cc := range ccs {
		s[i] = cc.Color
	}
	return s
}

// A ColorDistance contains a color and its distance from another color.
type ColorDistance struct {
	Color    color.Color
	Distance float64
}

// Distance returns the Euclidean distance between a and b.
func Distance(a, b color.Color) float64 {
	ra, ga, ba, _ := a.RGBA()
	rb, gb, bb, _ := b.RGBA()
	return math.Sqrt(math.Pow(float64(ra)-float64(rb), 2) +
		math.Pow(float64(ga)-float64(gb), 2) +
		math.Pow(float64(ba)-float64(bb), 2))
}

// NearestMatch returns the closest color in m to each color in base.
func NearestMatch(m image.Image, base color.Palette) color.Palette {
	cds := make([]ColorDistance, len(base))
	for i := range cds {
		cds[i].Distance = math.MaxFloat64
	}
	bounds := m.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := m.At(x, y)
			for i, cbase := range base {
				if distance := Distance(cbase, c); distance <= cds[i].Distance {
					cds[i].Color = c
					cds[i].Distance = distance
				}
			}
		}
	}
	cs := make(color.Palette, len(cds))
	for i, cd := range cds {
		cs[i] = cd.Color
	}
	return cs
}

// NewScheme creates a color scheme using the colors in m, following the color
// base color scheme.
func NewScheme(m image.Image, base color.Palette) color.Palette {
	s := BestMatch(m, base)
	var missing []int
	var missingColors color.Palette
	for place, c := range s {
		if c == nil {
			missing = append(missing, place)
			missingColors = append(missingColors, base[place])
		}
	}
	nearest := NearestMatch(m, missingColors)
	i := 0
	for _, match := range nearest {
		s[missing[i]] = match
		i++
	}
	return s
}

func main() {
	fbaseDesc := "the base color scheme file to use, or a built in color scheme:\n"
	for name := range Bases {
		fbaseDesc += fmt.Sprintf("\t\t%s\n", name)
	}
	fbase := flag.String("base", "vga", fbaseDesc)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [PATH]\n", os.Args[0])
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		flag.PrintDefaults()
		os.Exit(1)
	}
	base, ok := Bases[*fbase]
	if !ok {
		log.Fatal("base color scheme does not exist")
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
	for _, c := range NewScheme(m, base) {
		r, g, b, _ := c.RGBA()
		fmt.Printf("#%02x%02x%02x\n",
			uint8(r>>8),
			uint8(g>>8),
			uint8(b>>8))
	}
}
