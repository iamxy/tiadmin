package service

import (
	"flag"
	"github.com/pingcap/tiadmin/pkg"
	"regexp"
	"strconv"
	"strings"
)

const TiKV_SERVICE = "TiKV"

type TiKVService struct {
	service
}

func NewTiKVService() Service {
	return &TiKVService{
		service{
			svcName:      TiKV_SERVICE,
			version:      "1.0.0",
			executor:     []string{},
			command:      "bin/tikv-server",
			args:         []string{"-S", "raftkv", "--addr", "0.0.0.0:5551", "--advertise-addr", "$HOST_IP:5551", "--pd", "$PD_ADVERTISE_ADDR", "--store", "data", "--cluster-id", "1"},
			environments: map[string]string{},
			endpoints: map[string]pkg.Endpoint{
				"TIKV_ADDR": pkg.Endpoint{
					Port: pkg.Port(5551),
				},
				"TIKV_ADVERTISE_ADDR": pkg.Endpoint{
					Port: pkg.Port(5551),
				},
			},
		},
	}
}

func (s *TiKVService) ParseEndpointFromArgs(args []string) map[string]pkg.Endpoint {
	var res = make(map[string]pkg.Endpoint)
	argset := flag.NewFlagSet(TiKV_SERVICE, flag.ContinueOnError)
	argset.String("addr", "127.0.0.1:5551", "")
	argset.String("advertise-addr", "127.0.0.1:5551", "")
	argset.String("L", "debug", "")
	argset.String("store", "data", "")
	argset.String("S", "raftkv", "")
	argset.String("cluster-id", "1", "")
	argset.String("pd", "127.0.0.1:1234", "")
	if err := argset.Parse(args); err != nil {
		// handle error
		return s.endpoints
	}

	for k, v := range s.endpoints {
		switch k {
		case "TIKV_ADDR":
			if flag := argset.Lookup("addr"); flag != nil {
				addrParts := strings.Split(flag.Value.String(), ":")
				if len(addrParts) > 1 {
					if p, err := strconv.Atoi(addrParts[1]); err == nil {
						v.Port = pkg.Port(p)
					}
				}
			}
		case "TIKV_ADVERTISE_ADDR":
			if flag := argset.Lookup("advertise-addr"); flag != nil {
				addrParts := strings.Split(flag.Value.String(), ":")
				if len(addrParts) > 1 {
					if ok, _ := regexp.MatchString(`^(\d+)\.(\d+)\.(\d+)\.(\d+)$`, addrParts[0]); ok {
						v.IPAddr = addrParts[0]
					}
					if p, err := strconv.Atoi(addrParts[1]); err == nil {
						v.Port = pkg.Port(p)
					}
				}
			}
		}
		res[k] = v
	}
	return res
}
