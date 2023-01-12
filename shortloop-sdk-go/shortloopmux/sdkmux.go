package shortloopmux

import (
	"fmt"
	. "github.com/short-loop/shortloop-common-go/models/data"
	"github.com/short-loop/shortloop-sdk-go/buffer"
	"github.com/short-loop/shortloop-sdk-go/config"
	"github.com/short-loop/shortloop-sdk-go/sdklogger"
	"github.com/short-loop/shortloop-sdk-go/sdkversion"
	"github.com/short-loop/shortloop-sdk-go/shortloopfilter"
	"net/http"
)

type Options struct {
	ShortloopEndpoint string
	ApplicationName   string
	LoggingEnabled    bool
	LogLevel          string
}

type ShortloopMux struct {
	shortloopFilter *shortloopfilter.ShortloopFilter
}

func (shortloopMux *ShortloopMux) Filter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sf := shortloopMux.shortloopFilter

		var agentConfigLocal *AgentConfig = sf.AgentConfig

		if nil == agentConfigLocal {
			h.ServeHTTP(w, r)
			return
		}

		if !agentConfigLocal.GetCaptureApiSample() {
			h.ServeHTTP(w, r)
			return
		}

		var observedApi ObservedApi = sf.GetObservedApiFromRequest(r)
		if sf.IsBlackListedApi(observedApi, *agentConfigLocal) {
			h.ServeHTTP(w, r)
			return
		}

		nrw := NewResponseWriterWrapper(w)
		context := shortloopfilter.NewRequestResponseContext(nrw, r, sf.UserApplicationName)
		context.SetObservedApi(observedApi)
		context.SetAgentConfig(*agentConfigLocal)

		var apiConfig *ApiConfig = sf.GetApiConfig(observedApi, *agentConfigLocal)

		if apiConfig != nil {
			context.SetApiConfig(apiConfig)
			context.SetApiBufferKey(buffer.GetApiBufferKeyFromApiConfig(*context.GetApiConfig()))
			sdklogger.Logger.InfoF("ApiConfig found for observedApi: %+v\n", observedApi.GetUri())
			sf.ApiProcessor.ProcessRegisteredApi(context, func(canOffer bool) {
				if canOffer {
					h.ServeHTTP(nrw, context.GetHttpRequest())
				} else {
					h.ServeHTTP(w, context.GetHttpRequest())
				}
				return
			})
		} else {
			sdklogger.Logger.InfoF("ApiConfig not found for observedApi: %+v\n", observedApi.GetUri())
			context.SetApiBufferKey(buffer.GetApiBufferKeyFromObservedApi(context.GetObservedApi()))
			sf.ApiProcessor.ProcessDiscoveredApi(context, func(canOffer bool) {
				if canOffer {
					h.ServeHTTP(nrw, context.GetHttpRequest())
				} else {
					h.ServeHTTP(w, context.GetHttpRequest())
				}
				return
			})
		}

		// h.ServeHTTP(nrw, r)
	})
}

func Init(options Options) (*ShortloopMux, error) {
	shortloopMux := &ShortloopMux{shortloopfilter.CurrentShortloopFilter()}

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
	sdkversion.SdkType = "Go-Mux"
	sdklogger.Logger.Info("Initialized Shortloop SDK")
	sdklogger.Logger.InfoF("Initialized Shortloop SDK\napplication name: %v\nurl: %v\nagent id:%v\nSDK Version: %v.%v\nSDKType: %v\n",
		options.ApplicationName,
		options.ShortloopEndpoint,
		configManager.GetAgentId(),
		sdkversion.MAJOR_VERSION,
		sdkversion.MINOR_VERSION,
		sdkversion.SdkType)
	return shortloopMux, nil
}
