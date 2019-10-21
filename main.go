package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sort"
	"strconv"
)

func main() {
	a := os.Args[1:]
	if len(a) == 1 {
		histogram(a[0])
		return
	}
	if len(a) < 3 || len(a)%2 == 0 {
		panic("not enough arguments")
	}
	m, w, h := readImage(a[len(a)-1])
	a = a[:len(a)-1]
	colors := make([]color.RGBA, len(a))
	for i := range a {
		colors[i] = parseColor(a[i])
	}
	d := image.NewRGBA(image.Rectangle{Max: image.Point{w, h}})
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			c := color.RGBAModel.Convert(m.At(i, j)).(color.RGBA)
			for k := 0; k < len(colors); k += 2 {
				if colors[k] == c {
					c = colors[k+1]
					break
				}
			}
			d.Set(i, j, c)
		}
	}
	out, e := os.Create("out.png")
	fail(e)
	defer out.Close()
	fail(png.Encode(out, d))
}
func parseColor(s string) color.RGBA { // #RRGGBB #RRGGBBAA
	if len(s) == 7 {
		s = s + "FF"
	} else if len(s) != 9 || s[0] != '#' {
		panic("parse color")
	}
	u, e := strconv.ParseUint(s[1:], 16, 32)
	fail(e)
	return color.RGBA{uint8(u & 0xFF000000 >> 24), uint8(u & 0xFF0000 >> 16), uint8(u & 0xFF00 >> 8), uint8(u & 0xFF)}
}
func readImage(file string) (image.Image, int, int) {
	f, e := os.Open(file)
	fail(e)
	defer f.Close()
	m, e := png.Decode(f)
	fail(e)
	return m, m.Bounds().Dx(), m.Bounds().Dy()
}
func histogram(file string) {
	m, w, h := readImage(file)
	colormap := make(map[color.RGBA]uint64)
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			c := color.RGBAModel.Convert(m.At(i, j)).(color.RGBA)
			colormap[c]++
		}
	}
	if len(colormap) == 0 {
		panic("empty image")
	}
	colors := make([]color.RGBA, len(colormap))
	i := 0
	for c := range colormap {
		colors[i] = c
		i++
	}
	sort.Slice(colors, func(i, j int) bool { return colormap[colors[i]] < colormap[colors[j]] })
	for i := 0; i < len(colors); i++ {
		c := colors[i]
		fmt.Printf("#%02x%02x%02x%02x %d\n", c.R, c.G, c.B, c.A, colormap[c])
	}
}
func fail(e error) {
	if e != nil {
		panic(e)
	}
}
