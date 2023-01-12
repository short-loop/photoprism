package shortloopfilter

import "net/http"

type ResponseWriter interface {
	Status() int
	Size() int
	Written() bool
	Body() string
	WriteHeader(s int)
	Write([]byte) (int, error)
	GetHeader() http.Header
	BodyString() string
}
