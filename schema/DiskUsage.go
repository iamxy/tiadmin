package schema

import ()

type DiskUsage struct {
	Mount     string `json:"mount"`
	TotalSize int32  `json:"totalSize"`
	UsedSize  int32  `json:"usedSize"`
}
