package dds

import (
	"fmt"
	"image"
	"io"
	"math"
)

// DXTn decoding is based on https://github.com/kchapelier/decode-dxt.
// Documentation https://www.khronos.org/registry/OpenGL/extensions/EXT/EXT_texture_compression_s3tc.txt.

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
			// BGRA to RGBA re-order.
			for i := 0; i < len(p); i += 4 {
				p[i+0], p[i+2] = p[i+2], p[i+0]
			}
		}
	case pfRGB:
		b := make([]byte, 3*w)
		for y := 0; y != h; y++ {
			if _, err := io.ReadFull(r, b); err != nil {
				return nil, err
			}
			p := rgba.Pix[y*rgba.Stride : y*rgba.Stride+w*4]
			// BGRA to RGBA re-order.
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
	err := checkDivisibilityBy4(width, height)
	if err != nil {
		return nil, err
	}

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

func decodeDXT3(r io.Reader, header header) (image.Image, error) {
	width := int(header.width)
	height := int(header.height)
	err := checkDivisibilityBy4(width, height)
	if err != nil {
		return nil, err
	}

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
			alphaValues := []uint16{
				getUint16(buf[offset+6 : offset+6+2]),
				getUint16(buf[offset+4 : offset+4+2]),
				getUint16(buf[offset+2 : offset+2+2]),
				getUint16(buf[offset : offset+2]),
			}
			colorValues := interpolateColorValues(getUint16(buf[offset+8:offset+8+2]), getUint16(buf[offset+10:offset+10+2]), true)
			colorIndices := getUint32(buf[offset+12 : offset+12+4])

			for y := 0; y < 4; y++ {
				for x := 0; x < 4; x++ {
					pixelIndex := (3 - x) + (y * 4)
					rgbaIndex := (h*4+3-y)*width*4 + (w*4+x)*4
					colorIndex := (colorIndices >> uint((2 * (15 - pixelIndex)))) & 0x03
					rgba.Pix[rgbaIndex] = colorValues[colorIndex*4]
					rgba.Pix[rgbaIndex+1] = colorValues[colorIndex*4+1]
					rgba.Pix[rgbaIndex+2] = colorValues[colorIndex*4+2]
					rgba.Pix[rgbaIndex+3] = getAlphaValue(alphaValues, pixelIndex)
				}
			}

			offset += 16
		}
	}

	return rgba, nil
}

func decodeDXT5(r io.Reader, header header) (image.Image, error) {
	width := int(header.width)
	height := int(header.height)
	err := checkDivisibilityBy4(width, height)
	if err != nil {
		return nil, err
	}

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
			alphaValues := interpolateAlphaValues(uint8(buf[offset]), uint8(buf[offset+1]))
			alphaIndices := []uint16{
				getUint16(buf[offset+6 : offset+6+2]),
				getUint16(buf[offset+4 : offset+4+2]),
				getUint16(buf[offset+2 : offset+2+2]),
			}

			colorValues := interpolateColorValues(getUint16(buf[offset+8:offset+8+2]), getUint16(buf[offset+10:offset+10+2]), true)
			colorIndices := getUint32(buf[offset+12 : offset+12+4])

			for y := 0; y < 4; y++ {
				for x := 0; x < 4; x++ {
					pixelIndex := (3 - x) + (y * 4)
					rgbaIndex := (h*4+3-y)*width*4 + (w*4+x)*4
					colorIndex := (colorIndices >> uint((2 * (15 - pixelIndex)))) & 0x03
					rgba.Pix[rgbaIndex] = colorValues[colorIndex*4]
					rgba.Pix[rgbaIndex+1] = colorValues[colorIndex*4+1]
					rgba.Pix[rgbaIndex+2] = colorValues[colorIndex*4+2]
					rgba.Pix[rgbaIndex+3] = alphaValues[getAlphaIndex(alphaIndices, pixelIndex)]
				}
			}

			offset += 16
		}
	}

	return rgba, nil
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

func convert565ByteToRgb(b uint16) []uint8 {
	return []uint8{
		uint8(math.Round(float64(b>>11&31) * (255 / 31))),
		uint8(math.Round((float64(b>>5&63) * (255 / 63)))),
		uint8(math.Round(float64(b&31) * (255 / 31))),
	}
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

func interpolateAlphaValues(a0, a1 uint8) (alphaValues []uint8) {
	alphaValues = append(alphaValues, a0)
	alphaValues = append(alphaValues, a1)

	if a0 > a1 {
		alphaValues = append(alphaValues,
			uint8(math.Floor((6*float64(a0)+1*float64(a1))/7)),
			uint8(math.Floor((5*float64(a0)+2*float64(a1))/7)),
			uint8(math.Floor((4*float64(a0)+3*float64(a1))/7)),
			uint8(math.Floor((3*float64(a0)+4*float64(a1))/7)),
			uint8(math.Floor((2*float64(a0)+5*float64(a1))/7)),
			uint8(math.Floor((1*float64(a0)+6*float64(a1))/7)),
		)
	} else {
		alphaValues = append(alphaValues,
			uint8(math.Floor((4*float64(a0)+1*float64(a1))/5)),
			uint8(math.Floor((3*float64(a0)+2*float64(a1))/5)),
			uint8(math.Floor((2*float64(a0)+3*float64(a1))/5)),
			uint8(math.Floor((1*float64(a0)+4*float64(a1))/5)),
			0,
			255,
		)
	}

	return alphaValues
}

func getAlphaValue(alphaValue []uint16, pixelIndex int) uint8 {
	return extractBitsFromUin16Array(alphaValue, (4*(15-pixelIndex)), 4) * 17
}

func getAlphaIndex(alphaIndices []uint16, pixelIndex int) uint8 {
	return extractBitsFromUin16Array(alphaIndices, (3 * (15 - pixelIndex)), 3)
}

func extractBitsFromUin16Array(array []uint16, shift, length int) uint8 {
	height := len(array)
	heightm1 := height - 1
	width := 16
	rowS := (shift / width) | 0
	rowE := ((shift + length - 1) / width) | 0
	var result uint8

	if rowS == rowE {
		// all the requested bits are contained in a single uint16
		shiftS := uint(shift % width)
		result = uint8(array[heightm1-rowS]>>shiftS) & uint8(math.Pow(2, float64(length))-1)
	} else {
		// the requested bits are contained in two continuous uint16
		shiftS := uint(shift % width)
		shiftE := uint(width) - shiftS
		result = uint8(array[heightm1-rowS]>>shiftS) & uint8(math.Pow(2, float64(length))-1)
		result += uint8(array[heightm1-rowE]) & uint8(math.Pow(2, float64(length)-float64(shiftE))-1) << shiftE
	}

	return result
}

func checkDivisibilityBy4(w, h int) error {
	if w%4 != 0 || h%4 != 0 {
		return fmt.Errorf("DXT compressed image width and height must be multiple of four: %vx%v", w, h)
	}
	return nil
}
