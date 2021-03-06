package pkg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Protocol string
type Port int32

const (
	ProtocolHttp = Protocol("http")
	ProtocolUnix = Protocol("unix")
)

func (p Port) Value() int32 {
	return int32(p)
}

func (p Protocol) String() string {
	return string(p)
}

type Endpoint struct {
	Protocol Protocol
	IPAddr   string
	Port     Port
}

func (e Endpoint) String() string {
	var ip = "0.0.0.0"
	if len(e.IPAddr) > 0 {
		ip = e.IPAddr
	}
	if len(e.Protocol) > 0 {
		return fmt.Sprintf("%s://%s:%d", e.Protocol.String(), ip, e.Port.Value())
	} else {
		return fmt.Sprintf("%s:%d", ip, e.Port.Value())
	}
}

func TrimAddrs(addrs []string) []string {
	res := []string{}
	for _, addr := range addrs {
		parts := strings.Split(addr, "://")
		if len(parts) == 2 {
			res = append(res, parts[1])
		} else {
			res = append(res, addr)
		}
	}
	return res
}

func ParseEndpoint(str string) (Endpoint, error) {
	var res Endpoint
	parts := strings.Split(str, "://")
	if len(parts) == 2 {
		sparts := strings.Split(parts[1], ":")
		if len(sparts) < 2 {
			return res, errors.New(fmt.Sprintf("Illegal endpoint string: %s", str))
		}
		res.Protocol = Protocol(parts[0])
		res.IPAddr = sparts[0]
		if port, err := strconv.Atoi(sparts[1]); err != nil {
			return res, errors.New(fmt.Sprintf("Illegal endpoint string: %s", str))
		} else {
			res.Port = Port(port)
		}
	} else {
		sparts := strings.Split(parts[0], ":")
		if len(sparts) < 2 {
			return res, errors.New(fmt.Sprintf("Illegal endpoint string: %s", str))
		}
		res.IPAddr = sparts[0]
		if port, err := strconv.Atoi(sparts[1]); err != nil {
			return res, errors.New(fmt.Sprintf("Illegal endpoint string: %s", str))
		} else {
			res.Port = Port(port)
		}
	}
	return res, nil
}

func ParseEndpoints(slice []string) ([]Endpoint, error) {
	res := []Endpoint{}
	for _, s := range slice {
		ep, err := ParseEndpoint(s)
		if err != nil {
			return nil, err
		}
		res = append(res, ep)
	}
	return res, nil
}

func EndpointsToStrings(endpoints map[string]Endpoint) []string {
	res := []string{}
	for _, ep := range endpoints {
		res = append(res, ep.String())
	}
	return res
}
