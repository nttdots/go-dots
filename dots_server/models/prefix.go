package models

import (
	"net"
)

// Object to store CIDR
type Prefix struct {
	net *net.IPNet

	Addr      string
	PrefixLen int
}

/*
 * Convert to CIDR strings.
 */
func (p *Prefix) String() string {
	return p.net.String()
}

// Create new prefixes from CIDR format strings.
func NewPrefix(addrString string) (p Prefix, err error) {
	_, ipNet, err := net.ParseCIDR(addrString)
	if err != nil {
		return
	}
	sz, _ := ipNet.Mask.Size()
	p = Prefix{ipNet, ipNet.IP.String(), sz}

	return
}

// Obtain the broadcast address.
func (p *Prefix) LastIP() (ip net.IP) {
	ip = p.net.IP
	for i := 0; i < len(p.net.IP); i++ {
		ip[i] |= ^p.net.Mask[i]
	}
	return
}

// Check if this prefix includes the prefix.
func (p *Prefix) Includes(other *Prefix) (ret bool) {
	if p.net == nil || other.net == nil {
		panic(p)
	}
	otherIP := other.net.IP
	otherFirst := otherIP.Mask(other.net.Mask)
	otherLast := other.LastIP()

	return p.net.Contains(otherFirst) && p.net.Contains(otherLast)
}

// Check if this prefix includes the address(string).
func (p *Prefix) Validate(addr string) (ret bool) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return false
	}
	return p.net.Contains(ip)
}

func ConvertAddrStringToPrefix(addrString []string) []Prefix {
	var ret = make([]Prefix, len(addrString))
	count := 0

	for _, addr := range addrString {
		prefix, err := NewPrefix(addr)
		if err != nil {
			continue
		}
		ret[count] = prefix
		count++
	}

	return ret
}
