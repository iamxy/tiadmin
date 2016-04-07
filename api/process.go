package api

import (
	"encoding/json"
	"github.com/pingcap/tiadmin/schema"
	"strconv"
)

type ProcessController struct {
	baseController
}

func (c *ProcessController) FindAllProcesses() {
	procs := []schema.Process{}
	for _, val := range mockProcs {
		procs = append(procs, val)
	}
	c.Data["json"] = procs
	c.ServeJSON()
}

func (c *ProcessController) StartNewProcess() {
	var proc schema.Process
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &proc)
	if err != nil || len(proc.SvcName) == 0 || len(proc.MachID) == 0 {
		c.Abort("400")
	}
	if svc, ok := mockServices[proc.SvcName]; ok {
		if len(proc.Executor) == 0 {
			proc.Executor = svc.Executor
		}
		if len(proc.Command) > 0 {
			proc.Command = svc.Command
		}
		if len(proc.Args) == 0 {
			proc.Args = svc.Args
		}
		if len(proc.Environments) == 0 {
			proc.Environments = svc.Environments
		}
		if proc.Port == 0 {
			proc.Port = svc.Port
		}
		if len(proc.Protocol) == 0 {
			proc.Protocol = svc.Protocol
		}
	} else {
		c.Abort("400")
	}
	if mach, ok := mockHosts[proc.MachID]; ok {
		proc.PublicIP = mach.PublicIP
		proc.HostName = mach.HostName
		proc.HostMeta = mach.HostMeta
	} else {
		c.Abort("400")
	}
	proc.ProcID = strconv.Itoa(mockProcID)
	mockProcID++
	if len(proc.DesiredState) == 0 {
		proc.DesiredState = "started"
	}
	proc.CurrentState = "stopped"
	proc.IsAlive = false
	mockProcs[proc.ProcID] = proc
	c.Data["json"] = proc
	c.ServeJSON()
}

func (c *ProcessController) FindByHost() {
	machID := c.GetString("machID")
	if len(machID) == 0 {
		c.Abort("400")
	}
	if _, ok := mockHosts[machID]; ok {
		res := []schema.Process{}
		for _, proc := range mockProcs {
			if proc.MachID == machID {
				res = append(res, proc)
			}
		}
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		c.Abort("404")
	}
}

func (c *ProcessController) FindByService() {
	svcName := c.GetString("svcName")
	if len(svcName) == 0 {
		c.Abort("400")
	}
	if _, ok := mockServices[svcName]; ok {
		res := []schema.Process{}
		for _, proc := range mockProcs {
			if proc.SvcName == svcName {
				res = append(res, proc)
			}
		}
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		c.Abort("404")
	}
}

func (c *ProcessController) FindProcess() {
	procID := c.Ctx.Input.Param(":procID")
	if len(procID) == 0 {
		c.Abort("400")
	}
	if proc, ok := mockProcs[procID]; ok {
		c.Data["json"] = proc
		c.ServeJSON()
	} else {
		c.Abort("404")
	}
}

func (c *ProcessController) DestroyProcess() {
	procID := c.Ctx.Input.Param(":procID")
	if len(procID) == 0 {
		c.Abort("400")
	}
	if proc, ok := mockProcs[procID]; ok {
		delete(mockProcs, procID)
		c.Data["json"] = proc
		c.ServeJSON()
	} else {
		c.Abort("404")
	}
}

func (c *ProcessController) StartProcess() {
	procID := c.Ctx.Input.Param(":procID")
	if len(procID) == 0 {
		c.Abort("400")
	}
	if proc, ok := mockProcs[procID]; ok {
		proc.DesiredState = "started"
		mockProcs[procID] = proc
		c.Data["json"] = proc
		c.ServeJSON()
	} else {
		c.Abort("404")
	}
}

func (c *ProcessController) StopProcess() {
	procID := c.Ctx.Input.Param(":procID")
	if len(procID) == 0 {
		c.Abort("400")
	}
	if proc, ok := mockProcs[procID]; ok {
		proc.DesiredState = "stopped"
		mockProcs[procID] = proc
		c.Data["json"] = proc
		c.ServeJSON()
	} else {
		c.Abort("404")
	}
}
