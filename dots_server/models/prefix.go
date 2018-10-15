package models

import (
	"net"
	"strconv"
	"net/url"
	log "github.com/sirupsen/logrus"
)

// Object to store CIDR
type Prefix struct {
	net *net.IPNet

	Addr      string
	PrefixLen int
}

const IP_PREFIX_LENGTH int = 32

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

/*
 * Create new prefixes from FQDN format strings.
 * parameters:
 *  fqdn     fqdn input string
 * return:
 *  p        ip prefix of the fqdn
 *  err      error
 */
func NewPrefixFromFQDN(fqdn string) (p []Prefix, err error) {
	ips, err := net.LookupIP(fqdn)
	if err != nil {
		log.Warnf("Failed to look-up ip from fqdn: %+v", fqdn)
		return
	}

	p, err = NewPrefixFromIps(ips)
	return
}

/*
 * Create new prefixes from URI format strings.
 * parameters:
 *  uri      uri input string
 * return:
 *  p        ip prefix of the uri
 *  err      error
 */
func NewPrefixFromURI(uri string) (p []Prefix, err error) {
	url, err := url.Parse(uri)
	if err != nil {
		log.Warnf("Failed to parse uri: %+v to object", uri)
		return
	}

	ips, err := net.LookupIP(url.Hostname())
	if err != nil {
		log.Warnf("Failed to look-up ip from fqdn: %+v", url.Hostname())
		return
	}

	p, err = NewPrefixFromIps(ips)
	return
}

func NewPrefixFromIps(ips []net.IP) (p []Prefix, err error) {
	for _, ip := range ips {
		cdir := ip.String() + "/" + strconv.Itoa(IP_PREFIX_LENGTH)
		var ipNet *net.IPNet
		_, ipNet, err = net.ParseCIDR(cdir)
		if err != nil {
			return
		}
		sz, _ := ipNet.Mask.Size()
		p = append(p, Prefix{ ipNet, ipNet.IP.String(), sz })
	}
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
	// This is an ip-address
	if prefix.PrefixLen == IP_PREFIX_LENGTH {
		return false
	}
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

// Remove overlap prefix
func RemoveOverlapPrefix(prefixlst []Prefix) ([]Prefix){
    pl := make([]Prefix, len(prefixlst))

    copy(pl, prefixlst)
    currentIdex := 0
    checkIndex := 1

    for currentIdex < len(pl){

        if checkIndex >= len(pl){
            currentIdex++
            checkIndex = currentIdex + 1
            continue
        }

        if pl[currentIdex].Includes(&pl[checkIndex]){
            if checkIndex == len(pl) - 1{
                pl = pl[:checkIndex]
            } else if checkIndex < len(pl) - 1{
                pl = append(pl[:checkIndex], pl[checkIndex+1:]...)
            }else{
                break;
            }
            currentIdex = 0
            checkIndex = currentIdex + 1
            continue
        }

        if pl[checkIndex].Includes(&pl[currentIdex]){
			if currentIdex == 0{
				pl = pl[currentIdex+1:]
			} else if currentIdex > 0{
				pl = append(pl[:currentIdex], pl[currentIdex+1:]...)
			}
            currentIdex = 0
            checkIndex = currentIdex + 1
            continue
        }

        checkIndex++
    }

    return pl
}