package agent

import (
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/pkg"
	proc "github.com/pingcap/tiadmin/process"
	"github.com/pingcap/tiadmin/registry"
	"strings"
	"time"
)

const (
	// time between triggering reconciliation routine
	reconcileInterval = 5 * time.Second
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
	for {
		select {
		case <-stopc:
			log.Debug("Reconciler is exiting due to stop signal")
			return
		//case <-ar.clock.After(reconcileInterval):
		//	log.Debug("Trigger reconciling from tick")
		//	if err := ar.reconcile(); err != nil {
		//		log.Errorf("Reconcile failed, %v", err)
		//	}
		case event := <-ar.eStream.Next(reconcileInterval):
			if event.None() {
				log.Debug("Reconciling is triggered by tick")
			} else {
				log.Debugf("Reconciling is triggered by event, %v", event)
			}
			if err := ar.reconcile(); err != nil {
				log.Errorf("Failed to reconcile, %v", err)
			}
		}
	}
}

func (ar *AgentReconciler) reconcile() error {
	start := time.Now()
	toPublish, err := ar.doReconcile()
	if err != nil {
		return err
	}
	ar.agent.subscribe(toPublish)
	elapsed := time.Now().Sub(start)
	msg := fmt.Sprintf("Reconciling completed in %s", elapsed)
	if elapsed > reconcileInterval {
		log.Warning(msg)
	} else {
		log.Debug(msg)
	}
	return nil
}

func (ar *AgentReconciler) doReconcile() ([]string, error) {
	toPublish := make([]string, 0)
	allProcesses, err := ar.reg.Processes()
	if err != nil {
		return nil, err
	}
	targetProcesses, endpoints := prepareProcesses(allProcesses, ar.agent.Mach.ID())
	currentProcesses := ar.agent.ProcMgr.AllProcess()
	endpoints["ETCD_ADDR"] = ar.reg.GetEtcdAddrs()

	var checked = make(map[string]struct{}, 0)
	for procID, procStatus := range targetProcesses {
		process, ok := currentProcesses[procID]
		if ok {
			checked[procID] = struct{}{}
			if procStatus.DesiredState == proc.StateStarted && process.State() == proc.StateStopped {
				if err := ar.agent.ProcMgr.StartProcess(procID, endpoints); err != nil {
					log.Errorf("Failed to start local process, procID: %s", procID)
					return nil, err
				}
				toPublish = append(toPublish, procID)
			}
			if procStatus.DesiredState == proc.StateStopped && process.State() == proc.StateStarted {
				if err := ar.agent.ProcMgr.StopProcess(procID); err != nil {
					log.Errorf("Failed to stop local process, procID: %s", procID)
					return nil, err
				}
				toPublish = append(toPublish, procID)
			}
		} else {
			// local process not exists, create one
			proc, err := ar.agent.ProcMgr.CreateProcess(procStatus, endpoints)
			if err != nil {
				log.Errorf("Failed to create new local process, %v", procStatus)
				return nil, err
			}
			log.Infof("Create local process successfully, procID: %s, with state: %v", proc.GetProcID(), proc.State())
			toPublish = append(toPublish, procID)
		}
	}

	for procID, _ := range currentProcesses {
		if _, ok := checked[procID]; !ok {
			if err := ar.agent.ProcMgr.DestroyProcess(procID); err != nil {
				log.Errorf("Failed to destroy local process, procID: %s", procID)
				return nil, err
			}
			log.Infof("Destroy local process successfully, procID: %s", procID)
		}
	}
	return toPublish, nil
}

func prepareProcesses(allProcs map[string]*proc.ProcessStatus, machID string) (map[string]*proc.ProcessStatus, map[string]string) {
	procsOnMach := make(map[string]*proc.ProcessStatus)
	temp := make(map[string][]string)
	for k, v := range allProcs {
		if v.MachID == machID {
			procsOnMach[k] = v
		}
		for name, ep := range v.RunInfo.Endpoints {
			if _, ok := temp[name]; ok {
				temp[name] = append(temp[name], ep.String())
			} else {
				temp[name] = []string{ep.String()}
			}
		}
	}
	endpoints := make(map[string]string)
	for k, v := range temp {
		endpoints[k] = strings.Join(v, ",")
	}
	return procsOnMach, endpoints
}
