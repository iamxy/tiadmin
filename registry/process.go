package registry

import (
	"errors"
	"fmt"
	etcd "github.com/coreos/etcd/client"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/pkg"
	proc "github.com/pingcap/tiadmin/process"
	"path"
	"strings"
	"time"
)

const processPrefix = "process"

func (r *EtcdRegistry) Processes() (map[string]*proc.ProcessStatus, error) {
	key := r.prefixed(processPrefix)
	opts := &etcd.GetOptions{
		Recursive: true,
		Quorum:    true,
	}
	ctx, cancel := r.ctx()
	defer cancel()
	resp, err := r.kAPI.Get(ctx, key, opts)
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			e := errors.New(fmt.Sprintf("%s not found in etcd, cluster may not be properly bootstrapped", key))
			return nil, e
		}
		return nil, err
	}

	IDToProcess := make(map[string]*proc.ProcessStatus)
	for _, node := range resp.Node.Nodes {
		key := path.Base(node.Key)
		parts := strings.Split(key, "-")
		if len(parts) < 3 {
			e := errors.New(fmt.Sprintf("Node key[%s] is illegal, invalid key foramt of process", node.Key))
			return nil, e
		}
		procID := parts[0]
		machID := parts[1]
		svcName := parts[2]
		status, err := processStatusFromEtcdNode(procID, machID, svcName, node)
		if err != nil || status == nil {
			e := errors.New(fmt.Sprintf("Invalid process node, key[%s], error[%v]", node.Key, err))
			return nil, e
		}
		IDToProcess[procID] = status
	}
	return IDToProcess, nil
}

func (r *EtcdRegistry) Process(procID string) (*proc.ProcessStatus, error) {
	key := r.prefixed(processPrefix)
	opts := &etcd.GetOptions{
		Recursive: true,
		Quorum:    true,
	}
	ctx, cancel := r.ctx()
	defer cancel()
	resp, err := r.kAPI.Get(ctx, key, opts)
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			e := errors.New(fmt.Sprintf("Node[%s] not found in etcd, Ti-Cluster may not be properly bootstrapped", key))
			return nil, e
		}
		return nil, err
	}

	for _, node := range resp.Node.Nodes {
		key := path.Base(node.Key)
		parts := strings.Split(key, "-")
		if len(parts) < 3 {
			e := errors.New(fmt.Sprintf("Node key[%s] is illegal, invalid key foramt of process", node.Key))
			return nil, e
		}
		if procID != parts[0] {
			continue
		}
		machID := parts[1]
		svcName := parts[2]
		status, err := processStatusFromEtcdNode(procID, machID, svcName, node)
		if err != nil || status == nil {
			e := errors.New(fmt.Sprintf("Invalid process node, key[%s], error[%v]", node.Key, err))
			return nil, e
		}
		return status, nil
	}
	e := errors.New(fmt.Sprintf("No process found by procID[%s]", procID))
	return nil, e
}

func (r *EtcdRegistry) ProcessesOnMachine(machID string) (map[string]*proc.ProcessStatus, error) {
	key := r.prefixed(processPrefix)
	opts := &etcd.GetOptions{
		Recursive: true,
		Quorum:    true,
	}
	ctx, cancel := r.ctx()
	defer cancel()
	resp, err := r.kAPI.Get(ctx, key, opts)
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			e := errors.New(fmt.Sprintf("Node[%s] not found in etcd, Ti-Cluster may not be properly bootstrapped", key))
			return nil, e
		}
		return nil, err
	}

	procIDToProcess := make(map[string]*proc.ProcessStatus)
	for _, node := range resp.Node.Nodes {
		key := path.Base(node.Key)
		parts := strings.Split(key, "-")
		if len(parts) < 3 {
			e := errors.New(fmt.Sprintf("Node key[%s] is illegal, invalid key foramt of process", node.Key))
			return nil, e
		}
		procID := parts[0]
		if machID != parts[1] {
			continue
		}
		svcName := parts[2]
		status, err := processStatusFromEtcdNode(procID, machID, svcName, node)
		if err != nil || status == nil {
			e := errors.New(fmt.Sprintf("Invalid process node, key[%s], error[%v]", node.Key, err))
			return nil, e
		}
		procIDToProcess[procID] = status
	}
	return procIDToProcess, nil
}

func (r *EtcdRegistry) ProcessesOfService(svcName string) (map[string]*proc.ProcessStatus, error) {
	key := r.prefixed(processPrefix)
	opts := &etcd.GetOptions{
		Recursive: true,
		Quorum:    true,
	}
	ctx, cancel := r.ctx()
	defer cancel()
	resp, err := r.kAPI.Get(ctx, key, opts)
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			e := errors.New(fmt.Sprintf("Node[%s] not found in etcd, Ti-Cluster may not be properly bootstrapped", key))
			return nil, e
		}
		return nil, err
	}

	procIDToProcess := make(map[string]*proc.ProcessStatus)
	for _, node := range resp.Node.Nodes {
		key := path.Base(node.Key)
		parts := strings.Split(key, "-")
		if len(parts) < 3 {
			e := errors.New(fmt.Sprintf("Node key[%s] is illegal, invalid key foramt of process", node.Key))
			return nil, e
		}
		procID := parts[0]
		machID := parts[1]
		if svcName != parts[2] {
			continue
		}
		status, err := processStatusFromEtcdNode(procID, machID, svcName, node)
		if err != nil || status == nil {
			e := errors.New(fmt.Sprintf("Invalid process node, key[%s], error[%v]", node.Key, err))
			return nil, e
		}
		procIDToProcess[procID] = status
	}
	return procIDToProcess, nil
}

// The structure of node representing a process in etcd:
//   /root/process/{procID}-{machID}-{svcName}
//                  /desired-state
//                  /current-state
//                  /alive
//                  /object
//                  /endpoints/{endpoint}
func processStatusFromEtcdNode(procID, machID, svcName string, node *etcd.Node) (*proc.ProcessStatus, error) {
	if !node.Dir {
		return nil, errors.New(fmt.Sprintf("Invalid process node, not a etcd directory, key[%v]", node.Key))
	}
	status := &proc.ProcessStatus{
		ProcID:  procID,
		MachID:  machID,
		SvcName: svcName,
	}
	for _, n := range node.Nodes {
		key := path.Base(n.Key)
		switch key {
		case "desired-state":
			if state, err := parseProcessState(n.Value); err != nil {
				log.Errorf("Error parsing process state, procID: %s, %v", procID, err)
				return nil, err
			} else {
				status.DesiredState = state
			}
		case "current-state":
			if state, err := parseProcessState(n.Value); err != nil {
				log.Errorf("Error parsing process state, procID: %s, %v", procID, err)
				return nil, err
			} else {
				status.CurrentState = state
			}
		case "alive":
			status.IsAlive = true
		case "object":
			if err := unmarshal(n.Value, &status.RunInfo); err != nil {
				log.Errorf("Error unmarshaling RunInfo, procID: %s, %v", procID, err)
				return nil, err
			}
		case "endpoints":
			for _, epNode := range n.Nodes {
				if epNode.Value == "ok" {
					str := path.Base(epNode.Key)
					if endpoint, err := pkg.ParseEndpoint(str); err == nil {
						status.Endpoints = append(status.Endpoints, endpoint)
					} else {
						log.Errorf("Error parsing endpoint, procID: %s, %v", procID, err)
						return nil, err
					}
				}
			}
		}
	}
	return status, nil
}

func parseProcessState(state string) (proc.ProcessState, error) {
	switch state {
	case proc.StateStarted.String():
		return proc.StateStarted, nil
	case proc.StateStopped.String():
		return proc.StateStopped, nil
	default:
		return proc.StateStopped, errors.New(fmt.Sprintf("Illegal process state: %s", state))
	}
}

func (r *EtcdRegistry) UpdateProcessState(procID, machID, svcName string, state proc.ProcessState, isAlive bool, ttl time.Duration) error {
	// update the current-state of process
	if err := r.updateProcessCurrentState(procID, machID, svcName, state); err != nil {
		return err
	}
	// update the real alive state of process in etcd
	return r.refreshProcessAlive(procID, machID, svcName, isAlive, ttl)
}

func (r *EtcdRegistry) updateProcessCurrentState(procID, machID, svcName string, state proc.ProcessState) error {
	procKey := strings.Join([]string{procID, machID, svcName}, "-")
	currentStateKey := r.prefixed(processPrefix, procKey, "current-state")
	ctx, cancel := r.ctx()
	defer cancel()
	_, err := r.kAPI.Set(ctx, currentStateKey, state.String(), &etcd.SetOptions{
		PrevValue: state.Opposite().String(),
		PrevExist: etcd.PrevExist,
	})
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			// maybe process has been destroyed
			log.Warnf("Error updating process state of procID: %s, process node is gone, error: %v", procID, err)
			return nil
		} else if isEtcdError(err, etcd.ErrorCodeTestFailed) {
			log.Debugf("Process's current-state not changed in etcd, procID: %s, state: %s", procID, state.String())
			return nil
		} else {
			// other errors
			return err
		}
	}
	return nil
}

func (r *EtcdRegistry) refreshProcessAlive(procID, machID, svcName string, isAlive bool, ttl time.Duration) error {
	procKey := strings.Join([]string{procID, machID, svcName}, "-")
	aliveKey := r.prefixed(processPrefix, procKey, "alive")
	if isAlive {
		// try to touch alive node of process, if node not exists than create it
		if ok, err := r.touchProcessAlive(aliveKey, ttl); err != nil {
			return err
		} else if !ok {
			// create new node
			return r.createProcessAlive(aliveKey, ttl)
		}
	} else {
		// delete the alive state of process immediately
		ctx, cancel := r.ctx()
		defer cancel()
		_, err := r.kAPI.Delete(ctx, aliveKey, &etcd.DeleteOptions{})
		if err != nil && !isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			return err
		}
	}
	return nil
}

func (r *EtcdRegistry) touchProcessAlive(aliveKey string, ttl time.Duration) (bool, error) {
	ctx, cancel := r.ctx()
	defer cancel()
	_, err := r.kAPI.Set(ctx, aliveKey, "", &etcd.SetOptions{
		PrevExist: etcd.PrevExist,
		TTL:       ttl,
		Refresh:   true,
	})
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *EtcdRegistry) createProcessAlive(aliveKey string, ttl time.Duration) error {
	ctx, cancel := r.ctx()
	defer cancel()
	_, err := r.kAPI.Set(ctx, aliveKey, "", &etcd.SetOptions{
		TTL: ttl,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *EtcdRegistry) NewProcess(machID, svcName string, hostIP, hostName, hostRegion, hostIDC string,
	executor []string, command string, args []string, env map[string]string, port pkg.Port, protocol pkg.Protocol) error {
	// generate new process ID
	procID, err := r.GenerateProcID()
	if err != nil {
		e := fmt.Sprintf("Failed to generate new process ID, %v", err)
		log.Error(e)
		return errors.New(e)
	}
	procKey := strings.Join([]string{procID, machID, svcName}, "-")
	desiredState := proc.StateStarted
	currentState := proc.StateStopped
	object := &proc.ProcessRunInfo{
		HostIP:      hostIP,
		HostName:    hostName,
		HostRegion:  hostRegion,
		HostIDC:     hostIDC,
		Executor:    executor,
		Command:     command,
		Args:        args,
		Environment: env,
		Port:        port,
		Protocol:    protocol,
	}
	if err := r.mustCreateNode(r.prefixed(processPrefix, procKey), "", true); err != nil {
		e := fmt.Sprintf("Failed to create node of process, %s, %v", procKey, err)
		log.Error(e)
		return errors.New(e)
	}
	if err := r.createNode(r.prefixed(processPrefix, procKey, "desired-state"), desiredState.String(), false); err != nil {
		e := fmt.Sprintf("Failed to create desired-state of process node, %s, %v", procKey, err)
		log.Error(e)
		return errors.New(e)
	}
	if err := r.createNode(r.prefixed(processPrefix, procKey, "current-state"), currentState.String(), false); err != nil {
		e := fmt.Sprintf("Failed to create current-state of process node, %s, %v", procKey, err)
		log.Error(e)
		return errors.New(e)
	}
	if objstr, err := marshal(object); err == nil {
		if err := r.createNode(r.prefixed(processPrefix, procKey, "object"), objstr, false); err != nil {
			e := fmt.Sprintf("Failed to create RunInfo of process node, %s, %v, %v", procKey, object, err)
			log.Error(e)
			return errors.New(e)
		}
	} else {
		e := fmt.Sprintf("Error marshaling RunInfo, %v, %v", object, err)
		log.Errorf(e)
		return errors.New(e)
	}
	if err := r.createNode(r.prefixed(processPrefix, procKey, "endpoints"), "", true); err != nil {
		e := fmt.Sprintf("Failed to create endpoints of process node, %s, %v", procKey, err)
		log.Error(e)
		return errors.New(e)
	}
	return nil
}

func (r *EtcdRegistry) DeleteProcess(procID string) (*proc.ProcessStatus, error) {
	status, err := r.Process(procID)
	if err != nil {
		return nil, err
	}
	procKey := strings.Join([]string{status.ProcID, status.MachID, status.SvcName}, "-")
	if err := r.deleteNode(r.prefixed(processPrefix, procKey), true); err != nil {
		return nil, err
	}
	return status, nil
}

func (r *EtcdRegistry) UpdateProcessDesiredState(procID string, state proc.ProcessState) error {
	status, err := r.Process(procID)
	if err != nil {
		return err
	}
	procKey := strings.Join([]string{status.ProcID, status.MachID, status.SvcName}, "-")
	ctx, cancel := r.ctx()
	defer cancel()
	if _, err := r.kAPI.Set(ctx, r.prefixed(processPrefix, procKey, "desired-state"), state.String(), &etcd.SetOptions{
		PrevExist: etcd.PrevExist,
	}); err != nil {
		return err
	}
	return nil
}
