// Package dds provides a decoder for the DirectDraw surface format, which is compatible with the image package.
package dds

import (
	"dds/decoder"
	"dds/header"
	"fmt"
	"image"
	"image/color"
	"io"
)

// init registers the decoder for the dds image format
func init() {
	image.RegisterFormat("dds", "DDS ", Decode, DecodeConfig)
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	h, err := header.New(r)
	if err != nil {
		return image.Config{}, err
	}

	// set width and height
	c := image.Config{
		Width:  int(h.Width),
		Height: int(h.Height),
	}

	pf := h.PixelFlags
	hasAlpha := (pf&header.AlphaPixels == header.AlphaPixels) || (pf&header.Alpha == header.Alpha)
	hasRGB := (pf&header.FourCC == header.FourCC) || (pf&header.RGB == header.RGB)
	hasYUV := pf&header.YUV == header.YUV
	hasLuminance := pf&header.Luminance == header.Luminance

	switch {
	case hasLuminance && h.RgbBitCount == 8:
		c.ColorModel = color.GrayModel
	case hasAlpha && h.RgbBitCount == 8:
		c.ColorModel = color.AlphaModel
	case hasLuminance && h.RgbBitCount == 16:
		c.ColorModel = color.Gray16Model
	case hasAlpha && h.RgbBitCount == 16:
		c.ColorModel = color.AlphaModel
	case hasYUV && h.RgbBitCount == 24:
		c.ColorModel = color.YCbCrModel
	case hasRGB && h.RgbBitCount == 32:
		c.ColorModel = color.NRGBAModel
	case hasRGB && h.RgbBitCount == 64:
		c.ColorModel = color.NRGBA64Model
	default:
		return image.Config{}, fmt.Errorf("unrecognized image format: hasAlpha: %v, hasRGB: %v, hasYUV: %v, hasLuminance: %v, pf.flags: %x", hasAlpha, hasRGB, hasYUV, hasLuminance, pf)
	}

	return c, nil
}

func Decode(r io.Reader) (image.Image, error) {
	h, err := header.New(r)
	if err != nil {
		return nil, err
	}

	d, err := decoder.Find(h)
	if err != nil {
		return nil, err
	}

	return d.Decode(r)
}
