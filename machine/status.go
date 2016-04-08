package machine

type MachineStatus struct {
	MachID   string
	IsAlive  bool
	MachInfo MachineInfo
	MachStat MachineStat
}

type MachineInfo struct {
	HostName   string
	HostRegion string
	HostIDC    string
	PublicIP   string
}

type MachineStat struct {
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
