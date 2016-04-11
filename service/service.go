package service

import (
	"github.com/pingcap/tiadmin/pkg"
	"github.com/pingcap/tiadmin/registry"
)

var Registered map[string]Service

func RegisterServices() {
	Registered = make(map[string]Service)
	Registered[TiDB_SERVICE] = NewTiDBService()
	Registered[PD_SERVICE] = NewPDService()
	Registered[TiKV_SERVICE] = NewTiKVService()
}

func RegisterServciesFromEtcd(reg registry.Registry) {
	// TODO: implement it
}

type Service interface {
	Status() *ServiceStatus
	ParseEndpointFromArgs([]string) map[string]pkg.Endpoint
}

type service struct {
	svcName      string
	version      string
	executor     []string
	command      string
	args         []string
	environments map[string]string
	endpoints    map[string]pkg.Endpoint
}

func (s *service) Status() *ServiceStatus {
	return &ServiceStatus{
		SvcName:      s.svcName,
		Version:      s.version,
		Executor:     s.executor,
		Command:      s.command,
		Args:         s.args,
		Environments: s.environments,
		Endpoints:    s.endpoints,
	}
}
