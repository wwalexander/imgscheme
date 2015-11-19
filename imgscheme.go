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
	"os/signal"
	"sort"
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

// TODO: Implement RGB.String() (convert to a hex string)

const Size = 16

var (
	// TODO: Add more colorschemes
	SchemeVGA = color.Palette{
		RGB{0, 0, 0},
		RGB{170, 0, 0},
		RGB{0, 170, 0},
		RGB{170, 85, 0},
		RGB{0, 0, 170},
		RGB{170, 0, 170},
		RGB{0, 170, 170},
		RGB{170, 170, 170},
		RGB{85, 85, 85},
		RGB{255, 85, 85},
		RGB{85, 255, 85},
		RGB{255, 255, 85},
		RGB{85, 85, 255},
		RGB{255, 85, 255},
		RGB{85, 255, 255},
		RGB{255, 255, 255},
	}
)

type ByCount struct {
	Colors []color.Color
	// TODO: Use an array with an element for every possible RGB value and see if performance improves
	Counts map[color.Color]int
}

func NewByCount() ByCount {
	return ByCount{Counts: make(map[color.Color]int)}
}

func (c *ByCount) Len() int      { return len(c.Colors) }
func (c *ByCount) Swap(i, j int) { c.Colors[i], c.Colors[j] = c.Colors[j], c.Colors[i] }

func (c *ByCount) Less(i, j int) bool {
	icount, iok := c.Counts[c.Colors[i]]
	jcount, jok := c.Counts[c.Colors[j]]
	if iok && jok {
		return icount > jcount
	}
	return iok
}

func (c *ByCount) Insert(co color.Color) {
	count, ok := c.Counts[co]
	if !ok {
		c.Colors = append(c.Colors, co)
	}
	c.Counts[co] = count + 1
	sort.Sort(c)
}

func NewScheme(m image.Image, base color.Palette, c chan os.Signal) color.Palette {
	bounds := m.Bounds()
	var counts [Size]ByCount
	for i := range counts {
		counts[i] = NewByCount()
	}
ImageLoop:
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			select {
			case <-c:
				break ImageLoop
			default:
			}
			c := m.At(x, y)
			i := base.Index(c)
			counts[i].Insert(c)
		}
	}
	s := make(color.Palette, Size)
	for i := range counts {
		if len(counts[i].Colors) == 0 {
			// TODO: find a better behavior for unset colors
			// Perhaps search through all other ByCounts for the
			// nearest color - should counts be generated in one go
			// and then iterated over for each place?
			s[i] = RGB{0, 0, 0}
			continue
		}
		r, g, b, _ := counts[i].Colors[0].RGBA()
		c := RGB{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
		}
		s[i] = c
	}
	return s
}

func main() {
	// TODO: Add a flag to select the base colorscheme
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
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
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	fmt.Println(NewScheme(m, SchemeVGA, sigc))
}
