# ![tinyclient](tiny.jpg) tinyclient

#Features
* Support body in string,[]byte,io.Reader,io.ReadCloser,map,slice or struct types
* Support of redirection
* Support of logger injection
* Support of default error logger which can be overridden 
* Not yet auto user-agent info filling
* No Support XML as body, you can still send as []byte or string
* No Support Formdata, you can still configure your http.Request via request.HttpRequest
* No User agent info
* cert support





#Theory
* http.Request.Body is io.ReadCloser
* fmt.Printf(response.Body) would fail

## string, []byte,Reader to ReadCloser conversions

**[]byte("hello world")** string to byte

**string(b)** bytes to string

**bytes.NewReader([]byte("hello world")** string to Reader

**strings.NewReader("hello world")** string to Reader

**ioutil.NopCloser(strings.NewReader("hello world"))** string to ReadCloser

* after Go 1.16 use io.NopCloser(stringReader) 

## from ReadCloser to string
````
buf := new(bytes.Buffer)
buf.ReadFrom(r)
r.Close()
s := buf.String()
fmt.Println(s)
````

##Solution 1
so we convert it to a string by passing it through a buffer first. A 'costly' but useful process.

````
buf := new(bytes.Buffer)
buf.ReadFrom(response.Body)
newStr := buf.String()
fmt.Printf(newStr)
````

##Solution 2
````
b, _:= ioutil.ReadAll(req.Body)
b, err := ioutil.ReadAll(req.Body);
````