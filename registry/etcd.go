package registry

import (
	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"path"
	"time"
)

const DefaultKeyPrefix = "/_pingcap.com/tiadmin"

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

func (r *EtcdRegistry) createNode(key, val string, isDir bool) (err error) {
	opts := &etcd.SetOptions{
		PrevExist: etcd.PrevNoExist,
		Dir: isDir,
	}
	_, err = r.kAPI.Set(r.ctx(), key, val, opts)
	return
}

func (r *EtcdRegistry) deleteNode(key string, isDir bool) (err error) {
	opts := &etcd.DeleteOptions{
		Recursive: isDir, // weird ?
		Dir: isDir,
	}
	_, err = r.kAPI.Delete(r.ctx(), key, opts)
	return
}

func (r *EtcdRegistry) mustCreateNode(key, val string, isDir bool) (err error) {
	if err = r.createNode(key, val, isDir); err != nil {
		if isEtcdError(err, etcd.ErrorCodeNodeExist) {
			if err = r.deleteNode(key, isDir); err == nil {
				err = r.createNode(key, val, isDir)
			}
		}
	}
	return
}

