package process

import (
	"github.com/pingcap/tiadmin/pkg"
)

type ProcessState string

func (s ProcessState) String() string {
	return string(s)
}

const (
	StateStarted = ProcessState("StateStarted")
	StateStopped = ProcessState("StateStopped")
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
	HostName    string
	HostRegion  string
	HostIDC     string
	Executor    []string
	Command     string
	Args        []string
	Environment map[string]string
	Port        pkg.Port
	Protocol    pkg.Protocol
}
