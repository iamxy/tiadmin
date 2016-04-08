package machine

import (
	"github.com/jonboulle/clockwork"
	"github.com/ngaut/log"
	"math/rand"
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
			log.Debug("Trigger monitor routine after tick")
			if err := m.collect(); err != nil {
				log.Errorf("Collect statistics of this machine failed, %v", err)
			}
		}
	}
}

func randInt(min int, max int) int32 {
	return int32(min + rand.Intn(max-min))
}

func (m *machine) collect() error {
	stat := &MachineStat{
		UsageOfCPU:  randInt(0, 100),
		TotalMem:    1024 * 8,
		UsedMem:     randInt(1140, 4600),
		TotalSwp:    0,
		UsedSwp:     0,
		LoadAvg:     []float32{1.04, 1.46, 1.31},
		UsageOfDisk: []DiskUsage{},
	}
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	m.stat = stat
	return nil
}
