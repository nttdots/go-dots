package radius

import "testing"

func TestParseUserType(t *testing.T) {

	if p, _ := ParseServiceType("Login"); p != Login {
		t.Error("parse Login error")
	}

	if p, _ := ParseServiceType("Framed"); p != Framed {
		t.Error("parse Framed error")
	}

	if p, _ := ParseServiceType("Callback-Login"); p != CallbackLogin {
		t.Error("parse CallbackLogin error")
	}

	if p, _ := ParseServiceType("CallbackFramed"); p != CallbackFramed {
		t.Error("parse CallbackFramed error")
	}

	if p, _ := ParseServiceType("OUTBOUND"); p != Outbound {
		t.Error("parse Outbound error")
	}

	if p, _ := ParseServiceType("administrative"); p != Administrative {
		t.Error("parse Administrative error")
	}

	if p, _ := ParseServiceType("NAS Prompt"); p != NASPrompt {
		t.Error("parse NASPrompt error")
	}

	if p, _ := ParseServiceType("AuthenticateOnly"); p != AuthenticateOnly {
		t.Error("parse AuthenticateOnly error")
	}

	if p, _ := ParseServiceType("Callback NAS Prompt"); p != CallbackNASPrompt {
		t.Error("parse CallbackNASPrompt error")
	}

	if p, _ := ParseServiceType("AuthenticateOnly"); p != AuthenticateOnly {
		t.Error("parse AuthenticateOnly error")
	}

	if p, _ := ParseServiceType("CallCheck"); p != CallCheck {
		t.Error("parse CallCheck error")
	}

	if p, _ := ParseServiceType("CallbackAdministrative"); p != CallbackAdministrative {
		t.Error("parse CallbackAdministrative error")
	}

	if _, err := ParseServiceType("error"); err == nil {
		t.Error("parse invalid usertype error")
	}

}
