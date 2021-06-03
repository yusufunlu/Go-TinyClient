<p align="center"><img src="tiny.jpg"/>

<p align="center">
<h1 align="center">tinyclient</h1>
<p align="center">Simple HTTP client for Golang</p>

## Table of Contents
- [Features](#Features)
- [Prerequisites](#Prerequisites)
- [Installation](#Installation)
- [Usage](#Usage)
- [About](#about)
- [License](#license)


## Features
* Support body in string,[]byte,io.Reader,io.ReadCloser,map,slice or struct types
* Support of logger injection
* Support of redirection
* Support of default error logger which can be overridden 
* Support of client and system info sending as User-Agent
* Support of *http.Request access for edge case configuration
* Default SSL certificate verification is disabled, can be still overridden

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
You can use alias while importing

````
import tiny "github.com/yusufunlu/tinyclient"
...

client := tiny.NewClient()

````

