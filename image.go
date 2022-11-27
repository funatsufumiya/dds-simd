// Package dds provides a decoder for the DirectDraw surface format, which is compatible with the image package.
package dds

import (
	"dds/decoder"
	"dds/header"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
)

// init registers the decoder for the dds image format
func init() {
	image.RegisterFormat("dds", "DDS ", Decode, DecodeConfig)
}

var ErrUnsupported = errors.New("unsupported texture format")

func DecodeConfig(r io.Reader) (image.Config, error) {
	h, err := header.New().Read(r)
	if err != nil {
		return image.Config{}, err
	}

	// set width and height
	c := image.Config{
		Width:  int(h.Width),
		Height: int(h.Height),
	}

	switch pf, s := h.PixelFlags, h.RgbBitCount; {
	case pf.Is(header.DDPFFourCC):
		fmt.Println(h.FourCC, h.FourCCString)
		c.ColorModel = color.NRGBAModel
	case pf.Has(header.DDPFRGB): // because alpha is implicit
		if s <= 32 {
			c.ColorModel = color.NRGBAModel
		} else {
			c.ColorModel = color.NRGBA64Model
		}
	case pf.Is(header.DDPFYUV):
		c.ColorModel = color.NYCbCrAModel
	case pf.Is(header.DDPFLuminance):
		if s <= 8 {
			c.ColorModel = color.GrayModel
		} else {
			c.ColorModel = color.Gray16Model
		}
	case pf.Is(header.DDPFAlpha):
		if s <= 8 {
			c.ColorModel = color.AlphaModel
		} else {
			c.ColorModel = color.Alpha16Model
		}
	case pf.Is(header.DDPFLuminance | header.DDPFAlphaPixels):
		err = ErrUnsupported
	default:
		err = fmt.Errorf("unrecognized image format: pf.flags: %x", pf)
	}

	return c, err
}

func Decode(r io.Reader) (image.Image, error) {
	h, err := header.New().Read(r)
	if err != nil {
		return nil, err
	}

	d, err := decoder.Find(h)
	if err != nil {
		return nil, err
	}

	return d.Decode(r)
}
