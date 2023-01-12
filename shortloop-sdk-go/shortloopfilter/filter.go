package shortloopfilter

import (
	"bytes"
	"github.com/short-loop/shortloop-sdk-go/sdklogger"
	"net/http"

	. "github.com/short-loop/shortloop-common-go/models/data"
	"github.com/short-loop/shortloop-sdk-go/buffer"
	"github.com/short-loop/shortloop-sdk-go/config"
)

var currentShortloopFilter *ShortloopFilter = &ShortloopFilter{}

type RequestResponseContext struct {
	httpRequest                     *http.Request
	responseWriterWrapper           ResponseWriter
	applicationName                 string
	observedApi                     ObservedApi
	agentConfig                     AgentConfig
	apiConfig                       *ApiConfig
	apiBufferKey                    buffer.ApiBufferKey
	payloadCaptureAttempted         bool
	requestPayloadCaptureAttempted  bool
	responsePayloadCaptureAttempted bool
	latency                         int64
	// requestBody holds a reference to the original request.Body.
	requestPayload interface {
		Bytes() []byte
	}
}

func NewRequestResponseContext(responseWriterWrapper ResponseWriter, httpRequest *http.Request, applicationName string) RequestResponseContext {
	return RequestResponseContext{
		httpRequest:           httpRequest,
		responseWriterWrapper: responseWriterWrapper,
		applicationName:       applicationName,
	}
}

type ShortloopFilter struct {
	AgentConfig         *AgentConfig
	configManager       *config.Manager
	ApiProcessor        *ApiProcessor
	UserApplicationName string
}

func NewShortloopFilter(configManager *config.Manager, userApplicationName string) ShortloopFilter {
	return ShortloopFilter{
		configManager:       configManager,
		UserApplicationName: userApplicationName,
	}
}

func CurrentShortloopFilter() *ShortloopFilter {
	return currentShortloopFilter
}

func (sf *ShortloopFilter) Init() bool {
	sf.configManager.SubscribeToUpdates(sf)
	return true
}

// func Filter(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		nrw := NewResponseWriter(w)

// 		h.ServeHTTP(nrw, r)

// 		func() {
// 			if len(contextArray) < 20 {
// 				contextArray = append(contextArray, NewContext(nrw, r))
// 			}
// 		}()
// 	})
// }

// func (sf *ShortloopFilter) Filter(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		var agentConfigLocal AgentConfig = sf.agentConfig

// 		if reflect.DeepEqual(agentConfigLocal, AgentConfig{}) {
// 			h.ServeHTTP(w, r)
// 			return
// 		}

// 		if !agentConfigLocal.GetCaptureApiSample() {
// 			h.ServeHTTP(w, r)
// 			return
// 		}

// 		var observedApi ObservedApi = sf.getObservedApiFromRequest(r)
// 		nrw := NewResponseWriterWrapper(w)
// 		context := NewRequestResponseContext(nrw, r, sf.userApplicationName)
// 		context.SetObservedApi(observedApi)
// 		context.SetAgentConfig(agentConfigLocal)
// 		fmt.Println("context: ", context)

// 		var apiConfig ApiConfig = sf.getApiConfig(observedApi, agentConfigLocal)
// 		fmt.Println("apiConfig before registered api: ", apiConfig)

//if !reflect.DeepEqual(apiConfig, ApiConfig{}) {
//	context.SetApiConfig(apiConfig)
//	context.SetApiBufferKey(buffer.GetApiBufferKeyFromApiConfig(context.GetApiConfig()))
//	sf.apiProcessor.ProcessRegisteredApi(context, h)
//} else {
//	context.SetApiBufferKey(buffer.GetApiBufferKeyFromObservedApi(context.GetObservedApi()))
//	sf.apiProcessor.ProcessDiscoveredApi(context, h)
//}

// 		//h.ServeHTTP(nrw, r)
// 	})
// }

func (sf *ShortloopFilter) GetApiConfig(observedApi ObservedApi, agentConfigLocal AgentConfig) *ApiConfig {

	if agentConfigLocal.GetRegisteredApiConfigs() == nil {
		sdklogger.Logger.Info("No known APIs as per config, returning NULL")
		return nil
	}
	for _, apiConfig := range agentConfigLocal.GetRegisteredApiConfigs() {
		if observedApi.Matches(apiConfig) {
			return &apiConfig
		}
	}
	return nil
}

func (sf *ShortloopFilter) GetObservedApiFromRequest(r *http.Request) ObservedApi {
	return NewObservedApi(r.URL.Path, HTTPRequestMethod(r.Method))
}

func (sf *ShortloopFilter) OnSuccessfulConfigUpdate(agentConfig AgentConfig) {
	sf.AgentConfig = &agentConfig
}

func (sf *ShortloopFilter) OnErroneousConfigUpdate() {
	sf.AgentConfig = GetNoOpAgentConfig()
}

//var Shortloop ShortloopFilter
//
//func init() {
//	scm := NewConfigManager("http://localhost:8300", http.Client{}, "test-service")
//	scm.Init()
//	Shortloop = NewShortloopFilter(&scm, "test-service")
//	Shortloop.Init()
//}

func (sf *ShortloopFilter) GetUserApplicationName() string {
	return sf.UserApplicationName
}

func (sf *ShortloopFilter) SetUserApplicationName(userApplicationName string) {
	sf.UserApplicationName = userApplicationName
}

func (sf *ShortloopFilter) GetConfigManager() *config.Manager {
	return sf.configManager
}

func (sf *ShortloopFilter) SetConfigManager(configManager *config.Manager) {
	sf.configManager = configManager
}

func (sf *ShortloopFilter) GetApiProcessor() *ApiProcessor {
	return sf.ApiProcessor
}

func (sf *ShortloopFilter) SetApiProcessor(apiProcessor *ApiProcessor) {
	sf.ApiProcessor = apiProcessor
}

func (sf *ShortloopFilter) IsBlackListedApi(observedApi ObservedApi, agentConfig AgentConfig) bool {
	for _, blackListedApi := range agentConfig.BlackListRules {
		if blackListedApi.MatchUri(observedApi.GetUri(), observedApi.Method) {
			return true
		}
	}
	return false
}

func (rrc *RequestResponseContext) GetHttpRequest() *http.Request {
	return rrc.httpRequest
}

func (rrc *RequestResponseContext) GetResponseWriterWrapper() ResponseWriter {
	return rrc.responseWriterWrapper
}

func (rrc *RequestResponseContext) GetApplicationName() string {
	return rrc.applicationName
}

func (rrc *RequestResponseContext) GetObservedApi() ObservedApi {
	return rrc.observedApi
}

func (rrc *RequestResponseContext) GetAgentConfig() AgentConfig {
	return rrc.agentConfig
}

func (rrc *RequestResponseContext) GetApiConfig() *ApiConfig {
	return rrc.apiConfig
}

func (rrc *RequestResponseContext) GetPayloadCaptureAttempted() bool {
	return rrc.payloadCaptureAttempted
}

func (rrc *RequestResponseContext) GetRequestPayloadCaptureAttempted() bool {
	return rrc.requestPayloadCaptureAttempted
}

func (rrc *RequestResponseContext) GetResponsePayloadCaptureAttempted() bool {
	return rrc.responsePayloadCaptureAttempted
}

func (rrc *RequestResponseContext) GetLatency() int64 {
	return rrc.latency
}

func (rrc *RequestResponseContext) GetApiBufferKey() buffer.ApiBufferKey {
	return rrc.apiBufferKey
}

func (rrc *RequestResponseContext) SetHttpRequest(httpRequest *http.Request) {
	rrc.httpRequest = httpRequest
}

func (rrc *RequestResponseContext) SetResponseWriterWrapper(responseWriterWrapper ResponseWriter) {
	rrc.responseWriterWrapper = responseWriterWrapper
}

func (rrc *RequestResponseContext) SetApplicationName(applicationName string) {
	rrc.applicationName = applicationName
}

func (rrc *RequestResponseContext) SetObservedApi(observedApi ObservedApi) {
	rrc.observedApi = observedApi
}

func (rrc *RequestResponseContext) SetAgentConfig(agentConfig AgentConfig) {
	rrc.agentConfig = agentConfig
}

func (rrc *RequestResponseContext) SetApiConfig(apiConfig *ApiConfig) {
	rrc.apiConfig = apiConfig
}

func (rrc *RequestResponseContext) SetPayloadCaptureAttempted(payloadCaptureAttempted bool) {
	rrc.payloadCaptureAttempted = payloadCaptureAttempted
}

func (rrc *RequestResponseContext) SetRequestPayloadCaptureAttempted(requestPayloadCaptureAttempted bool) {
	rrc.requestPayloadCaptureAttempted = requestPayloadCaptureAttempted
}

func (rrc *RequestResponseContext) SetResponsePayloadCaptureAttempted(responsePayloadCaptureAttempted bool) {
	rrc.responsePayloadCaptureAttempted = responsePayloadCaptureAttempted
}

func (rrc *RequestResponseContext) SetLatency(latency int64) {
	rrc.latency = latency
}

func (rrc *RequestResponseContext) SetApiBufferKey(apiBufferKey buffer.ApiBufferKey) {
	rrc.apiBufferKey = apiBufferKey
}

func (rrc *RequestResponseContext) GetRequestPayload() []byte {
	return rrc.requestPayload.Bytes()
}

func (rrc *RequestResponseContext) SetRequestPayload(requestPayload *bytes.Buffer) {
	rrc.requestPayload = requestPayload
}
