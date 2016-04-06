package schema

import ()

type Host struct {
	MachID   string   `json:"machID"`
	HostName string   `json:"hostName"`
	HostMeta HostMeta `json:"hostMeta"`
	PublicIP string   `json:"publicIP"`
	IsAlive  bool     `json:"isAlive"`
	Machine  Machine  `json:"machine"`
}
