package service

import "github.com/pingcap/tiadmin/registry"

type TiDBService struct {
}

func NewTiDBService(reg registry.Registry) Service {
	return &TiDBService{}
}
