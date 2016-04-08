package agent

import (
	"errors"
	"fmt"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/machine"
	"github.com/pingcap/tiadmin/pkg"
	proc "github.com/pingcap/tiadmin/process"
	"github.com/pingcap/tiadmin/registry"
	"github.com/pingcap/tiadmin/service"
	"time"
)

type Agent struct {
	Reg       registry.Registry
	ProcMgr   proc.ProcMgr
	Mach      machine.Mach
	TTL       time.Duration
	publishch chan []string
}

func New(reg registry.Registry, pm proc.ProcMgr, m machine.Mach, ttl time.Duration) *Agent {
	return &Agent{
		Reg:       reg,
		ProcMgr:   pm,
		Mach:      m,
		TTL:       ttl,
		publishch: make(chan []string),
	}
}

func (a *Agent) Subscribe(procIDs []string) {
	if procIDs != nil && len(procIDs) > 0 {
		a.publishch <- procIDs
	}
}

func (a *Agent) Publish() chan []string {
	return a.publishch
}

func (a *Agent) StartNewProcess(status proc.ProcessStatus) error {
	machID := status.MachID
	svcName := status.SvcName
	var hostIP string
	var hostName string
	var hostRegion string
	var hostIDC string
	var executor []string
	var command string
	var args []string
	var env map[string]string
	var port pkg.Port
	var protocol pkg.Protocol

	// retrieve machine infomation from etcd
	if mach, err := a.Reg.Machine(machID); err == nil && mach != nil {
		hostIP = mach.MachInfo.PublicIP
		hostName = mach.MachInfo.HostName
		hostRegion = mach.MachInfo.HostRegion
		hostIDC = mach.MachInfo.HostIDC
	} else {
		return err
	}

	if svc, ok := service.Registered[svcName]; ok {
		ss, err := svc.Status()
		if err != nil {
			return err
		}
		executor = ss.Executor
		command = ss.Command
		if len(status.RunInfo.Args) > 0 {
			args = status.RunInfo.Args
		} else {
			args = ss.Args
		}
		if len(status.RunInfo.Environment) > 0 {
			env = status.RunInfo.Environment
		} else {
			env = ss.Environments
		}
	} else {
		e := fmt.Sprintf("Unregistered service: %s", svcName)
		log.Error(e)
		return errors.New(e)
	}

	if err := a.Reg.NewProcess(machID, svcName, hostIP, hostName, hostRegion, hostIDC,
		executor, command, args, env, port, protocol); err != nil {
		e := fmt.Sprintf("Create new process failed in etcd, %s, %s, %v", machID, svcName, err)
		log.Error(e)
		return errors.New(e)
	}
	return nil
}

func (a *Agent) DestroyProcess(procID string) error {
	_, err := a.Reg.DeleteProcess(procID)
	if err != nil {
		log.Errorf("Delete process failed in etcd, %s, %v", procID, err)
	}
	return err
}

func (a *Agent) StartProcess(procID string) error {
	err := a.Reg.UpdateProcessDesiredState(procID, proc.StateStarted)
	if err != nil {
		log.Errorf("Change desired state of process to started failed, %s, %v", procID, err)
	}
	return err
}

func (a *Agent) StopProcess(procID string) error {
	err := a.Reg.UpdateProcessDesiredState(procID, proc.StateStopped)
	if err != nil {
		log.Errorf("Change desired state of process to stopped failed, %s, %v", procID, err)
	}
	return err
}

func (a *Agent) ListAllProcesses() (res map[string]*proc.ProcessStatus, err error) {
	res, err = a.Reg.Processes()
	if err != nil {
		log.Errorf("List all processes failed, %v", err)
	}
	return
}

func (a *Agent) ListProcessesByMachID(machID string) (res map[string]*proc.ProcessStatus, err error) {
	res, err = a.Reg.ProcessesOnMachine(machID)
	if err != nil {
		log.Errorf("List processes on specified machine, %s, %v", machID, err)
	}
	return
}

func (a *Agent) ListProcessesBySvcName(svcName string) (res map[string]*proc.ProcessStatus, err error) {
	res, err = a.Reg.ProcessesOfService(svcName)
	if err != nil {
		log.Errorf("List processes of specified service, %s, %v", svcName, err)
	}
	return
}

func (a *Agent) ListProcess(procID string) (res *proc.ProcessStatus, err error) {
	res, err = a.Reg.Process(procID)
	if err != nil {
		log.Errorf("List specified process failed, %s, %v", procID, err)
	}
	return
}

func (a *Agent) ListMachines() (res map[string]*machine.MachineStatus, err error) {
	res, err = a.Reg.Machines()
	if err != nil {
		log.Errorf("List all machines in cluster failed, %v", err)
	}
	return
}

func (a *Agent) ListMachine(machID string) (res *machine.MachineStatus, err error) {
	res, err = a.Reg.Machine(machID)
	if err != nil {
		log.Errorf("List specified machines infomation failed, %s, %v", machID, err)
	}
	return
}
