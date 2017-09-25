package main

import (
	"strconv"

	"github.com/kirves/goradius"
	"github.com/nttdots/go-dots/dots_server/config"
	log "github.com/sirupsen/logrus"
)

func NewAuthenticator(aaa *config.AAA) chan<- interface{} {
	requestChan := make(chan interface{}, 100)
	var radiusAuth *goradius.AuthenticatorT
	if aaa.Enable {
		radiusAuth = goradius.Authenticator(aaa.Server, strconv.Itoa(aaa.Port), aaa.Secret)
	}

	go func() {
		for {
			switch req := requestChan.(type) {
			case RadiusAuthenticatorIdentifier:
				if radiusAuth != nil {
					b, err := radiusAuth.Authenticate(req.Username, req.Password, "1")
					if err != nil {
						log.WithError(err).Error("authenticate error.")
						req.Result <- false
					}
					req.Result <- b
				} else {
					log.WithField("request", req).Error("unknown authenticate request.")
					req.Result <- true
				}
			default:
				log.WithField("request", req).Error("unknown authenticate request.")
			}
		}
	}()
	return requestChan
}

type RadiusAuthenticatorIdentifier struct {
	Username string
	Password string
	Result   <-chan bool
}
