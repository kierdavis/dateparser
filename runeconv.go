package dateparser

import (
    "unicode/utf8"
)

func size(runes []rune) (size int) {
	for _, rune := range runes {
		size += utf8.RuneLen(rune)
	}
	return size
}

func encodeTo(bytes []byte, runes []rune) {
	pos := 0
	for _, rune := range runes {
		count := utf8.EncodeRune(bytes[pos:], rune)
		pos += count
	}
}

func encode(runes []rune) (bytes []byte) {
	bytes = make([]byte, size(runes))
	encodeTo(bytes, runes)
	return bytes
}

func decode(bytes []byte) (runes []rune) {
	pos := 0
	for pos < len(bytes) {
		rune, size := utf8.DecodeRune(bytes[pos:])
		runes = append(runes, rune)
		pos += size
	}
	
	return runes
}

