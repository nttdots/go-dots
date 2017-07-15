package models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/models"
)

var testAccessControlListEntry models.AccessControlListEntry
var testUpdAccessControlListEntry models.AccessControlListEntry

func accessControlListEntrySampleDataCreate() {
	// AccessControlListEntry
	testAccessControlListEntry = models.AccessControlListEntry{}
	testUpdAccessControlListEntry = models.AccessControlListEntry{}

	// matches create test data
	soruceIpv4Network, _ := models.NewPrefix("10.10.10.1/24")
	destinationIpv4Network, _ := models.NewPrefix("10.10.10.2/24")

	testAccessControlListEntry.AclName = "aclname1"
	testAccessControlListEntry.AclType = "ipv4"
	testAccessControlListEntry.AccessListEntries = &models.AccessListEntries{
		Ace: []models.Ace{{
			RuleName: "rule1",
			Matches: &models.Matches{
				SourceIpv4Network:      soruceIpv4Network,
				DestinationIpv4Network: destinationIpv4Network,
			},
			Actions: &models.Actions{
				Deny: []string{"deny"},
			},
		}},
	}

	// matches update test data
	updSourceIpv4Network, _ := models.NewPrefix("210.210.210.1/24")
	updDestinationIpv4Network, _ := models.NewPrefix("210.210.210.2/24")

	testUpdAccessControlListEntry.AclName = "aclname2"
	testUpdAccessControlListEntry.AclType = "ipv6"
	testUpdAccessControlListEntry.AccessListEntries = &models.AccessListEntries{
		Ace: []models.Ace{{
			RuleName: "rule2",
			Matches: &models.Matches{
				SourceIpv4Network:      updSourceIpv4Network,
				DestinationIpv4Network: updDestinationIpv4Network,
			},
			Actions: &models.Actions{
				Deny:      []string{"deny"},
				RateLimit: []string{"rate_limit1", "rate_limit2"},
			},
		}},
	}
}

func TestCreateAccessControlList(t *testing.T) {
	customer, err := models.GetCustomer(123)
	if err != nil {
		t.Errorf("GetCustomer err: %s", err)
	}
	_, err = models.CreateAccessControlList(&testAccessControlListEntry, &customer)
	if err != nil {
		t.Errorf("CreateAccessControlList err: %s", err)
	}
}

func TestGetAccessControlList(t *testing.T) {
	accessControlList, err := models.GetAccessControlList(123)
	if err != nil {
		t.Errorf("get accessControlList err: %s", err)
		return
	}

	if accessControlList.AclName != "aclname1" {
		t.Errorf("AclName got %s, want %s", accessControlList.AclName, testAccessControlListEntry.AclName)
		return
	}

	if accessControlList.AclType != "ipv4" {
		t.Errorf("AclType got %s, want %s", accessControlList.AclType, testAccessControlListEntry.AclType)
		return
	}

	for _, srcAce := range accessControlList.AccessListEntries.Ace {
		// RuleName check
		findFlag := false
		for _, testAce := range testAccessControlListEntry.AccessListEntries.Ace {
			if srcAce.RuleName == testAce.RuleName {
				findFlag = true
				break
			}

		}
		if !findFlag {
			t.Errorf("no RuleName data: %s", srcAce.RuleName)
		}

		// Matches.SourceIpv4Network check
		findFlag = false
		for _, testAce := range testAccessControlListEntry.AccessListEntries.Ace {
			if srcAce.Matches.SourceIpv4Network.Addr == testAce.Matches.SourceIpv4Network.Addr &&
				srcAce.Matches.SourceIpv4Network.PrefixLen == testAce.Matches.SourceIpv4Network.PrefixLen {
				findFlag = true
				break
			}

		}
		if !findFlag {
			t.Errorf("no SourceIpv4Network data: %s", srcAce.Matches.SourceIpv4Network)
		}

		// Matches.DestinationIpv4Network check
		findFlag = false
		for _, testAce := range testAccessControlListEntry.AccessListEntries.Ace {
			if srcAce.Matches.DestinationIpv4Network.Addr == testAce.Matches.DestinationIpv4Network.Addr &&
				srcAce.Matches.DestinationIpv4Network.PrefixLen == testAce.Matches.DestinationIpv4Network.PrefixLen {
				findFlag = true
				break
			}

		}
		if !findFlag {
			t.Errorf("no DestinationIpv4Network data: %s", srcAce.Matches.DestinationIpv4Network)
		}

		// Actions.Deny check
		for _, testAce := range testAccessControlListEntry.AccessListEntries.Ace {
			for _, srcDeny := range srcAce.Actions.Deny {
				findFlag = false
				for _, testDeny := range testAce.Actions.Deny {
					if srcDeny == testDeny {
						findFlag = true
						break
					}
				}
				if !findFlag {
					t.Errorf("no Actions.Deny data: %s", srcDeny)
				}
			}
		}

		// Actions.Permit check
		for _, testAce := range testAccessControlListEntry.AccessListEntries.Ace {
			for _, srcPermit := range srcAce.Actions.Permit {
				findFlag = false
				for _, testPermit := range testAce.Actions.Permit {
					if srcPermit == testPermit {
						findFlag = true
						break
					}
				}
				if !findFlag {
					t.Errorf("no Actions.Permit data: %s", srcPermit)
				}
			}
		}

		// Actions.RateLimit check
		for _, testAce := range testAccessControlListEntry.AccessListEntries.Ace {
			for _, srcRateLimit := range srcAce.Actions.RateLimit {
				findFlag = false
				for _, testRateLimit := range testAce.Actions.RateLimit {
					if srcRateLimit == testRateLimit {
						findFlag = true
						break
					}
				}
				if !findFlag {
					t.Errorf("no Actions.RateLimit data: %s", srcRateLimit)
				}
			}
		}
	}
}

func TestUpdateAccessControlList(t *testing.T) {
	customer, err := models.GetCustomer(127)
	if err != nil {
		t.Errorf("GetCustomer err: %s", err)
		return
	}
	err = models.UpdateAccessControlList(&testUpdAccessControlListEntry, &customer)
	if err != nil {
		t.Errorf("UpdateAccessControlList err: %s", err)
		return
	}

	accessControlList, err := models.GetAccessControlList(127)
	if err != nil {
		t.Errorf("get accessControlList err: %s", err)
		return
	}

	if accessControlList.AclName != testUpdAccessControlListEntry.AclName {
		t.Errorf("AclName got %s, want %s", accessControlList.AclName, testUpdAccessControlListEntry.AclName)
		return
	}

	if accessControlList.AclType != testUpdAccessControlListEntry.AclType {
		t.Errorf("AclType got %s, want %s", accessControlList.AclType, testUpdAccessControlListEntry.AclType)
		return
	}

	for _, srcAce := range accessControlList.AccessListEntries.Ace {
		// RuleName check
		findFlag := false
		for _, testAce := range testUpdAccessControlListEntry.AccessListEntries.Ace {
			if srcAce.RuleName == testAce.RuleName {
				findFlag = true
				break
			}

		}
		if !findFlag {
			t.Errorf("no RuleName data: %s", srcAce.RuleName)
		}

		// Matches.SourceIpv4Network check
		findFlag = false
		for _, testAce := range testUpdAccessControlListEntry.AccessListEntries.Ace {
			if srcAce.Matches.SourceIpv4Network.Addr == testAce.Matches.SourceIpv4Network.Addr &&
				srcAce.Matches.SourceIpv4Network.PrefixLen == testAce.Matches.SourceIpv4Network.PrefixLen {
				findFlag = true
				break
			}

		}
		if !findFlag {
			t.Errorf("no SourceIpv4Network data: %s", srcAce.Matches.SourceIpv4Network)
		}

		// Matches.DestinationIpv4Network check
		findFlag = false
		for _, testAce := range testUpdAccessControlListEntry.AccessListEntries.Ace {
			if srcAce.Matches.DestinationIpv4Network.Addr == testAce.Matches.DestinationIpv4Network.Addr &&
				srcAce.Matches.DestinationIpv4Network.PrefixLen == testAce.Matches.DestinationIpv4Network.PrefixLen {
				findFlag = true
				break
			}

		}
		if !findFlag {
			t.Errorf("no DestinationIpv4Network data: %s", srcAce.Matches.DestinationIpv4Network)
		}

		// Actions.Deny check
		for _, testAce := range testUpdAccessControlListEntry.AccessListEntries.Ace {
			for _, srcDeny := range srcAce.Actions.Deny {
				findFlag = false
				for _, testDeny := range testAce.Actions.Deny {
					if srcDeny == testDeny {
						findFlag = true
						break
					}
				}
				if !findFlag {
					t.Errorf("no Actions.Deny data: %s", srcDeny)
				}
			}
		}

		// Actions.Permit check
		for _, testAce := range testUpdAccessControlListEntry.AccessListEntries.Ace {
			for _, srcPermit := range srcAce.Actions.Permit {
				findFlag = false
				for _, testPermit := range testAce.Actions.Permit {
					if srcPermit == testPermit {
						findFlag = true
						break
					}
				}
				if !findFlag {
					t.Errorf("no Actions.Permit data: %s", srcPermit)
				}
			}
		}

		// Actions.RateLimit check
		for _, testAce := range testUpdAccessControlListEntry.AccessListEntries.Ace {
			for _, srcRateLimit := range srcAce.Actions.RateLimit {
				findFlag = false
				for _, testRateLimit := range testAce.Actions.RateLimit {
					if srcRateLimit == testRateLimit {
						findFlag = true
						break
					}
				}
				if !findFlag {
					t.Errorf("no Actions.RateLimit data: %s", srcRateLimit)
				}
			}
		}
	}
}
