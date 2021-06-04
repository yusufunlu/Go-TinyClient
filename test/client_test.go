package interview_accountapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	tiny "github.com/yusufunlu/tinyclient"
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
)

var desiredData = `{"success": true,"data": "done!"}`

func TestMethods(t *testing.T) {

	var table = []struct {
		name string
	}{
		{"get"},
		{"post"},
		{"put"},
		{"patch"},
		{"delete"},
	}

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			switch req.URL.String() {
			case "/get":
				t.Log("Get request executed")
				require.Equal(t, req.Method, "GET")
			case "/post":
				t.Log("Post request executed")
				require.Equal(t, req.Method, "POST")
			case "/put":
				t.Log("Put request executed")
				require.Equal(t, req.Method, "PUT")
			case "/patch":
				t.Log("Patch request executed")
				require.Equal(t, req.Method, "PATCH")
			case "/delete":
				t.Log("Delete request executed")
				require.Equal(t, req.Method, "DELETE")
			}
		}),
	)
	defer server.Close()

	client := tiny.NewClient()

	for _, row := range table {
		url := fmt.Sprintf("%s/"+row.name, server.URL)
		request := client.NewRequest().SetURL(url)

		var err error
		switch row.name {
		case "get":
			request.SetMethod(tiny.Get)
			_, err = client.Send(request)
		case "post":
			request.SetMethod(tiny.Post)
			_, err = client.Send(request)
		case "put":
			request.SetMethod(tiny.Put)
			_, err = client.Send(request)
		case "patch":
			request.SetMethod(tiny.Patch)
			_, err = client.Send(request)
		case "delete":
			request.SetMethod(tiny.Delete)
			_, err = client.Send(request)
		}
		require.NoError(t, err)
	}

}

func TestPostString(t *testing.T) {

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			time.Sleep(time.Second * 2)
			// Test request
			require.Equal(t, req.URL.String(), "/post")
			require.Equal(t, req.Method, "POST")
			require.Equal(t, req.Header.Get("Content-Type"), "application/json; charset=utf-8")
			require.Equal(t, req.Header.Get("Test-Header"), "this is a test")

			rw.Header().Set("Content-Type", req.Header.Get("Content-Type"))

			b, _ := ioutil.ReadAll(req.Body)

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

	response, err := client.Send(request)

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

	response, err := client.Send(request)

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

	response, err := client.Send(request)

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

	response, err := client.Send(request)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)

	resultBody, err := response.ReadBody()
	requestBodyString, err := json.Marshal(requestBody)
	require.Equal(t, string(resultBody), string(requestBodyString))
}

func TestPostExternalAdressSuccess(t *testing.T) {

	url := "https://httpbin.org/post"

	client := tiny.NewClient()
	request := client.NewRequest().SetURL(url).SetMethod("POST")
	response, err := client.Send(request)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://httpbin.org/post", http.StatusFound)
}

func TestPostRedirect(t *testing.T) {

	subPath := "/redirection/sub"

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request
			require.Equal(t, req.URL.String(), subPath)
			require.Equal(t, req.Header.Get("Content-Type"), "application/json; charset=utf-8")

			rw.Header().Set("Content-Type", req.Header.Get("Content-Type"))
			b, _ := ioutil.ReadAll(req.Body)
			req.Body.Close()

			_, err := rw.Write(b)
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	// Start a local HTTP server
	redirectServer := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			//http.Redirect(rw, req, server.URL+req.RequestURI, http.StatusFound)
			rw.Header().Set("Location", server.URL+req.RequestURI)
			rw.WriteHeader(301)
			rw.Write([]byte("Redirecting..."))
		}),
	)
	defer redirectServer.Close()

	url := fmt.Sprintf("%v%v", redirectServer.URL, subPath)

	client := tiny.NewClient().SetTimeout(30)
	request := client.NewRequest().SetURL(url).SetMethod("POST")
	request.SetContentType("application/json; charset=utf-8")

	response, err := client.Send(request)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
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

	response, err := client.Send(request)

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

	response, err := client.Send(request)

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

	response, err := client.Send(request)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
}

func TestPostDebugMode(t *testing.T) {

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			time.Sleep(time.Second * 2)
			// Test request
			require.Equal(t, req.URL.String(), "/post")
			require.Equal(t, req.Method, "POST")
			b, _ := ioutil.ReadAll(req.Body)
			_, err := rw.Write(b)
			require.NoError(t, err)
		}),
	)
	defer server.Close()

	url := fmt.Sprintf("%s/post", server.URL)
	client := tiny.NewClient().SetDebugMode(true)
	request := client.NewRequest().SetBody(desiredData).SetURL(url).SetMethod("POST")

	response, err := client.Send(request)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)
	resultBody, err := response.ReadBody()
	require.Equal(t, string(resultBody), desiredData)
}

func TestPostCancelled(t *testing.T) {

	// Start a local HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			time.Sleep(time.Second * 2)
			// Test request
			require.Equal(t, req.URL.String(), "/post")
			require.Equal(t, req.Method, "POST")
			require.Equal(t, req.Header.Get("Content-Type"), "application/json; charset=utf-8")
			require.Equal(t, req.Header.Get("Test-Header"), "this is a test")

			rw.Header().Set("Content-Type", req.Header.Get("Content-Type"))

			b, _ := ioutil.ReadAll(req.Body)

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
	cancel()
	client.SetContext(ctx)

	_, err := client.Send(request)

	require.Error(t, err)
}
