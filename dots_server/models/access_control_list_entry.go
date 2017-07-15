package models

type AccessControlListEntry struct {
	AclName           string
	AclType           string
	AccessListEntries *AccessListEntries
	Customer          *Customer
}

type AccessListEntries struct {
	Ace []Ace
}

type Ace struct {
	RuleName string
	Matches  *Matches
	Actions  *Actions
}

type Matches struct {
	SourceIpv4Network      Prefix
	DestinationIpv4Network Prefix
}

type Actions struct {
	Deny      []string
	Permit    []string
	RateLimit []string
}

func NewAccessControlListEntry(c *Customer) (s *AccessControlListEntry) {
	s = &AccessControlListEntry{
		"",
		"",
		&AccessListEntries{
			make([]Ace, 0),
		},
		c,
	}
	return
}

func NewAce() (a *Ace) {
	a = &Ace{
		RuleName: "",
		Matches: &Matches{
			SourceIpv4Network:      Prefix{},
			DestinationIpv4Network: Prefix{},
		},
		Actions: &Actions{
			Deny:      make([]string, 0),
			Permit:    make([]string, 0),
			RateLimit: make([]string, 0),
		},
	}
	return
}
