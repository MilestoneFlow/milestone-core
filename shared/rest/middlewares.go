package rest

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type logEntry struct {
	RequestID string              `json:"request_id"`
	Headers   map[string][]string `json:"headers"`
	Body      string              `json:"body"`
	IP        string              `json:"ip"`
	Caller    string              `json:"caller"`
	User      *interface{}        `json:"user"`
	Timestamp time.Time           `json:"timestamp"`
}

type RequestLogger struct {
	handler http.Handler
}

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return &RequestLogger{handler: next}
}

func NewRequestLogger(handler http.Handler) *RequestLogger {
	return &RequestLogger{handler: handler}
}

func (rl *RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log request headers
	headers := r.Header

	// Read request body
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// Log request ID
	requestID := middleware.GetReqID(r.Context())
	if requestID == "" {
		requestID = "unknown"
	}

	// Get client IP
	ip := r.RemoteAddr
	if ip == "" {
		ip = "unknown"
	}

	userData := r.Context().Value("user")

	// Create log entry
	entry := logEntry{
		RequestID: requestID,
		Headers:   headers,
		Body:      string(bodyBytes),
		IP:        ip,
		Caller:    r.URL.Path,
		User:      &userData,
		Timestamp: time.Now(),
	}

	// Marshal log entry to JSON
	logEntryBytes, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
	} else {
		yellow := "\033[33m"
		reset := "\033[0m"
		log.Printf("%s%s%s", yellow, string(logEntryBytes), reset)
	}

	// Continue with the request
	rl.handler.ServeHTTP(w, r)
}
