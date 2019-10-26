package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	a := os.Args[1:]
	grey := uint8(0)
	var colors []color.RGBA
	if len(a) == 1 {
		histogram(a[0])
		return
	} else if len(a) == 2 && a[0] == "h2" {
		h2(a[1])
		return
	} else if len(a) == 2 {
		grey = parseGrey(a[0])
		if grey == 0 {
			panic("grey must be > 0")
		}
	} else if len(a) > 2 && strings.HasSuffix(a[1], "#") {
		reshape(a)
		return
	} else if len(a) < 3 || len(a)%2 == 0 {
		panic("not enough arguments")
	} else {
		a = a[:len(a)-1]
		colors = make([]color.RGBA, len(a))
		for i := range a {
			colors[i] = parseColor(a[i])
		}
	}
	m, w, h := readImage(a[len(a)-1])
	d := image.NewRGBA(image.Rectangle{Max: image.Point{w, h}})
	var c color.Color
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			if grey == 0 {
				c = color.RGBAModel.Convert(m.At(i, j)).(color.RGBA)
				for k := 0; k < len(colors); k += 2 {
					if colors[k] == c {
						c = colors[k+1]
						break
					}
				}
			} else {
				rgba := color.RGBAModel.Convert(m.At(i, j)).(color.RGBA)
				if rgba.R <= grey {
					c = color.Black
				} else {
					c = color.White
				}
			}
			d.Set(i, j, c)
		}
	}
	writeImage(d)
}
func h2(s string) {
	m, w, h := readImage(s)
	d := image.NewRGBA(image.Rectangle{Max: image.Point{w, 2 * h}})
	draw.Draw(d, d.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			d.Set(i, 2*j, m.At(i, j))
		}
	}
	writeImage(d)
}
func parseGrey(s string) uint8 { // XX
	if len(s) != 2 {
		panic("parse grey")
	}
	u, e := strconv.ParseUint(s, 16, 32)
	fail(e)
	return uint8(u)
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
func writeImage(m image.Image) {
	out, e := os.Create("out.png")
	fail(e)
	defer out.Close()
	fail(png.Encode(out, m))
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
func reshape(a []string) {
	var rows, cols int
	var e error
	rows, e = strconv.Atoi(a[0])
	fail(e)
	cols, e = strconv.Atoi(strings.TrimSuffix(a[1], "#"))
	fail(e)
	a = a[2:]
	d := image.NewRGBA(image.Rectangle{})
	n := 0
	for i := 0; i < rows; i++ {
		dr := image.NewRGBA(image.Rectangle{})
		for j := 0; j < cols; j++ {
			m, _, _ := readImage(a[n%len(a)])
			dr = cat(dr, m, true)
			n++
		}
		d = cat(d, dr, false)
	}
	writeImage(d)
}
func cat(x, y image.Image, hor bool) *image.RGBA {
	var r image.Rectangle
	r.Max = image.Point{y.Bounds().Dx(), x.Bounds().Dy() + y.Bounds().Dy()}
	if hor {
		r.Max = image.Point{x.Bounds().Dx() + y.Bounds().Dx(), y.Bounds().Dy()}
	}
	m := image.NewRGBA(r)
	draw.Draw(m, x.Bounds(), x, image.ZP, draw.Src)
	r = y.Bounds().Add(image.Point{0, x.Bounds().Dy()})
	if hor {
		r = y.Bounds().Add(image.Point{x.Bounds().Dx(), 0})
	}
	draw.Draw(m, r, y, image.ZP, draw.Src)
	return m
}
func fail(e error) {
	if e != nil {
		panic(e)
	}
}
