package schema

import ()

type HostMeta struct {
	Region     string `json:"region,omitempty"`
	Datacenter string `json:"datacenter,omitempty"`
}
