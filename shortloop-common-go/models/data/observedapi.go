package data

import (
	"fmt"
	"reflect"
)

type ObservedApi struct {
	Uri    URI               `json:"uri"`
	Method HTTPRequestMethod `json:"method"`
}

func (o ObservedApi) GetUri() URI {
	return o.Uri
}

func (o ObservedApi) GetMethod() HTTPRequestMethod {
	return o.Method
}

func (o *ObservedApi) SetUri(uri URI) {
	o.Uri = uri
}

func (o *ObservedApi) SetMethod(method HTTPRequestMethod) {
	o.Method = method
}

func NewObservedApi(uri string, method HTTPRequestMethod) ObservedApi {
	return ObservedApi{Uri: GetNonTemplatedURI(uri), Method: method}
}

func (o *ObservedApi) Matches(apiConfig ApiConfig) bool {
	if reflect.DeepEqual(apiConfig, ApiConfig{}) {
		return false
	}
	if o.Method != apiConfig.GetMethod() {
		return false
	}

	return o.Uri.Equals(apiConfig.GetUri())
}

func (o ObservedApi) String() string {
	return fmt.Sprintf("ObservedApi{Uri=%s, Method=%s}", o.Uri, o.Method)
}
