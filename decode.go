package dds

import (
	"fmt"
	"image"
	"io"
	"math"
)

func decodeRGBA(r io.Reader, header header) (image.Image, error) {
	w := int(header.width)
	h := int(header.height)

	rgba := image.NewNRGBA(image.Rect(0, 0, w, h))
	if w == 0 || h == 0 {
		return rgba, nil
	}

	switch header.pixelFormat.flags {
	case pfAlphaPixels | pfRGB:
		for y := 0; y != h; y++ {
			p := rgba.Pix[y*rgba.Stride : y*rgba.Stride+w*4]
			if _, err := io.ReadFull(r, p); err != nil {
				return nil, err
			}
		}
	case pfRGB:
		b := make([]byte, 3*w)
		for y := 0; y != h; y++ {
			if _, err := io.ReadFull(r, b); err != nil {
				return nil, err
			}
			p := rgba.Pix[y*rgba.Stride : y*rgba.Stride+w*4]
			for i, j := 0, 0; i < len(p); i, j = i+4, j+3 {
				p[i+0] = b[j+2]
				p[i+1] = b[j+1]
				p[i+2] = b[j+0]
				p[i+3] = 0xFF
			}
		}
	}

	return rgba, nil
}

func decodeDXT1(r io.Reader, header header) (image.Image, error) {
	width := int(header.width)
	height := int(header.height)
	width4 := (width / 4) | 0
	height4 := (height / 4) | 0
	offset := 0

	rgba := image.NewNRGBA(image.Rect(0, 0, width, height))
	if width == 0 || height == 0 {
		return rgba, nil
	}

	buf := make([]byte, header.pitchOrLinearSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("reading image: %v", err)
	}

	for h := 0; h < height4; h++ {
		for w := 0; w < width4; w++ {
			colorValues := interpolateColorValues(getUint16(buf[offset:offset+2]), getUint16(buf[offset+2:offset+2+2]), true)
			colorIndices := getUint32(buf[offset+4 : offset+4+4])

			for y := 0; y < 4; y++ {
				for x := 0; x < 4; x++ {
					pixelIndex := (3 - x) + (y * 4)
					rgbaIndex := (h*4+3-y)*width*4 + (w*4+x)*4
					colorIndex := (colorIndices >> uint((2 * (15 - pixelIndex)))) & 0x03
					rgba.Pix[rgbaIndex] = colorValues[colorIndex*4]
					rgba.Pix[rgbaIndex+1] = colorValues[colorIndex*4+1]
					rgba.Pix[rgbaIndex+2] = colorValues[colorIndex*4+2]
					rgba.Pix[rgbaIndex+3] = colorValues[colorIndex*4+3]
				}
			}

			offset += 8
		}
	}

	return rgba, nil
}

func interpolateColorValues(v0, v1 uint16, isDxt1 bool) (colorValues []uint8) {
	c0 := convert565ByteToRgb(v0)
	c1 := convert565ByteToRgb(v1)

	colorValues = append(colorValues, c0...)
	colorValues = append(colorValues, 255)
	colorValues = append(colorValues, c1...)
	colorValues = append(colorValues, 255)

	if isDxt1 && v0 <= v1 {
		colorValues = append(colorValues,
			[]uint8{
				uint8(math.Round((float64(c0[0]) + float64(c1[0])) / 2)),
				uint8(math.Round((float64(c0[1]) + float64(c1[1])) / 2)),
				uint8(math.Round((float64(c0[2]) + float64(c1[2])) / 2)),
				255,

				0,
				0,
				0,
				0,
			}...,
		)
	} else {
		colorValues = append(colorValues,
			[]uint8{
				uint8(math.Round(2*float64(c0[0])+float64(c1[0])) / 3),
				uint8(math.Round(2*float64(c0[1])+float64(c1[1])) / 3),
				uint8(math.Round(2*float64(c0[2])+float64(c1[2])) / 3),
				255,

				uint8(math.Round(2*float64(c1[0])+float64(c0[0])) / 3),
				uint8(math.Round(2*float64(c1[1])+float64(c0[1])) / 3),
				uint8(math.Round(2*float64(c1[2])+float64(c0[2])) / 3),
				255,
			}...,
		)
	}

	return colorValues
}

func convert565ByteToRgb(b uint16) []uint8 {
	return []uint8{
		uint8(math.Round(float64(b>>11&31) * (255 / 31))),
		uint8(math.Round((float64(b>>5&63) * (255 / 63)))),
		uint8(math.Round(float64(b&31) * (255 / 31))),
	}
}

func getUint16(buf []byte) (n uint16) {
	n |= uint16(buf[0])
	n |= uint16(buf[1]) << 8
	return
}

func getUint32(buf []byte) (n uint32) {
	n |= uint32(buf[0])
	n |= uint32(buf[1]) << 8
	n |= uint32(buf[2]) << 16
	n |= uint32(buf[3]) << 24
	return
}
