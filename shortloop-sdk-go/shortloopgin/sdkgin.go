package shortloopgin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	. "github.com/short-loop/shortloop-common-go/models/data"
	"github.com/short-loop/shortloop-sdk-go/buffer"
	"github.com/short-loop/shortloop-sdk-go/config"
	"github.com/short-loop/shortloop-sdk-go/sdklogger"
	"github.com/short-loop/shortloop-sdk-go/sdkversion"
	"github.com/short-loop/shortloop-sdk-go/shortloopfilter"
	"net/http"
	"os"
)

type Options struct {
	ShortloopEndpoint string
	ApplicationName   string
	LoggingEnabled    bool
	LogLevel          string
}

type ShortloopGin struct {
	shortloopFilter *shortloopfilter.ShortloopFilter
}

func (shortloopGin *ShortloopGin) Filter() gin.HandlerFunc {
	return func(c *gin.Context) {
		sf := shortloopGin.shortloopFilter

		var agentConfigLocal *AgentConfig = sf.AgentConfig

		if nil == agentConfigLocal {
			c.Next()
			return
		}

		if !agentConfigLocal.GetCaptureApiSample() {
			c.Next()
			return
		}

		var observedApi ObservedApi = sf.GetObservedApiFromRequest(c.Request)
		if sf.IsBlackListedApi(observedApi, *agentConfigLocal) {
			c.Next()
			return
		}
		nrw := NewResponseWriterWrapper(c.Writer)
		//c.Writer = nrw
		// nrw := NewResponseWriterWrapper(w)
		context := shortloopfilter.NewRequestResponseContext(nrw, c.Request, sf.UserApplicationName)
		context.SetObservedApi(observedApi)
		context.SetAgentConfig(*agentConfigLocal)

		var apiConfig *ApiConfig = sf.GetApiConfig(observedApi, *agentConfigLocal)

		if apiConfig != nil {
			context.SetApiConfig(apiConfig)
			context.SetApiBufferKey(buffer.GetApiBufferKeyFromApiConfig(*context.GetApiConfig()))
			sdklogger.Logger.InfoF("ApiConfig found for observedApi: %+v\n", observedApi.GetUri())
			sf.ApiProcessor.ProcessRegisteredApi(context, func(canOffer bool, responsePayloadCaptureAttempted bool) {
				if canOffer {
					nrw.SetShouldCaptureBody(responsePayloadCaptureAttempted)
					c.Writer = nrw
				}
				c.Next()
				return
			})
		} else {
			sdklogger.Logger.InfoF("ApiConfig not found for observedApi: %+v\n", observedApi.GetUri())
			context.SetApiBufferKey(buffer.GetApiBufferKeyFromObservedApi(context.GetObservedApi()))
			sf.ApiProcessor.ProcessDiscoveredApi(context, func(canOffer bool) {
				if canOffer {
					c.Writer = nrw
				}
				c.Next()
				return
			})
		}

		// h.ServeHTTP(nrw, r)
	}
}

func Init(options Options) (*ShortloopGin, error) {

	if os.Getenv("GOARCH") == "386" {
		return nil, fmt.Errorf("32 bit Arch not supported by shortloop sdk")
	}

	shortloopGin := &ShortloopGin{shortloopfilter.CurrentShortloopFilter()}
	if options.ShortloopEndpoint == "" {
		return nil, fmt.Errorf("ShortloopEndpoint is required")
	}
	if options.ApplicationName == "" {
		return nil, fmt.Errorf("ApplicationName is required")
	}

	loggingEnabled := options.LoggingEnabled
	logLevel := "ERROR"
	if options.LogLevel != "" {
		logLevel = options.LogLevel
	}

	sdklogger.Logger.SetLoggingEnabled(loggingEnabled)
	sdklogger.Logger.SetLogLevel(sdklogger.GetLogLevel(logLevel))

	sdklogger.Logger.Info("Initializing Shortloop SDK")

	configManager := config.CurrentConfigManager()
	configManager.SetCtUrl(options.ShortloopEndpoint)
	configManager.SetUserApplicationName(options.ApplicationName)
	configManager.SetHttpClient(http.Client{})
	configManager.Init()

	discoveredApiManager := buffer.NewDiscoveredApiManager(configManager, http.Client{}, options.ShortloopEndpoint)
	discoveredApiManager.Init()
	registeredApiManager := buffer.NewRegisteredApiManager(configManager, http.Client{}, options.ShortloopEndpoint)
	registeredApiManager.Init()
	apiProcessor := shortloopfilter.NewApiProcessor(discoveredApiManager, registeredApiManager)

	shortloopFilter := shortloopfilter.CurrentShortloopFilter()
	shortloopFilter.SetUserApplicationName(options.ApplicationName)
	shortloopFilter.SetConfigManager(configManager)
	shortloopFilter.SetApiProcessor(apiProcessor)
	shortloopFilter.Init()
	sdkversion.SdkType = "Go-Gin"
	sdklogger.Logger.InfoF("Initialized Shortloop SDK\napplication name: %v\nurl: %v\nagent id:%v\nSDK Version: %v.%v\nSDKType: %v\n",
		options.ApplicationName,
		options.ShortloopEndpoint,
		configManager.GetAgentId(),
		sdkversion.MAJOR_VERSION,
		sdkversion.MINOR_VERSION,
		sdkversion.SdkType)
	return shortloopGin, nil
}
