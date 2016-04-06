package schema

import ()

type Service struct {
	SvcName      string        `json:"svcName"`
	Version      string        `json:"version"`
	Executor     []string      `json:"executor"`
	Command      string        `json:"command"`
	Args         []string      `json:"args"`
	Environments []Environment `json:"environments"`
	Port         int32         `json:"port"`
	Protocol     string        `json:"protocol"`
	Dependencies []string      `json:"dependencies"`
	Endpoints    []string      `json:"endpoints"`
}
