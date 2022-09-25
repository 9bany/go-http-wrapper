package request

type Target interface {
	GetMethod() string
	GetEndpoint() string
	GetBody() []byte
	GetHeader() Header
}
