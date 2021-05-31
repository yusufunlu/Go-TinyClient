package tinyclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Response struct holds response values of executed request.
type Response struct {
	Request    *Request
	Response   *http.Response
	bodyBytes  []byte
	size       int64
	receivedAt time.Time
}

// ReadBody reads the http.Response bodyBytes and assigns it to the r.Body
func (r *Response) ReadBody() ([]byte, error) {

	// If r.bodyBytes already set then return r.bodyBytes
	if len(r.bodyBytes) != 0 {
		return r.bodyBytes, nil
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

	// Read response bodyBytes
	b, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		//logger.Errorf("Can't read http.Response bodyBytes Error: %v!", err)
		fmt.Println(err)
		return nil, err
	}

	// Set response readBody
	r.bodyBytes = b

	// Close response bodyBytes
	err = r.Response.Body.Close()
	if err != nil {
		//logger.Errorf("Can't close http.Response body Error: %v!", err)
		return nil, err
	}

	return b, nil

}
