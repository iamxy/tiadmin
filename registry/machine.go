package registry

import "github.com/pingcap/tiadmin/machine"

const machinePrefix = "machine"

func (r *EtcdRegistry) Machine(machID string) (*machine.MachineStatus, error) {
	return nil, nil
}
