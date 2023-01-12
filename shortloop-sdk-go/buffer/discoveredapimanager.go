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

type DiscoveredApiManager struct {
	*AbstractManager
	httpClient http.Client
	ctUrl      string
}

func NewDiscoveredApiManager(configManager *config.Manager, httpClient http.Client, ctUrl string) *DiscoveredApiManager {
	m := NewAbstractManager(configManager)
	//m := &AbstractManager{configManager: configManager}
	dam := &DiscoveredApiManager{
		AbstractManager: m,
		httpClient:      httpClient,
		ctUrl:           ctUrl,
	}
	m.Manager = dam
	return dam
}

func (dam *DiscoveredApiManager) createWorker(newConfig *data.AgentConfig) ManagerWorker {
	return NewDiscoveredApiManagerWorker(newConfig, dam.httpClient, dam.ctUrl)
}

type DiscoveredApiManagerWorker struct {
	*AbstractManagerWorker
	httpClient http.Client
	semaphore  *semaphore.Weighted
}

func NewDiscoveredApiManagerWorker(config *data.AgentConfig, httpClient http.Client, ctUrl string) *DiscoveredApiManagerWorker {
	worker := NewAbstractManagerWorker(config, ctUrl)
	//worker := &AbstractManagerWorker{agentConfig: config, ctUrl: ctUrl}
	damw := &DiscoveredApiManagerWorker{
		AbstractManagerWorker: worker,
		httpClient:            httpClient,
		semaphore:             semaphore.NewWeighted(config.GetDiscoveryBufferSize()),
	}
	worker.ManagerWorker = damw
	//return damw
	return damw
}

func (damw *DiscoveredApiManagerWorker) Offer(apiBufferKey ApiBufferKey, apiSample data.APISample) bool {
	value, _ := damw.bufferMap.GetOrUpdate(apiBufferKey, func() interface{} {
		return NewSimpleBuffer(damw.GetOperatingConfig().GetDiscoveryBufferSizePerApi())
	})
	if value != nil {
		buffer := value.(SimpleBuffer)
		return buffer.Offer(&apiSample)
	} else {
		sdklogger.Logger.InfoF("buffer not found for api buffer key = %+v\n", apiBufferKey)
	}
	return false
}

func (damw *DiscoveredApiManagerWorker) CanOffer(apiBufferKey ApiBufferKey) bool {
	if !damw.GetOperatingConfig().GetCaptureApiSample() {
		return false
	}
	if damw.semaphore.TryAcquire(1) {
		var canOffer bool = false
		value, found := damw.bufferMap.Get(apiBufferKey)
		if found {
			buffer := value.(SimpleBuffer)
			canOffer = buffer.CanOffer()
		} else {
			var bufferLen int64 = damw.bufferMap.Len()
			canOffer = bufferLen < damw.GetOperatingConfig().GetDiscoveryBufferSize()
		}
		damw.semaphore.Release(1)
		return canOffer
	}
	return false
}

func (damw *DiscoveredApiManagerWorker) syncForKey(apiBufferKey ApiBufferKey) {

	defer func() {
		if err := recover(); err != nil {
			sdklogger.Logger.ErrorF("Panic in DiscoveredApiManagerWorker::syncForKey: %s\n", err)
		}
	}()

	value, found := damw.bufferMap.Get(apiBufferKey)
	if !found {
		sdklogger.Logger.ErrorF("Buffer not found for %+v\n", apiBufferKey.GetUri())
		return
	}
	buffer := value.(SimpleBuffer)
	var iterations int = buffer.GetContentCount()
	if iterations == 0 {
		damw.bufferMap.Delete(apiBufferKey)
		return
	}

	sdklogger.Logger.InfoF("Syncing %d samples of discovered api %+v\n", iterations, apiBufferKey.GetUri())

	var apiSamples []data.APISample
	for i := 0; i < iterations; i++ {
		apiSample := buffer.Poll()
		if apiSample == nil {
			damw.bufferMap.Delete(apiBufferKey)
			break
		}
		apiSamples = append(apiSamples, *apiSample)
	}
	jsonBodyContent, err := json.Marshal(apiSamples)
	if err != nil {
		sdklogger.Logger.ErrorF("Error marshalling api sample list: %s\n", err.Error())
		return
	}
	bodyReader := bytes.NewReader(jsonBodyContent)

	req, err := http.NewRequest("POST", damw.ctUrl+damw.GetUri(), bodyReader)
	if err != nil {
		sdklogger.Logger.ErrorF("Error creating request: %s\n", err.Error())
		return
	}

	//_, err = damw.httpClient.Do(req)
	//fmt.Println("fetched config: ", response, err)
	httpconnection.SendRequest(&damw.httpClient, req)
	//if err != nil {
	//	sdklogger.Logger.ErrorF("Error sending request: %s\n", err.Error())
	//	return
	//}
	sdklogger.Logger.InfoF("Synced %d samples for discovered api %+v\n", iterations, apiBufferKey.GetUri())
}
