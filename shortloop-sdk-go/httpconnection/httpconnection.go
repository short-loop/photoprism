package httpconnection

import (
	"github.com/short-loop/shortloop-sdk-go/sdklogger"
	"github.com/short-loop/shortloop-sdk-go/sdkversion"
	"net/http"
)

func SendRequest(httpClient *http.Client, httpRequest *http.Request) {
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set(sdkversion.MAJOR_VERSION_KEY, sdkversion.MAJOR_VERSION)
	httpRequest.Header.Set(sdkversion.MINOR_VERSION_KEY, sdkversion.MINOR_VERSION)
	httpRequest.Header.Set("sdkType", sdkversion.SdkType)

	response, err := httpClient.Do(httpRequest)
	if err != nil {
		sdklogger.Logger.ErrorF("Error sending request: %s\n", err.Error())
		return
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			sdklogger.Logger.ErrorF("Error closing connection: %s\n", err.Error())
			return
		}
	}()
	return
}
