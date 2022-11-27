package decoder

import (
	"dds/decoder/dxt"
	"dds/decoder/rgba"
	"dds/header"
	"fmt"
	"image"
	"io"
)

// Decoder is the default interface for actual decoding operations.
type Decoder interface {
	// Decode takes the header-less reader and tries to read an parse the image-data from it.
	Decode(io.Reader) (image.Image, error)
}

// Find takes a parsed header.Header and tries to find a fitting Decoder or returns an error.
func Find(h *header.Header) (d Decoder, err error) {
	switch cc := h.FourCCString(); cc {
	case "":
		if h.PixelFlags != header.AlphaPixels|header.RGB && h.PixelFlags != header.RGB {
			err = fmt.Errorf("unsupported pixel format %x", h.PixelFlags)
		} else {
			d = rgba.New(h)
		}

	case "DXT1", "DXT2", "DXT3", "DXT4", "DXT5":
		d, err = dxt.New(cc, int(h.Width), int(h.Height))

	default:
		err = fmt.Errorf("texture with compression '%v' is unsupported", cc)
	}

	return
}
