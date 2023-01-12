package data

import (
	"fmt"
	"time"
)

type AgentConfig struct {
	BufferSyncFreqInSec       int             `json:"bufferSyncFreqInSec"`
	CaptureApiSample          bool            `json:"captureApiSample"`
	ConfigFetchFreqInSec      int             `json:"configFetchFreqInSec"`
	RegisteredApiConfigs      []ApiConfig     `json:"registeredApiConfigs"`
	Timestamp                 time.Time       `json:"timestamp"`
	DiscoveryBufferSize       int64           `json:"discoveryBufferSize"`
	DiscoveryBufferSizePerApi int             `json:"discoveryBufferSizePerApi"`
	BlackListRules            []BlackListRule `json:"blackListRules"`
}

type ApiConfig struct {
	Uri                   URI               `json:"uri"`
	Method                HTTPRequestMethod `json:"method"`
	BufferSize            int               `json:"bufferSize"`
	CaptureSampleRequest  bool              `json:"captureSampleRequest"`
	CaptureSampleResponse bool              `json:"captureSampleResponse"`
}

func (a ApiConfig) GetUri() URI {
	return a.Uri
}

func (a ApiConfig) GetMethod() HTTPRequestMethod {
	return a.Method
}

func (a ApiConfig) GetBufferSize() int {
	return a.BufferSize
}

func (a ApiConfig) GetCaptureSampleRequest() bool {
	return a.CaptureSampleRequest
}

func (a ApiConfig) GetCaptureSampleResponse() bool {
	return a.CaptureSampleResponse
}

func (a *ApiConfig) SetUri(uri URI) {
	a.Uri = uri
}

func (a *ApiConfig) SetMethod(method HTTPRequestMethod) {
	a.Method = method
}

func (a *ApiConfig) SetBufferSize(bufferSize int) {
	a.BufferSize = bufferSize
}

func (a *ApiConfig) SetCaptureSampleRequest(captureSampleRequest bool) {
	a.CaptureSampleRequest = captureSampleRequest
}

func (a *ApiConfig) SetCaptureSampleResponse(captureSampleResponse bool) {
	a.CaptureSampleResponse = captureSampleResponse
}

func (a AgentConfig) GetBufferSyncFreqInSec() int {
	return a.BufferSyncFreqInSec
}

func (a AgentConfig) GetCaptureApiSample() bool {
	return a.CaptureApiSample
}

func (a AgentConfig) GetConfigFetchFreqInSec() int {
	return a.ConfigFetchFreqInSec
}

func (a AgentConfig) GetRegisteredApiConfigs() []ApiConfig {
	return a.RegisteredApiConfigs
}

func (a AgentConfig) GetTimestamp() time.Time {
	return a.Timestamp
}

func (a AgentConfig) GetDiscoveryBufferSize() int64 {
	return a.DiscoveryBufferSize
}

func (a AgentConfig) GetDiscoveryBufferSizePerApi() int {
	return a.DiscoveryBufferSizePerApi
}

func (a *AgentConfig) SetBufferSyncFreqInSec(bufferSyncFreqInSec int) {
	a.BufferSyncFreqInSec = bufferSyncFreqInSec
}

func (a *AgentConfig) SetCaptureApiSample(captureApiSample bool) {
	a.CaptureApiSample = captureApiSample
}

func (a *AgentConfig) SetConfigFetchFreqInSec(configFetchFreqInSec int) {
	a.ConfigFetchFreqInSec = configFetchFreqInSec
}

func (a *AgentConfig) SetRegisteredApiConfigs(registeredApiConfigs []ApiConfig) {
	a.RegisteredApiConfigs = registeredApiConfigs
}

func (a *AgentConfig) SetTimestamp(timestamp time.Time) {
	a.Timestamp = timestamp
}

func (a *AgentConfig) SetDiscoveryBufferSize(discoveryBufferSize int64) {
	a.DiscoveryBufferSize = discoveryBufferSize
}

func (a *AgentConfig) SetDiscoveryBufferSizePerApi(discoveryBufferSizePerApi int) {
	a.DiscoveryBufferSizePerApi = discoveryBufferSizePerApi
}

func GetNoOpAgentConfig() *AgentConfig {
	return &AgentConfig{
		BufferSyncFreqInSec:       120,
		CaptureApiSample:          false,
		ConfigFetchFreqInSec:      120,
		DiscoveryBufferSize:       0,
		DiscoveryBufferSizePerApi: 0,
		RegisteredApiConfigs:      []ApiConfig{},
		Timestamp:                 time.Time{},
	}
}

//
//func (a *ApiConfig) UnmarshalJSON(data []byte) error {
//	// Uri           URI
//	// Method        HTTPRequestMethod
//	// BufferSize    int
//	// CaptureSample bool
//
//	var tmpJson map[string]interface{}
//
//	if err := json.Unmarshal(data, &tmpJson); err != nil {
//		return err
//	}
//
//	uri, ok := tmpJson["Uri"].(string)
//	if !ok {
//		return errors.New("Uri is missing")
//	}
//
//	method, ok := tmpJson["Method"].(string)
//	if !ok {
//		return errors.New("Method is missing")
//	}
//
//	bufferSize, ok := tmpJson["BufferSize"].(float64)
//	if !ok {
//		return errors.New("BufferSize is missing")
//	}
//
//	captureSample, ok := tmpJson["CaptureSample"].(bool)
//	if !ok {
//		return errors.New("CaptureSample is missing")
//	}
//
//	var uriObj URI
//	if err := json.Unmarshal([]byte(uri), &uriObj); err != nil {
//		return err
//	}
//	a.SetUri(uriObj)
//	a.SetMethod(HTTPRequestMethod(method))
//	a.SetBufferSize(int(bufferSize))
//	a.SetCaptureSample(captureSample)
//
//	return nil
//}

//func (a *AgentConfig) UnmarshalJSON(data []byte) error {
//
//	var tmpJson map[string]interface{}
//
//	if err := json.Unmarshal(data, &tmpJson); err != nil {
//		return err
//	}
//	// fmt.Println(tmpJson, tmpJson["CaptureApiSample"])
//
//	bufferSyncFreqInSec, ok := tmpJson["BufferSyncFreqInSec"].(float64)
//	if !ok {
//		return errors.New("BufferSyncFreqInSec is missing")
//	}
//
//	captureApiSample, ok := tmpJson["CaptureApiSample"].(bool)
//	if !ok {
//		return errors.New("CaptureApiSample is missing")
//	}
//
//	configFetchFreqInSec, ok := tmpJson["ConfigFetchFreqInSec"].(float64)
//	if !ok {
//		return errors.New("ConfigFetchFreqInSec is missing")
//	}
//
//	discoveryBufferSize, ok := tmpJson["DiscoveryBufferSize"].(float64)
//	if !ok {
//		return errors.New("DiscoveryBufferSize is missing")
//	}
//
//	discoveryBufferSizePerApi, ok := tmpJson["DiscoveryBufferSizePerApi"].(float64)
//	if !ok {
//		return errors.New("DiscoveryBufferSizePerApi is missing")
//	}
//
//	registeredApiConfigs, ok := tmpJson["RegisteredApiConfigs"].([]interface{})
//	if !ok {
//		return errors.New("RegisteredApiConfigs is missing")
//	}
//
//	timestamp, ok := tmpJson["Timestamp"].(float64)
//	if !ok {
//		return errors.New("Timestamp is missing")
//	}
//
//	a.SetBufferSyncFreqInSec(int(bufferSyncFreqInSec))
//	a.SetCaptureApiSample(captureApiSample)
//	a.SetConfigFetchFreqInSec(int(configFetchFreqInSec))
//	a.SetDiscoveryBufferSize(int(discoveryBufferSize))
//	a.SetDiscoveryBufferSizePerApi(int(discoveryBufferSizePerApi))
//	a.SetTimestamp(time.Unix(int64(timestamp), 0))
//
//	var apiConfigs []ApiConfig
//	for _, apiConfig := range registeredApiConfigs {
//
//		apiConfigMap := apiConfig.(map[string]interface{})
//
//		uri, ok := apiConfigMap["Uri"].(map[string]interface{})
//		if !ok {
//			return errors.New("Uri is missing")
//		}
//
//		method, ok := apiConfigMap["Method"].(string)
//		if !ok {
//			return errors.New("Method is missing")
//		}
//
//		bufferSize, ok := apiConfigMap["BufferSize"].(float64)
//		if !ok {
//			return errors.New("BufferSize is missing")
//		}
//
//		captureSample, ok := apiConfigMap["CaptureSample"].(bool)
//		if !ok {
//			return errors.New("CaptureSample is missing")
//		}
//
//		uriPath, ok := uri["UriPath"].(string)
//		if !ok {
//			return errors.New("UriPath is not string")
//		}
//		hasPathVariable, ok := uri["HasPathVariable"].(bool)
//		if !ok {
//			return errors.New("HasPathVariable is not bool")
//		}
//
//		var u URI
//
//		u.SetURIPath(uriPath)
//		u.SetHasPathVariable(hasPathVariable)
//
//		var apiConfigObj ApiConfig
//
//		apiConfigObj.SetUri(u)
//		apiConfigObj.SetMethod(HTTPRequestMethod(method))
//		apiConfigObj.SetBufferSize(int(bufferSize))
//		apiConfigObj.SetCaptureSample(captureSample)
//
//		apiConfigs = append(apiConfigs, apiConfigObj)
//	}
//	a.SetRegisteredApiConfigs(apiConfigs)
//
//	return nil
//	// for key, value := range tmpJson {
//	// 	switch key {
//	// 	case "BufferSyncFreqInSec":
//	// 		a.SetBufferSyncFreqInSec(int(value.(float64)))
//	// 	case "CaptureApiSample":
//	// 		a.SetCaptureApiSample(value.(bool))
//	// 	case "ConfigFetchFreqInSec":
//	// 		a.SetConfigFetchFreqInSec(int(value.(float64)))
//	// 	case "DiscoveryBufferSize":
//	// 		a.SetDiscoveryBufferSize(int(value.(float64)))
//	// 	case "DiscoveryBufferSizePerApi":
//	// 		a.SetDiscoveryBufferSizePerApi(int(value.(float64)))
//	// 	case "RegisteredApiConfigs":
//	// 		var apiConfigs []ApiConfig
//	// 		for _, apiConfig := range value.([]interface{}) {
//	// 			var apiConfigObj ApiConfig
//	// 			for key, value := range apiConfig.(map[string]interface{}) {
//	// 				switch key {
//	// 				case "Uri":
//	// 					var Uri URI
//	// 					for key, value := range value.(map[string]interface{}) {
//	// 						switch key {
//	// 						case "HasPathVariable":
//	// 							Uri.SetHasPathVariable(value.(bool))
//	// 						case "UriPath":
//	// 							Uri.SetURIPath(value.(string))
//	// 						}
//	// 					}
//	// 					apiConfigObj.SetUri(Uri)
//	// 				case "Method":
//	// 					apiConfigObj.SetMethod(HTTPRequestMethod(value.(string)))
//	// 				case "BufferSize":
//	// 					apiConfigObj.SetBufferSize(int(value.(float64)))
//	// 				case "CaptureSample":
//	// 					apiConfigObj.SetCaptureSample(value.(bool))
//	// 				}
//	// 			}
//	// 			apiConfigs = append(apiConfigs, apiConfigObj)
//	// 		}
//	// 		a.SetRegisteredApiConfigs(apiConfigs)
//	// 	case "Timestamp":
//	// 		a.SetTimestamp(time.Unix(int64(value.(float64)), 0))
//	// 	}
//	// }
//}

func (a AgentConfig) String() string {
	return fmt.Sprintf("AgentConfig{BufferSyncFreqInSec=%d, CaptureApiSample=%t, ConfigFetchFreqInSec=%d, RegisteredApiConfigs=%v, Timestamp=%v, DiscoveryBufferSize=%d, DiscoveryBufferSizePerApi=%d}", a.BufferSyncFreqInSec, a.CaptureApiSample, a.ConfigFetchFreqInSec, a.RegisteredApiConfigs, a.Timestamp, a.DiscoveryBufferSize, a.DiscoveryBufferSizePerApi)
}
