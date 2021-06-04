# Task Conclusion
* A client library made by Go which is suitable for use in another software project.
* Create, Fetch and Delete operations tested with integration tests on accounts resource
* Success and fail scenarios has been tested
* All tests start with docker-compose up
## About me
I used to have almost zero experience on Golang except reading blog posts on Go performance
As a ex-C and currenct Java focused fullstack developer Go seemed to me so humble and powerfull with good module system 
It's module system as easy as node, it
## Inspired Projects
* https://github.com/sendgrid/rest
* https://github.com/google/go-github
* https://github.com/go-resty/resty
* https://github.com/messagebird/go-rest-api
* https://github.com/bozd4g/go-http-client
* https://github.com/golang-standards/project-layout
## Decisions

* I wanted tinyclient's users not to be overwhelmed with converting between string, byte, io.Reader, struct etc...
So tinyclient accept interface and handle the conversions inside
**parseRequestBody** function at **request.go** inspired from resty's middleware
Understanding the reader, buffer, readcloser, byte arrays and their creating methods such as new(bytes.Buffer) , &bytes.Buffer{}, bytes.NewBuffer(bytes), bytes.NewReader(bytes)  took some time. I investigated so many codes and blogs with considering the immutable data and using less memory. 

* With Java experience I was looking for SLF4J alike abstraction middleware and planning to write my own logger (like SLF4J and Log4j relationship). So logger would be extendible, compatible with dependency inversion and interface segregation principles. I finally noticed log.Logger  go built in struct which provides my desires
As a demonstration, tinyclient has error and info loggers. error logger is always enable while debug logger depends on debugmode bool variable.

* I did read builtin http.Client support 10 redirection at most. I found some caveats on stackoverflow such as ``request.HttpRequest.ContentLength > 0 && request.HttpRequest.GetBody == nil`` and ``redirection needs reading the body more than once`` and implemented with logs

* Collection data from clients is important nowadays for debug and analytics. I have added sending some system and client on User-Agent. This feature is controlled by debugMode too.

* I know there is always a edge case we missed. Accessibility of *http.Request   from outside can solve those cases.

* While I am working with Java some certificate problems such as ``javax.net.ssl.SSLHandshakeException: sun.security.validator.ValidatorException: PKIX path building failed: sun.security.provider.certpath.SunCertPathBuilderException: unable to find valid certification path to requested target`` always overwhelming us per 6 months. So I decided to disable default SSL certificate verification 

* I can say context and managing the lifecycle of objects is so important after my docker, Spring context and Android context experiences. Context is like a handlebar of which ones live upon that. Every tinyclient is can be managed by different contexts. Builtin http.Request support context injection which is more flexible but I took my way.

* 

