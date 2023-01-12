package shortloopmux

import (
	"bytes"
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter

	// Status returns the status code of the response or 0 if the response has
	// not been written
	Status() int
	// Written returns whether or not the ResponseWriter has been written.
	Written() bool
	// Size returns the size of the response body.
	Size() int
	// Body returns the response body
	Body() string
}

type ResponseWriterWrapper struct {
	http.ResponseWriter
	body   *bytes.Buffer
	status int
	size   int
}

func (rw *ResponseWriterWrapper) Status() int {
	return rw.status
}

func (rw *ResponseWriterWrapper) Size() int {
	return rw.size
}

func (rw *ResponseWriterWrapper) Written() bool {
	return rw.status != 0
}

func (rw *ResponseWriterWrapper) Body() string {
	return rw.body.String()
}

func NewResponseWriterWrapper(rw http.ResponseWriter) *ResponseWriterWrapper {
	nrw := &ResponseWriterWrapper{
		ResponseWriter: rw,
		body:           bytes.NewBufferString(""),
	}
	return nrw
}

func (rw *ResponseWriterWrapper) WriteHeader(s int) {
	if rw.Written() {
		return
	}
	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *ResponseWriterWrapper) Write(b []byte) (int, error) {
	//if !rw.Written() {
	//	// The status will be StatusOK if WriteHeader has not been called yet
	//	rw.WriteHeader(http.StatusOK)
	//}
	size, err := rw.body.Write(b)
	if err != nil {
		return size, err
	}
	size, err = rw.ResponseWriter.Write(b)
	//_, err = io.Copy(rw.ResponseWriter, rw.body)
	rw.size += size
	return size, err
}

func (rw *ResponseWriterWrapper) GetHeader() http.Header {
	return rw.Header()
}

func (rw *ResponseWriterWrapper) BodyString() string {
	return rw.body.String()
}
