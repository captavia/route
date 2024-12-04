package main

import (
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"testing"
)

var pathKeys [1000]string // random /paths/of/parts keys
const partsPerKey = 3     // (e.g. /a/b/c has parts /a, /b, /c)
const bytesPerPart = 10

func init() {
	// path keys
	for i := 0; i < len(pathKeys); i++ {
		var key string
		for j := 0; j < partsPerKey; j++ {
			key += "/"
			part := make([]byte, bytesPerPart)
			if _, err := rand.Read(part); err != nil {
				panic("error generating random byte slice")
			}
			key += string(part)
		}
		pathKeys[i] = key
	}
}

func TestHttpSegmenter1(t *testing.T) {
	var cases = map[string][]string{
		"":              {""},
		"/":             {"/"},
		"//":            {"/", "/"},
		"hello/world":   {"hello", "/world"},
		"/hello/world":  {"/hello", "/world"},
		"/hello/world/": {"/hello", "/world", "/"},
		"/a/b/c":        {"/a", "/b", "/c"},
		"a/":            {"a", "/"},
	}

	for path, segments := range cases {
		var _seg = make([]string, 0, len(segments))
		for p := range PathSegmenterWithDelimiter('/', path) {
			_seg = append(_seg, p)
		}

		assert.Equal(t, segments, _seg)
	}
}

func BenchmarkSegment1(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for p := range PathSegmenterWithDelimiter('/', pathKeys[i%len(pathKeys)]) {
			var _ = p
		}
	}
}
