package schema

import ()

type Service struct {
	SvcName      string        `json:"svcName,omitempty"`
	Version      string        `json:"version,omitempty"`
	Executor     []string      `json:"executor,omitempty"`
	Command      string        `json:"command,omitempty"`
	Args         []string      `json:"args,omitempty"`
	Environments []Environment `json:"environments,omitempty"`
	Port         int32         `json:"port,omitempty"`
	Protocol     string        `json:"protocol,omitempty"`
	Dependencies []string      `json:"dependencies,omitempty"`
	Endpoints    []string      `json:"endpoints,omitempty"`
}
