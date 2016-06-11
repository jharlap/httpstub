Easily configure http test servers for stubbing external dependencies, for testing in Go.

API doc: http://godoc.org/github.com/jharlap/httpstub

```go
ts := httpstub.New().WithDefaultContentType(ctJSON)
defer ts.Close()

ts.Path("/user/*/name").WithBody(`{"id":"a1","name":"Alice"}`)
ts.Path("/user/*/xml").WithContentType(ctXML).WithBody(`<user id="a1"><name>Alice</name></user>`)
ts.Path("/user").WithBody(`{"id":"a1","name":"Alice","gender":"f"}`)

client := mine{a3rdPartyServerURL: ts.URL}
client.DoSomething() // that makes HTTP requests to the 3rd party server
```

