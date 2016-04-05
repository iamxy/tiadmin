package registry

import (
	etcd "github.com/coreos/etcd/client"
	"github.com/pingcap/tiadmin/pkg"
)

type EtcdEventStream struct {
}

func NewEtcdEventStream(kapi etcd.KeysAPI, keyPrefix string) pkg.EventStream {
	return &EtcdEventStream{}
}

func (e *EtcdEventStream) Next(stopc <-chan struct{}) chan pkg.Event {
	return nil
}
