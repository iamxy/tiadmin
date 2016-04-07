package api

import "github.com/pingcap/tiadmin/schema"

type ServiceController struct {
	baseController
}

func (c *ServiceController) AllServices() {
	services := []schema.Service{}
	for _, val := range mockServices {
		services = append(services, val)
	}
	c.Data["json"] = services
	c.ServeJSON()
}

func (c *ServiceController) Service() {
	svcName := c.Ctx.Input.Param(":svcName")
	if len(svcName) == 0 {
		c.Abort("400")
	}
	svc, ok := mockServices[svcName]
	if !ok {
		c.Abort("404")
	}
	c.Data["json"] = svc
	c.ServeJSON()
}
