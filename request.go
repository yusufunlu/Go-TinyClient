package tinyclient

import (
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	URL        string
	Method     string
	useSSL     bool
	Token      string
	AuthScheme string
	QueryParam url.Values
	FormData   url.Values
	Header     http.Header
	Time       time.Time
	Body       interface{}
	Result     interface{}
	Error      interface{}
	Cookies    []*http.Cookie
	//HttpRequest is exposed because of missing cases of this client wrapper and so a professional user can handle for this edge
	HttpRequest      *http.Request
	setContentLength bool
	bodyBytes        []byte
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
	r.Header.Set(header, value)
	return r
}

// SetHeaders sets request headers
func (r *Request) SetHeaders(headers map[string]string) *Request {
	if len(headers) > 0 {
		for k, v := range headers {
			r.Header.Set(k, v)
		}
	}
	return r
}

// SetContentType sets content type of request
func (r *Request) SetContentType(contentType string) *Request {
	r.Header.Set(ContentType, contentType)
	return r
}
