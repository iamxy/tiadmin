package agent

import (
	"github.com/jonboulle/clockwork"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/registry"
	"time"
)

func NewProcessStatePublisher(reg registry.Registry, ag *Agent, ttl time.Duration) *ProcessStatePublisher {
	return &ProcessStatePublisher{
		reg:   reg,
		agent: ag,
		clock: clockwork.NewRealClock(),
		ttl:   ttl,
	}
}

type ProcessStatePublisher struct {
	reg   registry.Registry
	agent *Agent
	clock clockwork.Clock
	ttl   time.Duration
}

func (p *ProcessStatePublisher) Run(stopc <-chan struct{}) {
	for {
		select {
		case <-stopc:
			log.Debug("ProcessStatePublisher is exiting due to stop signal")
			return
		case <-p.clock.After(p.ttl / 2):
			log.Debug("Trigger ProcessStatePublisher after tick")
			if err := p.doPublishAll(); err != nil {
				log.Errorf("ProcessStatePublisher failed, %v", err)
			}
		case pub := <-p.agent.publish():
			log.Debug("Trigger ProcessStatePublisher by event of publish")
			if err := p.doPublish(pub); err != nil {
				log.Errorf("ProcessStatePublisher failed, %v", err)
			}
		}
	}
}

func (p *ProcessStatePublisher) doPublish(pub []string) error {
	for _, procID := range pub {
		process := p.agent.ProcMgr.FindByProcID(procID)
		if process == nil {
			log.Warnf("Local process not found while publishing the state to etcd, procID: %s", procID)
			continue
		}
		if err := p.reg.UpdateProcessState(procID, p.agent.Mach.ID(), process.GetSvcName(), process.State(),
			process.IsActive(), p.ttl); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProcessStatePublisher) doPublishAll() error {
	for procID, process := range p.agent.ProcMgr.AllProcess() {
		if err := p.reg.UpdateProcessState(procID, p.agent.Mach.ID(), process.GetSvcName(), process.State(),
			process.IsActive(), p.ttl); err != nil {
			return err
		}
	}
	return nil
}
