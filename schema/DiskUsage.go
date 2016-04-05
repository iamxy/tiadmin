package schema

import ()

type DiskUsage struct {
	Mount     string `json:"mount,omitempty"`
	TotalSize int32  `json:"totalSize,omitempty"`
	UsedSize  int32  `json:"usedSize,omitempty"`
}
