package mimedb

import (
	"slices"
	"testing"
)

func TestTypeByExtension(t *testing.T) {
	tests := []struct {
		ext string
		typ string
	}{
		// https://github.com/golang/go/blob/e7015c9327c4d755651ed3de3fd34fd99a479924/src/mime/type.go#L60-L77
		{".avif", "image/avif"},
		{".css", "text/css; charset=utf-8"},
		{".gif", "image/gif"},
		{".htm", "text/html; charset=utf-8"},
		{".html", "text/html; charset=utf-8"},
		{".jpeg", "image/jpeg"},
		{".jpg", "image/jpeg"},
		{".js", "text/javascript; charset=utf-8"},
		{".json", "application/json"},
		{".mjs", "text/javascript; charset=utf-8"},
		{".pdf", "application/pdf"},
		{".png", "image/png"},
		{".svg", "image/svg+xml"},
		{".wasm", "application/wasm"},
		{".webp", "image/webp"},
		{".xml", "text/xml; charset=utf-8"},
	}

	for _, test := range tests {
		typ := TypeByExtension(test.ext)
		if typ != test.typ {
			t.Errorf("TypeByExtension(%q) = %q, want %q", test.ext, typ, test.typ)
		}
	}
}

func TestExtensionsByType(t *testing.T) {
	tests := []struct {
		typ string
		ext []string
	}{
		// https://github.com/golang/go/blob/e7015c9327c4d755651ed3de3fd34fd99a479924/src/mime/type.go#L60-L77
		{"image/avif", []string{".avif"}},
		{"text/css", []string{".css"}},
		{"image/gif", []string{".gif"}},
		{"text/html", []string{".htm", ".html", ".shtml"}},
		{"image/jpeg", []string{".jpe", ".jpeg", ".jpg"}},
		{"application/javascript", []string{".js", ".mjs"}},
		{"application/json", []string{".json", ".map"}},
		{"application/pdf", []string{".pdf"}},
		{"image/png", []string{".png"}},
		{"image/svg+xml", []string{".svg", ".svgz"}},
		{"application/wasm", []string{".wasm"}},
		{"image/webp", []string{".webp"}},
		{"text/xml", []string{".xml"}},
	}

	for _, test := range tests {
		ext, err := ExtensionsByType(test.typ)
		if err != nil {
			t.Errorf("ExtensionsByType(%q) = %v", test.typ, err)
		}
		if !slices.Equal(ext, test.ext) {
			t.Errorf("ExtensionsByType(%q) = %v, want %v", test.typ, ext, test.ext)
		}
	}
}
