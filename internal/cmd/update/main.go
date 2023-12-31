package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

const releaseURL = "https://github.com/jshttp/mime-db/releases/latest"
const databaseURL = "https://raw.githubusercontent.com/jshttp/mime-db/%s/db.json"

func main() {
	// get the latest release
	ctx := context.Background()
	latest, err := getLatestRelease(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("latest release:", latest)

	// get the database
	db, err := getDatabase(ctx, latest)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("database size:", len(db))

	// for compatibility with Go's mime package.
	db["text/html"].Charset = "UTF-8"
	db["application/json"].Charset = ""

	// format extensions
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf, "// Code generated by internal/cmd/update/main.go; DO NOT EDIT.")
	fmt.Fprintln(buf, "")
	fmt.Fprintln(buf, "package mimedb")
	fmt.Fprintln(buf, "")
	fmt.Fprintf(buf, "// The data is generated from [mime-db] version %s\n", latest)
	fmt.Fprintf(buf, "// [mime-db]: https://github.com/jshttp/mime-db\n")
	fmt.Fprintln(buf, "var extensions = map[string][]string{")
	types := make([]string, 0, len(db))
	for k := range db {
		types = append(types, k)
	}
	slices.Sort(types)

	for _, k := range types {
		v := db[k]
		if len(v.Extensions) == 0 {
			continue
		} else {
			fmt.Fprintf(buf, "\t%q: {", k)
			slices.Sort(v.Extensions)
			for i, ext := range v.Extensions {
				if i > 0 {
					fmt.Fprint(buf, ", ")
				}
				fmt.Fprintf(buf, "%q", "."+ext)
			}
			fmt.Fprintln(buf, "},")
		}
	}
	fmt.Fprintln(buf, "}")
	fmt.Fprintln(buf, "")

	// format mimeTypes and mimeTypesLower
	mimeTypes := make(map[string][]string, len(db))
	mimeTypesLower := make(map[string][]string, len(db))
	for k, v := range db {
		typ := k
		if v.Charset == "UTF-8" {
			typ += "; charset=utf-8"
		}
		for _, ext := range v.Extensions {
			lower := strings.ToLower(ext)
			mimeTypes[ext] = append(mimeTypes[ext], typ)
			mimeTypesLower[lower] = append(mimeTypesLower[lower], typ)
		}
	}

	exts := make([]string, 0, len(mimeTypes))
	for k := range mimeTypes {
		exts = append(exts, k)
	}
	slices.Sort(exts)

	fmt.Fprintln(buf, "var mimeTypes = map[string]string{")
	for _, ext := range exts {
		if len(mimeTypes[ext]) == 0 {
			continue
		}
		fmt.Fprintf(buf, "\t%q: %q,\n", "."+ext, selectType(ext, mimeTypes[ext]))
	}
	fmt.Fprintln(buf, "}")
	fmt.Fprintln(buf, "")

	fmt.Fprintln(buf, "var mimeTypesLower = map[string]string{")
	for _, ext := range exts {
		if len(mimeTypes[ext]) == 0 {
			continue
		}
		fmt.Fprintf(buf, "\t%q: %q,\n", "."+ext, selectType(ext, mimeTypes[ext]))
	}
	fmt.Fprintln(buf, "}")

	source, err := format.Source(buf.Bytes())
	if err != nil {
		log.Println(buf.String())
		log.Fatal(err)
	}

	if err := os.WriteFile("mimedb_generated.go", source, 0644); err != nil {
		log.Fatal(err)
	}
}

func selectType(ext string, types []string) string {
	// special cases
	switch ext {
	case "mp3":
		return "audio/mpeg"
	case "xml":
		return "text/xml; charset=utf-8"
	case "wav":
		return "audio/wav"
	case "js":
		return "text/javascript; charset=utf-8"
	case "mjs":
		return "text/javascript; charset=utf-8"
	}

	if len(types) == 1 {
		return types[0]
	}

	standardTypes := make([]string, 0, len(types))
	for _, typ := range types {
		if strings.Contains(typ, "/x-") {
			// it's a vendor-specific type, skip it
			continue
		}
		standardTypes = append(standardTypes, typ)
	}
	if len(standardTypes) == 1 {
		return standardTypes[0]
	}

	specialTypes := make([]string, 0, len(types))
	for _, typ := range types {
		if strings.HasSuffix(typ, "+json") || strings.HasSuffix(typ, "+xml") {
			specialTypes = append(specialTypes, typ)
		}
	}
	if len(specialTypes) == 1 {
		return specialTypes[0]
	}

	log.Println("warning: multiple types:", ext, types)

	if len(specialTypes) > 1 {
		slices.Sort(specialTypes)
		return specialTypes[0]
	}

	if len(standardTypes) > 1 {
		slices.Sort(standardTypes)
		return standardTypes[0]
	}

	slices.Sort(types)
	return types[0]
}

type release struct {
	TagName string `json:"tag_name"`
}

func getLatestRelease(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, releaseURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var r release
	if err := json.Unmarshal(data, &r); err != nil {
		return "", err
	}

	return r.TagName, nil
}

type entry struct {
	Source       string   `json:"source"`
	Extensions   []string `json:"extensions"`
	Compressible bool     `json:"compressible"`
	Charset      string   `json:"charset"`
}

func getDatabase(ctx context.Context, tag string) (map[string]*entry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(databaseURL, tag), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var v map[string]*entry
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return v, nil
}
