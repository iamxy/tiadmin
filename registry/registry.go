package registry

import (
	"encoding/json"
	"fmt"
	"github.com/pingcap/tiadmin/config"
	"github.com/pingcap/tiadmin/machine"
	"github.com/pingcap/tiadmin/pkg"
	proc "github.com/pingcap/tiadmin/process"
	"time"
)

// Registry interface defined a set of operations to access a distributed key value store,
// which always organizes data as directory structure, simlilar to a file system.
// now we implemented a registry driving ETCD as backend
type Registry interface {
	// Check whether tiadmin registry is bootstrapped normally
	IsBootstrapped(*config.Config) bool
	// Initialize the basic directory structure of tiadmin registry
	Bootstrap() error
	// Get Infomation of machine in cluster by the given machID
	Machine(machID string) (*machine.MachineStatus, error)
	// Retrieve all machines in Ti-Cluster,
	// return a map of machID to machineStatus
	Machines() (map[string]*machine.MachineStatus, error)
	// Create new machine node in etcd
	NewMachine(machID, hostName, hostRegion, hostIDC, publicIP string) error
	// Update statistic info of machine and refresh the TTL of alive state in etcd
	MachineHeartbeat(machID string, machStat *machine.MachineStat, ttl time.Duration) error
	// Return the status of process with specified procID
	Process(procID string) (*proc.ProcessStatus, error)
	// Retrieve all processes in Ti-Cluster,
	// with either running or stopped state
	// return a map of procID to processStatus
	Processes() (map[string]*proc.ProcessStatus, error)
	// Retrieve processes which scheduled at the specified host by given machID
	// return a map of procID to status infomation of process
	ProcessesOnMachine(machID string) (map[string]*proc.ProcessStatus, error)
	// Retrieve all processes instantiated from the specified service
	// return a map of procID to status infomation of process
	ProcessesOfService(svcName string) (map[string]*proc.ProcessStatus, error)
	// Create new process node of specified service in etcd
	NewProcess(machID, svcName string, hostIP, hostName, hostRegion, hostIDC string,
		executor []string, command string, args []string, env map[string]string, port pkg.Port, protocol pkg.Protocol) error
	// Destroy the process, normally the process should be in stopped state
	DeleteProcess(procID string) (*proc.ProcessStatus, error)
	// Update process desirede state in etcd
	UpdateProcessDesiredState(procID string, state proc.ProcessState) error
	// Update process current state in etcd, notice that isAlive is real run state of the local process
	UpdateProcessState(procID, machID, svcName string, state proc.ProcessState, isAlive bool, ttl time.Duration) error
}

func marshal(obj interface{}) (string, error) {
	encoded, err := json.Marshal(obj)
	if err == nil {
		return string(encoded), nil
	}
	return "", fmt.Errorf("unable to JSON-serialize object: %s", err)
}

func unmarshal(val string, obj interface{}) error {
	err := json.Unmarshal([]byte(val), &obj)
	if err == nil {
		return nil
	}
	return fmt.Errorf("unable to JSON-deserialize object: %s", err)
}
