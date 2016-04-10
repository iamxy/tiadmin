package registry

import (
	etcd "github.com/coreos/etcd/client"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/pkg"
	"golang.org/x/net/context"
	"path"
	"strings"
	"time"
)

const (
	jobPrefix = "job"

	// Occurs when any Process's target state is touched
	ProcessTargetStateChangeEvent = pkg.Event("ProcessTargetStateChangeEvent")
	// Occurs when any Machine's state is touched
	MachineStateChangeEvent = pkg.Event("MachineStateChangeEvent")
)

type EtcdEventStream struct {
	watcher    etcd.Watcher
	rootPrefix string
}

func NewEtcdEventStream(kapi etcd.KeysAPI, keyPrefix string) pkg.EventStream {
	key := path.Join(keyPrefix, jobPrefix)
	opts := &etcd.WatcherOptions{
		AfterIndex: 0,
		Recursive:  true,
	}
	return &EtcdEventStream{
		watcher:    kapi.Watcher(key, opts),
		rootPrefix: keyPrefix,
	}
}

// Next returns a channel which will emit an Event as soon as one of interest occurs
func (es *EtcdEventStream) Next(timeout time.Duration) chan pkg.Event {
	evchan := make(chan pkg.Event)
	go func() {
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		res, err := es.watcher.Next(ctx)
		if err != nil {
			if err == context.DeadlineExceeded {
				close(evchan)
			} else {
				log.Errorf("Some failure encountered while waiting for next etcd event, %v", err)
			}
		} else {
			if ev, ok := parse(res, es.rootPrefix); ok {
				evchan <- ev
			}
		}
		return
	}()
	return evchan
}

func parse(res *etcd.Response, prefix string) (ev pkg.Event, ok bool) {
	if res == nil || res.Node == nil {
		return
	}
	if !strings.HasPrefix(res.Node.Key, path.Join(prefix, jobPrefix)) {
		return
	}
	switch path.Base(res.Node.Key) {
	case "process-state":
		ev = ProcessTargetStateChangeEvent
		ok = true
	case "machine-state":
		ev = MachineStateChangeEvent
		ok = true
	}
	return
}
