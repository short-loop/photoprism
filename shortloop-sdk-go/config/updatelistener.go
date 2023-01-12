package config

import (
	"github.com/short-loop/shortloop-common-go/models/data"
)

type UpdateListener interface {
	OnSuccessfulConfigUpdate(agentConfig data.AgentConfig)
	OnErroneousConfigUpdate()
}
