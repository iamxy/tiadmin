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
	res.memFree = mem.MemFree + mem.Buffers + mem.Cached
	res.memUsed = mem.MemTotal - res.memFree
	res.swapUsed = mem.SwapUsed
	res.swapFree = mem.SwapFree
	return res
}
