// Package dxt holds decoder for all kind DXT based dds encodings
// Specification: https://www.khronos.org/registry/OpenGL/extensions/EXT/EXT_texture_compression_s3tc.txt.
package dxt

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"

	. "github.com/funatsufumiya/dds-simd/decoder/dxt/internal"
)

type (
	Decoder struct {
		strategy
		bounds image.Point
	}

	strategy interface {
		New(bounds image.Rectangle) draw.Image
		BlockSize() byte
		DecodeBlock(buffer []byte)
		Pixel(index byte) color.Color
		PixelBlock() [16]color.Color
	}
)

func New(fourCC string, width, height int) (*Decoder, error) {
	decoder := &Decoder{bounds: image.Pt(width, height)}

	switch fourCC {
	case "DXT1":
		decoder.strategy = new(dxt1)
	case "DXT3":
		decoder.strategy = new(dxt3)
	case "DXT5":
		decoder.strategy = new(dxt5)
	default:
		return nil, fmt.Errorf("DXT type '%s' not supported", fourCC)
	}

	return decoder, nil
}

// Decode decodes from r and returns a new image.Image as before.
func (d *Decoder) Decode(r io.Reader) (image.Image, error) {
       bounds := image.Rectangle{Max: d.bounds}
       img := d.New(bounds)
       if bounds.Empty() {
	       return img, nil
       }
       err := d.DecodeTo(r, img)
       if err != nil {
	       return nil, err
       }
       return img, nil
}

// DecodeTo decodes from r and writes the result into dst (must be at least bounds size).
// This allows memory reuse and avoids unnecessary allocations.
func (d *Decoder) DecodeTo(r io.Reader, dst draw.Image) error {
       bounds := image.Rectangle{Max: d.bounds}
       if bounds.Empty() {
	       return nil
       }
       rd := NewReader(r, d.BlockSize())
       for h := 0; h < d.bounds.Y; h += 4 {
	       for w := 0; w < d.bounds.X; w += 4 {
		       buffer, err := rd.Read()
		       if err != nil {
			       return err
		       }
		       d.DecodeBlock(buffer)
		       blockColors := d.PixelBlock()
		       for y := 0; y < 4; y++ {
			       for x := 0; x < 4; x++ {
				       pxIndex := x + y*4
				       if w+x < d.bounds.X && h+y < d.bounds.Y {
					       dst.Set(w+x, h+y, blockColors[pxIndex])
				       }
			       }
		       }
	       }
       }
       return nil
}
