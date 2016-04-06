package schema

import ()

type Machine struct {
	MachID      string      `json:"machID"`
	UsageOfCPU  int32       `json:"usageOfCPU"`
	TotalMem    int32       `json:"totalMem"`
	UsedMem     int32       `json:"usedMem"`
	TotalSwp    int32       `json:"totalSwp"`
	UsedSwp     int32       `json:"usedSwp"`
	LoadAvg     []float32   `json:"loadAvg"`
	UsageOfDisk []DiskUsage `json:"usageOfDisk"`
}
