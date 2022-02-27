package drawing

import (
	"fmt"
	"math/rand"
	"os"
	"unicode"
)

var charset []rune

// GetCharset returns the static set of the characters will appear. It returns the
// set of runes used for drawing the screen.
func GetCharset() []rune {
	if len(charset) == 0 {
		calcCharset(true)
	}

	return charset
}

func calcCharset(refresh bool) {
	if refresh || len(charset) == 0 {
		excluded := configureExclusions()
		charset = charsetFrom(Tables, excluded)
	}
}

// DumpCharset is intended to be used for debugging and will dump the charset
// into a text file specified with the path.
func DumpCharset(path string) {
	err := os.WriteFile(path, []byte(string(GetCharset())), 0644)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func SelectRandomN(n int) []rune {
	calcCharset(false)
	var result []rune
	for i := 0; i < n; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return result
}

var greekSlice = unicode.RangeTable{
	R16: unicode.Greek.R16[7:9],
}

var latinSlice = unicode.RangeTable{
	R16: unicode.Latin.R16[:2],
}

var numberSlice = unicode.RangeTable{
	R16: unicode.Number.R16[:],
}

var Tables = []*unicode.RangeTable{
	&greekSlice,
	&latinSlice,
	&numberSlice,
	unicode.Runic,
}

// Used to build a unicode.RangeTable of character ranges excluded from the
// Tables variable specified above
func configureExclusions() unicode.RangeTable {
	var excludedRanges16 []unicode.Range16
	var excludedRanges32 []unicode.Range32

	// Trimming out number sets that don't look good in the terminal
	excludedRanges16 = append(excludedRanges16, unicode.Number.R16[1:5]...)
	excludedRanges16 = append(excludedRanges16, unicode.Number.R16[9:10]...)
	excludedRanges16 = append(excludedRanges16, unicode.Number.R16[14])
	excludedRanges16 = append(excludedRanges16, unicode.Number.R16[18:]...)

	// Trimming out runic chars that don't render
	excludedRanges16 = append(excludedRanges16, unicode.Runic.R16[1])

	return unicode.RangeTable{R16: excludedRanges16, R32: excludedRanges32}
}

func runesFromR16(r16 unicode.Range16, excluded []unicode.Range16) []rune {
	var result []rune
	for i := r16.Lo; i <= r16.Hi; i = i + r16.Stride {
		if isExcluded16(i, excluded) {
			continue
		}
		result = append(result, rune(i))
	}

	return result
}

func runesFromR32(r32 unicode.Range32, excluded []unicode.Range32) []rune {
	var result []rune
	for i := r32.Lo; i <= r32.Hi; i += r32.Stride {
		if isExcluded32(i, excluded) {
			continue
		}
		result = append(result, rune(i))
	}

	return result
}

// getRunes converts a *unicode.RangeTable of ranges into a slice of runes, with
// range exclusions specified with a unicode.RangeTable. Implemented such that if
// *ranges == excluded then then the returned slice would be empty.
func getRunes(ranges *unicode.RangeTable, excluded unicode.RangeTable) []rune {
	var allRunes16 []rune
	for _, runeRange16 := range ranges.R16 {
		allRunes16 = append(allRunes16, runesFromR16(runeRange16, excluded.R16)...)
	}

	var allRunes32 []rune
	for _, runeRange32 := range ranges.R32 {
		allRunes32 = append(allRunes32, runesFromR32(runeRange32, excluded.R32)...)
	}

	return append(allRunes16, allRunes32...)
}

// charsetFrom is just a wrapper around getRunes that allows a slice of table pointers to be
// supplied with the exclusions.
func charsetFrom(tables []*unicode.RangeTable, excluded unicode.RangeTable) []rune {
	var charset []rune
	for _, table := range tables {
		charset = append(charset, getRunes(table, excluded)...)
	}

	return charset
}

// isExcluded16 will check if the character at code point i is in the sixteen bit
// portion of the exclusions
func isExcluded16(i uint16, excluded []unicode.Range16) bool {
	for _, r16 := range excluded {
		if (i-r16.Lo)%r16.Stride == 0 && i >= r16.Lo && i <= r16.Hi {
			return true
		}
	}

	return false
}

// isExcluded32 is the 32-bit equivalent of isExcluded16
func isExcluded32(i uint32, excluded []unicode.Range32) bool {
	for _, r32 := range excluded {
		if (i-r32.Lo)%r32.Stride == 0 && i >= r32.Lo && i <= r32.Hi {
			return true
		}
	}
	return false
}
