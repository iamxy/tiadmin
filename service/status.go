package service

import "github.com/pingcap/tiadmin/pkg"

type ServiceStatus struct {
	SvcName      string
	Version      string
	Executor     []string
	Command      string
	Args         []string
	Environments map[string]string
	Endpoints    map[string]pkg.Endpoint
}

type TiDBPerfMetrics struct {
	TPS         int64  `json:"tps"`
	Connections int    `json:"connections"`
	Version     string `json:"version"`
}
