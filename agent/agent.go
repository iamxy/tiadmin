package agent

import (
	"github.com/pingcap/tiadmin/machine"
	proc "github.com/pingcap/tiadmin/process"
	"github.com/pingcap/tiadmin/registry"
	"time"
)

type Agent struct {
	Reg     registry.Registry
	ProcMgr proc.ProcMgr
	Mach    machine.Mach
	TTL     time.Duration
}

func New(reg registry.Registry, pm proc.ProcMgr, m machine.Mach, ttl time.Duration) *Agent {
	return &Agent{
		Reg:     reg,
		ProcMgr: pm,
		Mach:    m,
		TTL:     ttl,
	}
}
