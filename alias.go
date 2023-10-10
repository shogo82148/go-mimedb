package mimedb

import (
	"mime"
)

// WordEncoder is an alias of [mime.WordEncoder].
type WordEncoder = mime.WordEncoder

// WordDecoder is an alias of [mime.WordDecoder].
type WordDecoder = mime.WordDecoder

// FormatMediaType is an alias of [mime.FormatMediaType].
func FormatMediaType(t string, param map[string]string) string {
	return mime.FormatMediaType(t, param)
}

// ParseMediaType is an alias of [mime.ParseMediaType].
func ParseMediaType(v string) (mediatype string, params map[string]string, err error) {
	return mime.ParseMediaType(v)
}
