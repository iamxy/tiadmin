package service

import (
	"github.com/pingcap/tiadmin/pkg"
	proc "github.com/pingcap/tiadmin/process"
	"github.com/pingcap/tiadmin/registry"
)

type Service interface {
	UpdateConfig() error
	Status() (ServiceStatus, error)
	Endpoints() ([]pkg.Endpoint, error)
	RunUpProcess(args ...string) (proc.ProcessStatus, error)
	KillProcess(procID string) error
	TriggerStartProcess(procID string) error
	TriggerStopProcess(procID string) error
	ListProcess() []proc.ProcessStatus
}

type ServiceManager map[string]Service

func RegisterServices(reg registry.Registry) ServiceManager {
	svcMgr := make(ServiceManager)
	svcMgr["tidb"] = NewTiDBService(reg)
	return svcMgr
}
