package dxt

import (
	"image/color"

	. "github.com/funatsufumiya/dds-simd/decoder/dxt/internal"
)

type dxt1 struct {
	ColorDecoder
}

func (*dxt1) BlockSize() byte {
	return 8
}

func (d *dxt1) DecodeBlock(buffer []byte) {
	d.BlockColor(buffer[0:8:8])
}

func (d *dxt1) Pixel(index byte) color.Color {
	return d.PixelColor(index)
}

// PixelBlock returns a 4x4 block of colors (16 pixels) for the current block.
func (d *dxt1) PixelBlock() [16]color.Color {
	var out [16]color.Color
	for i := 0; i < 16; i++ {
		out[i] = d.Pixel(byte(i))
	}
	return out
}
