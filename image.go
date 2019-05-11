/*
Copyright 2017 Luke Granger-Brown

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package dds provides a decoder for the DirectDraw surface format,
// which is compatible with the standard image package.
//
// It should normally be used by importing it with a blank name, which
// will cause it to register itself with the image package:
//  import _ "github.com/lukegb/dds"
package dds

import (
	"fmt"
	"image"
	"image/color"
	"io"
)

func init() {
	image.RegisterFormat("dds", "DDS ", Decode, DecodeConfig)
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	h, err := readHeader(r)
	if err != nil {
		return image.Config{}, err
	}

	// set width and height
	c := image.Config{
		Width:  int(h.width),
		Height: int(h.height),
	}

	pf := h.pixelFormat
	hasAlpha := (pf.flags&pfAlphaPixels == pfAlphaPixels) || (pf.flags&pfAlpha == pfAlpha)
	hasRGB := (pf.flags&pfFourCC == pfFourCC) || (pf.flags&pfRGB == pfRGB)
	hasYUV := (pf.flags&pfYUV == pfYUV)
	hasLuminance := (pf.flags&pfLuminance == pfLuminance)
	switch {
	case hasRGB && pf.rgbBitCount == 32:
		c.ColorModel = color.NRGBAModel
	case hasRGB && pf.rgbBitCount == 64:
		c.ColorModel = color.NRGBA64Model
	case hasYUV && pf.rgbBitCount == 24:
		c.ColorModel = color.YCbCrModel
	case hasLuminance && pf.rgbBitCount == 8:
		c.ColorModel = color.GrayModel
	case hasLuminance && pf.rgbBitCount == 16:
		c.ColorModel = color.Gray16Model
	case hasAlpha && pf.rgbBitCount == 8:
		c.ColorModel = color.AlphaModel
	case hasAlpha && pf.rgbBitCount == 16:
		c.ColorModel = color.AlphaModel
	default:
		return image.Config{}, fmt.Errorf("unrecognized image format: hasAlpha: %v, hasRGB: %v, hasYUV: %v, hasLuminance: %v, pf.flags: %x", hasAlpha, hasRGB, hasYUV, hasLuminance, pf.flags)
	}

	return c, nil
}

func Decode(r io.Reader) (image.Image, error) {
	h, err := readHeader(r)
	if err != nil {
		return nil, err
	}

	switch fourccToString(h.pixelFormat.fourCC) {
	case "\x00\x00\x00\x00":
		if h.pixelFormat.flags != pfAlphaPixels|pfRGB && h.pixelFormat.flags != pfRGB {
			return nil, fmt.Errorf("unsupported pixel format %x", h.pixelFormat.flags)
		}
		return decodeRGBA(r, h)
	case "DXT1":
		return decodeDXT1(r, h)
	case "DXT3":
		return decodeDXT3(r, h)
	case "DXT5":
		return decodeDXT5(r, h)
	default:
		return nil, fmt.Errorf("image data is compressed with %v; this compression is unsupported", fourccToString(h.pixelFormat.fourCC))
	}
}
