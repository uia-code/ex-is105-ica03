package moustache

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype/raster"
)

// pt returns the raster.Point corresponding to the pixel
// position (x,y)
func pt(x, y int) fixed.Point26_6 {
	return fixed.Point26_6{fixed.Int26_6(x << 6),
		fixed.Int26_6(y << 6)}
}

// rgba returns and RGBA version of the imaage,
// making a copy only if necessary.
func rgba(m image.Image) *image.RGBA {
	if r, ok := m.(*image.RGBA); ok {
		return r
	}
	b := m.Bounds()
	r := image.NewRGBA(b)
	draw.Draw(r, b, m, image.ZP, draw.Src)
	return r
}

// moustache draws a moustache of the specified size and
// droops onto the image m and returns the result.
// It may overwrite the original.
func Moustache(m image.Image, x, y, size, droopFactor int) image.Image {
	mrgba := rgba(m) // Create specialized RGBA image from m (original image)

	p := raster.NewRGBAPainter(mrgba)
	p.SetColor(color.RGBA{0, 0, 0, 255}) // black?

	w, h := m.Bounds().Dx(), m.Bounds().Dy()
	r := raster.NewRasterizer(w, h)
	var (
		mag   = fixed.Int26_6((10 + size) << 6)
		width = pt(20, 0).Mul(mag)
		mid   = pt(x, y)
		droop = pt(0, droopFactor).Mul(mag)
		left  = mid.Sub(width).Add(droop)
		right = mid.Add(width).Add(droop)
		bow   = pt(0, 5).Mul(mag).Sub(droop)
		curlx = pt(10, 0).Mul(mag)
		curly = pt(0, 2).Mul(mag)
		risex = pt(2, 0).Mul(mag)
		risey = pt(0, 5).Mul(mag)
	)
	r.Start(left)
	r.Add3(
		mid.Sub(curlx).Add(curly),
		mid.Sub(risex).Sub(risey),
		mid,
	)
	r.Add3(
		mid.Add(risex).Sub(risey),
		mid.Add(curlx).Add(curly),
		right,
	)
	r.Add2(
		mid.Add(bow),
		left,
	)
	r.Rasterize(p)

	return mrgba
}
