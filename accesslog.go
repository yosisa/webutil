package webutil

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

type LogWriter interface {
	WriteLog(*AccessLog)
}

type AccessLog struct {
	Time       time.Time
	Elapsed    time.Duration
	RemoteAddr string
	Status     int
	Size       int
	Request    *http.Request
}

func (a *AccessLog) String() string {
	return fmt.Sprintf("%s [%s] %s %s %d %d %v",
		a.RemoteAddr, a.Time.Format("02/Jan/2006:15:04:05 -0700"),
		a.Request.Method, a.Request.RequestURI, a.Status, a.Size, a.Elapsed)
}

func Logger(h http.Handler, logger LogWriter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := WrapResponseWriter(w)
		h.ServeHTTP(lw, r)
		elapsed := time.Since(start)

		remoteAddr := r.Header.Get("X-Forwarded-For")
		if remoteAddr == "" {
			remoteAddr = r.RemoteAddr
		}
		if port := r.Header.Get("X-Forwarded-Port"); port != "" {
			remoteAddr = net.JoinHostPort(remoteAddr, port)
		}

		logger.WriteLog(&AccessLog{
			Time:       start,
			Elapsed:    elapsed,
			RemoteAddr: remoteAddr,
			Status:     lw.Status,
			Size:       lw.Size,
			Request:    r,
		})
	})
}

func NewConsoleLogWriter(w io.Writer) *ConsoleLogWriter {
	return &ConsoleLogWriter{w: w}
}

type ConsoleLogWriter struct {
	w io.Writer
	m sync.Mutex
}

func (w *ConsoleLogWriter) WriteLog(l *AccessLog) {
	w.m.Lock()
	defer w.m.Unlock()
	fmt.Fprintln(w.w, l)
}

func (w *ConsoleLogWriter) Swap(new io.Writer) (old io.Writer) {
	w.m.Lock()
	defer w.m.Unlock()
	old, w.w = w.w, new
	return
}
