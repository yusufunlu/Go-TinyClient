package tinyclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
)

var (
	ContentType     = "Content-Type"
	plainTextType   = "text/plain; charset=utf-8"
	JsonContentType = "application/json; charset=utf-8"
	formContentType = "application/x-www-form-urlencoded"

	jsonCheck = regexp.MustCompile(`(?i:(application|text)/(json|.*\+json|json\-.*)(;|$))`)
)

const (
	ClientName        = "tinyclient"
	ClientVersion     = "1.0.0"
	httpClientTimeout = 15 * time.Second
)

type Client struct {
	HTTPClient  *http.Client // The HTTP client to send requests on.
	Cookies     []*http.Cookie
	ctx         context.Context
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	debugMode   bool
}

func (client *Client) SetContext(ctx context.Context) *Client {
	client.ctx = ctx
	return client
}

// NewClient creates a new TinyClient object.
func NewClient() *Client {
	transport := http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	return &Client{
		InfoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		HTTPClient: &http.Client{
			Timeout:   httpClientTimeout,
			Transport: transport,
		},
	}
}

func (client *Client) SetTimeout(timeout time.Duration) *Client {
	client.HTTPClient.Timeout = timeout * time.Second
	return client
}

// newRequest method creates a new request instance, it will be used for Get, Post, Put, Delete, Patch
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

func (client *Client) SetDebugMode(debugMode bool) *Client {
	client.debugMode = debugMode
	return client
}

func (client *Client) Send(request *Request) (*Response, error) {

	err := request.parseRequestBody()
	if err != nil {
		request.client.ErrorLogger.Println(err)
		return nil, err
	}
	err = client.fillHttpRequest(request)
	if err != nil {
		request.client.ErrorLogger.Println(err)
		return nil, err
	}
	if request.HttpRequest.ContentLength > 0 && request.HttpRequest.GetBody == nil {
		err := errors.New("request.GetBody cannot be nil because it prevents redirection when content length>0")
		client.ErrorLogger.Println(err)
		return nil, err
	}

	if client.debugMode {
		var headerString string
		if headerBytes, err := json.Marshal(request.HttpRequest.Header); err != nil {
			headerString = "Could not Marshal Request Headers"
			client.ErrorLogger.Println(err)
			return nil, err
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
				"------------------------------------------------------------------------------\n"
		client.InfoLogger.Printf(requestLogString)
	}

	res, err := client.HTTPClient.Do(request.HttpRequest)

	if err != nil {
		client.ErrorLogger.Println(err)
		return nil, err
	}

	response := &Response{
		client:     client,
		Response:   res,
		Request:    request,
		ReceivedAt: time.Now(),
	}

	if client.debugMode {
		var elapsedDuration time.Duration
		elapsedDuration = response.ReceivedAt.Sub(request.SentAt)

		var responseHeaderString string
		if headerBytes, err := json.Marshal(res.Header); err != nil {
			responseHeaderString = "Could not Marshal Req Headers"
			client.ErrorLogger.Println(err)
			return nil, err
		} else {
			responseHeaderString = string(headerBytes)
		}

		responseBytes, err := response.ReadBody()
		if err != nil {
			client.ErrorLogger.Println(err)
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
				"------------------------------------------------------------------------------\n"

		client.InfoLogger.Printf(responseLogString)
	}

	return response, nil
}

func (client *Client) fillHttpRequest(r *Request) (err error) {

	r.SentAt = time.Now()
	//Set request Body
	r.HttpRequest.Body = ioutil.NopCloser(bytes.NewReader(r.bodyBytes))

	r.HttpRequest.ContentLength = int64(len(r.bodyBytes))

	// Set request URL
	URL, err := r.generateURL()
	if err != nil {
		client.ErrorLogger.Println(err)
		return err
	}
	r.HttpRequest.URL = URL

	// Set request method
	r.HttpRequest.Method = string(r.Method)
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

	if client.debugMode {
		hostStat, _ := host.Info()
		cpuStat, _ := cpu.Info()

		infoString := fmt.Sprintf("%v/%v; %v; %v; %v",
			ClientName,
			ClientVersion,
			hostStat.Platform,
			cpuStat[0].ModelName,
			hostStat.Hostname,
		)
		if err == nil {
			r.HttpRequest.Header.Set("User-Agent", infoString)
		}
	}

	if err != nil {
		client.ErrorLogger.Println(err)
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

func IsJSONType(isJsonString string) bool {
	return jsonCheck.MatchString(isJsonString)
}
