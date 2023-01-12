package config

import (
	"encoding/json"
	"github.com/short-loop/shortloop-common-go/models/data"
	"github.com/short-loop/shortloop-sdk-go/executor"
	"github.com/short-loop/shortloop-sdk-go/sdklogger"
	"github.com/short-loop/shortloop-sdk-go/sdkversion"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var currentConfigManager *Manager = &Manager{}

type Manager struct {
	ctUrl               string
	httpClient          http.Client
	userApplicationName string
	scheduledExecutor   executor.ScheduledExecutor
	agentId             string
	updateListeners     []UpdateListener
	sync.Mutex
}

func CurrentConfigManager() *Manager {
	return currentConfigManager
}

func (m *Manager) Init() {
	m.agentId = strconv.Itoa(rand.Intn(10000))
	m.scheduledExecutor = executor.NewScheduledExecutor()
	m.scheduleConfigRefresh(60)
	sdklogger.Logger.Info("ConfigManager initialized and scheduler started to fetch config")
}

func (m *Manager) scheduleConfigRefresh(timeInSecs int) {
	m.scheduledExecutor.Schedule(
		func() {
			agentConfig, errCode := m.fetchConfig()
			if agentConfig == nil || errCode != NIL {
				m.scheduleConfigRefresh(300)
				m.onUnSuccessfulConfigFetch(errCode)
			} else {
				m.scheduleConfigRefresh(agentConfig.GetConfigFetchFreqInSec())
				m.onSuccessfulConfigFetch(*agentConfig)
			}
		}, time.Duration(timeInSecs)*time.Second)
}

func (m *Manager) getUri() string {
	return "/api/v1/agent-config"
}

func (m *Manager) fetchConfig() (*data.AgentConfig, ErrorCode) {
	defer func() {
		if err := recover(); err != nil {
			sdklogger.Logger.ErrorF("Panic while fetching config: %s\n", err)
		}
	}()
	m.Lock()
	defer m.Unlock()

	req, err := http.NewRequest("GET", m.ctUrl+m.getUri(), nil)
	req.Header.Set(sdkversion.MAJOR_VERSION_KEY, sdkversion.MAJOR_VERSION)
	req.Header.Set(sdkversion.MINOR_VERSION_KEY, sdkversion.MINOR_VERSION)
	req.Header.Set("sdkType", sdkversion.SdkType)
	q := req.URL.Query()
	q.Add("appName", m.userApplicationName)
	q.Add("agentId", m.agentId)
	req.URL.RawQuery = q.Encode()

	response, err := m.httpClient.Do(req)
	// fmt.Println("fetched config: ", response, err)
	if err != nil {
		sdklogger.Logger.ErrorF("Error while fetching config from CT: %+v\n", err)
		return nil, TIMEOUT
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			sdklogger.Logger.ErrorF("Error while closing response body: %+v\n", err)
		}
	}(response.Body)
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		// fmt.Println("bodyBytes: ", string(bodyBytes), "err: ", err, "len: ", len(bodyBytes))
		if err != nil {
			sdklogger.Logger.ErrorF("Error while reading response body of /agent-config: %+v\n", err)
			return nil, PARSE_ERROR
		}
		if bodyBytes == nil || len(bodyBytes) == 0 {
			sdklogger.Logger.ErrorF("Empty response body received from /agent-config: %+v\n", err)
			return nil, INVALID_CONFIG
		}
		var agentConfig data.AgentConfig = data.AgentConfig{}
		err = json.Unmarshal(bodyBytes, &agentConfig)
		if err != nil {
			sdklogger.Logger.ErrorF("Error while unmarshalling response body of /agent-config: %+v\n", err)
			return nil, PARSE_ERROR
		}
		return &agentConfig, NIL
	} else {
		sdklogger.Logger.ErrorF("Error while fetching config from CT, Status Code: %+v\n", response.StatusCode)
		return nil, TIMEOUT
	}
}

func (m *Manager) onSuccessfulConfigFetch(agentConfig data.AgentConfig) {
	sdklogger.Logger.InfoF("onSuccessfulConfigFetch config from CT: %+v\n", agentConfig)
	for _, listener := range m.updateListeners {
		listener.OnSuccessfulConfigUpdate(agentConfig)
	}
}

func (m *Manager) onUnSuccessfulConfigFetch(errCode ErrorCode) {
	sdklogger.Logger.ErrorF("onUnSuccessfulConfigFetch config from CT: %+v\n", errCode)
	for _, listener := range m.updateListeners {
		listener.OnErroneousConfigUpdate()
	}
}

func (m *Manager) shutdown() {
	sdklogger.Logger.Info("Shutting down ConfigManager")
	err := m.scheduledExecutor.Shutdown()
	if err != nil {
		sdklogger.Logger.ErrorF("Error while shutting down ConfigManager: %+v\n", err)
		return
	}
}

func (m *Manager) SubscribeToUpdates(listener UpdateListener) bool {
	m.updateListeners = append(m.updateListeners, listener)
	return true
}

type ErrorCode string

const (
	TIMEOUT        ErrorCode = "TIMEOUT"
	PARSE_ERROR    ErrorCode = "PARSE_ERROR"
	INVALID_CONFIG ErrorCode = "INVALID_CONFIG"
	NIL            ErrorCode = "NIL"
)

func (m *Manager) SetCtUrl(ctUrl string) {
	m.ctUrl = ctUrl
}

func (m *Manager) SetHttpClient(httpClient http.Client) {
	m.httpClient = httpClient
}

func (m *Manager) SetUserApplicationName(userApplicationName string) {
	m.userApplicationName = userApplicationName
}

func (m *Manager) GetAgentId() string {
	return m.agentId
}
