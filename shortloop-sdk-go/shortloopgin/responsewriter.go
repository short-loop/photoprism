package shortloopgin

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseWriterWrapper struct {
	gin.ResponseWriter
	body              *bytes.Buffer
	status            int
	size              int
	shouldCaptureBody bool
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

func NewResponseWriterWrapper(rw gin.ResponseWriter) *ResponseWriterWrapper {
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
	size, err := rw.ResponseWriter.Write(b)
	if err != nil {
		return size, err
	}
	rw.size += size
	if rw.shouldCaptureBody {
		size, err = rw.body.Write(b)
		if err != nil {
			return size, err
		}
	}
	return size, err
}

func (rw *ResponseWriterWrapper) GetHeader() http.Header {
	return rw.Header()
}

func (rw *ResponseWriterWrapper) BodyString() string {
	return rw.body.String()
}

func (rw *ResponseWriterWrapper) ShouldCaptureBody() bool {
	return rw.shouldCaptureBody
}

func (rw *ResponseWriterWrapper) SetShouldCaptureBody(shouldCaptureBody bool) {
	rw.shouldCaptureBody = shouldCaptureBody
}
