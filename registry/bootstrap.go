package registry

import (
	etcd "github.com/coreos/etcd/client"
	"github.com/ngaut/log"
)

const bootstrapPrefix = "bootstrapped"

func (r *EtcdRegistry) IsBootstrapped() (bool, error) {
	key := r.prefixed(bootstrapPrefix)
	opts := &etcd.GetOptions{
		Quorum: true,
	}
	resp, err := r.kAPI.Get(r.ctx(), key, opts)
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			// not bootstrapped
			log.Infof("The etcd registry of tiadmin not bootstrapped yet")
			return false, nil
		}
		return false, err
	} else {
		if resp.Node.Dir {
			log.Fatalf("Node[%s] is a directory in etcd, which's unexpected", key)
		}
		// already bootstrapped
		return true, nil
	}
}

func (r *EtcdRegistry) Bootstrap() (err error) {
	if err = r.mustCreateNode(r.prefixed(processPrefix), "", true); err != nil {
		return
	}
	if err = r.mustCreateNode(r.prefixed(machinePrefix), "", true); err != nil {
		return
	}
	if err = r.mustCreateNode(r.prefixed(jobPrefix), "", true); err != nil {
		return
	}
	if err = r.mustCreateNode(r.prefixed(bootstrapPrefix), "bootstrapped", false); err != nil {
		return
	}
	return
}
