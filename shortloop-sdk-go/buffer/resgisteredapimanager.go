package buffer

import (
	"bytes"
	"encoding/json"
	"github.com/short-loop/shortloop-common-go/models/data"
	"github.com/short-loop/shortloop-sdk-go/config"
	"github.com/short-loop/shortloop-sdk-go/httpconnection"
	"github.com/short-loop/shortloop-sdk-go/sdklogger"
	"golang.org/x/sync/semaphore"
	"net/http"
)

type RegisteredApiManager struct {
	*AbstractManager
	httpClient http.Client
	ctUrl      string
}

func NewRegisteredApiManager(configManager *config.Manager, httpClient http.Client, ctUrl string) *RegisteredApiManager {
	m := NewAbstractManager(configManager)
	//m := &AbstractManager{configManager: configManager}
	ram := &RegisteredApiManager{
		AbstractManager: m,
		httpClient:      httpClient,
		ctUrl:           ctUrl,
	}
	m.Manager = ram
	return ram
}

func (ram *RegisteredApiManager) createWorker(newConfig *data.AgentConfig) ManagerWorker {
	return NewRegisteredApiManagerWorker(newConfig, ram.httpClient, ram.ctUrl)
}

type RegisteredApiManagerWorker struct {
	*AbstractManagerWorker
	httpClient http.Client
	semaphore  *semaphore.Weighted
}

func NewRegisteredApiManagerWorker(config *data.AgentConfig, httpClient http.Client, ctUrl string) *RegisteredApiManagerWorker {
	worker := NewAbstractManagerWorker(config, ctUrl)
	//worker := &AbstractManagerWorker{agentConfig: config, ctUrl: ctUrl}
	ramw := &RegisteredApiManagerWorker{
		AbstractManagerWorker: worker,
		httpClient:            httpClient,
	}
	ramw.semaphore = semaphore.NewWeighted(ramw.GetRegisteredApiCountToCapture())
	worker.ManagerWorker = ramw
	return ramw
}

func (ramw *RegisteredApiManagerWorker) GetRegisteredApiCountToCapture() int64 {
	agentConfig := ramw.GetOperatingConfig()
	if len(agentConfig.GetRegisteredApiConfigs()) == 0 {
		return 0
	}
	var totalApis int64 = 0
	for _, registeredApiConfig := range agentConfig.GetRegisteredApiConfigs() {
		totalApis += int64(registeredApiConfig.GetBufferSize())
	}
	return totalApis
}

func (ramw *RegisteredApiManagerWorker) Offer(apiBufferKey ApiBufferKey, apiSample data.APISample) bool {
	value, _ := ramw.bufferMap.GetOrUpdate(apiBufferKey, func() interface{} {
		return NewSimpleBuffer(ramw.GetRegisteredApiBufferSize(apiBufferKey))
	})
	if value != nil {
		buffer := value.(SimpleBuffer)
		return buffer.Offer(&apiSample)
	} else {
		sdklogger.Logger.ErrorF("RegisteredApiManagerWorker.Offer: buffer is nil for uri: %+v\n", apiBufferKey.GetUri())
	}
	return false
}

func (ramw *RegisteredApiManagerWorker) GetRegisteredApiBufferSize(apiBufferKey ApiBufferKey) int {
	agentConfig := ramw.GetOperatingConfig()
	if len(agentConfig.GetRegisteredApiConfigs()) == 0 {
		return 0
	}
	for _, registeredApiConfig := range agentConfig.GetRegisteredApiConfigs() {
		if registeredApiConfig.GetMethod() == apiBufferKey.GetMethod() && registeredApiConfig.GetUri().Equals(apiBufferKey.GetUri()) {
			return registeredApiConfig.GetBufferSize()
		}
	}
	return 0
}

func (ramw *RegisteredApiManagerWorker) CanOffer(apiBufferKey ApiBufferKey) bool {
	if !ramw.GetOperatingConfig().GetCaptureApiSample() {
		return false
	}
	if ramw.GetRegisteredApiBufferSize(apiBufferKey) == 0 {
		return false
	}
	if ramw.semaphore.TryAcquire(1) {
		value, found := ramw.bufferMap.Get(apiBufferKey)
		var canOffer bool = false
		if !found {
			canOffer = true
		} else {
			buffer := value.(SimpleBuffer)
			canOffer = buffer.CanOffer()
		}
		ramw.semaphore.Release(1)
		return canOffer
	}
	return false
}

func (ramw *RegisteredApiManagerWorker) syncForKey(apiBufferKey ApiBufferKey) {

	defer func() {
		if err := recover(); err != nil {
			sdklogger.Logger.ErrorF("Panic in RegisteredApiManagerWorker::syncForKey: %s\n", err)
		}
	}()

	value, found := ramw.bufferMap.Get(apiBufferKey)
	if !found {
		sdklogger.Logger.ErrorF("Buffer not found for %+v\n", apiBufferKey.GetUri())
		return
	}
	buffer := value.(SimpleBuffer)
	var iterations int = buffer.GetContentCount()
	sdklogger.Logger.InfoF("Syncing %d samples of registered api %+v\n", iterations, apiBufferKey.GetUri())
	for i := 0; i < iterations; i++ {
		apiSample := buffer.Poll()
		if apiSample == nil {
			ramw.bufferMap.Delete(apiBufferKey)
			break
		}

		apiSampleList := [1]data.APISample{*apiSample}
		jsonBodyContent, err := json.Marshal(apiSampleList)
		if err != nil {
			sdklogger.Logger.ErrorF("Error marshalling api sample list: %s\n", err.Error())
			return
		}

		bodyReader := bytes.NewReader(jsonBodyContent)

		req, err := http.NewRequest("POST", ramw.ctUrl+ramw.GetUri(), bodyReader)
		if err != nil {
			sdklogger.Logger.ErrorF("Error creating request: %s\n", err.Error())
			return
		}
		//_, err = ramw.httpClient.Do(req)
		httpconnection.SendRequest(&ramw.httpClient, req)

		//if err != nil {
		//	sdklogger.Logger.ErrorF("Error sending request: %s\n", err.Error())
		//	return
		//}
	}
	sdklogger.Logger.InfoF("Synced %d samples for registered api %+v\n", iterations, apiBufferKey.GetUri())
}
