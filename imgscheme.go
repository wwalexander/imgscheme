package main

import (
	"bufio"
	"errors"
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
	"strconv"
	"strings"
)

// ParseChannel parses a channel from a hex triplet into a uint8.
func ParseChannel(channel string) (uint8, error) {
	ch, err := strconv.ParseUint(channel, 16, 8)
	if err != nil {
		return 0, err
	}
	return uint8(ch), nil
}

// ParseTriplet parses a hex triplet into an RGB.
func ParseTriplet(triplet string) (RGB, error) {
	length := len(triplet)
	if triplet[0] != '#' || length != 7 {
		log.Println(triplet[0], length)
		fmt.Print(triplet)
		return RGB{}, errors.New("malformed hex triplet in base color scheme")
	}
	var c RGB
	var err error
	c.R, err = ParseChannel(triplet[1:3])
	if err != nil {
		return RGB{}, err
	}
	c.G, err = ParseChannel(triplet[3:5])
	if err != nil {
		return RGB{}, err
	}
	c.B, err = ParseChannel(triplet[5:7])
	if err != nil {
		return RGB{}, err
	}
	return c, nil
}

// ReadBase reads a base color scheme from the file located at path.
func ReadBase(path string) (color.Palette, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(file)
	var s color.Palette
	for line, err := r.ReadString('\n'); err == nil; line, err = r.ReadString('\n') {
		line = strings.TrimSuffix(line, "\n")
		c, err := ParseTriplet(line)
		if err != nil {
			return nil, err
		}
		s = append(s, c)
	}
	return s, nil
}

// ReadImage reads an image.Image from the file located at path.
func ReadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	m, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// RGB is a 24-bit RGB color.
type RGB struct {
	R, G, B uint8
}

// RGBA satisfies the color.Color interface for RGB.
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
		RGB{0x00, 0xaa, 0xaa},
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

// Colors returns a color.Palette containing the colors in m, and a map from
// each color to the number of times it appears in m.
func Colors(m image.Image) (colors color.Palette, counts map[color.Color]int) {
	bounds := m.Bounds()
	colors = make(color.Palette, 0, bounds.Max.X*bounds.Max.Y)
	counts = make(map[color.Color]int)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := m.At(x, y)
			colors = append(colors, c)
			count := counts[c]
			counts[c] = count + 1
		}
	}
	return colors, counts
}

// A ColorCount contains a color and the number of times it appears in an image.
type ColorCount struct {
	Color color.Color
	Count int
}

// NewScheme creates a color scheme using the colors in m, following the base
// color scheme
func NewScheme(m image.Image, base color.Palette) color.Palette {
	p, counts := Colors(m)
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
		if cc.Count == 0 {
			s[i] = p.Convert(base[i])
		} else {
			s[i] = cc.Color
		}
	}
	return s
}

func main() {
	fbaseDesc := "the base color scheme file to use, or a built-in color scheme:\n"
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
		var err error
		base, err = ReadBase(*fbase)
		if err != nil {
			log.Fatal(err)
		}
	}
	m, err := ReadImage(args[0])
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
