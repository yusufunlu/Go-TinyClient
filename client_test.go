package tinyclient_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
	tiny "tinyclient"

	"github.com/stretchr/testify/require"
)

var desiredData = `{"success": true,"data": "done!"}`

func TestPostString(t *testing.T) {

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			time.Sleep(time.Second * 5)
			// Test request
			require.Equal(t, req.URL.String(), "/post")
			require.Equal(t, req.Method, "POST")
			require.Equal(t, req.Header.Get("Content-Type"), "application/json; charset=utf-8")
			require.Equal(t, req.Header.Get("Test-Header"), "this is a test")

			rw.Header().Set("Content-Type", req.Header.Get("Content-Type"))

			b, _ := ioutil.ReadAll(req.Body)

			fmt.Println("request body: ", string(b))

			_, err := rw.Write(b)
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	url := fmt.Sprintf("%s/post", server.URL)

	client := tiny.NewClient().SetTimeout(30)

	request := client.NewRequest().SetBody(desiredData).SetURL(url).SetMethod("POST")
	request.AddHeaders(map[string]string{"Test-Header": "this is a test"})
	request.SetContentType("application/json; charset=utf-8")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(time.Second * 2)
		println("Cancel")
		cancel()
	}()

	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
	resultBody, err := response.ReadBody()
	require.Equal(t, string(resultBody), desiredData)
}

func TestPostByte(t *testing.T) {

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request
			require.Equal(t, req.URL.String(), "/post")
			require.Equal(t, req.Method, "POST")
			require.Equal(t, req.Header.Get("Content-Type"), "application/json; charset=utf-8")
			require.Equal(t, req.Header.Get("Test-Header"), "this is a test")

			rw.Header().Set("Content-Type", req.Header.Get("Content-Type"))

			b, _ := ioutil.ReadAll(req.Body)

			fmt.Println("request body: ", string(b))

			_, err := rw.Write(b)
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	url := fmt.Sprintf("%s/post", server.URL)

	client := tiny.NewClient().SetTimeout(30)

	desiredDataBytes := []byte(desiredData)
	request := client.NewRequest().SetBody(desiredDataBytes).SetURL(url).SetMethod("POST")
	request.AddHeaders(map[string]string{"Test-Header": "this is a test"})
	request.SetContentType("application/json; charset=utf-8")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
	resultBody, err := response.ReadBody()
	require.Equal(t, string(resultBody), desiredData)
}

func TestPostReader(t *testing.T) {

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request
			require.Equal(t, req.URL.String(), "/post")
			require.Equal(t, req.Method, "POST")
			require.Equal(t, req.Header.Get("Content-Type"), "application/json; charset=utf-8")
			require.Equal(t, req.Header.Get("Test-Header"), "this is a test")

			rw.Header().Set("Content-Type", req.Header.Get("Content-Type"))

			b, _ := ioutil.ReadAll(req.Body)
			req.Body.Close()
			fmt.Println("request body: ", string(b))

			_, err := rw.Write(b)
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	url := fmt.Sprintf("%s/post", server.URL)

	client := tiny.NewClient().SetTimeout(30)

	desiredDataReader := strings.NewReader(desiredData)
	request := client.NewRequest().SetBody(desiredDataReader).SetURL(url).SetMethod("POST")
	request.AddHeaders(map[string]string{"Test-Header": "this is a test"})
	request.SetContentType("application/json; charset=utf-8")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
	resultBody, err := response.ReadBody()
	require.Equal(t, string(resultBody), desiredData)
}

func TestPostJSONMapSuccess(t *testing.T) {

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request
			require.Equal(t, req.URL.String(), "/post")
			require.Equal(t, req.Method, "POST")
			require.Equal(t, req.Header.Get("Content-Type"), "application/json; charset=utf-8")
			require.Equal(t, req.Header.Get("Test-Header"), "this is a test")

			rw.Header().Set("Content-Type", req.Header.Get("Content-Type"))

			b, _ := ioutil.ReadAll(req.Body)
			req.Body.Close()
			fmt.Println("request body: ", string(b))

			_, err := rw.Write(b)
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	url := fmt.Sprintf("%s/post", server.URL)

	client := tiny.NewClient().SetTimeout(30)

	requestBody := map[string]interface{}{"success": true, "data": "done!"}
	request := client.NewRequest().SetBody(requestBody).SetURL(url).SetMethod("POST")
	request.AddHeaders(map[string]string{"Test-Header": "this is a test"})
	request.SetContentType("application/json; charset=utf-8")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)

	resultBody, err := response.ReadBody()
	requestBodyString, err := json.Marshal(requestBody)
	require.Equal(t, string(resultBody), string(requestBodyString))
}

func TestPostHTTPBinAdressSuccess(t *testing.T) {

	url := "https://httpbin.org/post"

	client := tiny.NewClient().SetTimeout(30)

	request := client.NewRequest().SetURL(url).SetMethod("POST")
	request.SetContentType("application/json; charset=utf-8")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)

	resultBody, err := response.ReadBody()

	fmt.Println(string(resultBody))
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://httpbin.org/post", http.StatusFound)
}

func TestPostRedirect(t *testing.T) {

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request
			require.Equal(t, req.URL.String(), "/post")
			require.Equal(t, req.Header.Get("Content-Type"), "application/json; charset=utf-8")

			rw.Header().Set("Content-Type", req.Header.Get("Content-Type"))

			b, _ := ioutil.ReadAll(req.Body)
			req.Body.Close()
			fmt.Println("request body: ", string(b))

			_, err := rw.Write(b)
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	// Start a local HTTP server
	redirectServer := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			http.Redirect(rw, req, server.URL+req.RequestURI, http.StatusFound)

		}),
	)
	defer redirectServer.Close()

	url := fmt.Sprintf("%s/post", redirectServer.URL)

	client := tiny.NewClient().SetTimeout(30)
	request := client.NewRequest().SetURL(url).SetMethod("POST")
	request.SetContentType("application/json; charset=utf-8")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)

	resultBody, err := response.ReadBody()

	fmt.Println(string(resultBody))

}

func TestPostLoggerInject(t *testing.T) {

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request
			require.Equal(t, req.URL.String(), "/post")
			require.Equal(t, req.Method, "POST")
			require.Equal(t, req.Header.Get("Content-Type"), "application/json; charset=utf-8")
			require.Equal(t, req.Header.Get("Test-Header"), "this is a test")

			rw.Header().Set("Content-Type", req.Header.Get("Content-Type"))

			b, _ := ioutil.ReadAll(req.Body)

			fmt.Println("request body: ", string(b))

			_, err := rw.Write(b)
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	url := fmt.Sprintf("%s/post", server.URL)

	infoLog, err := os.OpenFile("infos.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	//errorLog, err := os.OpenFile("errors.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	infoLogger := log.New(infoLog, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	//errorLogger:= log.New(errorLog, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	client := tiny.NewClient().SetTimeout(30)
	client.InfoLogger = infoLogger
	//client.ErrorLogger = errorLogger

	request := client.NewRequest().SetBody(desiredData).SetURL(url).SetMethod(tiny.Post)
	request.AddHeaders(map[string]string{"Test-Header": "this is a test"})
	request.SetContentType("application/json; charset=utf-8")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
	resultBody, err := response.ReadBody()
	require.Equal(t, string(resultBody), desiredData)
}

func TestGet(t *testing.T) {

	desiredResult := runtime.FuncForPC(reflect.ValueOf(TestGet).Pointer()).Name()

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request
			require.Equal(t, req.URL.String(), "/get")
			require.Equal(t, req.Method, "GET")
			require.Equal(t, req.Header.Get(tiny.ContentType), tiny.JsonContentType)
			require.Equal(t, req.Header.Get("Test-Header"), "this is a test")

			rw.Header().Set("Content-Type", "application/json")

			_, err := rw.Write([]byte(desiredResult))
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	url := fmt.Sprintf("%s/get", server.URL)

	client := tiny.NewClient().SetTimeout(30)

	request := client.NewRequest().SetURL(url).SetMethod("GET")
	request.AddHeaders(map[string]string{"Test-Header": "this is a test"})
	request.SetContentType("application/json; charset=utf-8")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
}

func TestGetQueryParams(t *testing.T) {

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request
			require.Equal(t, req.URL.Path, "/get")
			require.Equal(t, req.Method, "GET")
			require.Equal(t, req.Header.Get(tiny.ContentType), tiny.JsonContentType)
			require.Equal(t, req.Header.Get("Test-Header"), "this is a test")

			params, err := url.ParseQuery(req.URL.RawQuery)
			require.NoError(t, err)
			if err != nil {
				return
			}

			require.Equal(t, params.Get("param1"), "value1")
			require.Equal(t, params.Get("param2"), "value2")
			require.Equal(t, params.Get("param3"), "value3")
			require.Equal(t, params.Get("param4"), "value4")

			rw.Header().Set("Content-Type", "application/json")

			_, err = rw.Write([]byte(desiredData))
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	url := fmt.Sprintf("%s/get", server.URL)

	client := tiny.NewClient().SetTimeout(30)
	infoLogger := log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	client.InfoLogger = infoLogger

	request := client.NewRequest().SetURL(url).SetMethod("GET").
		AddHeaders(map[string]string{"Test-Header": "this is a test"}).
		SetContentType("application/json; charset=utf-8").
		AddQueryParam("param1", "value1").
		AddQueryParams(map[string]string{"param2": "value2", "param3": "value3"})

	request.QueryParams["param4"] = "value4"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
}

func TestGetDebugMode(t *testing.T) {

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request
			require.Equal(t, req.URL.String(), "/get")
			require.Equal(t, req.Method, "GET")
			require.Equal(t, req.Header.Get(tiny.ContentType), tiny.JsonContentType)
			require.Equal(t, req.Header.Get("Test-Header"), "this is a test")

			rw.Header().Set("Content-Type", "application/json")

			_, err := rw.Write([]byte(desiredData))
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	url := fmt.Sprintf("%s/get", server.URL)

	client := tiny.NewClient().SetTimeout(30)
	client.SetDebugMode(true)

	request := client.NewRequest().SetURL(url).SetMethod("GET")
	request.AddHeaders(map[string]string{"Test-Header": "this is a test"})
	request.SetContentType("application/json; charset=utf-8")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
}