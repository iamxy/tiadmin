package api

import (
	"github.com/pingcap/tiadmin/schema"
	"github.com/pingcap/tiadmin/server"
)

type MonitorController struct {
	baseController
}

func (c *MonitorController) TiDBPerformanceMetrics() {
	// TODO: implement it
	var status = server.Agent.ShowTiDBRealPerfermance()
	c.Data["json"] = schema.PerfMetrics{
		Tps:   int32(status.TPS),
		Qps:   int32(0),
		Iops:  int32(0),
		Conns: int32(status.Connections),
	}
	c.ServeJSON()
}

func (c *MonitorController) TiKVStorageMetrics() {
	// TODO: implement it
	c.Data["json"] = schema.StorageMetrics{
		Usage:    int64(randInt(535, 565)),
		Capacity: 8192,
	}
	c.ServeJSON()
}
