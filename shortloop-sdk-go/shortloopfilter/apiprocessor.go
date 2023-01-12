package shortloopfilter

import (
	"bytes"
	. "github.com/short-loop/shortloop-common-go/models/data"
	"github.com/short-loop/shortloop-sdk-go/buffer"
	"github.com/short-loop/shortloop-sdk-go/sdklogger"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ApiProcessor struct {
	DiscoveredApiManager *buffer.DiscoveredApiManager
	RegisteredApiManager *buffer.RegisteredApiManager
}

func NewApiProcessor(discoveredApiManager *buffer.DiscoveredApiManager, registeredApiManager *buffer.RegisteredApiManager) *ApiProcessor {
	return &ApiProcessor{
		DiscoveredApiManager: discoveredApiManager,
		RegisteredApiManager: registeredApiManager,
	}
}

func (ap *ApiProcessor) ProcessDiscoveredApi(context RequestResponseContext, next func(canOffer bool)) {
	var worker buffer.ManagerWorker = ap.DiscoveredApiManager.GetWorker()
	if worker == nil {
		sdklogger.Logger.Error("BufferManagerWorker is nil inside DiscoveredApiProcessor")
		next(false)
		return
	}

	var canOffer bool = worker.CanOffer(context.GetApiBufferKey())

	context.SetPayloadCaptureAttempted(false)

	next(canOffer)
	// h.ServeHTTP(context.GetResponseWriterWrapper(), context.GetHttpRequest())

	if canOffer {
		ap.tryOffering(context, worker)
	}
}

func (ap *ApiProcessor) ProcessRegisteredApi(context RequestResponseContext, next func(canOffer bool, responsePayloadCaptureAttempted bool)) {
	var worker buffer.ManagerWorker = ap.RegisteredApiManager.GetWorker()
	if worker == nil {
		sdklogger.Logger.Error("BufferManagerWorker is nil inside RegisteredApiProcessor")
		next(false, false)
		return
	}
	var canOffer bool = worker.CanOffer(context.GetApiBufferKey())

	context.SetPayloadCaptureAttempted(true)
	var requestPayloadCaptureAttempted bool = false
	var responsePayloadCaptureAttempted bool = false

	if canOffer {
		requestPayloadCaptureAttempted = ap.shouldCaptureSampleRequest(context)
		if requestPayloadCaptureAttempted {
			context.SetRequestPayload(ap.wrapRequest(context.GetHttpRequest()))
		}
		responsePayloadCaptureAttempted = ap.shouldCaptureSampleResponse(context)
		context.SetRequestPayloadCaptureAttempted(requestPayloadCaptureAttempted)
		context.SetResponsePayloadCaptureAttempted(responsePayloadCaptureAttempted)
	}

	var startTime time.Time
	var shouldComputeLatency bool = canOffer
	if shouldComputeLatency {
		startTime = time.Now()
	}
	// h.ServeHTTP(context.GetResponseWriterWrapper(), context.GetHttpRequest())
	next(canOffer, responsePayloadCaptureAttempted)

	if shouldComputeLatency {
		context.SetLatency(time.Since(startTime).Milliseconds())
	}

	if canOffer {
		ap.tryOffering(context, worker)
	}

}

func (ap *ApiProcessor) tryOffering(context RequestResponseContext, worker buffer.ManagerWorker) {

	defer func() {
		if err := recover(); err != nil {
			sdklogger.Logger.ErrorF("Panic in tryOffering: %s\n", err)
		}
	}()

	var apiSample *APISample = ap.getBufferEntryForApiSample(context)
	sdklogger.Logger.InfoF("trying to offer apiSample of %s to buffer\n", apiSample.GetRawUri())
	if apiSample != nil {
		worker.Offer(context.GetApiBufferKey(), *apiSample)
	}
}

func (ap *ApiProcessor) getBufferEntryForApiSample(context RequestResponseContext) *APISample {
	var apiSample APISample = APISample{}

	apiSample.SetApplicationName(context.GetApplicationName())
	if nil != context.GetApiConfig() {
		apiSample.SetMethod(context.GetApiConfig().GetMethod())
	} else {
		apiSample.SetMethod(context.GetObservedApi().GetMethod())
	}
	apiSample.SetRawUri(context.GetObservedApi().GetUri().GetURIPath())
	apiSample.SetParameters(ap.getParameters(context.GetHttpRequest()))
	apiSample.SetRequestHeaders(ap.getRequestHeaders(context.GetHttpRequest()))
	apiSample.SetResponseHeaders(ap.getResponseHeaders(context.GetResponseWriterWrapper()))
	apiSample.SetLatency(context.GetLatency())

	hostname, port := ap.getHostAndPort(context.GetHttpRequest())
	apiSample.SetHostName(hostname)
	apiSample.SetPort(port)
	apiSample.SetScheme(context.GetHttpRequest().URL.Scheme)

	if context.GetRequestPayloadCaptureAttempted() {
		apiSample.SetRequestPayload(string(context.GetRequestPayload()))
	}

	apiSample.SetStatusCode(context.GetResponseWriterWrapper().Status())
	if context.GetResponsePayloadCaptureAttempted() {
		apiSample.SetResponsePayload(context.GetResponseWriterWrapper().BodyString())
	}

	apiSample.SetPayloadCaptureAttempted(context.GetPayloadCaptureAttempted())
	apiSample.SetRequestPayloadCaptureAttempted(context.GetRequestPayloadCaptureAttempted())
	apiSample.SetResponsePayloadCaptureAttempted(context.GetResponsePayloadCaptureAttempted())

	return &apiSample
}

func (ap *ApiProcessor) getHostAndPort(request *http.Request) (hostname string, port int) {
	parts := strings.Split(request.Host, ":")
	hostname = parts[0]
	port, _ = strconv.Atoi(parts[1])
	return
}

func (ap *ApiProcessor) shouldCaptureSampleRequest(context RequestResponseContext) bool {
	if nil != context.GetApiConfig() && !context.GetApiConfig().GetCaptureSampleRequest() {
		return false
	}
	return true
}

func (ap *ApiProcessor) shouldCaptureSampleResponse(context RequestResponseContext) bool {
	if nil != context.GetApiConfig() && !context.GetApiConfig().GetCaptureSampleResponse() {
		return false
	}
	return true
}

func (ap *ApiProcessor) getRequestHeaders(request *http.Request) map[string]string {
	if request == nil {
		return map[string]string{}
	}

	var headers map[string]string = map[string]string{}

	for key := range request.Header {
		headers[key] = request.Header.Get(key)
	}
	return headers
}

func (ap *ApiProcessor) getResponseHeaders(responseWriterWrapper ResponseWriter) map[string]string {
	if responseWriterWrapper == nil {
		return map[string]string{}
	}

	var headers map[string]string = map[string]string{}

	for key := range responseWriterWrapper.GetHeader() {
		headers[key] = responseWriterWrapper.GetHeader().Get(key)
	}
	return headers
}

func (ap *ApiProcessor) getParameters(request *http.Request) map[string][]string {
	return request.URL.Query()
}

func (ap *ApiProcessor) wrapRequest(request *http.Request) *bytes.Buffer {
	buf := &bytes.Buffer{}
	// Don't buffer if there is no body.
	if request.Body == nil || request.Body == http.NoBody {
		return buf
	}
	request.Body = readCloser{
		Reader: io.TeeReader(request.Body, buf),
		Closer: request.Body,
	}
	return buf
}

// readCloser combines an io.Reader and an io.Closer to implement io.ReadCloser.
type readCloser struct {
	io.Reader
	io.Closer
}
