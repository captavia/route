package route

import (
	"iter"
	"strings"
)

type Segmenter func(delimiter rune, path string) iter.Seq[string]

type Handler[T Context] func(ctx T)

type Router[T Context] struct {
	delimiter rune

	children  map[string]*Router[T]
	handler   Handler[T]
	isDynamic bool
	paramName string

	// todo: closure have performance problem
	defaultSegmenter Segmenter
}

type Context interface {
	NotMatch()
	WithParam(map[string]string)
}

func NewRouter[T Context](opts ...RouterOpt[T]) *Router[T] {
	var router = &Router[T]{
		delimiter:        '/',
		defaultSegmenter: PathSegmenterWithDelimiter,
	}

	for _, opt := range opts {
		opt(router)
	}

	return router
}

type RouterOpt[T Context] func(*Router[T])

func WithDelimiter[T Context](d rune) RouterOpt[T] {
	return func(r *Router[T]) {
		r.delimiter = d
	}
}

func WithSegmenter[T Context](fn Segmenter) RouterOpt[T] {
	return func(r *Router[T]) {
		r.defaultSegmenter = fn
	}
}

func (r *Router[T]) Handle(path string, handler Handler[T]) *Router[T] {
	correct := r
	for part := range r.defaultSegmenter(r.delimiter, path) {
		trimmed := strings.TrimPrefix(part, string(r.delimiter))
		isDynamic := strings.HasPrefix(trimmed, ":") || (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(part, "}"))
		var paramName = ""
		if isDynamic {
			if strings.HasPrefix(part, ":") {
				paramName = trimmed[1:] // :id -> id
				part = "*"
			} else {
				paramName = trimmed[1 : len(part)-1] // {id} -> id
				part = "*"
			}
		}
		if correct.children == nil {
			correct.children = make(map[string]*Router[T])
		}
		if _, exists := r.children[part]; !exists {
			correct.children[part] = &Router[T]{
				isDynamic: isDynamic,
				paramName: paramName,
			}
		}
		correct = correct.children[part]
	}
	correct.handler = handler
	return correct
}

func (r *Router[T]) Serve(path string, ctx func() T) {
	var c = ctx()
	correct := r
	var params map[string]string
	for part := range r.defaultSegmenter(r.delimiter, path) {
		switch {
		case correct.children[part] != nil:
			correct = correct.children[part]
		case correct.children["*"] != nil:
			if params == nil {
				params = make(map[string]string)
			}
			correct = correct.children["*"]
			if correct.isDynamic {
				params[correct.paramName] = strings.TrimPrefix(part, string(r.delimiter)) // 提取参数
			}
		default:
			c.NotMatch()
			return
		}
	}
	c.WithParam(params)
	correct.handler(c)
}

func (r *Router[T]) Delete(path string) {
	parent := r
	correct := r
	lastPart := ""
	for part := range r.defaultSegmenter(r.delimiter, path) {
		var _ = part
		switch {
		case correct.children[part] != nil:
			parent = correct
			lastPart = part
			correct = correct.children[part]
		case correct.children["*"] != nil:
			parent = correct
			lastPart = part
			correct = correct.children["*"]
		default:

		}
	}
	delete(parent.children, lastPart)
}
