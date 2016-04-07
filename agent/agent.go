package agent

import (
	"github.com/pingcap/tiadmin/machine"
	proc "github.com/pingcap/tiadmin/process"
	"github.com/pingcap/tiadmin/registry"
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
