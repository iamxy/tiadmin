package service

import (
	"flag"
	"github.com/pingcap/tiadmin/pkg"
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
			command:      "target/release/tikv-server",
			args:         []string{"-S", "raftkv", "--addr", "$HOST_IP:5551", "--pd", "$PD_ADDR", "-s", "data", "--cluster-id", "1", "-L", "debug"},
			environments: map[string]string{},
			endpoints: map[string]pkg.Endpoint{
				"TIKV_ADDR": pkg.Endpoint{
					Port: pkg.Port(5551),
				},
			},
		},
	}
}

func (s *TiKVService) ParseEndpointFromArgs(args []string) map[string]pkg.Endpoint {
	res := make(map[string]pkg.Endpoint)
	argset := flag.NewFlagSet(TiKV_SERVICE, flag.ExitOnError)
	argset.String("S", "raftkv", "")
	argset.String("addr", "127.0.0.1:5551", "")
	argset.String("pd", "127.0.0.1:1234", "")
	argset.String("s", "data", "")
	argset.String("cluster-id", "TiCluster", "")
	argset.String("L", "debug", "")
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
					if p, err := strconv.Atoi(addrParts[1]); err != nil {
						v.Port = pkg.Port(p)
					}
				}
			}
		}
		res[k] = v
	}
	return res
}
