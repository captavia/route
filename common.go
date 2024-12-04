package route

import (
	"iter"
	"strings"
)

func PathSegmenterWithDelimiter(s rune, path string) iter.Seq[string] {
	return func(yield func(string) bool) {
		var start = 0
		if len(path) == 0 {
			yield(path)
			return
		}
		end := strings.IndexRune(path[start+1:], s)
		for ; end != -1; end = strings.IndexRune(path[start+1:], s) {
			end++
			yield(path[start : start+end])
			start += end
		}
		yield(path[start:])
	}
}

func PathSegmenter(s rune, path string) iter.Seq[string] {
	return func(yield func(string) bool) {
		var start, end = 0, 0
		if len(path) == 0 {
			yield(path)
			return
		}

		if end = strings.IndexRune(path, s); end == 0 {
			start++
			end = strings.IndexRune(path[start:], s)
		}

		for ; end != -1; end = strings.IndexRune(path[start:], s) {
			yield(path[start : start+end])
			start += end + 1
		}
		yield(path[start:])
	}
}
