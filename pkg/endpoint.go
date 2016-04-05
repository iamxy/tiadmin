package pkg

import "fmt"

type Protocol string
type Port int

const (
	ProtocolHttp = Protocol("http")
	ProtocolUnix = Protocol("unix")
)

type Endpoint struct {
	Protocol Protocol
	// TODO: This should allow hostname or IP
	IPAddress string
	Port      Port
}

func (e Endpoint) String() string {
	return e.Protocol + e.IPAddress + ":" + e.Port
	return fmt.Sprintf("%s://%s:%d", e.Protocol, e.IPAddress, e.Port)
}
