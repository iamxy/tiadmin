package service

import (
	"flag"
	"github.com/pingcap/tiadmin/pkg"
	"strconv"
)

const TiDB_SERVICE = "TiDB"

type TiDBService struct {
	service
}

func NewTiDBService() Service {
	return &TiDBService{
		service{
			svcName:      TiDB_SERVICE,
			version:      "1.0.0",
			executor:     []string{},
			command:      "tidb-server",
			args:         []string{"-L", "info", "-path", "/tmp/tidb", "-P", "4000"},
			environments: map[string]string{},
			endpoints: map[string]pkg.Endpoint{
				"TIDB_ADDR": pkg.Endpoint{
					Protocol: pkg.Protocol("mysql"),
					Port:     pkg.Port(4000),
				},
				"TIDB_STATUS_ADDR": pkg.Endpoint{
					Protocol: pkg.Protocol("http"),
					Port:     pkg.Port(10080),
				},
			},
		},
	}
}

func (s *TiDBService) ParseEndpointFromArgs(args []string) map[string]pkg.Endpoint {
	res := make(map[string]pkg.Endpoint)
	argset := flag.NewFlagSet(TiDB_SERVICE, flag.ExitOnError)
	argset.String("store", "goleveldb", "registered store name, [memory, goleveldb, hbase, boltdb, tikv]")
	argset.String("path", "/tmp/tidb", "tidb storage path")
	argset.String("L", "debug", "log level: info, debug, warn, error, fatal")
	argset.String("P", "4000", "mp server port")
	argset.String("status", "10080", "tidb server status port")
	argset.Int("lease", 1, "schema lease seconds, very dangerous to change only if you know what you do")
	if err := argset.Parse(args); err != nil {
		// handle error
		return s.endpoints
	}

	for k, v := range s.endpoints {
		switch k {
		case "TIDB_ADDR":
			if flag := argset.Lookup("P"); flag != nil {
				if p, err := strconv.Atoi(flag.Value.String()); err != nil {
					v.Port = pkg.Port(p)
				}
			}
		case "TIDB_STATUS_ADDR":
			if flag := argset.Lookup("status"); flag != nil {
				if p, err := strconv.Atoi(flag.Value.String()); err != nil {
					v.Port = pkg.Port(p)
				}
			}
		}
		res[k] = v
	}
	return res
}
