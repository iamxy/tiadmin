package registry

import (
	"errors"
	"fmt"
	etcd "github.com/coreos/etcd/client"
	proc "github.com/pingcap/tiadmin/process"
	"github.com/prometheus/common/log"
	"path"
	"strings"
)

const ProcessPrefix = "process"

func (r *EtcdRegistry) Processes() (map[string]*proc.ProcessStatus, error) {
	key := r.prefixed(ProcessPrefix)
	opts := &etcd.GetOptions{
		Recursive: true,
		Quorum:    true,
	}
	resp, err := r.kAPI.Get(r.ctx(), key, opts)
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			log.Fatal("Node key[%s] not found in etcd, Ti-Cluster may not be properly bootstrapped")
		}
		return nil, err
	}

	procIDToProcess := map[string]*proc.ProcessStatus{}
	for _, node := range resp.Node.Nodes {
		_, key := path.Split(node.Key)
		parts := strings.Split(key, "-")
		if len(parts) < 3 {
			log.Errorf("Node key[%v] is illegal, invalid process key parts", node.Key)
			continue
		}
		procID := parts[0]
		status, err := r.processStatusFromEtcdNode(procID, node)
		if err != nil || status == nil {
			log.Errorf("Invalid process node, key[%v], error[%v]", node.Key, err)
			continue
		}
		procIDToProcess[procID] = status
	}

	return procIDToProcess, nil
}

// The structure of node representing a process in etcd directory:
//   /root/process/{procID}-{machID}-{svcName}
//                  /desired_state
//                  /current_state
//                  /alive
//                  /object
//                  /endpoints/{endpoint}
func (r *EtcdRegistry) processStatusFromEtcdNode(procID string, node *etcd.Node) (*proc.ProcessStatus, error) {
	if !node.Dir {
		return nil, errors.New(fmt.Sprintf("Invalid process node, not a etcd directory, key[%v]", node.Key))
	}
	status := &proc.ProcessStatus{}
	for _, n := range node.Nodes {
		_, key := path.Split(n.Key)
		switch key {
		case "desired_state":
			status.DesiredState = n.Value
		case "current_state":
			status.CurrentState = n.Value
		case "alive":
			status.IsAlive = true
		case "object":
			var ri proc.ProcessRunInfo
			if err := unmarshal(n.Value, &ri); err != nil {
				log.Errorf("error unmarshaling ProcessRunInfo(procID: %s): %v", procID, err)
				return nil
			}
		case "endpoints":
			for _, ep := range n.Nodes {
				if ep.Value == "ok" {
					_, endpoint := path.Split(ep.Key)
					status.Endpoints = append(status.Endpoints, endpoint)
				}
			}
		}
	}
	return nil, nil
}