package tinyclient

import (
	"net/http"
	"time"
)

// Response struct holds response values of executed request.
type Response struct {
	Request  *Request
	Response *http.Response

	body       []byte
	size       int64
	receivedAt time.Time
}