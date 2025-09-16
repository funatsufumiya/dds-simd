package internal

import (
	"github.com/pehringer/simd"
)

// WeightedSIMD applies the weighted average to two byte slices using SIMD.

func WeightedSIMD(w0, w1 float32, v0, v1, out []float32) int {
	// out[i] = (w0*v0[i] + w1*v1[i]) / (w0 + w1)
	n := len(v0)
	if n > len(v1) {
		n = len(v1)
	}
	if n > len(out) {
		n = len(out)
	}
	if n == 0 {
		return 0
	}
	w0s := make([]float32, n)
	w1s := make([]float32, n)
	for i := range w0s {
		w0s[i] = w0
		w1s[i] = w1
	}
	simd.MulFloat32(w0s, v0, out)
	simd.MulFloat32(w1s, v1, w1s)
	simd.AddFloat32(out, w1s, out)
	denom := w0 + w1
	for i := 0; i < n; i++ {
		out[i] = out[i] / denom
	}
	return n
}

// ExtractIndexSIMD extracts indices from a byte slice using SIMD (for batch processing).
// This is a stub for illustration; actual SIMD bit extraction is non-trivial and may need custom implementation.
func ExtractIndexSIMD(bytes []byte, indices, lengths []byte, out []byte) int {
	// Not implemented: SIMD bit extraction is complex and not directly supported by simd package.
	// Fallback to scalar for now.
	n := len(indices)
	if n > len(lengths) {
		n = len(lengths)
	}
	if n > len(out) {
		n = len(out)
	}
	for i := 0; i < n; i++ {
		out[i] = ExtractIndex(bytes, indices[i], lengths[i])
	}
	return n
}
