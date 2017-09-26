package main

import (
	"github.com/kirves/goradius"
	common "github.com/nttdots/go-dots/dots_common"
	log "github.com/sirupsen/logrus"
)

func main() {
	common.SetUpLogger()

	auth := goradius.Authenticator("127.0.0.1", "1812", "testing123")

	b, err := auth.Authenticate("client1", "password1", "")
	log.WithField("result", b).WithError(err).Info("auth result")
}
