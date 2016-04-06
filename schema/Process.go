package schema

import ()

type Process struct {
	ProcID       string        `json:"procID"`
	SvcName      string        `json:"svcName"`
	MachID       string        `json:"machID"`
	DesiredState string        `json:"desiredState"`
	CurrentState string        `json:"currentState"`
	IsAlive      bool          `json:"isAlive"`
	Endpoints    []string      `json:"endpoints"`
	Executor     []string      `json:"executor"`
	Command      string        `json:"command"`
	Args         []string      `json:"args"`
	Environments []Environment `json:"environments"`
	PublicIP     string        `json:"publicIP"`
	HostName     string        `json:"hostName"`
	HostMeta     HostMeta      `json:"hostMeta"`
	Port         int32         `json:"port"`
	Protocol     string        `json:"protocol"`
}
