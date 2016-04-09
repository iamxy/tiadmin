package service

import (
	"github.com/pingcap/tiadmin/registry"
)

var Registered map[string]Service

func RegisterServices() {
	Registered = make(map[string]Service)
	Registered[TiDB_SERVICE] = NewTiDBService()
}

func RegisterServciesFromEtcd(reg registry.Registry) {
	// TODO: implement it
}

type Service interface {
	Status() *ServiceStatus
}

type service struct {
	svcName      string
	version      string
	executor     []string
	command      string
	args         []string
	environments map[string]string
}

func (s *service) Status() *ServiceStatus {
	return &ServiceStatus{
		SvcName:      s.svcName,
		Version:      s.version,
		Executor:     s.executor,
		Command:      s.command,
		Args:         s.args,
		Environments: s.environments,
	}
}
