package agent

import (
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/machine"
	"github.com/pingcap/tiadmin/pkg"
	proc "github.com/pingcap/tiadmin/process"
	"github.com/pingcap/tiadmin/registry"
	"time"
)

const (
	// time between triggering reconciliation routine
	ReconcileInterval = 5 * time.Second
)

func NewReconciler(reg registry.Registry, es pkg.EventStream, ag *Agent) *AgentReconciler {
	return &AgentReconciler{
		reg:     reg,
		eStream: es,
		agent:   ag,
		clock:   clockwork.NewRealClock(),
	}
}

type AgentReconciler struct {
	reg     registry.Registry
	eStream pkg.EventStream
	agent   *Agent
	clock   clockwork.Clock
}

func (ar *AgentReconciler) Run(stopc <-chan struct{}) {
	// execute reconciling once immediately
	if err := ar.reconcile(); err != nil {
		log.Fatalf("Reconciling run failed at first time, %v", err)
	}

	// reconciling loop
	for {
		select {
		case <-stopc:
			log.Debug("Reconciler is exiting due to stop signal")
			return
		case <-ar.clock.After(ReconcileInterval):
			log.Debug("Trigger reconciling from tick")
			if err := ar.reconcile(); err != nil {
				log.Errorf("Reconcile failed, %v", err)
			}
		case <-ar.eStream.Next(stopc):
			log.Debug("Trigger reconciling fome etcd watcher")
			if err := ar.reconcile(); err != nil {
				log.Errorf("Reconcile failed, %v", err)
			}
		}
	}
}

func (ar *AgentReconciler) reconcile() error {
	start := time.Now()

	toPublish, err := doReconcile(ar.reg, ar.eStream, ar.agent.Mach, ar.agent.ProcMgr)
	if err != nil {
		return err
	}
	ar.agent.Subscribe(toPublish)

	elapsed := time.Now().Sub(start)
	msg := fmt.Sprintf("Reconciling completed in %s", elapsed)
	if elapsed > ReconcileInterval {
		log.Warning(msg)
	} else {
		log.Debug(msg)
	}
	return nil
}

func doReconcile(reg registry.Registry, es pkg.EventStream, mach machine.Machine, procMgr proc.ProcMgr) ([]string, error) {
	// collect the procs which state changes and needed to be published to etcd
	toPublish := make([]string, 0)
	targetProcesses, err := reg.ProcessesOnMachine(mach.ID())
	if err != nil {
		return toPublish, err
	}
	currentProcesses := procMgr.AllProcess()

	for procID, procStatus := range targetProcesses {
		process := procMgr.FindByProcID(procID)
		if process == nil {
			// local process not exists, create new one
			if _, err := procMgr.CreateProcess(procStatus); err != nil {
				log.Errorf("Failed to create new local process, %v", procStatus)
				return toPublish, err
			}
			log.Infof("Create new local process successfully, procID: %s", procID)
			toPublish = append(toPublish, procID)
		} else {
			delete(currentProcesses, procID)
			if procStatus.DesiredState == proc.StateStarted && process.State() == proc.StateStopped {
				if err := procMgr.StartProcess(procID); err != nil {
					log.Errorf("Failed to start local process, procID: %s", procID)
					return toPublish, err
				}
				toPublish = append(toPublish, procID)
			}
			if procStatus.DesiredState == proc.StateStopped && process.State() == proc.StateStarted {
				if err := procMgr.StopProcess(procID); err != nil {
					log.Errorf("Failed to stop local process, procID: %s", procID)
					return toPublish, err
				}
				toPublish = append(toPublish, procID)
			}
		}
	}

	for procID, _ := range currentProcesses {
		if err := procMgr.DestroyProcess(procID); err != nil {
			log.Errorf("Failed to destroy local process, procID: %s", procID)
			return toPublish, err
		}
		//toPublish = append(toPublish, procID)
	}
	return toPublish, nil
}
