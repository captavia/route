package route

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkRouterSimple(b *testing.B) {
	mux := NewHttpMux()
	mux.router.Handle("/hello/world/", func(ctx *httpContext) {
		ctx.w.Write([]byte("hello world"))
	})

	req := httptest.NewRequest("GET", "/hello/world/", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkRouterWithParam(b *testing.B) {
	mux := NewHttpMux()
	mux.router.Handle("/hello/world/:id", func(ctx *httpContext) {
		ctx.w.Write([]byte("hello world"))
	})

	req := httptest.NewRequest("GET", "/hello/world/123", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkRouterStd(b *testing.B) {
	router := http.NewServeMux()
	router.HandleFunc("/hello/world/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	req := httptest.NewRequest("GET", "/hello/world/123", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}
