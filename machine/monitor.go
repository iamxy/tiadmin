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
		select {
		case <-stopc:
			log.Debug("Machine monitor is exiting due to stop signal")
			return
		case <-clock.After(monitorInterval):
			// log.Debug("Trigger monitor routine after tick")
			m.collect()
			//if err := m.collect(); err != nil {
			//	log.Errorf("Collect statistics of this machine failed, %v", err)
			//}
		}
	}
}

func (m *machine) collect() {
	updateCpuStat()
	memInfo := memInfo()
	load := loadAvg()
	stat := &MachineStat{
		UsageOfCPU:  100.0 - cpuIdle(),
		TotalMem:    (memInfo.memFree + memInfo.memUsed) / 1024 / 1024,
		UsedMem:     memInfo.memUsed / 1024 / 1024,
		TotalSwp:    (memInfo.swapFree + memInfo.swapUsed) / 1024 / 1024,
		UsedSwp:     memInfo.swapUsed / 1024 / 1024,
		LoadAvg:     []float64{load.Avg1min, load.Avg5min, load.Avg15min},
		UsageOfDisk: diskInfo(),
	}
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	m.stat = stat
}
