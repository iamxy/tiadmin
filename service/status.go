package service

import "github.com/pingcap/tiadmin/pkg"

type ServiceStatus struct {
	SvcName      string
	Version      string
	Executor     []string
	Command      string
	Args         []string
	Environments map[string]string
	Endpoints    []pkg.Endpoint
}
