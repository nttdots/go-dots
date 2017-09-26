package main

import (
	"strconv"

	"github.com/kirves/goradius"
	"github.com/nttdots/go-dots/dots_server/config"
	log "github.com/sirupsen/logrus"
)

type Authenticator struct {
	enable              bool
	radiusAuthenticator *goradius.AuthenticatorT
}

func (a *Authenticator) CheckClient(clientName, password, nasId string) (bool, error) {
	if !a.enable {
		return true, nil
	}

	log.WithFields(log.Fields{
		"clientName": clientName,
		"password":   password,
		"nasId":      nasId,
	}).Debug("check client")

	return a.radiusAuthenticator.Authenticate(clientName, password, nasId)
}

func NewAuthenticator(aaa *config.AAA) *Authenticator {
	if !aaa.Enable {
		return &Authenticator{
			enable:              false,
			radiusAuthenticator: nil,
		}
	}

	radiusAuth := goradius.Authenticator(aaa.Server, strconv.Itoa(aaa.Port), aaa.Secret)
	return &Authenticator{
		enable:              true,
		radiusAuthenticator: radiusAuth,
	}
}
