package machine

type MachineStatus struct {
	MachID         string
	HostName       string
	HostRegion     string
	HostDatacenter string
	PublicIP       string
	IsAlive        bool
	metrics        MachineMetrics
}

type MachineMetrics struct {
	UsageOfCPU  int32
	TotalMem    int32
	UsedMem     int32
	TotalSwp    int32
	UsedSwp     int32
	LoadAvg     []float32
	UsageOfDisk []DiskUsage
}

type DiskUsage struct {
	Mount     string
	TotalSize int32
	UsedSize  int32
}
