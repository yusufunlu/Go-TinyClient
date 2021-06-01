package tinyclient

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	client *Client
	URL    string
	Method string
	//bodyBytes is not exposed, it is for holding Body interface internally as bytes and reading the body without spoiling HttpRequest.GetBody
	bodyBytes []byte
	Body      interface{}
	Cookies   []*http.Cookie
	//HttpRequest is exposed because of missing cases of this client wrapper and so a professional user can handle for this edge
	HttpRequest *http.Request
	Headers     map[string]string
	QueryParams map[string]string
	FormData    url.Values
	SentAt      time.Time
	useSSL      bool
}

func (r *Request) SetBody(body interface{}) *Request {
	r.Body = body
	return r
}

func (r *Request) SetMethod(method string) *Request {
	r.Method = method
	return r
}

func (r *Request) SetURL(url string) *Request {
	r.URL = url
	return r
}

func (r *Request) SetHeader(header, value string) *Request {
	r.Headers[header] = value
	return r
}

// SetHeaders sets request headers
func (r *Request) SetHeaders(headers map[string]string) *Request {
	if len(headers) > 0 {
		for k, v := range headers {
			r.Headers[k] = v
		}
	}
	return r
}

// SetContentType sets content type of request
func (r *Request) SetContentType(contentType string) *Request {
	r.Headers[ContentType] = contentType
	return r
}

// ReadBody reads already set bodyBytes field from Request which is wrapper of *http.Request
func (r *Request) ReadBody() ([]byte, error) {

	if len(r.bodyBytes) != 0 {
		return r.bodyBytes, nil
	}

	err := errors.New("bodyBytes is empty")
	r.client.ErrorLogger.Println(err)
	return nil, err
}
