package tinyclient_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
			require.Equal(t, req.Header.Get("Content-Type"), "application/json")
			require.Equal(t, req.Header.Get("Test-Header"), "this is a test")

			rw.Header().Set("Content-Type", "application/json")

			b, _ := ioutil.ReadAll(req.Body)

			fmt.Println("request body: ", string(b))

			_, err := rw.Write(b)
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	url := fmt.Sprintf("%s/post", server.URL)

	client := tiny.CreateClient().SetTimeout(30)

	request := client.NewRequest().SetBody(desiredData).SetURL(url).SetMethod("POST")
	request.SetHeaders(map[string]string{"Test-Header": "this is a test"})
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
	resultBody, err := response.Body()
	require.Equal(t, string(resultBody), desiredData)
}

func TestPostByte(t *testing.T) {

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request
			require.Equal(t, req.URL.String(), "/post")
			require.Equal(t, req.Method, "POST")
			require.Equal(t, req.Header.Get(tiny.ContentType), tiny.JsonContentType)
			require.Equal(t, req.Header.Get("Test-Header"), "this is a test")

			rw.Header().Set("Content-Type", "application/json")

			b, _ := ioutil.ReadAll(req.Body)

			fmt.Println("request body: ", string(b))

			_, err := rw.Write(b)
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	url := fmt.Sprintf("%s/post", server.URL)

	client := tiny.CreateClient().SetTimeout(30)

	desiredDataBytes := []byte(desiredData)
	request := client.NewRequest().SetBody(desiredDataBytes).SetURL(url).SetMethod("POST")
	request.SetHeaders(map[string]string{"Test-Header": "this is a test"})
	request.SetContentType("application/json; charset=utf-8")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
	resultBody, err := response.Body()
	require.Equal(t, string(resultBody), desiredData)
}

func TestPostReader(t *testing.T) {

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request
			require.Equal(t, req.URL.String(), "/post")
			require.Equal(t, req.Method, "POST")
			require.Equal(t, req.Header.Get(tiny.ContentType), tiny.JsonContentType)
			require.Equal(t, req.Header.Get("Test-Header"), "this is a test")

			rw.Header().Set("Content-Type", "application/json")

			b, _ := ioutil.ReadAll(req.Body)
			req.Body.Close()
			fmt.Println("request body: ", string(b))

			_, err := rw.Write(b)
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	url := fmt.Sprintf("%s/post", server.URL)

	client := tiny.CreateClient().SetTimeout(30)

	desiredDataReader := strings.NewReader(desiredData)
	request := client.NewRequest().SetBody(desiredDataReader).SetURL(url).SetMethod("POST")
	request.SetHeaders(map[string]string{"Test-Header": "this is a test"})
	request.SetContentType("application/json; charset=utf-8")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
	resultBody, err := response.Body()
	require.Equal(t, string(resultBody), desiredData)
}

func TestPostJSONMapSuccess(t *testing.T) {

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request
			require.Equal(t, req.URL.String(), "/post")
			require.Equal(t, req.Method, "POST")
			require.Equal(t, req.Header.Get(tiny.ContentType), tiny.JsonContentType)
			require.Equal(t, req.Header.Get("Test-Header"), "this is a test")

			rw.Header().Set("Content-Type", "application/json")

			b, _ := ioutil.ReadAll(req.Body)
			req.Body.Close()
			fmt.Println("request body: ", string(b))

			_, err := rw.Write(b)
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	url := fmt.Sprintf("%s/post", server.URL)

	client := tiny.CreateClient().SetTimeout(30)

	requestBody := map[string]interface{}{"username": "testuser", "password": "testpass"}
	request := client.NewRequest().SetBody(requestBody).SetURL(url).SetMethod("POST")
	request.SetHeaders(map[string]string{"Test-Header": "this is a test"})
	request.SetContentType("application/json; charset=utf-8")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
	resultBody, err := response.Body()

	requestBodyString, err := json.Marshal(requestBody)

	require.Equal(t, string(resultBody), string(requestBodyString))
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

	client := tiny.CreateClient().SetTimeout(30)

	request := client.NewRequest().SetBody(desiredResult).SetURL(url).SetMethod("GET")
	request.SetHeaders(map[string]string{"Test-Header": "this is a test"})
	request.SetContentType("application/json; charset=utf-8")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
}
