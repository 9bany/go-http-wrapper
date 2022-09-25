# go-http-wrapper

## How to use

### 1.  Create an target 

```go
type exampleTarget struct {}

func (target *exampleTarget) GetMethod() string {
	return "POST" // http.PostMethod
}

func (target *exampleTarget) GetEndpoint() string {
	return "localhost:3000"
}

func (target *exampleTarget) GetBody() []byte {
	return nil
}

func (target *exampleTarget) GetHeader() Header {
	return map[string]string{}
}
```

### Use target with request wrapper

```go
// Data structure mapping from json response
type Data struct {
    Msg string `json:"data"`
}

target := &exampleTarget{}
requestor := NewRequest[Data](target)

resp, err := requestor.Execute(context.Background())
fmt.Println(resp)
fmt.Println(err)

```