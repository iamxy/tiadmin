package schema

import ()

type Host struct {
	MachID   string   `json:"machID,omitempty"`
	HostName string   `json:"hostName,omitempty"`
	HostMeta HostMeta `json:"hostMeta,omitempty"`
	PublicIP string   `json:"publicIP,omitempty"`
	IsAlive  bool     `json:"isAlive,omitempty"`
	Machine  Machine  `json:"machine,omitempty"`
}
