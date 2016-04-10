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
	Mach      machine.Machine
	TTL       time.Duration
	publishch chan []string
}

func (a *Agent) subscribe(procIDs []string) {
	if procIDs != nil && len(procIDs) > 0 {
		a.publishch <- procIDs
	}
}

func (a *Agent) publish() chan []string {
	return a.publishch
}

func New(reg registry.Registry, pm proc.ProcMgr, m machine.Machine, ttl time.Duration) *Agent {
	return &Agent{
		Reg:       reg,
		ProcMgr:   pm,
		Mach:      m,
		TTL:       ttl,
		publishch: make(chan []string),
	}
}

func (a *Agent) StartNewProcess(machID, svcName string, runinfo *proc.ProcessRunInfo) error {
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
	if mach, err := a.Reg.Machine(machID); err == nil {
		hostIP = mach.MachInfo.PublicIP
		hostName = mach.MachInfo.HostName
		hostRegion = mach.MachInfo.HostRegion
		hostIDC = mach.MachInfo.HostIDC
		// check if the target machine is offline
		if !mach.IsAlive {
			e := fmt.Sprintf("Should not start new processes on a offline host, machID: %s, svcName: %s", machID, svcName)
			log.Error(e)
			return errors.New(e)
		}
	} else {
		return err
	}

	if svc, ok := service.Registered[svcName]; ok {
		ss := svc.Status()
		if len(runinfo.Executor) > 0 {
			executor = runinfo.Executor
		} else {
			executor = ss.Executor
		}
		if len(runinfo.Command) > 0 {
			command = runinfo.Command
		} else {
			command = ss.Command
		}
		if len(runinfo.Args) > 0 {
			args = runinfo.Args
		} else {
			args = ss.Args
		}
		if len(runinfo.Environment) > 0 {
			env = runinfo.Environment
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

func (a *Agent) ListAllMachines() (res map[string]*machine.MachineStatus, err error) {
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

func (a *Agent) BirthCry() error {
	status := a.Mach.Status()
	if err := a.Reg.RegisterMachine(status.MachID, status.MachInfo.HostName, status.MachInfo.HostRegion,
		status.MachInfo.HostIDC, status.MachInfo.PublicIP); err != nil {
		log.Errorf("Register machine status into etcd failed, %v", err)
		return err
	}
	return nil
}
