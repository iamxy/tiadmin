package service

import (
	"github.com/pingcap/tiadmin/pkg"
	proc "github.com/pingcap/tiadmin/process"
	"github.com/pingcap/tiadmin/registry"
)

type TiDBService struct {
}

func NewTiDBService(reg registry.Registry) Service {
	return &TiDBService{}
}

func (s *TiDBService) UpdateConfig() error {
	return nil
}
func (s *TiDBService) Status() (ServiceStatus, error) {
	var status ServiceStatus
	return status, nil
}
func (s *TiDBService) Endpoints() ([]pkg.Endpoint, error) {
	return nil, nil
}
func (s *TiDBService) RunUpProcess(args ...string) (proc.ProcessStatus, error) {
	var status proc.ProcessStatus
	return status, nil
}
func (s *TiDBService) KillProcess(procID string) error {
	return nil
}
func (s *TiDBService) TriggerStartProcess(procID string) error {
	return nil
}
func (s *TiDBService) TriggerStopProcess(procID string) error {
	return nil
}
func (s *TiDBService) ListProcess() ([]proc.ProcessStatus, error) {
	return nil, nil
}
