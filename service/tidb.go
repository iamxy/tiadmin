package service

import (
	"encoding/json"
	"flag"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/pkg"
	proc "github.com/pingcap/tiadmin/process"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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
			command:      "bin/tidb-server",
			args:         []string{"-L", "info", "--store", "tikv", "--path", "$ETCD_ADDR/pd?cluster=1", "-P", "4000", "--lease", "1"},
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
	var res = make(map[string]pkg.Endpoint)
	argset := flag.NewFlagSet(TiDB_SERVICE, flag.ContinueOnError)
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
				if p, err := strconv.Atoi(flag.Value.String()); err == nil {
					v.Port = pkg.Port(p)
				}
			}
		case "TIDB_STATUS_ADDR":
			if flag := argset.Lookup("status"); flag != nil {
				if p, err := strconv.Atoi(flag.Value.String()); err == nil {
					v.Port = pkg.Port(p)
				}
			}
		}
		res[k] = v
	}
	return res
}

func (s *TiDBService) RetrieveRealPerformance(allProcesses map[string]*proc.ProcessStatus) *TiDBPerfMetrics {
	var res = &TiDBPerfMetrics{}
	for _, proc := range allProcesses {
		if proc.SvcName != TiDB_SERVICE || !proc.IsAlive {
			continue
		}
		if endpoint, ok := proc.RunInfo.Endpoints["TIDB_STATUS_ADDR"]; ok {
			if status, err := s.fetchTiDBStatusFromHttp(endpoint); err == nil {
				// accumulating
				res.Connections += status.Connections
				res.TPS += status.TPS
			}
		}
	}
	return res
}

func (s *TiDBService) fetchTiDBStatusFromHttp(addr pkg.Endpoint) (*TiDBPerfMetrics, error) {
	url := addr.String() + "/status"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warnf("Fetch TiDB status error while building http request, %v", err)
		return nil, err
	}
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("postman-token", "e5d36c00-595e-d099-7966-97dd984afbd7")
	client := &http.Client{
		Timeout: time.Duration(1) * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Warnf("Fetch TiDB status error while doing http request, %v", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Warnf("Fetch TiDB status error while reading response body, %v", err)
		return nil, err
	}
	var status = &TiDBPerfMetrics{}
	if err := json.Unmarshal(body, status); err != nil {
		log.Warnf("Fetch TiDB status error while unmarshal response body, %v", err)
		return nil, err
	}
	return status, nil
}
