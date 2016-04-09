package api

import (
	"github.com/pingcap/tiadmin/pkg"
	"github.com/pingcap/tiadmin/schema"
	"github.com/pingcap/tiadmin/service"
)

type ServiceController struct {
	baseController
}

func (c *ServiceController) AllServices() {
	res := []*schema.Service{}
	for _, svc := range service.Registered {
		status := svc.Status()
		s := &schema.Service{
			SvcName:      status.SvcName,
			Version:      status.Version,
			Executor:     status.Executor,
			Command:      status.Command,
			Args:         status.Args,
			Environments: transformMapToEnvironments(status.Environments),
			Endpoints:    pkg.EndpointsToStrings(status.Endpoints),
		}
		res = append(res, s)
	}
	c.Data["json"] = res
	c.ServeJSON()
}

func (c *ServiceController) Service() {
	svcName := c.Ctx.Input.Param(":svcName")
	if len(svcName) == 0 {
		c.Abort("400")
	}
	if svc, ok := service.Registered[svcName]; ok {
		status := svc.Status()
		c.Data["json"] = &schema.Service{
			SvcName:      status.SvcName,
			Version:      status.Version,
			Executor:     status.Executor,
			Command:      status.Command,
			Args:         status.Args,
			Environments: transformMapToEnvironments(status.Environments),
			Endpoints:    pkg.EndpointsToStrings(status.Endpoints),
		}
		c.ServeJSON()
	} else {
		c.Abort("404")
	}
}
