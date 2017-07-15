package models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/models"
)

func TestNewPrefix(t *testing.T) {
	v4, err := models.NewPrefix("192.168.0.0/25")
	if err != nil {
		t.Errorf("ipv4_parse error %s", err)
	}
	var expects interface{}

	expects = "192.168.0.0"
	if v4.Addr != expects {
		t.Errorf("ipv4_addr got %s, want %s", v4.Addr, expects)
	}
	expects = "192.168.0.127"
	if v4.LastIP().String() != expects {
		t.Errorf("ipv4_last_ip got %s, want %s", v4.LastIP().String(), expects)
	}
	expects = 25
	if v4.PrefixLen != expects {
		t.Errorf("ipv4_prefix got %s, want %s", v4.PrefixLen, expects)
	}
	if v4.Validate("192.168.0.192") {
		t.Errorf("ipV4_include, got %s, want %s", true, false)
	}
	if !v4.Validate("192.168.0.64") {
		t.Errorf("ipV4_include, got %s, want %s", false, true)
	}

	v6, err := models.NewPrefix("2001:db8:abcd:3f01::/64")
	if err != nil {
		t.Errorf("ipv6_parse error %s", err)
	}
	expects = "2001:db8:abcd:3f01::"
	if v6.Addr != expects {
		t.Errorf("ipv6_addr got %s, want %s", v6.Addr, expects)
	}
	expects = 64
	if v6.PrefixLen != expects {
		t.Errorf("ipv6_addr got %s, want %s", v6.PrefixLen, expects)
	}
	expects = "2001:db8:abcd:3f01:ffff:ffff:ffff:ffff"
	if v6.LastIP().String() != expects {
		t.Errorf("ipv4_last_ip got %s, want %s", v6.LastIP().String(), expects)
	}

}
