package registry

import (
	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"path"
	"time"
)

const DefaultKeyPrefix = "/_pingcap.com/tidb-admin"

// EtcdRegistry implement the Registry interface and uses etcd as backend
type EtcdRegistry struct {
	kAPI       etcd.KeysAPI
	keyPrefix  string
	reqTimeout time.Duration
}

func NewEtcdRegistry(kapi etcd.KeysAPI, keyPrefix string, reqTimeout time.Duration) Registry {
	return &EtcdRegistry{
		kAPI:       kapi,
		keyPrefix:  keyPrefix,
		reqTimeout: reqTimeout,
	}
}

func (r *EtcdRegistry) ctx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), r.reqTimeout)
	return ctx
}

func (r *EtcdRegistry) prefixed(p ...string) string {
	return path.Join(r.keyPrefix, path.Join(p...))
}

func isEtcdError(err error, code int) bool {
	eerr, ok := err.(etcd.Error)
	return ok && eerr.Code == code
}
