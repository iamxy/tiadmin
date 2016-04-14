package machine

import (
	"github.com/toolkits/nux"
)

func diskInfo() []DiskUsage {
	var res = []DiskUsage{}
	mountPoints, err := nux.ListMountPoint()
	if err != nil {
		return res
	}
	for idx := range mountPoints {
		var du *nux.DeviceUsage
		du, err = nux.BuildDeviceUsage(mountPoints[idx][0], mountPoints[idx][1], mountPoints[idx][2])
		if err == nil {
			res = append(res, DiskUsage{
				Mount:     du.FsFile,
				TotalSize: du.BlocksAll / 1024 / 1024,
				UsedSize:  du.BlocksUsed / 1024 / 1024,
			})
		}
	}
	return res
}
