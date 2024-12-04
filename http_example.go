package route

import (
	"net/http"
	"sync"
)

type HttpMux struct {
	router *Router[*httpContext]
	pool   *sync.Pool
}

func NewHttpMux() *HttpMux {
	return &HttpMux{
		router: NewRouter[*httpContext](),
		pool: &sync.Pool{
			New: func() interface{} {
				return new(httpContext)
			},
		},
	}
}

func (c *HttpMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := c.pool.Get().(*httpContext)
	c.router.Serve(r.RequestURI, func() *httpContext {
		ctx.w = w
		ctx.r = r
		return ctx
	})
	c.pool.Put(ctx)
}

type httpContext struct {
	w http.ResponseWriter
	r *http.Request
}

func (c *httpContext) NotMatch() {
	http.NotFound(c.w, c.r)
}

func (c *httpContext) WithParam(param map[string]string) {

}
