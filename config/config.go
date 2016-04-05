package config

import (
	"errors"
	"flag"
	"github.com/pingcap/tiadmin/pkg"
	"github.com/pingcap/tiadmin/registry"
	"github.com/rakyll/globalconf"
	"os"
	"strings"
)

const (
	// TTL to use with all state pushed to Registry
	DefaultTTL = "30s"
	// If an environment variable with the EnvPrefix is set, it will take precedence over values
	// in the configuration file. Command line flags will override the environment variables.
	EnvConfigPrefix = "TIDB_"
	// First try to load configuration file in $(PWD), if not exist then check /etc/tidbadm/tidbadm.conf
	DefaultConfigFile = "/tidbadm.conf"
	DefaultConfigDir  = "/etc/tidbadm"
)

type Config struct {
	EtcdServers         []string
	EtcdKeyPrefix       string
	EtcdRequestTimeout  int
	MonitorLoopInterval int
	HostIP              string
	HostMetadata        string
	HostStateTTL        string
	TokenLimit          int
	LogLevel            string
}

func ParseFlag() (*Config, error) {
	etcdServers := flag.String("etcd_servers", "http://127.0.0.1:2379,http://127.0.0.1:4001", "List of etcd endpoints")
	etcdKeyPrefix := flag.String("etcd_key_prefix", registry.DefaultKeyPrefix, "namespace for tidb-admin registry in etcd")
	etcdRequestTimeout := flag.Int("etcd_request_timeout", 1000, "Amount of time in milliseconds to allow a single etcd request before considering it failed.")
	monitorLoopInterval := flag.Int("monitor_loop_interval", 2000, "Interval at which the monitor should check and report the cluster status in etcd periodically.")
	hostIP := flag.String("host_ip", "", "IP address that this host should publish")
	hostMetaData := flag.String("host_metadata", "", "List of key-value metadata to assign to this host")
	hostStateTTL := flag.String("host_state_ttl", DefaultTTL, "TTL in seconds of host state in etcd")
	tokenLimit := flag.Int("token_limit", 100, "Maximum number of entries per page returned from API requests")
	logLevel := flag.String("log_level", "debug", "Assigned log level: info, debug, warn, error, fatal")

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
		HostMetadata:        *hostMetaData,
		HostStateTTL:        *hostStateTTL,
		TokenLimit:          *tokenLimit,
		LogLevel:            *logLevel,
	}
	return cfg, nil
}

func pathToConfigFile() (string, error) {
	wd, err := pkg.GetRootDir()
	if err == nil {
		if path, err := checkFileExist(wd + DefaultConfigFile); err == nil {
			return path, nil
		}
		if path, err := checkFileExist(wd + "/conf" + DefaultConfigFile); err == nil {
			return path, nil
		}
		if path, err := checkFileExist(wd + "../conf" + DefaultConfigFile); err == nil {
			return path, nil
		}
		if path, err := checkFileExist(os.Getwd() + DefaultConfigFile); err == nil {
			return path, nil
		}
		if path, err := checkFileExist(os.Getwd() + "/conf" + DefaultConfigFile); err == nil {
			return path, nil
		}
	}
	if path, err := checkFileExist(DefaultConfigDir + DefaultConfigFile); err == nil {
		return path, nil
	}
	return "", errors.New("Not configuration file found")
}

func checkFileExist(filepath string) (string, error) {
	fi, err := os.Stat(filepath)
	if err != nil || fi.IsDir() {
		return "", err
	}
	return filepath, nil
}

func (c *Config) Metadata() map[string]string {
	meta := make(map[string]string, 0)
	for _, pair := range strings.Split(c.HostMetadata, ",") {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		meta[key] = val
	}
	return meta
}
