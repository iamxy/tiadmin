package machine

import "github.com/toolkits/nux"

type Loadavg struct {
	Avg1min  float64
	Avg5min  float64
	Avg15min float64
}

func loadAvg() *Loadavg {
	var res = &Loadavg{}
	load, err := nux.LoadAvg()
	if err != nil {
		return res
	}
	res.Avg1min = load.Avg1min
	res.Avg5min = load.Avg5min
	res.Avg15min = load.Avg15min
	return res
}
