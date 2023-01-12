package buffer

import (
	"github.com/short-loop/shortloop-common-go/models/data"
	"github.com/short-loop/shortloop-sdk-go/config"
	"github.com/short-loop/shortloop-sdk-go/sdklogger"
)

type Manager interface {
	GetWorker() ManagerWorker
	GetUri() string
	createWorker(newConfig *data.AgentConfig) ManagerWorker
	OnSuccessfulConfigUpdate(newAgentConfig data.AgentConfig)
	OnErroneousConfigUpdate()
	Init() bool
	Shutdown() bool
	isRefreshNeeded(olderConfig *data.AgentConfig, newConfig *data.AgentConfig) bool
}

type AbstractManager struct {
	Manager
	configManager *config.Manager
	dummyWorker   ManagerWorker
	worker        ManagerWorker
}

func NewAbstractManager(configManager *config.Manager) *AbstractManager {
	am := AbstractManager{
		configManager: configManager,
	}
	am.dummyWorker = NewNoOpManagerWorker()
	am.worker = am.dummyWorker
	return &am
}

func (am *AbstractManager) GetWorker() ManagerWorker {
	return am.worker
}

func (am *AbstractManager) GetUri() string {
	return "/data-ingestion/api-sample"
}

func (am *AbstractManager) OnSuccessfulConfigUpdate(newAgentConfig data.AgentConfig) {
	if am.isRefreshNeeded(am.worker.GetOperatingConfig(), &newAgentConfig) {
		sdklogger.Logger.Info("Refreshing worker")
		oldManagerWorker := am.worker
		am.worker = am.createWorker(&newAgentConfig)
		oldManagerWorker.Shutdown()
	}
}

func (am *AbstractManager) OnErroneousConfigUpdate() {
	oldManagerWorker := am.worker
	am.worker = am.dummyWorker
	oldManagerWorker.Shutdown()
}

func (am *AbstractManager) Init() bool {
	return am.configManager.SubscribeToUpdates(am)
}

func (am *AbstractManager) Shutdown() bool {
	if am.worker != nil {
		am.worker.Shutdown()
		am.worker = nil
	}
	return true
}

func (am *AbstractManager) isRefreshNeeded(olderConfig *data.AgentConfig, newConfig *data.AgentConfig) bool {
	if newConfig.GetTimestamp().IsZero() && olderConfig.GetTimestamp().IsZero() {
		return false
	}
	return newConfig.GetTimestamp().After(olderConfig.GetTimestamp())
}
