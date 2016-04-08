package agent

import (
	"github.com/jonboulle/clockwork"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/registry"
	"time"
)

func NewAgentHeartbeat(reg registry.Registry, ag *Agent, ttl time.Duration) *AgentHeartbeat {
	return &AgentHeartbeat{
		reg:   reg,
		agent: ag,
		clock: clockwork.NewRealClock(),
		ttl:   ttl,
	}
}

type AgentHeartbeat struct {
	reg   registry.Registry
	agent *Agent
	clock clockwork.Clock
	ttl   time.Duration
}

func (h *AgentHeartbeat) Run(stopc <-chan struct{}) {
	for {
		select {
		case <-stopc:
			log.Debug("AgentHeartbeat is exiting due to stop signal")
			return
		case <-h.clock.After(h.ttl / 2):
			log.Debug("Trigger heartbeat after tick")
			if err := h.heartBeat(); err != nil {
				log.Errorf("Heartbeat refresh state to etcd failed, %v", err)
			}
		}
	}
}

func (h *AgentHeartbeat) heartBeat() error {
	if err := h.reg.RefreshMachine(h.agent.Mach.ID(), h.agent.Mach.Status().MachStat, h.ttl); err != nil {
		log.Errorf("Refresh machine state in ectd meet error, %v", err)
		return err
	}
	return nil
}
