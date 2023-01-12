package data

import "fmt"

type APISample struct {
	RawUri                          string              `json:"rawUri"`
	ApplicationName                 string              `json:"applicationName"`
	HostName                        string              `json:"hostName"`
	Method                          HTTPRequestMethod   `json:"method"`
	Parameters                      map[string][]string `json:"parameters"`
	RequestHeaders                  map[string]string   `json:"requestHeaders"`
	ResponseHeaders                 map[string]string   `json:"responseHeaders"`
	StatusCode                      int                 `json:"statusCode"`
	RequestPayload                  string              `json:"requestPayload"`
	ResponsePayload                 string              `json:"responsePayload"`
	UncaughtExceptionMessage        string              `json:"uncaughtExceptionMessage"`
	PayloadCaptureAttempted         bool                `json:"payloadCaptureAttempted"`
	RequestPayloadCaptureAttempted  bool                `json:"requestPayloadCaptureAttempted"`
	ResponsePayloadCaptureAttempted bool                `json:"responsePayloadCaptureAttempted"`
	Latency                         int64               `json:"latency"`
	Port                            int                 `json:"port"`
	Scheme                          string              `json:"scheme"`
}

func (a *APISample) GetRawUri() string {
	return a.RawUri
}

func (a *APISample) GetApplicationName() string {
	return a.ApplicationName
}

func (a *APISample) GetHostName() string {
	return a.HostName
}

func (a *APISample) GetMethod() HTTPRequestMethod {
	return a.Method
}

func (a *APISample) GetParameters() map[string][]string {
	return a.Parameters
}

func (a *APISample) GetRequestHeaders() map[string]string {
	return a.RequestHeaders
}

func (a *APISample) GetResponseHeaders() map[string]string {
	return a.ResponseHeaders
}

func (a *APISample) GetStatusCode() int {
	return a.StatusCode
}

func (a *APISample) GetRequestPayload() string {
	return a.RequestPayload
}

func (a *APISample) GetResponsePayload() string {
	return a.ResponsePayload
}

func (a *APISample) GetUncaughtExceptionMessage() string {
	return a.UncaughtExceptionMessage
}

func (a *APISample) GetPayloadCaptureAttempted() bool {
	return a.PayloadCaptureAttempted
}

func (a *APISample) GetRequestPayloadCaptureAttempted() bool {
	return a.RequestPayloadCaptureAttempted
}

func (a *APISample) GetResponsePayloadCaptureAttempted() bool {
	return a.ResponsePayloadCaptureAttempted
}

func (a *APISample) GetLatency() int64 {
	return a.Latency
}

func (a *APISample) GetPort() int {
	return a.Port
}

func (a *APISample) GetScheme() string {
	return a.Scheme
}

func (a *APISample) SetRawUri(rawUri string) {
	a.RawUri = rawUri
}

func (a *APISample) SetApplicationName(applicationName string) {
	a.ApplicationName = applicationName
}

func (a *APISample) SetHostName(hostName string) {
	a.HostName = hostName
}

func (a *APISample) SetMethod(method HTTPRequestMethod) {
	a.Method = method
}

func (a *APISample) SetParameters(parameters map[string][]string) {
	a.Parameters = parameters
}

func (a *APISample) SetRequestHeaders(requestHeaders map[string]string) {
	a.RequestHeaders = requestHeaders
}

func (a *APISample) SetResponseHeaders(responseHeaders map[string]string) {
	a.ResponseHeaders = responseHeaders
}

func (a *APISample) SetStatusCode(statusCode int) {
	a.StatusCode = statusCode
}

func (a *APISample) SetRequestPayload(requestPayload string) {
	a.RequestPayload = requestPayload
}

func (a *APISample) SetResponsePayload(responsePayload string) {
	a.ResponsePayload = responsePayload
}

func (a *APISample) SetUncaughtExceptionMessage(uncaughtExceptionMessage string) {
	a.UncaughtExceptionMessage = uncaughtExceptionMessage
}

func (a *APISample) SetPayloadCaptureAttempted(payloadCaptureAttempted bool) {
	a.PayloadCaptureAttempted = payloadCaptureAttempted
}

func (a *APISample) SetRequestPayloadCaptureAttempted(requestPayloadCaptureAttempted bool) {
	a.RequestPayloadCaptureAttempted = requestPayloadCaptureAttempted
}

func (a *APISample) SetResponsePayloadCaptureAttempted(responsePayloadCaptureAttempted bool) {
	a.ResponsePayloadCaptureAttempted = responsePayloadCaptureAttempted
}

func (a *APISample) SetLatency(latency int64) {
	a.Latency = latency
}

func (a *APISample) SetPort(port int) {
	a.Port = port
}

func (a *APISample) SetScheme(scheme string) {
	a.Scheme = scheme
}

func (a APISample) String() string {
	return fmt.Sprintf("APISample{RawUri='%s', ApplicationName='%s', HostName='%s', Method=%s, Parameters=%v, RequestHeaders=%v, ResponseHeaders=%v, StatusCode=%d, RequestPayload='%s', ResponsePayload='%s', UncaughtExceptionMessage='%s', PayloadCaptureAttempted='%s', RequestPayloadCaptureAttempted=%t, ResponsePayloadCaptureAttempted=%t, Latency=%d}", a.RawUri, a.ApplicationName, a.HostName, a.Method, a.Parameters, a.RequestHeaders, a.ResponseHeaders, a.StatusCode, a.RequestPayload, a.ResponsePayload, a.UncaughtExceptionMessage, a.PayloadCaptureAttempted, a.RequestPayloadCaptureAttempted, a.ResponsePayloadCaptureAttempted, a.Latency)
}
