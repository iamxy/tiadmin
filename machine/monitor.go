package machine

import (
	"github.com/jonboulle/clockwork"
	"github.com/ngaut/log"
	"time"
)

const (
	// time between triggering monitor routine
	monitorInterval = 1 * time.Second
)

func (m *machine) Monitor(stopc <-chan struct{}) {
	var clock = clockwork.NewRealClock()
	for {
		begin := time.Now()
		select {
		case <-stopc:
			log.Debug("Machine monitor is exiting due to stop signal")
			return
		case <-clock.After(monitorInterval):
			// log.Debug("Trigger monitor routine after tick")
			m.collect(begin)
		}
	}
}

func (m *machine) collect(begin time.Time) {
	updateCpuStat()
	memInfo := memInfo()
	load := loadAvg()
	offset := time.Now().Sub(begin) - monitorInterval
	stat := &MachineStat{
		UsageOfCPU:  100.0 - cpuIdle(),
		TotalMem:    memInfo.memFree + memInfo.memUsed,
		UsedMem:     memInfo.memUsed,
		TotalSwp:    memInfo.swapFree + memInfo.swapUsed,
		UsedSwp:     memInfo.swapUsed,
		LoadAvg:     []float64{load.Avg1min, load.Avg5min, load.Avg15min},
		UsageOfDisk: diskInfo(),
		ClockOffset: offset.Seconds(),
	}
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	m.stat = stat
}
