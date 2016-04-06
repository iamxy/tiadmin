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
	// Occurs when any Process's target state is touched
	ProcessTargetStateChangeEvent = pkg.Event("ProcessTargetStateChangeEvent")
	// Occurs when any Machine's target state is touched
	MachineTargetStateChangeEvent = pkg.Event("MachineTargetStateChangeEvent")

	jobPrefix = "job"
)

type EtcdEventStream struct {
	kAPI       etcd.KeysAPI
	rootPrefix string
}

func NewEtcdEventStream(kapi etcd.KeysAPI, keyPrefix string) pkg.EventStream {
	return &EtcdEventStream{kapi, keyPrefix}
}

// Next returns a channel which will emit an Event as soon as one of interest occurs
func (es *EtcdEventStream) Next(stop <-chan struct{}) chan pkg.Event {
	evchan := make(chan pkg.Event)
	key := path.Join(es.rootPrefix, jobPrefix)
	go func() {
		for {
			select {
			case <-stop:
				log.Debugf("Gracefully closing etcd watch loop: key[%s]", key)
				return
			default:
				opts := &etcd.WatcherOptions{
					AfterIndex: 0,
					Recursive:  true,
				}
				watcher := es.kAPI.Watcher(key, opts)
				log.Debugf("Creating etcd watcher: %s", key)
				res, err := watcher.Next(context.Background())
				if err != nil {
					log.Debugf("etcd watcher %v returned error: %v", key, err)
				} else {
					if ev, ok := parse(res, es.rootPrefix); ok {
						evchan <- ev
						return
					}
				}
			}
			// Let's not slam the etcd server in the event that we know
			// an unexpected error occurred.
			time.Sleep(time.Second)
		}
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
		ev = MachineTargetStateChangeEvent
		ok = true
	}
	return
}
