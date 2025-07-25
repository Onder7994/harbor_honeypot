package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var staticExt = map[string]struct{}{
	".js": {}, ".css": {}, ".png": {}, ".jpg": {}, ".jpeg": {}, ".gif": {}, ".svg": {},
	".ico": {}, ".map": {}, ".json": {}, ".woff": {}, ".woff2": {}, ".ttf": {}, ".eot": {},
}

var logFile *os.File

type LogEntry struct {
	Time        string              `json:"time"`
	Method      string              `json:"method"`
	Path        string              `json:"path"`
	RemoteAddr  string              `json:"remote_addr"`
	UserAgent   string              `json:"user_agent"`
	Status      int                 `json:"status"`
	Fingerprint string              `json:"fingerprint,omitempty"`
	Query       map[string][]string `json:"query_params"`
	Form        map[string][]string `json:"body,omitempty"`
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	if lrw.wroteHeader {
		return
	}
	lrw.statusCode = code
	lrw.wroteHeader = true
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if !lrw.wroteHeader {
		lrw.WriteHeader(http.StatusOK)
	}
	return lrw.ResponseWriter.Write(b)
}

func InitLogger() {
	log_file_path := os.Getenv("APP_LOG_FILE_PATH")
	if log_file_path == "" {
		log_file_path = "/var/log/honeypot/harbor_honeypot.json"
	}

	var err error
	logFile, err = os.OpenFile(log_file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

func isStaticPath(p string) bool {
	ext := strings.ToLower(path.Ext(p))
	_, ok := staticExt[ext]
	return ok
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var formData map[string][]string
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err == nil {
				formData = r.PostForm
			}
		}
		fingerprint := ""
		if cookie, err := r.Cookie("fingerprint"); err == nil {
			fingerprint = cookie.Value
		}

		lrw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lrw, r)
		status := lrw.statusCode
		if status == 0 {
			status = http.StatusOK
		}

		if (r.Method == http.MethodGet || r.Method == http.MethodHead) &&
			(status == http.StatusNotModified || isStaticPath(r.URL.Path)) {
			return
		}

		entry := LogEntry{
			Time:        time.Now().Format(time.RFC3339),
			Method:      r.Method,
			Path:        r.URL.Path,
			RemoteAddr:  r.RemoteAddr,
			UserAgent:   r.UserAgent(),
			Status:      lrw.statusCode,
			Fingerprint: fingerprint,
			Query:       r.URL.Query(),
			Form:        formData,
		}
		jsonEntry, err := json.Marshal(entry)
		if err == nil {
			logFile.Write(jsonEntry)
			logFile.Write([]byte("\n"))
		}
	})
}
