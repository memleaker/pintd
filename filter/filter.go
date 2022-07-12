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

var denyaddr = make(map[string][]Addr, 0)

// Add deny address for pintd config.
func AddDenyAddrs(cfg *config.PintdConfig) {
	for _, section := range cfg.Redirects {
		AddDenyAddr(section.Denyaddr, section.SectionName)
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

// filter address, if matched return true, if no matched return false.
func FilterAddr(addr string, section string) bool {
	ip, mask := ParseIpAddr(addr)

	for _, deny := range denyaddr[section] {
		if ip&deny.mask == deny.ip&mask {
			return true
		}
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
