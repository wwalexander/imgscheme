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

func parseChannel(channel string) (uint8, error) {
	ch, err := strconv.ParseUint(channel, 16, 8)
	if err != nil {
		return 0, err
	}
	return uint8(ch), nil
}

func parseTriplet(triplet string) (RGB, error) {
	length := len(triplet)
	if triplet[0] != '#' || length != 7 {
		return RGB{}, errors.New("malformed hex triplet in base color scheme")
	}
	var c RGB
	var err error
	c.R, err = parseChannel(triplet[1:3])
	if err != nil {
		return RGB{}, err
	}
	c.G, err = parseChannel(triplet[3:5])
	if err != nil {
		return RGB{}, err
	}
	c.B, err = parseChannel(triplet[5:7])
	if err != nil {
		return RGB{}, err
	}
	return c, nil
}

func readBase(path string) (color.Palette, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return color.Palette{}, err
	}
	r := bufio.NewReader(file)
	var s color.Palette
	for line, err := r.ReadString('\n'); err == nil; line, err = r.ReadString('\n') {
		line = strings.TrimSuffix(line, "\n")
		c, err := parseTriplet(line)
		if err != nil {
			return color.Palette{}, err
		}
		s = append(s, c)
	}
	return s, nil
}

func readImage(path string) (image.Image, error) {
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

// RGBA returns the color channels in c and a maximal alpha value.
func (c RGB) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	return r, g, b, math.MaxUint16
}

func colors(m image.Image) (colors color.Palette, counts map[color.Color]int) {
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

type colorCount struct {
	color color.Color
	count int
}

// NewScheme creates a color scheme using the colors in m, following the base
// color scheme.
func NewScheme(m image.Image, base color.Palette) color.Palette {
	p, counts := colors(m)
	ccs := make([]colorCount, len(base))
	for c, count := range counts {
		// TODO: Try organizing colors by hue instead of Euclidean
		// distance - maybe use luminance for black/white
		i := base.Index(c)
		if count > ccs[i].count {
			ccs[i].color = c
			ccs[i].count = count
		}
	}
	s := make(color.Palette, len(ccs))
	for i, cc := range ccs {
		if cc.count == 0 {
			s[i] = p.Convert(base[i])
		} else {
			s[i] = cc.color
		}
	}
	return s
}

// SchemeVGA is the standard VGA color scheme.
var SchemeVGA = color.Palette{
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
}

const usage = `usage: imgscheme [-base path] [path]

imgscheme generates a terminal color scheme from the image at the named path. It
attempts to make the generated colors correspond to prominent colors in the
image, and also close to the corresponding colors in a base color scheme. A base
color scheme file can be specified using the -base flag; this file should
contain a newline-separated list of hex triplets. The generated colors will be
output in the same order as the base color scheme. If the -base flag is not set,
imgscheme will use the standard VGA color scheme; the first 8 triplets of the
output will be the 8 normal colors in order (black, red, green, yellow, blue,
magenta, cyan, and white) and the next 8 triplets will be the 8 bold colors in
the same order.`

func main() {
	fbase := flag.String("base", "", "the base color scheme file to use")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	base := SchemeVGA
	if *fbase != "" {
		var err error
		base, err = readBase(*fbase)
		if err != nil {
			log.Fatal(err)
		}
	}
	m, err := readImage(args[0])
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
