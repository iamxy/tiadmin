package registry

import (
	"fmt"
	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"path"
	"strconv"
	"time"
)

const maxProcessID = "max-process-id"

// EtcdRegistry implement the Registry interface and uses etcd as backend
type EtcdRegistry struct {
	kAPI       etcd.KeysAPI
	keyPrefix  string
	reqTimeout time.Duration
	etcdAddrs  string
}

func NewEtcdRegistry(kapi etcd.KeysAPI, keyPrefix string, reqTimeout time.Duration, etcdAddrs string) Registry {
	return &EtcdRegistry{
		kAPI:       kapi,
		keyPrefix:  keyPrefix,
		reqTimeout: reqTimeout,
		etcdAddrs:  etcdAddrs,
	}
}

func (r *EtcdRegistry) ctx() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), r.reqTimeout)
	return ctx, cancel
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
		Dir:       isDir,
	}
	ctx, cancel := r.ctx()
	defer cancel()
	_, err = r.kAPI.Set(ctx, key, val, opts)
	return
}

func (r *EtcdRegistry) deleteNode(key string, isDir bool) (err error) {
	opts := &etcd.DeleteOptions{
		Recursive: isDir, // weird ?
		Dir:       isDir,
	}
	ctx, cancel := r.ctx()
	defer cancel()
	_, err = r.kAPI.Delete(ctx, key, opts)
	return
}

func (r *EtcdRegistry) mustCreateNode(key, val string, isDir bool) (err error) {
	err = r.createNode(key, val, isDir)
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeNodeExist) {
			err = r.deleteNode(key, isDir)
			if err == nil {
				err = r.createNode(key, val, isDir)
			}
		}
	}
	return
}

func (r *EtcdRegistry) GenerateProcID() (string, error) {
	for {
		currProcID, err := r.getCurrentProcessID()
		if err != nil {
			return "", err
		}
		if suc, err := r.tryIncreaseMaxProcessID(currProcID); suc && err == nil {
			return currProcID, nil
		} else if err != nil {
			return "", err
		}
		// try failed, next loop
	}
}

func (r *EtcdRegistry) getCurrentProcessID() (string, error) {
	ctx, cancel := r.ctx()
	defer cancel()
	resp, err := r.kAPI.Get(ctx, r.prefixed(maxProcessID), &etcd.GetOptions{
		Quorum: true,
	})
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			panic(fmt.Sprintf("Max-Process-ID not exists in etcd, maybe not normally bootstrapped, %v", err))
		}
		return "", err
	}
	return resp.Node.Value, nil
}

func (r *EtcdRegistry) tryIncreaseMaxProcessID(currProcID string) (bool, error) {
	id, err := strconv.Atoi(currProcID)
	if err != nil {
		panic(fmt.Sprintf("Illegal value of Max-Process-ID stored in etcd, %v", err))
	}

	nextProcID := strconv.Itoa(id + 1)

	ctx, cancel := r.ctx()
	defer cancel()
	if _, err := r.kAPI.Set(ctx, r.prefixed(maxProcessID), nextProcID, &etcd.SetOptions{
		PrevExist: etcd.PrevExist,
		PrevValue: currProcID,
	}); err != nil {
		if isEtcdError(err, etcd.ErrorCodeTestFailed) {
			// try failed
			return false, nil
		}
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			panic(fmt.Sprintf("Max-Process-ID not exists in etcd, maybe not normally bootstrapped, %v", err))
		}
		return false, err
	} else {
		// success
		return true, nil
	}
}

func (r *EtcdRegistry) GetEtcdAddrs() string {
	return r.etcdAddrs
}
