package api

import "github.com/pingcap/tiadmin/schema"

// Version API
type VersionController struct {
	baseController
}

func (c *VersionController) VersionInfo() {
	c.Data["json"] = schema.Version{
		Version:      "1.0.0",
		BuildUTCTime: "2016-01-19 08:12:47",
	}
	c.ServeJSON()
}
