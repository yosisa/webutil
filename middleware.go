package webutil

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
)

func Recoverer(h http.Handler, out io.Writer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 65535)
				n := runtime.Stack(buf, false)
				fmt.Fprintf(out, "panic: %+v\n%s", err, buf[:n])
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
