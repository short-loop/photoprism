package buffer

import (
	"github.com/short-loop/shortloop-common-go/models/data"
	"github.com/short-loop/shortloop-sdk-go/executor"
	"github.com/short-loop/shortloop-sdk-go/sdklogger"
	"time"
)

type ManagerWorker interface {
	GetOperatingConfig() *data.AgentConfig
	Offer(apiBufferKey ApiBufferKey, apiSample data.APISample) bool
	CanOffer(apiBufferKey ApiBufferKey) bool
	Shutdown() bool
	cleanUpBufferMap()
	getUri()
	syncForKey(apiBufferKey ApiBufferKey)
	syncForKeys()
}

type AbstractManagerWorker struct {
	ManagerWorker
	scheduledExecutor executor.ScheduledExecutor
	bufferMap         BufferMap
	agentConfig       *data.AgentConfig
	ctUrl             string
}

func (m *AbstractManagerWorker) GetOperatingConfig() *data.AgentConfig {
	return m.agentConfig
}

func NewAbstractManagerWorker(agentConfig *data.AgentConfig, ctUrl string) *AbstractManagerWorker {
	m := AbstractManagerWorker{
		scheduledExecutor: executor.NewScheduledExecutor(),
		bufferMap:         BufferMap{},
		ctUrl:             ctUrl,
		agentConfig:       agentConfig,
	}
	//fmt.Println("NewAbstractManagerWorker: ", m)
	m.scheduledExecutor.ScheduleAtFixedRate(m.syncForKeys, time.Duration(agentConfig.GetBufferSyncFreqInSec())*time.Second)
	return &m
}

func (m *AbstractManagerWorker) Shutdown() bool {

	defer func() {
		if err := recover(); err != nil {
			sdklogger.Logger.ErrorF("Panic in Worker::Shutdown: %s\n", err)
		}
	}()

	sdklogger.Logger.Info("Shutting down ApiSample BufferManagerWorker")
	err := m.scheduledExecutor.Shutdown()
	if err != nil {
		sdklogger.Logger.ErrorF("Error shutting down scheduled executor: %s", err.Error())
		return false
	}
	m.cleanUpBufferMap()
	return true
}

func (m *AbstractManagerWorker) cleanUpBufferMap() {
	sdklogger.Logger.Info("Cleaning up buffer map...")
	m.syncForKeys()
	m.bufferMap.Range(func(key, value interface{}) bool {
		buffer := value.(SimpleBuffer)
		buffer.Clear()
		m.bufferMap.Delete(key)
		return true
	})
}

func (m *AbstractManagerWorker) GetUri() string {
	return "/api/v1/data-ingestion/api-sample"
}

func (m *AbstractManagerWorker) syncForKeys() {
	defer func() {
		if err := recover(); err != nil {
			sdklogger.Logger.ErrorF("Panic in syncForKeys: %s\n", err)
		}
	}()

	sdklogger.Logger.Info("Syncing apis...")
	keys := m.bufferMap.getKeys()
	for _, key := range keys {
		m.syncForKey(key.(ApiBufferKey))
	}
}

type NoOpManagerWorker struct {
	*AbstractManagerWorker
	agentConfig *data.AgentConfig
}

func NewNoOpManagerWorker() *NoOpManagerWorker {
	worker := &AbstractManagerWorker{agentConfig: data.GetNoOpAgentConfig()}
	nmw := &NoOpManagerWorker{
		AbstractManagerWorker: worker,
		agentConfig:           data.GetNoOpAgentConfig(),
	}
	nmw.ManagerWorker = worker
	return nmw
}

func (m *NoOpManagerWorker) Init() bool {
	return true
}

func (m *NoOpManagerWorker) GetOperatingConfig() *data.AgentConfig {
	return m.agentConfig
}

func (m *NoOpManagerWorker) Offer(apiBufferKey ApiBufferKey, apiSample data.APISample) bool {
	return false
}

func (m *NoOpManagerWorker) CanOffer(apiBufferKey ApiBufferKey) bool {
	return false
}

func (m *NoOpManagerWorker) Shutdown() bool {
	return true
}

func (m *NoOpManagerWorker) syncForKey(apiBufferKey ApiBufferKey) {
	return
}
