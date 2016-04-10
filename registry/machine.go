package registry

import (
	"errors"
	"fmt"
	etcd "github.com/coreos/etcd/client"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/machine"
	"path"
	"time"
)

const machinePrefix = "machine"

func (r *EtcdRegistry) Machine(machID string) (*machine.MachineStatus, error) {
	ctx, cancel := r.ctx()
	defer cancel()
	resp, err := r.kAPI.Get(ctx, r.prefixed(machinePrefix, machID), &etcd.GetOptions{
		Recursive: true,
		Quorum:    true,
	})
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			// not found
			e := fmt.Sprintf("Machine not found in etcd, machID: %s, %v", machID, err)
			log.Error(e)
			return nil, errors.New(e)
		}
		return nil, err
	}
	status, err := machineStatusFromEtcdNode(machID, resp.Node)
	if err != nil || status == nil {
		e := errors.New(fmt.Sprintf("Invalid machine node, machID[%s], error[%v]", machID, err))
		return nil, e
	}
	return status, nil
}

func (r *EtcdRegistry) Machines() (map[string]*machine.MachineStatus, error) {
	key := r.prefixed(machinePrefix)
	ctx, cancel := r.ctx()
	defer cancel()
	resp, err := r.kAPI.Get(ctx, r.prefixed(machinePrefix), &etcd.GetOptions{
		Recursive: true,
		Quorum:    true,
	})
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			e := errors.New(fmt.Sprintf("%s not found in etcd, cluster may not be properly bootstrapped", key))
			return nil, e
		}
		return nil, err
	}
	IDToMachine := make(map[string]*machine.MachineStatus)
	for _, node := range resp.Node.Nodes {
		machID := path.Base(node.Key)
		status, err := machineStatusFromEtcdNode(machID, node)
		if err != nil || status == nil {
			e := errors.New(fmt.Sprintf("Invalid machine node, machID[%s], error[%v]", node.Key, err))
			return nil, e
		}
		IDToMachine[machID] = status
	}
	return IDToMachine, nil
}

// The structure of node representing machine in etcd:
//   /root/machine/{machID}
//                  /object
//                  /alive
//                  /statistic
func machineStatusFromEtcdNode(machID string, node *etcd.Node) (*machine.MachineStatus, error) {
	status := &machine.MachineStatus{
		MachID: machID,
	}
	for _, n := range node.Nodes {
		key := path.Base(n.Key)
		switch key {
		case "object":
			if err := unmarshal(n.Value, &status.MachInfo); err != nil {
				log.Errorf("Error unmarshaling MachInfo, machID: %s, %v", machID, err)
				return nil, err
			}
		case "alive":
			status.IsAlive = true
		case "statistic":
			if err := unmarshal(n.Value, &status.MachStat); err != nil {
				log.Errorf("Error unmarshaling MachStat, machID: %s, %v", machID, err)
				return nil, err
			}
		}
	}
	return status, nil
}

func (r *EtcdRegistry) RegisterMachine(machID, hostName, hostRegion, hostIDC, publicIP string) error {
	if exists, err := r.checkMachineExists(machID); err != nil {
		return err
	} else if !exists {
		// not found then create a new machine node
		return r.createMachine(machID, hostName, hostRegion, hostIDC, publicIP)
	}

	// found it, update host infomation of the machine
	machInfo := &machine.MachineInfo{
		HostName:   hostName,
		HostRegion: hostRegion,
		HostIDC:    hostIDC,
		PublicIP:   publicIP,
	}
	return r.updateMeachineInfo(machID, machInfo)
}

func (r *EtcdRegistry) checkMachineExists(machID string) (bool, error) {
	ctx, cancel := r.ctx()
	defer cancel()
	_, err := r.kAPI.Get(ctx, r.prefixed(machinePrefix, machID), &etcd.GetOptions{
		Quorum: true,
	})
	if err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *EtcdRegistry) updateMeachineInfo(machID string, machInfo *machine.MachineInfo) error {
	object, err := marshal(machInfo)
	if err != nil {
		e := fmt.Sprintf("Error marshaling MachineInfo, %v, %v", object, err)
		log.Errorf(e)
		return errors.New(e)
	}
	ctx, cancel := r.ctx()
	defer cancel()
	key := r.prefixed(machinePrefix, machID, "object")
	if _, err := r.kAPI.Set(ctx, key, object, &etcd.SetOptions{}); err != nil {
		e := fmt.Sprintf("Failed to update MachInfo in etcd, %s, %v, %v", machID, object, err)
		log.Error(e)
		return errors.New(e)
	}
	return nil
}

func (r *EtcdRegistry) createMachine(machID, hostName, hostRegion, hostIDC, publicIP string) error {
	object := &machine.MachineInfo{
		HostName:   hostName,
		HostRegion: hostRegion,
		HostIDC:    hostIDC,
		PublicIP:   publicIP,
	}
	statobj := &machine.MachineStat{
		UsageOfCPU:  0,
		TotalMem:    0,
		UsedMem:     0,
		TotalSwp:    0,
		UsedSwp:     0,
		LoadAvg:     make([]float32, 0),
		UsageOfDisk: make([]machine.DiskUsage, 0),
	}
	if err := r.mustCreateNode(r.prefixed(machinePrefix, machID), "", true); err != nil {
		e := fmt.Sprintf("Failed to create node of machine, %s, %v", machID, err)
		log.Error(e)
		return errors.New(e)
	}
	if objstr, err := marshal(object); err == nil {
		if err := r.createNode(r.prefixed(machinePrefix, machID, "object"), objstr, false); err != nil {
			e := fmt.Sprintf("Failed to create MachInfo of machine node, %s, %v, %v", machID, object, err)
			log.Error(e)
			return errors.New(e)
		}
	} else {
		e := fmt.Sprintf("Error marshaling MachineInfo, %v, %v", object, err)
		log.Errorf(e)
		return errors.New(e)
	}
	if statstr, err := marshal(statobj); err == nil {
		if err := r.createNode(r.prefixed(machinePrefix, machID, "statistic"), statstr, false); err != nil {
			e := fmt.Sprintf("Failed to create MachStat of machine node, %s, %v, %v", machID, statobj, err)
			log.Error(e)
			return errors.New(e)
		}
	} else {
		e := fmt.Sprintf("Error marshaling MachineStat, %v, %v", statobj, err)
		log.Errorf(e)
		return errors.New(e)
	}
	return nil
}

func (r *EtcdRegistry) RefreshMachine(machID string, machStat machine.MachineStat, ttl time.Duration) error {
	if err := r.refreshMachineStatistic(machID, &machStat); err != nil {
		return nil
	}
	if err := r.refreshMachineAlive(machID, ttl); err != nil {
		return nil
	}
	return nil
}

func (r *EtcdRegistry) refreshMachineStatistic(machID string, machStat *machine.MachineStat) error {
	object, err := marshal(machStat)
	if err != nil {
		e := fmt.Sprintf("Error marshaling MachineStat, %v, %v", machStat, err)
		log.Errorf(e)
		return errors.New(e)
	}
	key := r.prefixed(machinePrefix, machID, "statistic")
	ctx, cancel := r.ctx()
	defer cancel()
	if _, err := r.kAPI.Set(ctx, key, object, &etcd.SetOptions{}); err != nil {
		e := fmt.Sprintf("Failed to update machine statistic node of machine in etcd, %s, %v", machID, err)
		log.Error(e)
		return errors.New(e)
	}
	return nil
}

func (r *EtcdRegistry) refreshMachineAlive(machID string, ttl time.Duration) error {
	aliveKey := r.prefixed(machinePrefix, machID, "alive")
	// try to touch alive state of machine, update ttl
	ctx, cancel := r.ctx()
	defer cancel()
	if _, err := r.kAPI.Set(ctx, aliveKey, "", &etcd.SetOptions{
		PrevExist: etcd.PrevExist,
		TTL:       ttl,
		Refresh:   true,
	}); err != nil {
		if isEtcdError(err, etcd.ErrorCodeKeyNotFound) {
			return r.createMachineAlive(machID, ttl)
		}
		return err
	}
	return nil
}

func (r *EtcdRegistry) createMachineAlive(machID string, ttl time.Duration) error {
	aliveKey := r.prefixed(machinePrefix, machID, "alive")
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
