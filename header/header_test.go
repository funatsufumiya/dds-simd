package header

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadHeader(t *testing.T) {
	var data = [128]byte{'D', 'D', 'S', ' '}
	for i := byte(1); i < 32; i++ {
		data[(i * 4)] = i
	}

	// expected values during parsing
	data[1*4] = 124
	data[2*4] = 7
	data[2*4+1] = 16
	data[19*4] = 32

	expected := &Header{
		TextureHeader: TextureHeader{
			TextureFlags:      4103,
			Height:            3,
			Width:             4,
			PitchOrLinearSize: 5,
			Depth:             6,
			MipMapCount:       7,
		},
		PixelFormatHeader: PixelFormatHeader{
			PixelFlags:  20,
			FourCC:      21,
			RgbBitCount: 22,
			RBitMask:    23,
			GBitMask:    24,
			BBitMask:    25,
			ABitMask:    26,
		},
		Caps: [4]uint32{27, 28, 29, 30},
	}

	rd := bytes.NewReader(data[:])
	h, err := New(rd)
	assert.NoError(t, err)
	assert.Equal(t, expected, h)
}
