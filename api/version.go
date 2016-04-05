package api

// Version API
type VersionController struct {
	baseController
}

func (c *VersionController) VersionInfo() {
	c.Ctx.WriteString("tidb-admin version 1.0 !")
}
