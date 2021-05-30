package tinyclient

import (
	"fmt"
	"io/ioutil"
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

// Body returns the response body
func (r *Response) Body() ([]byte, error) {
	body, err := r.readBody()
	if err != nil {
		//logger.Errorf("Can not read http.Response body Error: %v", err)
		return nil, err
	}
	return body, nil
}

// readBody reads the http.Response body and assigns it to the r.Body
func (r *Response) readBody() ([]byte, error) {

	// If r.body already set then return r.body
	if len(r.body) != 0 {
		return r.body, nil
	}

	// Check if Response.resp (*http.Response) is nil
	if r.Response == nil {
		err := fmt.Errorf("http.Response is nil")
		//logger.Errorf("%v", err)
		return nil, err
	}

	// Check if Response.resp.Body (*http.Response.Body) is nil
	if r.Response.Body == nil {
		err := fmt.Errorf("http.Response's Body is nil")
		//logger.Errorf("%v", err)
		return nil, err
	}

	// Read response body
	b, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		//logger.Errorf("Can't read http.Response body Error: %v!", err)
		fmt.Println(err)
		return nil, err
	}

	// Set response readBody
	r.body = b

	// Close response body
	err = r.Response.Body.Close()
	if err != nil {
		//logger.Errorf("Can't close http.Response body Error: %v!", err)
		return nil, err
	}

	return b, nil

}
