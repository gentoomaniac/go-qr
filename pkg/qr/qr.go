package qr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	baseSize = 21

	Numeric      = 0b0001
	Alphanumeric = 0b0010
	Binary       = 0b0100
	Kanji        = 0b1000
)

func New(version int, data []byte, mode uint8) (error, *QrCode) {
	if version < 1 || version > 40 {
		return errors.New("invalid version"), nil
	}

	// qr code side length, 21 in version 1, increases by 4 each version increment

	code := &QrCode{
		codeData: initialiseCodeData(version),
		data:     data,
		mode:     mode,
		version:  version,
	}

	if err := setMode(code.codeData, mode); err != nil {
		log.Error().Uint8("mode", mode).Msg("")
		return err, nil
	}
	if err := setDataLength(code.codeData, data); err != nil {
		log.Error().Uint8("mode", mode).Msg("")
		return err, nil
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
	addTimingPattern(data)

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

	// always dark bit
	data[len(data)-8][8] = 1
}

func addTimingPattern(data [][]int8) {
	offset := len(data) - 8
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

func setMode(data [][]int8, mode uint8) error {
	if mode == Numeric || mode == Alphanumeric || mode == Binary || mode == Kanji {
		dimension := len(data) - 1
		// FIXME: a Hack to set the right value. Remove when switching back to [][]bool
		data[dimension][dimension] = int8(mode&Kanji) / Kanji
		data[dimension][dimension-1] = int8(mode&Binary) / Binary
		data[dimension-1][dimension] = int8(mode&Alphanumeric) / Alphanumeric
		data[dimension-1][dimension-1] = int8(mode&Numeric) / Numeric
	} else {
		return errors.New("invalid mode specified")
	}

	return nil
}

func setDataLength(data [][]int8, payload []byte) error {
	return nil
}

type QrCode struct {
	codeData [][]int8
	data     []byte
	mode     uint8
	version  int
}

func (q *QrCode) ToString() []string {
	asText := []string{}
	asText = append(asText, strings.Repeat("⬚", len(q.codeData)+6))
	asText = append(asText, fmt.Sprintf("⬚%s⬚", strings.Repeat("🮋", len(q.codeData)+4)))
	asText = append(asText, fmt.Sprintf("⬚%s⬚", strings.Repeat("🮋", len(q.codeData)+4)))
	for _, line := range q.codeData {
		qrLine := ""
		for _, cell := range line {
			switch cell {
			case 1:
				qrLine += "⬚"
			case 0:
				qrLine += "🮋"

			default:
				qrLine += "·"
			}
		}
		asText = append(asText, fmt.Sprintf("⬚🮋🮋%s🮋🮋⬚", qrLine))
	}
	asText = append(asText, fmt.Sprintf("⬚%s⬚", strings.Repeat("🮋", len(q.codeData)+4)))
	asText = append(asText, fmt.Sprintf("⬚%s⬚", strings.Repeat("🮋", len(q.codeData)+4)))
	asText = append(asText, strings.Repeat("⬚", len(q.codeData)+6))
	return asText
}
