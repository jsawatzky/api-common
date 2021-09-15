package internal

import "net/http"

type ResponseRecorder struct {
	http.ResponseWriter
	staus        int
	responseSize int
}

func RecordResponse(w http.ResponseWriter) *ResponseRecorder {
	if rr, ok := w.(*ResponseRecorder); ok {
		return rr
	}

	return &ResponseRecorder{
		ResponseWriter: w,
		staus:          http.StatusOK,
	}
}

func (r *ResponseRecorder) Status() int {
	return r.staus
}

func (r *ResponseRecorder) ResponseSize() int {
	return r.responseSize
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseSize = size
	return size, err
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.staus = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
