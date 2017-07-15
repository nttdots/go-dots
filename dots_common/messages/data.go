package messages

import (
	"fmt"
)

/*
 * struct that express CreateIdetifier on the data channel.
 */
type CreateIdentifier struct {
	Identifier Identifier `json:"identifier" cbor:"identifier"`
}

/*
 * struct to store Identifiers
 */
type Identifier struct {
	Alias []Alias `json:"alias" cbor:"alias"`
}

/*
 * struct to store Aliases.
 */
type Alias struct {
	// Name of the alias.  This is a mandatory attribute.
	AliasName string `json:"alias-name" cbor:"alias-name"`
	// IP addresses are separated by commas.  This is an optional attribute.
	Ip []string `json:"ip" cbor:"ip"`
	// Prefixes are separated by commas.  This is an optional attribute.
	Prefix []string `json:"prefix" cbor:"prefix"`
	// The port range, lower-port for lower port number and upper-port for upper port number.
	// For TCP, UDP, SCTP, or DCCP: the range of ports (e.g., 80 to 8080).
	// This is an optional attribute.
	PortRange []PortRange `json:"port-range" cbor:"port-range"`
	// Internet Protocol numbers.  This is an optional attribute.
	TrafficProtocol []int `json:"traffic-protocol" cbor:"traffic-protocol"`
	// FQDN
	FQDN []string `json:"FQDN" cbor:"FQDN"`
	// URI
	URI []string `json:"URI" cbor:"URI"`
	// E.164";
	E164 []string `json:"E.164" cbor:"E.164"`
}

/*
 * TargetPortRange
 */
type PortRange struct {
	LowerPort int `json:"lower-port" cbor:"lower-port"`
	UpperPort int `json:"upper-port" cbor:"upper-port"`
}

/*
 * convert given Identifiers to strings.
 */
func (m *CreateIdentifier) String() (result string) {
	result = "\n"
	for key, alias := range m.Identifier.Alias {
		result += fmt.Sprintf("   \"%s[%d]\":\n", "alias", key+1)
		result += fmt.Sprintf("     \"%s\": %s\n", "alias-name", alias.AliasName)
		if alias.Ip != nil {
			for k, v := range alias.Ip {
				result += fmt.Sprintf("     \"%s[%d]\": %s\n", "ip", k+1, v)
			}
		}
		if alias.Prefix != nil {
			for k, v := range alias.Prefix {
				result += fmt.Sprintf("     \"%s[%d]\": %s\n", "prefix", k+1, v)
			}
		}
		if alias.PortRange != nil {
			for k, v := range alias.PortRange {
				result += fmt.Sprintf("     \"%s[%d]\":\n", "port-range", k+1)
				result += fmt.Sprintf("       \"%s\": %d\n", "lower-port", v.LowerPort)
				result += fmt.Sprintf("       \"%s\": %d\n", "upper-port", v.UpperPort)
			}
		}
		if alias.FQDN != nil {
			for k, v := range alias.FQDN {
				result += fmt.Sprintf("     \"%s[%d]\": %s\n", "FQDN", k+1, v)
			}
		}
		if alias.URI != nil {
			for k, v := range alias.URI {
				result += fmt.Sprintf("     \"%s[%d]\": %s\n", "URI", k+1, v)
			}
		}
		if alias.E164 != nil {
			for k, v := range alias.E164 {
				result += fmt.Sprintf("     \"%s[%d]\": %s\n", "E.164", k+1, v)
			}
		}
	}
	return
}

/*
 * struct to store InstallFilteringRule.
 */
type InstallFilteringRule struct {
	AccessLists AccessLists `json:"access-lists" cbor:"access-lists"`
}

/*
 * struct to store AccessLists
 */
type AccessLists struct {
	Acl []Acl `json:"acl" cbor:"acl"`
}

/*
 * struct to store ACL.
 */
type Acl struct {
	// The name of access-list.  This is a mandatory attribute.
	AclName string `json:"acl-name" cbor:"acl-name"`
	// Indicates the primary intended type of match criteria (e.g.  IPv4, IPv6).  This is a mandatory attribute.
	AclType           string            `json:"acl-type" cbor:"acl-type"`
	AccessListEntries AccessListEntries `json:"access-list-entries" cbor:"access-list-entries"`
}

/*
 * struct to store AccessListEntries
 */
type AccessListEntries struct {
	Ace []Ace `json:"ace" cbor:"ace"`
}

/*
 * struct to store ACL rules.
 */
type Ace struct {
	RuleName string  `json:"rule-name" cbor:"rule-name"`
	Matches  Matches `json:"matches" cbor:"matches"`
	// deny" or "permit" or "rate-limit".
	// "permit" action is used to white-list traffic.
	// "deny" action is used to black-list traffic.
	// "rate-limit" action is used to rate-limit traffic, the allowed traffic rate is represented in bytes per second
	// indicated in IEEE floating point format [IEEE.754.1985].
	// If actions attribute is not specified in the request then the default action is "deny".
	// This is an optional attribute.
	Actions Actions `json:"actions" cbor:"actions"`
}

/*
 * struct to match ACL rules.
 */
type Matches struct {
	// The source IPv4 prefix.  This is an optional attribute.
	SourceIpv4Network string `json:"source-ipv4-network" cbor:"source-ipv4-network"`
	// The destination IPv4 prefix.  This is an optional attribute.
	DestinationIpv4Network string `json:"destination-ipv4-network" cbor:"destination-ipv4-network"`
}

/*
 * struct to hold ACL actions.
 */
type Actions struct {
	Deny      []string `json:"deny" cbor:"deny"`
	Permit    []string `json:"permit" cbor:"permit"`
	RateLimit []string `json:"rate-limit" cbor:"rate-limit"`
}

/*
 * convert InstallFilteringRule to strings.
 */
func (m *InstallFilteringRule) String() (result string) {
	result = "\n"
	for key, acl := range m.AccessLists.Acl {
		result += fmt.Sprintf("   \"%s[%d]\":\n", "acl", key+1)
		result += fmt.Sprintf("     \"%s\": %s\n", "acl-name", acl.AclName)
		result += fmt.Sprintf("     \"%s\": %s\n", "acl-type", acl.AclType)
		if acl.AccessListEntries.Ace != nil {
			for k, v := range acl.AccessListEntries.Ace {
				result += fmt.Sprintf("     \"%s[%d]\":\n", "ace", k+1)
				result += fmt.Sprintf("       \"%s\": %s\n", "rule-name", v.RuleName)
				result += fmt.Sprintf("       \"%s\":\n", "matches")
				result += fmt.Sprintf("         \"%s\": %s\n", "source-ipv4-network", v.Matches.SourceIpv4Network)
				result += fmt.Sprintf("         \"%s\": %s\n", "destination-ipv4-network", v.Matches.DestinationIpv4Network)
				result += fmt.Sprintf("       \"%s\":\n", "actions")
				result += fmt.Sprintf("         \"%s\": %s\n", "deny", v.Actions.Deny)
				result += fmt.Sprintf("         \"%s\": %s\n", "permit", v.Actions.Permit)
				result += fmt.Sprintf("         \"%s\": %s\n", "rate-limit", v.Actions.RateLimit)
			}
		}
	}
	return
}
