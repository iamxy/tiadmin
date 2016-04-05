package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/pkg"
	"github.com/pingcap/tiadmin/registry"
	"github.com/rakyll/globalconf"
	"os"
)

const (
	// TTL to use with all state pushed to Registry
	DefaultTTL = "30s"
	// If an environment variable with the EnvPrefix is set, it will take precedence over values
	// in the configuration file. Command line flags will override the environment variables.
	EnvConfigPrefix = "TIADM_"
	// First try to load configuration file in $(PWD), if not exist then check /etc/tiadm/tiadm.conf
	DefaultConfigFile = "/tiadm.conf"
	DefaultConfigDir  = "/etc/tiadm"
)

type Config struct {
	EtcdServers         []string
	EtcdKeyPrefix       string
	EtcdRequestTimeout  int
	MonitorLoopInterval int
	HostIP              string
	HostRegion          string
	HostIDC             string
	HostStateTTL        string
	TokenLimit          int
	IsMock              bool
}

func ParseFlag() (*Config, error) {
	etcdServers := flag.String("etcd_servers", "http://127.0.0.1:2379,http://127.0.0.1:4001", "List of etcd endpoints")
	etcdKeyPrefix := flag.String("etcd_key_prefix", registry.DefaultKeyPrefix, "Namespace for tiadmin registry in etcd")
	etcdRequestTimeout := flag.Int("etcd_request_timeout", 1000, "Amount of time in milliseconds to allow a single etcd request before considering it failed.")
	monitorLoopInterval := flag.Int("monitor_loop_interval", 2000, "Interval at which the monitor should check and report the cluster status in etcd periodically.")
	hostIP := flag.String("host_ip", "", "IP address that this host should publish")
	hostRegion := flag.String("host_region", "", "Geographical region where this node was set")
	hostIDC := flag.String("host_idc", "", "IDC identifier this node belongs to")
	hostStateTTL := flag.String("host_ttl", DefaultTTL, "TTL in seconds of host state in etcd")
	tokenLimit := flag.Int("token_limit", 100, "Maximum number of entries per page returned from API requests")
	isMock := flag.Bool("mock", false, "Whether to privide mock APIs for test")

	opts := globalconf.Options{EnvPrefix: EnvConfigPrefix}
	if file, err := pathToConfigFile(); err == nil {
		opts.Filename = file
	}
	if gconf, err := globalconf.NewWithOptions(&opts); err == nil {
		gconf.ParseAll()
	} else {
		return nil, err
	}

	var etcdEndpoints pkg.StringSlice
	etcdEndpoints.Set(*etcdServers)

	cfg := &Config{
		EtcdServers:         etcdEndpoints,
		EtcdKeyPrefix:       *etcdKeyPrefix,
		EtcdRequestTimeout:  *etcdRequestTimeout,
		MonitorLoopInterval: *monitorLoopInterval,
		HostIP:              *hostIP,
		HostRegion:          *hostRegion,
		HostIDC:             *hostIDC,
		HostStateTTL:        *hostStateTTL,
		TokenLimit:          *tokenLimit,
		IsMock:              *isMock,
	}
	return cfg, nil
}

func pathToConfigFile() (string, error) {
	cd, err := pkg.GetCmdDir()
	if err != nil {
		log.Errorf("get command directory error, %v", err)
		return "", err
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Errorf("get work directory error, %v", err)
		return "", err
	}
	if path, err := checkFileExist(cd + DefaultConfigFile); err == nil {
		return path, nil
	}
	if path, err := checkFileExist(cd + "/conf" + DefaultConfigFile); err == nil {
		return path, nil
	}
	if path, err := checkFileExist(cd + "../conf" + DefaultConfigFile); err == nil {
		return path, nil
	}
	if path, err := checkFileExist(wd + DefaultConfigFile); err == nil {
		return path, nil
	}
	if path, err := checkFileExist(wd + "/conf" + DefaultConfigFile); err == nil {
		return path, nil
	}
	if path, err := checkFileExist(wd + "../conf" + DefaultConfigFile); err == nil {
		return path, nil
	}
	if path, err := checkFileExist(DefaultConfigDir + DefaultConfigFile); err == nil {
		return path, nil
	}
	return "", errors.New("Not configuration file found")
}

func checkFileExist(filepath string) (string, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return "", err
	}
	if fi.IsDir() {
		return "", errors.New(fmt.Sprintf("filepath: %s, is a directory, not a file", filepath))
	}
	return filepath, nil
}
