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
		args = ss.Args
		env = ss.Environments
	} else {
		e := fmt.Sprintf("Unregistered service: %s", svcName)
		log.Error(e)
		return errors.New(e)
	}

	if err := a.Reg.NewProcess(machID, svcName, hostIP, hostName, hostRegion, hostIDC,
		executor, command, args, env, port, protocol); err != nil {
		e := fmt.Sprintf("Create new process failed in etcd, %v", err)
		log.Error(e)
		return errors.New(e)
	}
	return nil
}
