package buffer

import (
	"github.com/short-loop/shortloop-common-go/models/data"
)

type ApiBufferKey struct {
	uri    data.URI
	method data.HTTPRequestMethod
}

func GetApiBufferKeyFromObservedApi(observedApi data.ObservedApi) ApiBufferKey {
	return ApiBufferKey{uri: observedApi.GetUri(), method: observedApi.GetMethod()}
}

func GetApiBufferKeyFromApiConfig(apiConfig data.ApiConfig) ApiBufferKey {
	return ApiBufferKey{uri: apiConfig.GetUri(), method: apiConfig.GetMethod()}
}

func (abk ApiBufferKey) Equals(object interface{}) bool {
	if object == nil {
		return false
	}
	other, ok := object.(ApiBufferKey)
	if !ok {
		return false
	}
	return abk.method == other.GetMethod() && abk.uri.Equals(other.GetUri())
}

func (abk ApiBufferKey) GetUri() data.URI {
	return abk.uri
}

func (abk ApiBufferKey) GetMethod() data.HTTPRequestMethod {
	return abk.method
}

func (abk *ApiBufferKey) SetUri(uri data.URI) {
	abk.uri = uri
}

func (abk *ApiBufferKey) SetMethod(method data.HTTPRequestMethod) {
	abk.method = method
}
