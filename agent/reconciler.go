package agent

import (
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/pkg"
	proc "github.com/pingcap/tiadmin/process"
	"github.com/pingcap/tiadmin/registry"
	"time"
)

const (
	// time between triggering reconciliation routine
	ReconcileInterval = 5 * time.Second
)

func NewReconciler(reg registry.Registry, evtStream pkg.EventStream, procMgr *proc.ProcessManager) *AgentReconciler {
	return &AgentReconciler{
		reg:       reg,
		evtStream: evtStream,
		procMgr:   procMgr,
		clock:     clockwork.NewRealClock(),
	}
}

type AgentReconciler struct {
	reg       registry.Registry
	evtStream pkg.EventStream
	procMgr   *proc.ProcessManager
	clock     clockwork.Clock
}

func (ar *AgentReconciler) Run(stopc <-chan struct{}) {
	// When starting up, reconcile once immediately
	ar.reconcile()
	// Execute periodic reconciling
	for {
		ticker := ar.clock.After(ReconcileInterval)
		select {
		case <-stopc:
			log.Debug("Reconciler exiting due to stop signal")
			return
		case <-ticker:
			log.Debug("Reconciler tick")
			ar.reconcile()
		case <-ar.evtStream.Next(stopc):
			log.Debug("Reconciler triggered")
			ar.reconcile()
		}
	}
}

func (ar *AgentReconciler) reconcile() {
	start := time.Now()

	elapsed := time.Now().Sub(start)
	msg := fmt.Sprintf("AgentReconciler completed reconciliation in %s", elapsed)
	if elapsed > ReconcileInterval {
		log.Warning(msg)
	} else {
		log.Debug(msg)
	}
}
