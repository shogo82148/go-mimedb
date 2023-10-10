package mimedb

import "testing"

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
