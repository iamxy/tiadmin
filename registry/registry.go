package registry

import (
	"encoding/json"
	"fmt"
	proc "github.com/pingcap/tiadmin/process"
)

// Registry interface defined a set of operations to access a distributed key value store,
// which always organizes data as directory structure, simlilar to a file system.
// now we implemented a registry driving ETCD as backend
type Registry interface {
	// Return the status of process with specified procID
	Process(procID string) (*proc.ProcessStatus, error)
	// Retrieve all processes in Ti-Cluster,
	// with either running or stopped state
	// return a map of procID to status info of process
	Processes() (map[string]*proc.ProcessStatus, error)
	// Retrieve processes which scheduled at the specified host by given machID
	// return a map of procID to status infomation of process
	ProcessesOnHost(machID string) (map[string]*proc.ProcessStatus, error)
	// Retrieve all processes instantiated from the specified service
	// return a map of procID to status infomation of process
	ProcessesOfService(svcName string) (map[string]*proc.ProcessStatus, error)
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
