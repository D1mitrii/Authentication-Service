package middlewares

import (
	"fmt"
	"net/http"
)

type StatusResponseWriter struct {
	http.ResponseWriter
	Status int
}

func NewStatusResponseWriter(w http.ResponseWriter) *StatusResponseWriter {
	return &StatusResponseWriter{
		ResponseWriter: w,
		Status:         http.StatusOK,
	}
}

func (r *StatusResponseWriter) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *StatusResponseWriter) StatusString() string {
	return fmt.Sprintf("%d", r.Status)
}
