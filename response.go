package request

type Response[T any] struct {
	Status     string
	StatusCode int
	Data       T
}
