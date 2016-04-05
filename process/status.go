package process

import (
	"github.com/pingcap/tiadmin/pkg"
)

type ProcessStatus struct {
	ProcID       string
	SvcName      string
	MachID       string
	DesiredState ProcessState
	CurrentState ProcessState
	IsAlive      bool
	Endpoints    []pkg.Endpoint
	RunInfo      ProcessRunInfo
}

type ProcessRunInfo struct {
	HostIP      string
	Executor    []string
	Command     string
	Args        []string
	Environment map[string]string
	Port        pkg.Port
	Protocol    pkg.Protocol
}
