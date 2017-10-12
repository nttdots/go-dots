package main

import (
	"context"
	"fmt"

	"github.com/nttdots/go-dots/dots_server/config"

	"net"

	dots_radius "github.com/nttdots/go-dots/dots_server/radius"
	log "github.com/sirupsen/logrus"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

type Authenticator struct {
	Enable     bool
	ServerAddr string
	Secret     string
	NASAddress net.IP
}

func (a *Authenticator) CheckClient(clientName, realm, password string, checkType dots_radius.ServiceType) (bool, error) {
	if !a.Enable {
		return true, nil
	}

	log.WithFields(log.Fields{
		"clientName": clientName,
		"realm":      realm,
		"password":   password,
	}).Debug("check client")

	var userName string
	if realm == "" {
		userName = clientName
	} else {
		userName = fmt.Sprintf("%s@%s", clientName, realm)
	}

	radiusPacket := radius.New(radius.CodeAccessRequest, []byte(a.Secret))
	rfc2865.UserName_SetString(radiusPacket, userName)
	rfc2865.UserPassword_SetString(radiusPacket, password)
	rfc2865.NASIPAddress_Set(radiusPacket, a.NASAddress)
	rfc2865.ServiceType_Add(radiusPacket, rfc2865.ServiceType(uint32(checkType)))

	response, err := radius.Exchange(context.Background(), radiusPacket, a.ServerAddr)
	if err != nil {
		return false, err
	}

	serviceTypes, err := rfc2865.ServiceType_Gets(response)
	if err != nil || response.Code != radius.CodeAccessAccept {
		return false, err
	}

	for _, a := range serviceTypes {
		if uint32(a) == uint32(checkType) {
			return true, nil
		}
	}
	return false, nil
}

func NewAuthenticator(aaa *config.AAA) *Authenticator {
	if !aaa.Enable {
		return &Authenticator{
			Enable: false,
		}
	}

	return &Authenticator{
		Enable:     true,
		ServerAddr: fmt.Sprintf("%s:%d", aaa.Host, aaa.Port),
		Secret:     aaa.Secret,
		NASAddress: aaa.ClientIPAddr,
	}
}
