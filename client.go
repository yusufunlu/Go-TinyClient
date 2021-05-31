package tinyclient

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var (
	ContentType     = "Content-Type"
	plainTextType   = "text/plain; charset=utf-8"
	JsonContentType = "application/json; charset=utf-8"
	formContentType = "application/x-www-form-urlencoded"

	jsonCheck = regexp.MustCompile(`(?i:(application|text)/(json|.*\+json|json\-.*)(;|$))`)
	xmlCheck  = regexp.MustCompile(`(?i:(application|text)/(xml|.*\+xml)(;|$))`)
)

type Client struct {
	HTTPClient *http.Client // The HTTP client to send requests on.
	DebugLog   *log.Logger  // Optional logger for debugging purposes.
	Cookies    []*http.Cookie
	ctx        context.Context
}

const (
	ClientVersion     = "1.0.0"
	httpClientTimeout = 15 * time.Second
)

// CreateClient creates a new TinyClient object.
func CreateClient() *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: httpClientTimeout,
		},
	}
}

func (client *Client) SetTimeout(timeout time.Duration) *Client {
	client.HTTPClient.Timeout = timeout * time.Second
	return client
}

// newRequest method creates a new request instance, it will be used for Get, Post, Put, Delete, Patch, Head, Options, etc.
func (client *Client) NewRequest() *Request {
	return &Request{
		QueryParams: map[string]string{},
		Headers:     map[string]string{},
		Cookies:     make([]*http.Cookie, 0),
		HttpRequest: &http.Request{
			Header: make(http.Header),
		},
		FormData: url.Values{},
	}
}

func (client *Client) Send(request *Request, ctx context.Context) (*Response, error) {

	parseRequestBody(request)
	client.fillHttpRequest(request)

	if request.HttpRequest.ContentLength > 0 && request.HttpRequest.GetBody == nil {
		return nil, errors.New("request.GetBody cannot be nil because it prevents redirection when content length>0")
	}

	client.ctx = ctx

	res, err := client.HTTPClient.Do(request.HttpRequest)

	if err != nil {
		return nil, err
	}

	response := &Response{
		Response:   res,
		Request:    request,
		receivedAt: time.Now(),
	}

	return response, nil
}

//parseRequestBody logics can't be in Request because of checking contentType
func parseRequestBody(r *Request) (err error) {
	contentType := r.Headers[ContentType]
	if r.Body == nil {
		return
	}
	kind := reflect.TypeOf(r.Body).Kind()

	//reader case can be used for sending received request to another server
	if reader, ok := r.Body.(io.Reader); ok {
		r.bodyBytes, err = io.ReadAll(reader)
	} else if b, ok := r.Body.([]byte); ok {
		r.bodyBytes = b
	} else if s, ok := r.Body.(string); ok {
		r.bodyBytes = []byte(s)
	} else if IsJSONType(contentType) &&
		(kind == reflect.Struct || kind == reflect.Map || kind == reflect.Slice) {
		b, err := json.Marshal(r.Body)
		r.bodyBytes = b
		if err != nil {
			return err
		}
	} else if IsXMLType(contentType) && (kind == reflect.Struct) {
		b, err = xml.Marshal(r.Body)
		r.bodyBytes = b
		if err != nil {
			return
		}
	}

	fmt.Println("request body bytes: ", string(r.bodyBytes))

	return
}

func (client *Client) fillHttpRequest(r *Request) (err error) {

	//Set request Body
	r.HttpRequest.Body = ioutil.NopCloser(bytes.NewReader(r.bodyBytes))
	r.HttpRequest.Body = ioutil.NopCloser(bytes.NewBuffer(r.bodyBytes))
	//redirection need reading the body more than once
	r.HttpRequest.GetBody = func() (io.ReadCloser, error) {
		return ioutil.NopCloser(strings.NewReader("deneme")), nil
	}
	r.HttpRequest.ContentLength = int64(len(r.bodyBytes))

	// Set request URL
	URL, err := GenerateURL(r.URL, r.useSSL)
	if err != nil {
		return err
	}
	r.HttpRequest.URL = URL

	// Set request method
	r.HttpRequest.Method = r.Method
	// Add headers into http request
	for key, value := range r.Headers {
		r.HttpRequest.Header.Set(key, value)
	}

	// Add cookies from client instance into http request
	for _, cookie := range client.Cookies {
		r.HttpRequest.AddCookie(cookie)
	}

	// Add cookies from request instance into http request
	for _, cookie := range r.Cookies {
		r.HttpRequest.AddCookie(cookie)
	}

	// Use context if it was specified
	if client.ctx != nil {
		r.HttpRequest = r.HttpRequest.WithContext(client.ctx)
	}

	if err != nil {
		return err
	}

	// assign get body func for the underlying raw request instance
	r.HttpRequest.GetBody = func() (io.ReadCloser, error) {
		if r.bodyBytes != nil {
			return ioutil.NopCloser(bytes.NewReader(r.bodyBytes)), nil
		}
		return nil, nil
	}

	return
}

// IsJSONType method is to check JSON content type or not
func IsJSONType(ct string) bool {
	return jsonCheck.MatchString(ct)
}

// IsXMLType method is to check XML content type or not
func IsXMLType(ct string) bool {
	return xmlCheck.MatchString(ct)
}

func GenerateURL(address string, useSSL bool) (*url.URL, error) {
	v := strings.ToLower(address)
	if strings.HasPrefix(v, "http://") {
		address = strings.TrimPrefix(address, "http://")
	} else if strings.HasPrefix(v, "https://") {
		address = strings.TrimPrefix(address, "https://")
	}

	// Generate prefix
	prefix := "http://"
	if useSSL {
		prefix = "https://"
	}

	// Merge prefix and URL
	URL := prefix + address

	// Parse URL
	parsedURL, err := url.Parse(URL)
	if err != nil {
		//logger.Errorf("URL parsing error: %s, %v", URL, err)
		return nil, err
	}

	return parsedURL, nil
}

// ReadBody reads already set bodyBytes field from Request which is wrapper of *http.Request
func (r *Request) ReadBody() ([]byte, error) {

	if len(r.bodyBytes) != 0 {
		return r.bodyBytes, nil
	}

	err := errors.New("bodyBytes is empty")
	return nil, err
}
