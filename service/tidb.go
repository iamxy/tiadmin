package service

import ()

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
			command:      "bin/tidb_server",
			args:         []string{"-L", "info", "-path", "/tmp/tidb", "-P", "4000"},
			environments: map[string]string{},
		},
	}
}
