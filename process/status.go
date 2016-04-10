package process

import (
	"github.com/pingcap/tiadmin/pkg"
)

type ProcessState string

const (
	StateStarted = ProcessState("StateStarted")
	StateStopped = ProcessState("StateStopped")
)

func (s ProcessState) String() string {
	return string(s)
}

func (s ProcessState) Opposite() ProcessState {
	if s == StateStarted {
		return StateStopped
	} else {
		return StateStarted
	}
}

type ProcessStatus struct {
	ProcID       string
	SvcName      string
	MachID       string
	DesiredState ProcessState
	CurrentState ProcessState
	IsAlive      bool
	RunInfo      ProcessRunInfo
}

type ProcessRunInfo struct {
	HostIP      string
	HostName    string
	HostRegion  string
	HostIDC     string
	Executor    []string
	Command     string
	Args        []string
	Environment map[string]string
	Endpoints   map[string]pkg.Endpoint
}
