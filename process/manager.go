package process

import (
	"errors"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/pkg"
	"sync"
)

type ProcMgr interface {
	CreateProcess(*ProcessStatus) (Proc, error)
	DestroyProcess(string) error
	StartProcess(string) error
	StopProcess(string) error
	AllProcess() map[string]Proc
	AllActiveProcess() map[string]Proc
	TotalProcess() int
	TotalActiveProcess() int
	FindByProcID(string) Proc
	FindBySvcName(string) map[string]Proc
}

type ProcessManager struct {
	procs   map[string]Proc
	rwMutex sync.RWMutex
}

func NewProcessManager() ProcMgr {
	return &ProcessManager{
		procs: make(map[string]Proc),
	}
}

func (pm *ProcessManager) CreateProcess(target *ProcessStatus) (Proc, error) {
	metadata := make(map[string]string) // reserve
	var pwd string
	if root, err := pkg.GetRootDir(); err != nil {
		log.Errorf("GetRootDir failed, error: %v", err)
		return nil, err
	} else {
		pwd = root
	}

	if proc, err := NewProcess(target.ProcID, target.SvcName, target.RunInfo.Executor, target.RunInfo.Command,
		target.RunInfo.Args, "", "", target.RunInfo.Environment, metadata, pwd); err == nil {
		if target.DesiredState == StateStarted {
			if err := proc.Start(); err != nil {
				log.Errorf("Failed to start local process, procID: %s, error: %v", target.ProcID, err)
				return nil, err
			}
		}
		pm.rwMutex.Lock()
		defer pm.rwMutex.Unlock()
		pm.procs[target.ProcID] = proc
		return proc, nil
	} else {
		log.Errorf("Failed to create new local process, procID: %s, error: %v", target.ProcID, err)
		return nil, err
	}
}

func (pm *ProcessManager) DestroyProcess(procID string) (err error) {
	pm.rwMutex.RLock()
	proc, ok := pm.procs[procID]
	pm.rwMutex.RUnlock()
	if !ok {
		return errors.New("Failed to destroy a local process, the procID not exists: " + procID)
	}

	if proc.State() == StateStarted {
		if err := proc.Stop(); err != nil {
			log.Errorf("Failed to stop local process, procID: %s, error: %v", procID, err)
			return err
		}
	}

	pm.rwMutex.Lock()
	delete(pm.procs, procID)
	pm.rwMutex.Unlock()
	return nil
}

func (pm *ProcessManager) StartProcess(procID string) error {
	pm.rwMutex.RLock()
	proc, ok := pm.procs[procID]
	pm.rwMutex.RUnlock()
	if !ok {
		return errors.New("Failed to start a local process, the procID not exists: " + procID)
	}
	if proc.State() == StateStopped {
		if err := proc.Start(); err != nil {
			log.Errorf("Failed to start local process, procID: %s, error: %v", procID, err)
			return err
		}
	} else {
		log.Warnf("Process is already started, no need to be stopped, procID: %s", procID)
	}
	return nil
}

func (pm *ProcessManager) StopProcess(procID string) error {
	pm.rwMutex.RLock()
	proc, ok := pm.procs[procID]
	pm.rwMutex.RUnlock()
	if !ok {
		return errors.New("Failed to start a local process, the procID not exists: " + procID)
	}
	if proc.State() == StateStarted {
		if err := proc.Stop(); err != nil {
			log.Errorf("Failed to stop local process, procID: %s, error: %v", procID, err)
			return err
		}
	} else {
		log.Warnf("Process is already started, no need to be stopped, procID: %s", procID)
	}
	return nil
}

func (pm *ProcessManager) AllProcess() map[string]Proc {
	res := make(map[string]Proc)
	pm.rwMutex.RLock()
	defer pm.rwMutex.RUnlock()
	for k, v := range pm.procs {
		res[k] = v
	}
	return res
}

func (pm *ProcessManager) AllActiveProcess() map[string]Proc {
	res := make(map[string]Proc)
	pm.rwMutex.RLock()
	defer pm.rwMutex.RUnlock()
	for k, v := range pm.procs {
		if v.IsActive() {
			res[k] = v
		}
	}
	return res
}

func (pm *ProcessManager) TotalProcess() int {
	pm.rwMutex.RLock()
	defer pm.rwMutex.RUnlock()
	return len(pm.procs)
}

func (pm *ProcessManager) TotalActiveProcess() int {
	res := 0
	pm.rwMutex.RLock()
	defer pm.rwMutex.RUnlock()
	for _, v := range pm.procs {
		if v.IsActive() {
			res++
		}
	}
	return res
}

func (pm *ProcessManager) FindByProcID(procID string) Proc {
	pm.rwMutex.RLock()
	defer pm.rwMutex.RUnlock()
	if proc, ok := pm.procs[procID]; ok {
		return proc
	} else {
		return nil
	}
}

func (pm *ProcessManager) FindBySvcName(svcName string) map[string]Proc {
	res := make(map[string]Proc)
	pm.rwMutex.RLock()
	defer pm.rwMutex.RUnlock()
	for k, v := range pm.procs {
		if v.GetSvcName() == svcName {
			res[k] = v
		}
	}
	return res
}
