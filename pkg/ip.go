package pkg

import (
	"github.com/docker/libcontainer/netlink"
	"github.com/ngaut/log"
	"net"
	"strconv"
	"strings"
)

// Method 1 to get local IP addr
func IntranetIP() (ips []string, err error) {
	ips = make([]string, 0)
	ifaces, e := net.Interfaces()
	if e != nil {
		return ips, e
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		// ignore docker and warden bridge
		if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "w-") {
			continue
		}
		addrs, e := iface.Addrs()
		if e != nil {
			return ips, e
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ipStr := ip.String()
			if IsIntranet(ipStr) {
				ips = append(ips, ipStr)
			}
		}
	}
	return ips, nil
}

func IsIntranet(ipStr string) bool {
	if strings.HasPrefix(ipStr, "10.") || strings.HasPrefix(ipStr, "192.168.") {
		return true
	}
	if strings.HasPrefix(ipStr, "172.") {
		// 172.16.0.0-172.31.255.255
		arr := strings.Split(ipStr, ".")
		if len(arr) != 4 {
			return false
		}
		second, err := strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			return false
		}
		if second >= 16 && second <= 31 {
			return true
		}
	}
	return false
}

// Method 2 to get local IP addr
func GetLocalIP() (got string) {
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
