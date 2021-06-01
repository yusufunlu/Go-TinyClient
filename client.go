package tinyclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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
)

type Client struct {
	HTTPClient  *http.Client // The HTTP client to send requests on.
	Cookies     []*http.Cookie
	ctx         context.Context
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

const (
	ClientVersion     = "1.0.0"
	httpClientTimeout = 15 * time.Second
)

// CreateClient creates a new TinyClient object.
func CreateClient() *Client {
	return &Client{
		ErrorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
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
		client:      client,
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

	client.parseRequestBody(request)
	client.fillHttpRequest(request)

	if request.HttpRequest.ContentLength > 0 && request.HttpRequest.GetBody == nil {
		err := errors.New("request.GetBody cannot be nil because it prevents redirection when content length>0")
		client.ErrorLogger.Println(err)
		return nil, err
	}

	client.ctx = ctx

	if client.InfoLogger != nil {
		var headerString string
		if headerBytes, err := json.Marshal(request.HttpRequest.Header); err != nil {
			headerString = "Could not Marshal Request Headers"
		} else {
			headerString = string(headerBytes)
		}

		requestLogString :=
			"\n==============================================================================\n" +
				"~~~ HTTP REQUEST ~~~\n" +
				fmt.Sprintf("%s  %s\n", request.HttpRequest.Method, request.HttpRequest.URL) +
				fmt.Sprintf("HOST   : %s\n", request.HttpRequest.URL.Host) +
				fmt.Sprintf("HEADERS:\n%s\n", headerString) +
				fmt.Sprintf("BODY   :\n%v\n", string(request.bodyBytes)) +
				"==============================================================================\n"
		client.InfoLogger.Printf(requestLogString)
	}

	res, err := client.HTTPClient.Do(request.HttpRequest)

	if err != nil {
		return nil, err
	}

	response := &Response{
		Response:   res,
		Request:    request,
		ReceivedAt: time.Now(),
	}

	if client.InfoLogger != nil {
		var elapsedDuration time.Duration
		elapsedDuration = response.ReceivedAt.Sub(request.SentAt)

		var responseHeaderString string
		if headerBytes, err := json.Marshal(res.Header); err != nil {
			responseHeaderString = "Could not Marshal Req Headers"
		} else {
			responseHeaderString = string(headerBytes)
		}

		responseBytes, err := response.ReadBody()
		if err != nil {
			return nil, err
		}

		responseLogString :=
			"\n==============================================================================\n" +
				"~~~ HTTP RESPONSE ~~~\n" +
				fmt.Sprintf("STATUS       : %s\n", res.Status) +
				fmt.Sprintf("PROTO        : %s\n", res.Proto) +
				fmt.Sprintf("RECEIVED AT  : %v\n", response.ReceivedAt) +
				fmt.Sprintf("TIME DURATION: %v\n", elapsedDuration) +
				fmt.Sprintf("RESPONSE BODY: %v\n", string(responseBytes)) +
				fmt.Sprintf("HEADERS:\n%s\n", responseHeaderString) +
				"==============================================================================\n"

		client.InfoLogger.Printf(responseLogString)
	}

	return response, nil
}

//parseRequestBody logics can't be in Request because of checking contentType
func (client *Client) parseRequestBody(r *Request) (err error) {
	contentType := r.Headers[ContentType]
	if r.Body == nil {
		return
	}
	kind := reflect.TypeOf(r.Body).Kind()

	//http.Request.Body is io.ReadCloser and implements io.Reader too
	//a server can put http.Request.Body into Request.Body
	//it can be any other stream too
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
	}

	return
}

func (client *Client) fillHttpRequest(r *Request) (err error) {

	r.SentAt = time.Now()
	//Set request Body
	r.HttpRequest.Body = ioutil.NopCloser(bytes.NewReader(r.bodyBytes))

	r.HttpRequest.ContentLength = int64(len(r.bodyBytes))

	// Set request URL
	URL, err := client.generateURL(r.URL, r.useSSL)
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

	//redirection needs reading the body more than once
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

func (client *Client) generateURL(address string, useSSL bool) (*url.URL, error) {
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
		client.ErrorLogger.Println(err)
		return nil, err
	}

	return parsedURL, nil
}
