package api

import "github.com/pingcap/tiadmin/schema"

type MonitorController struct {
	baseController
}

func (c *MonitorController) TiDBPerformanceMetrics() {
	// TODO: implement it
	c.Data["json"] = schema.PerfMetrics{
		Tps:   int32(randInt(20, 100)),
		Qps:   int32(randInt(65, 300)),
		Iops:  int32(randInt(80, 550)),
		Conns: int32(randInt(3, 10)),
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
