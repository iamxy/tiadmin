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
	UsageOfCPU  float64
	TotalMem    uint64
	UsedMem     uint64
	TotalSwp    uint64
	UsedSwp     uint64
	LoadAvg     []float64
	UsageOfDisk []DiskUsage
	ClockOffset float64
}

type DiskUsage struct {
	Mount     string
	TotalSize uint64
	UsedSize  uint64
}
