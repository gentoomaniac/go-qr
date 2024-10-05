package qr

import (
	"errors"
)

const (
	baseSize = 21
)

func New(version int, data string) (error, *QrCode) {
	if version < 1 || version > 40 {
		return errors.New("invalid version"), nil
	}

	// qr code side length, 21 in version 1, increases by 4 each version increment

	code := &QrCode{
		codeData: initialiseCodeData(version),
		data:     data,
		version:  version,
	}

	return nil, code
}

func initialiseCodeData(version int) [][]int8 {
	dimension := 21 + (version-1)*4
	data := make([][]int8, dimension)
	for index := range data {
		data[index] = make([]int8, dimension)
		for cI := range data[index] {
			data[index][cI] = -1
		}
	}

	addFinderPattern(data)

	return data
}

func addFinderPattern(data [][]int8) {
	finderPattern := [8][8]int8{
		{1, 1, 1, 1, 1, 1, 1, 0},
		{1, 0, 0, 0, 0, 0, 1, 0},
		{1, 0, 1, 1, 1, 0, 1, 0},
		{1, 0, 1, 1, 1, 0, 1, 0},
		{1, 0, 1, 1, 1, 0, 1, 0},
		{1, 0, 0, 0, 0, 0, 1, 0},
		{1, 1, 1, 1, 1, 1, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
	}
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			data[x][y] = finderPattern[x][y]
		}
	}

	offset := len(data) - 8
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			data[y+offset][x] = finderPattern[x][7-y]
		}
	}

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			data[y][x+offset] = finderPattern[7-x][y]
		}
	}

	alignmentPattern := [5][5]int8{
		{1, 1, 1, 1, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 1, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 1, 1, 1, 1},
	}
	aOffset := len(data) - len(alignmentPattern) - 4
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			data[y+aOffset][x+aOffset] = alignmentPattern[y][x]
		}
	}

	// timing pattern
	patternBit := int8(1)
	for x := 8; x < offset; x++ {
		data[6][x] = patternBit
		if patternBit == 1 {
			patternBit = 0
		} else {
			patternBit = 1
		}
	}

	// timing pattern
	patternBit = int8(1)
	for y := 8; y < offset; y++ {
		data[y][6] = patternBit
		if patternBit == 1 {
			patternBit = 0
		} else {
			patternBit = 1
		}
	}
}

type QrCode struct {
	codeData [][]int8
	data     string
	version  int
}

func (q *QrCode) ToString() []string {
	asText := []string{}
	for _, line := range q.codeData {
		qrLine := ""
		for _, cell := range line {
			switch cell {
			case 1:
				qrLine += "🮋"
			case 0:
				qrLine += " "

			default:
				qrLine += "·"
			}
		}
		asText = append(asText, qrLine)
	}
	return asText
}
