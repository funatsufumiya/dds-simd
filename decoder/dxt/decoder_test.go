package dxt

import (
	"bytes"
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecoderBatchDXT1(t *testing.T) {
	// 4x4 DXT1 block (dummy data, not a real image)
	data := []byte{
		0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00, // color block
	}
	d, err := New("DXT1", 4, 4)
	assert.NoError(t, err)
	img, err := d.Decode(bytes.NewReader(data))
	assert.NoError(t, err)
	assert.Equal(t, image.Rect(0, 0, 4, 4), img.Bounds())
}

func TestDecoderBatchDXT5(t *testing.T) {
	// 4x4 DXT5 block (dummy data, not a real image)
	data := []byte{
		0x00, 0xFF, // alpha0, alpha1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // alpha indices
		0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00, // color block
	}
	d, err := New("DXT5", 4, 4)
	assert.NoError(t, err)
	img, err := d.Decode(bytes.NewReader(data))
	assert.NoError(t, err)
	assert.Equal(t, image.Rect(0, 0, 4, 4), img.Bounds())
}
