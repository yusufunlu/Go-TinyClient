# ![tinyclient](tiny.jpg) tinyclient

#Features
* Support body in string,[]byte,io.Reader,io.ReadCloser,map,slice or struct types
* Support of redirection
* Support of logger injection
* Support of default error logger which can be overridden 
* Support of client and system info sending as User-Agent
* No special function support of XML as body, you can still send as []byte or string
* No special support of formdata, you can still configure your http.Request via request.HttpRequest
* Default SSL certificate verification is disabled, can be still overridden

