package models

import (
	"net"
	"strconv"
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

// IsBroadCast reports whether ip is a broadcast address.
func (prefix *Prefix) IsBroadCast() bool {
	targetPrefix := prefix.net.IP.String()
	if ip4 := prefix.net.IP.To4(); ip4 != nil {
		return targetPrefix == prefix.LastIP().String()
	}
	return false
}

// IsMulticast reports whether ip is a multicast address.
func (prefix *Prefix) IsMulticast() bool {
	return prefix.net.IP.IsMulticast()
}

// IsLoopback reports whether ip is a loopback address.
func (prefix *Prefix) IsLoopback() bool {
	return prefix.net.IP.IsLoopback()
}

// Check validate range ip address
func (prefix Prefix) CheckValidRangeIpAddress(addressRangePrefixes AddressRange) (bool, []string){
	if !addressRangePrefixes.Validate(prefix) {
		var addressRange []string
        for _,address := range addressRangePrefixes.Prefixes {
          addressRange = append(addressRange, address.Addr+"/"+ strconv.Itoa(address.PrefixLen))
        }
		return false, addressRange
	}
	return true, nil
}