<p align="center"><img src="tiny.jpg"/>

<p align="center">
<h1 align="center">tinyclient</h1>
<p align="center">Simple HTTP client for Golang</p>

## Table of Contents
- [Features](#Features)
- [Prerequisites](#Prerequisites)
- [Installation](#Installation)
- [Usage](#Usage)
- [Detailed Explanation](#explanation)
- [Story](#story)


## Features
* Support body in string,[]byte,io.Reader,io.ReadCloser,map,slice or struct types
* Support of logger injection with default error logger and opt-in info logger
* Support of redirection
* Support of client and system info sending as User-Agent
* Support of *http.Request access for edge case configuration
* Default SSL certificate verification is disabled, can be still overridden
* Context injection
## Prerequisites
Go version 1.13.X, 1.14.X, 1.15.X and 1.16.X
## Installation
Repo is private so best installing way is from local. You can use ``git clone https://github.com/yusufunlu/tinyclient`` to copy to your local
```
<home>/
 |-- tinyclient/
 |-- hello/
    |-- put into go.mod 
        replace github.com/yusufunlu/tinyclient => ../tinyclient
    |-- put into hello.go 
        import (tiny "github.com/yusufunlu/tinyclient")
```
go to ``hello`` folder, use ``go mod tidy`` to synchronize the hello module's dependencies. It will add like ``require example.com/greetings v0.0.0-00010101000000-000000000000`` to ``go.mod``


## Usage
You can use alias while importing. Following code is a quite enhanced example 

````
import tiny "github.com/yusufunlu/tinyclient"
import "fmt"

func main() {
	url := "https://httpbin.org/post"
	client := tiny.NewClient().SetDebugMode(true)
	request := client.NewRequest().SetURL(url).SetMethod("POST")
	response, err := client.Send(request)
}
````

When you execute above code, console output will be similar to following: because debugMode is set
````
INFO: 2021/06/04 01:11:39 client.go:119: 
==============================================================================
~~~ HTTP REQUEST ~~~
POST  http://httpbin.org/post
HOST   : httpbin.org
HEADERS:
{"User-Agent":["tinyclient/1.0.0; Microsoft Windows 10 Pro for Workstations; Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz; LAPTOP-E6QR5KGQ"]}
BODY   :

------------------------------------------------------------------------------
INFO: 2021/06/04 01:11:40 client.go:162: 
==============================================================================
~~~ HTTP RESPONSE ~~~
STATUS       : 200 OK
PROTO        : HTTP/1.1
RECEIVED AT  : 2021-06-04 01:11:40.0605231 +0300 +03 m=+1.390798201
TIME DURATION: 1.3859438s
RESPONSE BODY: {
  "args": {}, 
  "data": "", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Content-Length": "0", 
    "Host": "httpbin.org", 
    "User-Agent": "tinyclient/1.0.0; Microsoft Windows 10 Pro for Workstations; Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz; LAPTOP-E6QR5KGQ", 
    "X-Amzn-Trace-Id": "Root=1-60b9539c-7cfd49b57d6efdfa2be25e90"
  }, 
  "json": null, 
  "origin": "94.54.16.185", 
  "url": "http://httpbin.org/post"
}

HEADERS:
{"Access-Control-Allow-Credentials":["true"],"Access-Control-Allow-Origin":["*"],"Connection":["keep-alive"],"Content-Length":["462"],"Content-Type":["application/json"],"Date":["Thu, 03 Jun 2021 22:11:40 GMT"],"Server":["gunicorn/19.9.0"]}
------------------------------------------------------------------------------
````
## Detailed Explanation

You need to create client first
``client := tiny.NewClient()``

Then create request object
``request := client.NewRequest()``

Use SetBody function for setting body in string,[]byte,io.Reader,io.ReadCloser,map,slice or struct types
``request.SetBody("test")``
``request.SetBody([]byte("test"))``
``request.SetBody(map[string]interface{}{"success": true})``
``request.SetBody(strings.NewReader("test"))``

Reach underlying ***http.Request** via request.HttpRequest and change it
``request.HttpRequest.URL = someString``

Set debug mode for a client
``client.SetDebugMode(true)``

Inject context from outside to client and therefore request. Context will be injected to request when used **client.Send(request)**
````
ctx, cancel := context.WithCancel(context.Background())
cancel()
client.SetContext(ctx)
````
Client has default ErrorLogger. You override it or inject InfoLogger
````
infoLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
client.InfoLogger = infoLogger
````
Default client timeout is 15 sec, you can change it via `client.SetTimeout(30)` 

Client and Request use builder pattern so you can chain functions
Get request with query parameters with function chain
````
request := client.NewRequest().SetURL(url).SetMethod("GET").
    AddQueryParam("param1", "value1").
    AddQueryParams(map[string]string{"param2": "value2", "param3": "value3"})
````

Creating request object via which client object doesn't effect request for now. But it will in next releases
````
request := client.NewRequest()
````

<a name="story"></a>
# Story
* I am new to Go 
- [Story of development process and technical decisions ](STORY.md)

## License 
Copyright (c) 2021 Yusuf Unlu

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE. 
