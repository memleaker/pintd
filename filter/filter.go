package filter

import (
	"fmt"
	"pintd/config"
	"strings"
)

type Addr struct {
	ip   uint32
	mask uint32
}

var admitaddr = make(map[string][]Addr, 0)
var denyaddr = make(map[string][]Addr, 0)

// init filter
func InitFilter(cfg *config.PintdConfig) {
	for _, section := range cfg.Redirects {
		AddAdmitAddr(section.Admitaddr, section.SectionName)
		AddDenyAddr(section.Denyaddr, section.SectionName)
	}
}

// Add Admit address for single redirect config.
func AddAdmitAddr(admit []string, section string) {
	admitaddr[section] = make([]Addr, 0)

	for _, addr := range admit {
		ip, mask := ParseIpAddr(addr)
		admitaddr[section] = append(admitaddr[section], Addr{ip, mask})
	}
}

// Add Deny address for single redirect config.
func AddDenyAddr(deny []string, section string) {
	denyaddr[section] = make([]Addr, 0)

	for _, addr := range deny {
		ip, mask := ParseIpAddr(addr)
		denyaddr[section] = append(denyaddr[section], Addr{ip, mask})
	}
}

func Match(ip, mask uint32, addrs []Addr) bool {
	for _, addr := range addrs {
		if ip&addr.mask == addr.ip&mask {
			return true
		}
	}

	return false
}

// filter address, return true means access deny.
func DenyAccess(addr string, section string) bool {
	ip, mask := ParseIpAddr(addr)

	if len(admitaddr[section]) > 0 {
		admit := Match(ip, mask, admitaddr[section])
		if !admit {
			return !admit
		}

		deny := Match(ip, mask, denyaddr[section])
		return deny
	}

	if len(denyaddr[section]) > 0 {
		return Match(ip, mask, denyaddr[section])
	}

	return false
}

// parse ip addr to ip and mask.
func ParseIpAddr(addr string) (uint32, uint32) {
	var (
		ip   uint32 = 0
		mask uint32 = 32
		arr  [4]uint32
	)

	if ipstr, maskstr, ok := strings.Cut(addr, "/"); ok {
		addr = ipstr

		fmt.Sscanf(maskstr, "%d", &mask)
	}

	fmt.Sscanf(addr, "%d.%d.%d.%d", &arr[0], &arr[1], &arr[2], &arr[3])

	ip |= arr[0] << 24
	ip |= arr[1] << 16
	ip |= arr[2] << 8
	ip |= arr[3]

	mask = (0xffffffff << (32 - mask))

	return ip, mask
}
