package api

import (
	"encoding/json"
	"github.com/pingcap/tiadmin/schema"
)

type HostController struct {
	baseController
}

func (c *HostController) FindAllHosts() {
	hosts := []schema.Host{}
	for _, val := range mockHosts {
		hosts = append(hosts, val)
	}
	c.Data["json"] = hosts
	c.ServeJSON()
}

func (c *HostController) FindHost() {
	machID := c.Ctx.Input.Param(":machID")
	if len(machID) == 0 {
		c.Abort("400")
	}
	host, ok := mockHosts[machID]
	if !ok {
		c.Abort("404")
	}
	c.Data["json"] = host
	c.ServeJSON()
}

func (c *HostController) SetHostMetaInfo() {
	machID := c.Ctx.Input.Param(":machID")
	if len(machID) == 0 {
		c.Abort("400")
	}
	host, ok := mockHosts[machID]
	if !ok {
		c.Abort("404")
	}
	var meta schema.HostMeta
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &meta)
	if err != nil || len(meta.Region) == 0 || len(meta.Datacenter) == 0 {
		c.Abort("400")
	}
	host.HostMeta = meta
	mockHosts[machID] = host
	c.Data["json"] = host
	c.ServeJSON()
}
