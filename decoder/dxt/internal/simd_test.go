package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightedSIMD(t *testing.T) {
	w0, w1 := float32(2), float32(1)
	v0 := []float32{10, 20, 30, 40}
	v1 := []float32{1, 2, 3, 4}
	out := make([]float32, 4)
	WeightedSIMD(w0, w1, v0, v1, out)
	// out[i] = (w0*v0[i] + w1*v1[i]) / (w0 + w1)
	expected := []float32{
		(2*10 + 1*1) / 3,
		(2*20 + 1*2) / 3,
		(2*30 + 1*3) / 3,
		(2*40 + 1*4) / 3,
	}
	for i := range out {
		assert.InDelta(t, expected[i], out[i], 0.01)
	}
}

func TestExtractIndexSIMD(t *testing.T) {
	in := []byte{0b00100001, 0b10000100}
	indices := []byte{0, 1, 2, 3}
	lengths := []byte{4, 4, 4, 4}
	out := make([]byte, 4)
	ExtractIndexSIMD(in, indices, lengths, out)
	expected := []byte{0b0001, 0b0010, 0b0100, 0b1000}
	assert.Equal(t, expected, out)
}
