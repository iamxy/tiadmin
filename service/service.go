package service

import (
	"github.com/coreos/fleet/registry"
)

var Registered map[string]Service

func RegisterServices() {
	Registered = make(map[string]Service)
	Registered[TiDB_SERVICE] = NewTiDB()
}

func RegisterServciesFromEtcd(reg registry.Registry) {
	// TODO: implement it
}

type Service interface {
	Status() (*ServiceStatus, error)
}

type service struct {
	svcName      string
	version      string
	executor     []string
	command      string
	args         []string
	environments map[string]string
}

func (s *service) Status() (*ServiceStatus, error) {
	return &ServiceStatus{
		SvcName:      s.svcName,
		Version:      s.version,
		Executor:     s.executor,
		Command:      s.command,
		Args:         s.args,
		Environments: s.environments,
	}, nil
}
