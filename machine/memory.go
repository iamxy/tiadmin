package machine

import "github.com/toolkits/nux"

type Mem struct {
	memFree  uint64
	memUsed  uint64
	swapUsed uint64
	swapFree uint64
}

func memInfo() *Mem {
	var res = &Mem{}
	mem, err := nux.MemInfo()
	if err != nil {
		return res
	}
	res.memFree = (mem.MemFree + mem.Buffers + mem.Cached) / 1024 / 1024
	res.memUsed = (mem.MemTotal - mem.MemFree - mem.Buffers - mem.Cached) / 1024 / 1024
	res.swapUsed = mem.SwapUsed / 1024 / 1024
	res.swapFree = mem.SwapFree / 1024 / 1024
	return res
}
