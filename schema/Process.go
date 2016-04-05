package schema

import ()

type Process struct {
	ProcID       string        `json:"procID,omitempty"`
	SvcName      string        `json:"svcName,omitempty"`
	MachID       string        `json:"machID,omitempty"`
	DesiredState string        `json:"desiredState,omitempty"`
	CurrentState string        `json:"currentState,omitempty"`
	IsAlive      bool          `json:"isAlive,omitempty"`
	Endpoints    []string      `json:"endpoints,omitempty"`
	Executor     []string      `json:"executor,omitempty"`
	Command      string        `json:"command,omitempty"`
	Args         []string      `json:"args,omitempty"`
	Environments []Environment `json:"environments,omitempty"`
	PublicIP     string        `json:"publicIP,omitempty"`
	HostName     string        `json:"hostName,omitempty"`
	HostMeta     HostMeta      `json:"hostMeta,omitempty"`
	Port         int32         `json:"port,omitempty"`
	Protocol     string        `json:"protocol,omitempty"`
}
