package registry

import (
	"errors"
	"fmt"
	etcd "github.com/coreos/etcd/client"
	"github.com/pingcap/tiadmin/pkg"
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
			e := errors.New(fmt.Sprintf("Node[%s] not found in etcd, Ti-Cluster may not be properly bootstrapped", key))
			return nil, e
		}
		return nil, err
	}

	procIDToProcess := make(map[string]*proc.ProcessStatus)
	for _, node := range resp.Node.Nodes {
		_, key := path.Split(node.Key)
		parts := strings.Split(key, "-")
		if len(parts) < 3 {
			e := errors.New(fmt.Sprintf("Node key[%s] is illegal, invalid key foramt of process", node.Key))
			return nil, e
		}
		procID := parts[0]
		status, err := r.processStatusFromEtcdNode(procID, node)
		if err != nil || status == nil {
			e := errors.New(fmt.Sprintf("Invalid process node, key[%s], error[%v]", node.Key, err))
			return nil, e
		}
		procIDToProcess[procID] = status
	}
	return procIDToProcess, nil
}

func (r *EtcdRegistry) Process(procID string) (*proc.ProcessStatus, error) {
	key := r.prefixed(ProcessPrefix)
	opts := &etcd.GetOptions{
		Recursive: true,
		Quorum:    true,
	}
	resp, err := r.kAPI.Get(r.ctx(), key, opts)
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			e := errors.New(fmt.Sprintf("Node[%s] not found in etcd, Ti-Cluster may not be properly bootstrapped", key))
			return nil, e
		}
		return nil, err
	}

	for _, node := range resp.Node.Nodes {
		_, key := path.Split(node.Key)
		parts := strings.Split(key, "-")
		if len(parts) < 3{
			e := errors.New(fmt.Sprintf("Node key[%s] is illegal, invalid key foramt of process", node.Key))
			return nil, e
		}
		if procID != parts[0] {
			continue
		}
		status, err := r.processStatusFromEtcdNode(procID, node)
		if err != nil || status == nil {
			e := errors.New(fmt.Sprintf("Invalid process node, key[%s], error[%v]", node.Key, err))
			return nil, e
		}
		return status, nil
	}
	e := errors.New(fmt.Sprintf("No process found by procID[%s]", procID))
	return nil, e
}

func (r *EtcdRegistry) ProcessesOnHost(machID string) (map[string]*proc.ProcessStatus, error) {
	key := r.prefixed(ProcessPrefix)
	opts := &etcd.GetOptions{
		Recursive: true,
		Quorum:    true,
	}
	resp, err := r.kAPI.Get(r.ctx(), key, opts)
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			e := errors.New(fmt.Sprintf("Node[%s] not found in etcd, Ti-Cluster may not be properly bootstrapped", key))
			return nil, e
		}
		return nil, err
	}

	procIDToProcess := make(map[string]*proc.ProcessStatus)
	for _, node := range resp.Node.Nodes {
		_, key := path.Split(node.Key)
		parts := strings.Split(key, "-")
		if len(parts) < 3{
			e := errors.New(fmt.Sprintf("Node key[%s] is illegal, invalid key foramt of process", node.Key))
			return nil, e
		}
		procID := parts[0]
		if machID != parts[1] {
			continue
		}
		status, err := r.processStatusFromEtcdNode(procID, node)
		if err != nil || status == nil {
			e := errors.New(fmt.Sprintf("Invalid process node, key[%s], error[%v]", node.Key, err))
			return nil, e
		}
		procIDToProcess[procID] = status
	}
	return procIDToProcess, nil
}

func (r *EtcdRegistry) ProcessesOfService(svcName string) (map[string]*proc.ProcessStatus, error) {
	key := r.prefixed(ProcessPrefix)
	opts := &etcd.GetOptions{
		Recursive: true,
		Quorum:    true,
	}
	resp, err := r.kAPI.Get(r.ctx(), key, opts)
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			e := errors.New(fmt.Sprintf("Node[%s] not found in etcd, Ti-Cluster may not be properly bootstrapped", key))
			return nil, e
		}
		return nil, err
	}

	procIDToProcess := make(map[string]*proc.ProcessStatus)
	for _, node := range resp.Node.Nodes {
		_, key := path.Split(node.Key)
		parts := strings.Split(key, "-")
		if len(parts) < 3{
			e := errors.New(fmt.Sprintf("Node key[%s] is illegal, invalid key foramt of process", node.Key))
			return nil, e
		}
		procID := parts[0]
		if svcName != parts[2] {
			continue
		}
		status, err := r.processStatusFromEtcdNode(procID, node)
		if err != nil || status == nil {
			e := errors.New(fmt.Sprintf("Invalid process node, key[%s], error[%v]", node.Key, err))
			return nil, e
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
			if state, err := r.parseProcessState(n.Value); err != nil {
				log.Errorf("error parsing process state, procID: %s, %v", procID, err)
				return nil, err
			} else {
				status.DesiredState = state
			}
		case "current_state":
			if state, err := r.parseProcessState(n.Value); err != nil {
				log.Errorf("error parsing process state, procID: %s, %v", procID, err)
				return nil, err
			} else {
				status.CurrentState = state
			}
		case "alive":
			status.IsAlive = true
		case "object":
			var ri proc.ProcessRunInfo
			if err := unmarshal(n.Value, &ri); err != nil {
				log.Errorf("error unmarshaling ProcessRunInfo(procID: %s): %v", procID, err)
				return nil, err
			}
		case "endpoints":
			for _, epNode := range n.Nodes {
				if epNode.Value == "ok" {
					_, str := path.Split(epNode.Key)
					if endpoint, err := pkg.ParseEndpoint(str); err != nil {
						log.Errorf("error parsing endpoint, procID: %s, %v", procID, err)
						return nil, err
					} else {
						status.Endpoints = append(status.Endpoints, endpoint)
					}
				}
			}
		}
	}
	return nil, nil
}

func (r *EtcdRegistry) parseProcessState(state string) (proc.ProcessState, error) {
	switch state {
	case "started":
		return proc.StateStarted, nil
	case "Started":
		return proc.StateStarted, nil
	case "stopped":
		return proc.StateStopped, nil
	case "Stopped":
		return proc.StateStopped, nil
	default:
		return proc.StateStopped, errors.New(fmt.Sprintf("Illegal process state: %s", state))
	}
}
