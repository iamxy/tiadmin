package process

import "sync"

type ProcessManager struct {
	procs   map[string]Proc
	rwMutex sync.Mutex
}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		procs: make(map[string]Proc),
	}
}

func (pm *ProcessManager) CreateProcess(svcName string, procID string) (*Process, error) {
	return nil, nil
}

func (pm *ProcessManager) DestroyProcess(procID string) error {
	return nil
}

func (pm *ProcessManager) StartProcess(procID string) error {
	return nil
}

func (pm *ProcessManager) StopProcess(procID string) error {
	return nil
}

func (pm *ProcessManager) AllProcess() []Proc {
	return nil
}

func (pm *ProcessManager) AllActiveProcess() []Proc {
	return nil
}

func (pm *ProcessManager) TotalProcess() int {
	return nil
}

func (pm *ProcessManager) TotalActiveProcess() int {
	return nil
}

func (pm *ProcessManager) FindByProcID(procID string) Proc {
	return nil
}

func (pm *ProcessManager) FindBySvcName(svcName string) []Proc {
	return nil
}
