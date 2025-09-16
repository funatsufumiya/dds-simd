package dxt

import (
	"image/color"

	. "github.com/robroyd/dds/decoder/dxt/internal"
)

type dxt5 struct {
	ColorDecoder
	alphaValues  [8]byte
	alphaIndices []byte
}

func (*dxt5) BlockSize() byte {
	return 16
}

func (d *dxt5) DecodeBlock(buffer []byte) {
	d.alphaValues = d.interpolateAlphaValues(buffer[0:2:2])
	d.alphaIndices = buffer[2:8:8]
	d.BlockColor(buffer[8:16:16])
}

func (d *dxt5) Pixel(index byte) color.Color {
	alphaIndex := ExtractIndex(d.alphaIndices, index, 3)
	alpha := d.alphaValues[alphaIndex]
	return d.PixelAlpha(index, alpha)
}

// PixelBlock returns a 4x4 block of colors (16 pixels) for the current block.
func (d *dxt5) PixelBlock() [16]color.Color {
	var out [16]color.Color
	for i := 0; i < 16; i++ {
		out[i] = d.Pixel(byte(i))
	}
	return out
}

func (d *dxt5) interpolateAlphaValues(a0 []byte) (av [8]byte) {
	av[0] = a0[0]
	av[1] = a0[1]

	v0 := float32(a0[0])
	v1 := float32(a0[1])
	out := make([]float32, 6)
	if a0[0] <= a0[1] {
		w0 := []float32{4, 3, 2, 1}
		w1 := []float32{1, 2, 3, 4}
		v0s := make([]float32, 4)
		v1s := make([]float32, 4)
		for i := 0; i < 4; i++ {
			v0s[i] = v0
			v1s[i] = v1
		}
		WeightedSIMD(w0[0], w1[0], v0s, v1s, out[0:1])
		WeightedSIMD(w0[1], w1[1], v0s, v1s, out[1:2])
		WeightedSIMD(w0[2], w1[2], v0s, v1s, out[2:3])
		WeightedSIMD(w0[3], w1[3], v0s, v1s, out[3:4])
		av[2] = byte(out[0] + 0.5)
		av[3] = byte(out[1] + 0.5)
		av[4] = byte(out[2] + 0.5)
		av[5] = byte(out[3] + 0.5)
		av[6] = 0
		av[7] = 255
	} else {
		w0 := []float32{6, 5, 4, 3, 2, 1}
		w1 := []float32{1, 2, 3, 4, 5, 6}
		v0s := make([]float32, 6)
		v1s := make([]float32, 6)
		for i := 0; i < 6; i++ {
			v0s[i] = v0
			v1s[i] = v1
		}
		WeightedSIMD(w0[0], w1[0], v0s, v1s, out[0:1])
		WeightedSIMD(w0[1], w1[1], v0s, v1s, out[1:2])
		WeightedSIMD(w0[2], w1[2], v0s, v1s, out[2:3])
		WeightedSIMD(w0[3], w1[3], v0s, v1s, out[3:4])
		WeightedSIMD(w0[4], w1[4], v0s, v1s, out[4:5])
		WeightedSIMD(w0[5], w1[5], v0s, v1s, out[5:6])
		av[2] = byte(out[0] + 0.5)
		av[3] = byte(out[1] + 0.5)
		av[4] = byte(out[2] + 0.5)
		av[5] = byte(out[3] + 0.5)
		av[6] = byte(out[4] + 0.5)
		av[7] = byte(out[5] + 0.5)
	}
	return
}
