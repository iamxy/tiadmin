package schema

import ()

type Machine struct {
	MachID      string      `json:"machID,omitempty"`
	UsageOfCPU  int32       `json:"usageOfCPU,omitempty"`
	TotalMem    int32       `json:"totalMem,omitempty"`
	UsedMem     int32       `json:"usedMem,omitempty"`
	TotalSwp    int32       `json:"totalSwp,omitempty"`
	UsedSwp     int32       `json:"usedSwp,omitempty"`
	LoadAvg     []float32   `json:"loadAvg,omitempty"`
	UsageOfDisk []DiskUsage `json:"usageOfDisk,omitempty"`
}
