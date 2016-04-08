package schema

import ()

type PerfMetrics struct {
	Tps   int32 `json:"tps"`
	Qps   int32 `json:"qps"`
	Iops  int32 `json:"iops"`
	Conns int32 `json:"conns"`
}
