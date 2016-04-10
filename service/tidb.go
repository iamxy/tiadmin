package service

import (
	"github.com/pingcap/tiadmin/pkg"
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
