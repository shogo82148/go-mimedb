package mimedb

import (
	"errors"
	"fmt"
	"mime"
	"slices"
	"sort"
	"strings"
)

//go:generate go run internal/cmd/update/main.go

// AddExtensionType is not supported by this package.
// MIE type database is read-only.
func AddExtensionType(ext, typ string) error {
	if !strings.HasPrefix(ext, ".") {
		return fmt.Errorf("mimedb: extension %q missing leading dot", ext)
	}
	return errors.New("mimedb: adding extension is not supported")
}

// TypeByExtension returns the MIME type associated with the file extension ext.
// The extension ext should begin with a leading dot, as in ".html".
// When ext has no associated type, TypeByExtension returns "".
//
// Extensions are looked up first case-sensitively, then case-insensitively.
//
// Text types have the charset parameter set to "utf-8" by default.
func TypeByExtension(ext string) string {
	// Case-sensitive lookup.
	if v, ok := mimeTypes[ext]; ok {
		return v
	}

	// Case-insensitive lookup.
	// Optimistically assume a short ASCII extension and be
	// allocation-free in that case.
	var buf [10]byte
	lower := buf[:0]
	const utf8RuneSelf = 0x80 // from utf8 package, but not importing it.
	for i := 0; i < len(ext); i++ {
		c := ext[i]
		if c >= utf8RuneSelf {
			// Slow path.
			return mimeTypesLower[strings.ToLower(ext)]
		}
		if 'A' <= c && c <= 'Z' {
			lower = append(lower, c+('a'-'A'))
		} else {
			lower = append(lower, c)
		}
	}
	return mimeTypesLower[string(lower)]
}

// ExtensionsByType returns the extensions known to be associated with the MIME
// type typ. The returned extensions will each begin with a leading dot, as in
// ".html". When typ has no associated extensions, ExtensionsByType returns an
// nil slice.
func ExtensionsByType(typ string) ([]string, error) {
	justType, _, err := mime.ParseMediaType(typ)
	if err != nil {
		return nil, err
	}

	s, ok := extensions[justType]
	if !ok {
		return nil, nil
	}
	ret := slices.Clone(s)
	sort.Strings(ret)
	return ret, nil
}
