package radius

import "testing"

func TestParseUserType(t *testing.T) {

	if p, _ := ParseUserType("Login"); p != Login {
		t.Error("parse Login error")
	}

	if p, _ := ParseUserType("Framed"); p != Framed {
		t.Error("parse Framed error")
	}

	if p, _ := ParseUserType("Callback-Login"); p != CallbackLogin {
		t.Error("parse CallbackLogin error")
	}

	if p, _ := ParseUserType("CallbackFramed"); p != CallbackFramed {
		t.Error("parse CallbackFramed error")
	}

	if p, _ := ParseUserType("OUTBOUND"); p != Outbound {
		t.Error("parse Outbound error")
	}

	if p, _ := ParseUserType("administrative"); p != Administrative {
		t.Error("parse Administrative error")
	}

	if p, _ := ParseUserType("NAS Prompt"); p != NASPrompt {
		t.Error("parse NASPrompt error")
	}

	if p, _ := ParseUserType("AuthenticateOnly"); p != AuthenticateOnly {
		t.Error("parse AuthenticateOnly error")
	}

	if p, _ := ParseUserType("Callback NAS Prompt"); p != CallbackNASPrompt {
		t.Error("parse CallbackNASPrompt error")
	}

	if p, _ := ParseUserType("AuthenticateOnly"); p != AuthenticateOnly {
		t.Error("parse AuthenticateOnly error")
	}

	if p, _ := ParseUserType("CallCheck"); p != CallCheck {
		t.Error("parse CallCheck error")
	}

	if p, _ := ParseUserType("CallbackAdministrative"); p != CallbackAdministrative {
		t.Error("parse CallbackAdministrative error")
	}

	if _, err := ParseUserType("error"); err == nil {
		t.Error("parse invalid usertype error")
	}

}
