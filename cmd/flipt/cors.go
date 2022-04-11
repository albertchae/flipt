package main

import (
	"net/http"
)

func CorsVaryOriginHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &serverWriter{
			w: w,
		}

		next.ServeHTTP(sw, r)
	})
}

type serverWriter struct {
	w           http.ResponseWriter
	wroteHeader bool
}

func (s *serverWriter) Header() http.Header {
	return s.w.Header()
}

func (s *serverWriter) WriteHeader(code int) {
	if s.wroteHeader == false {
		s.Header().Add("Vary", "Origin")
		s.wroteHeader = true
	}
	s.w.WriteHeader(code)
}

func (s *serverWriter) Write(b []byte) (int, error) {
	if s.wroteHeader == false {
		// We hit this case if user never calls WriteHeader (default 200)
		s.w.Header().Add("Vary", "Origin")
		s.wroteHeader = true
	}
	return s.w.Write(b)
}
