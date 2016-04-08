package machine

import (
	"crypto/sha1"
	"fmt"
	"github.com/docker/libcontainer/netlink"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/config"
	"github.com/pingcap/tiadmin/pkg"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"sync"
)

const (
	shortIDLen    = 8
	machineIDFile = ".machineID"
)

type Machine interface {
	ID() string
	ShortID() string
	MatchID(ID string) bool
	Status() *MachineStatus
	Monitor(<-chan struct{})
}

type machine struct {
	machID     string
	hostName   string
	hostRegion string
	hostIDC    string
	publicIP   string
	stat       *MachineStat
	rwMutex    sync.RWMutex
}

func NewMachineFromConfig(cfg *config.Config) (Machine, error) {
	machID, err := readLocalMachineID()
	if err != nil {
		log.Errorf("Read local machine ID error, %v", err)
		return nil, err
	}

	var publicIP string
	if len(cfg.HostIP) > 0 {
		publicIP = cfg.HostIP
	} else {
		publicIP = getLocalIP()
	}
	var hostName string
	if len(cfg.HostName) > 0 {
		hostName = cfg.HostName
	} else {
		hostName = publicIP
	}

	mach := &machine{
		machID:     machID,
		hostName:   hostName,
		hostRegion: cfg.HostRegion,
		hostIDC:    cfg.HostIDC,
		publicIP:   publicIP,
		stat: &MachineStat{
			LoadAvg:     []float32{},
			UsageOfDisk: []DiskUsage{},
		},
	}
	return mach, nil
}

// IsLocalMachineID returns whether the given machine ID is equal to that of the local machine
func IsLocalMachineID(mID string) bool {
	m, err := readLocalMachineID()
	return err == nil && m == mID
}

func readLocalMachineID() (string, error) {
	root, err := pkg.GetRootDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(root, machineIDFile)
	if _, err := pkg.CheckFileExist(fullPath); err != nil {
		return generateLocalMachineID(fullPath)
	} else {
		// read the machine ID from file
		hash, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return "", err
		}
		machID := fmt.Sprintf("%X", hash)
		if len(machID) == 0 {
			return generateLocalMachineID(fullPath)
		}
		return machID, nil
	}
}

// generate a new machine ID, and save it to file
func generateLocalMachineID(fullPath string) (string, error) {
	t := sha1.New()
	rand64 := string(pkg.KRand(64, pkg.KC_RAND_KIND_ALL))
	log.Debugf("Generated a string of 64 rand bytes, %s", rand64)
	io.WriteString(t, rand64)
	hash := t.Sum(nil)
	if err := ioutil.WriteFile(fullPath, hash, os.ModePerm); err != nil {
		return "", err
	}
	machID := fmt.Sprintf("%X", hash)
	return machID, nil
}

func getLocalIP() (got string) {
	iface := getDefaultGatewayIface()
	if iface == nil {
		return
	}
	addrs, err := iface.Addrs()
	if err != nil || len(addrs) == 0 {
		return
	}
	for _, addr := range addrs {
		// Attempt to parse the address in CIDR notation
		// and assert that it is IPv4 and global unicast
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}
		if !usableAddress(ip) {
			continue
		}
		got = ip.String()
		break
	}
	return
}

func usableAddress(ip net.IP) bool {
	return ip.To4() != nil && ip.IsGlobalUnicast()
}

func getDefaultGatewayIface() *net.Interface {
	log.Debug("Attempting to retrieve IP route info from netlink")
	routes, err := netlink.NetworkGetRoutes()
	if err != nil {
		log.Debugf("Unable to detect default interface: %v", err)
		return nil
	}
	if len(routes) == 0 {
		log.Debug("Netlink returned zero routes")
		return nil
	}
	for _, route := range routes {
		if route.Default {
			if route.Iface == nil {
				log.Debugf("Found default route but could not determine interface")
			}
			log.Debugf("Found default route with interface %v", route.Iface.Name)
			return route.Iface
		}
	}
	log.Debugf("Unable to find default route")
	return nil
}

func (m *machine) ID() string {
	return m.machID
}

func (m *machine) ShortID() string {
	if len(m.machID) <= shortIDLen {
		return m.machID
	}
	return m.machID[0:shortIDLen]
}

func (m *machine) MatchID(ID string) bool {
	return m.machID == ID || m.ShortID() == ID
}

func (m *machine) Status() *MachineStatus {
	return &MachineStatus{
		MachID:  m.machID,
		IsAlive: true,
		MachInfo: MachineInfo{
			HostName:   m.hostName,
			HostRegion: m.hostRegion,
			HostIDC:    m.hostIDC,
			PublicIP:   m.publicIP,
		},
		MachStat: m.getMachineStat(),
	}
}

func (m *machine) getMachineStat() MachineStat {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	return *m.stat
}
