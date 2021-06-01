package tinyclient

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// Method is a type
type Method string

// Supported HTTP methods, for preventing typos
const (
	Get    Method = "GET"
	Post   Method = "POST"
	Put    Method = "PUT"
	Patch  Method = "PATCH"
	Delete Method = "DELETE"
)

type Request struct {
	client *Client
	//HttpRequest is exposed because of missing cases of this client wrapper and so a professional user can handle for this edge
	HttpRequest *http.Request
	URL         string
	Method      Method
	Body        interface{}
	//bodyBytes is not exposed, it is for holding Body interface internally as bytes and reading the body without spoiling HttpRequest.GetBody
	bodyBytes   []byte
	Cookies     []*http.Cookie
	Headers     map[string]string
	QueryParams map[string]string
	FormData    url.Values
	SentAt      time.Time
	useSSL      bool
}

func (request *Request) SetBody(body interface{}) *Request {
	request.Body = body
	return request
}

func (request *Request) SetMethod(method Method) *Request {
	request.Method = method
	return request
}

func (request *Request) SetURL(url string) *Request {
	request.URL = url
	return request
}

func (request *Request) AddHeader(header, value string) *Request {
	request.Headers[header] = value
	return request
}

// AddHeaders sets request headers
func (request *Request) AddHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		request.Headers[k] = v
	}
	return request
}

// SetContentType sets content type of request
func (request *Request) SetContentType(contentType string) *Request {
	request.Headers[ContentType] = contentType
	return request
}

// ReadBody reads already set bodyBytes field from Request which is wrapper of *http.Request
func (request *Request) ReadBody() ([]byte, error) {

	if len(request.bodyBytes) != 0 {
		return request.bodyBytes, nil
	}

	err := errors.New("bodyBytes is empty")
	request.client.ErrorLogger.Println(err)
	return nil, err
}

func (request *Request) AddQueryParam(param, value string) *Request {
	request.QueryParams[param] = value
	return request
}

func (request *Request) AddQueryParams(params map[string]string) *Request {
	for k, v := range params {
		request.AddQueryParam(k, v)
	}
	return request
}

//parseRequestBody logics can't be in Request because of checking contentType
func (request *Request) parseRequestBody() (err error) {
	contentType := request.Headers[ContentType]
	if request.Body == nil {
		return
	}
	kind := reflect.TypeOf(request.Body).Kind()

	//http.Request.Body is io.ReadCloser and implements io.Reader too
	//a server can put http.Request.Body into Request.Body
	//it can be any other stream too
	if reader, ok := request.Body.(io.Reader); ok {
		request.bodyBytes, err = io.ReadAll(reader)
	} else if b, ok := request.Body.([]byte); ok {
		request.bodyBytes = b
	} else if s, ok := request.Body.(string); ok {
		request.bodyBytes = []byte(s)
	} else if IsJSONType(contentType) &&
		(kind == reflect.Struct || kind == reflect.Map || kind == reflect.Slice) {
		b, err := json.Marshal(request.Body)
		request.bodyBytes = b
		if err != nil {
			return err
		}
	}

	return
}

func (request *Request) generateURL() (*url.URL, error) {
	v := strings.ToLower(request.URL)
	if strings.HasPrefix(v, "http://") {
		request.URL = strings.TrimPrefix(request.URL, "http://")
	} else if strings.HasPrefix(v, "https://") {
		request.URL = strings.TrimPrefix(request.URL, "https://")
	}

	// Generate prefix
	prefix := "http://"
	if request.useSSL {
		prefix = "https://"
	}

	if len(request.QueryParams) != 0 {
		//Create new URL with query params
		params := url.Values{}
		for key, value := range request.QueryParams {
			params.Add(key, value)
		}
		request.URL = request.URL + "?" + params.Encode()
	}

	// Merge prefix and URL
	URL := prefix + request.URL

	// Parse URL
	parsedURL, err := url.Parse(URL)
	if err != nil {
		request.client.ErrorLogger.Println(err)
		return nil, err
	}

	return parsedURL, nil
}
