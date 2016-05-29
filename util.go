package webutil

import "net/http"

type LoggedResponseWriter struct {
	http.ResponseWriter
	wroteHeader bool
	Status      int
	Size        int
}

func (w *LoggedResponseWriter) WriteHeader(code int) {
	w.Status = code
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *LoggedResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(b)
	w.Size += n
	return n, err
}

func WrapResponseWriter(w http.ResponseWriter) *LoggedResponseWriter {
	if lw, ok := w.(*LoggedResponseWriter); ok {
		return lw
	}
	return &LoggedResponseWriter{ResponseWriter: w}
}
