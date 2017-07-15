package models

import (
	log "github.com/sirupsen/logrus"
	"github.com/nttdots/go-dots/dots_server/db_models"
)

// Connector object to the mitigation devices.
type DeviceConnector interface {
	connect(*LoginProfile) (err error)
	executeCommand(string) (err error)
}

type SshConnector struct {
	connected bool
}

/*
 * Establish SSH connections to the mitigation devices.
 *
 * parameter:
 *  login_profile login account information
 * return:
 *  err error
*/
func (d *SshConnector) connect(login_profile LoginProfile) (err error) {
	log.Infof("SshConnector.connect profile: %+v", login_profile)
	d.connected = true
	return
}

/*
 * Execute commands to the mitigation devices via established SSH connections.
 *
 * parameter:
 *  command command strings to be executed.
 * return:
 *  err error
 */
func (d *SshConnector) executeCommand(command string) (err error) {
	if d.connected {
		log.Infof("SshConnector.execute command: %s", command)
	}
	return
}

type ApiConnector struct {
}

func (d *ApiConnector) connect(login_profile LoginProfile) (err error) {
	return
}

/*
 * Execute commands to the mitigation devices via Device APIs(REST, gRPC...)

parameter:
  command command strings to be executed.
return:
  err error
*/
func (d ApiConnector) executeCommand(command string) (err error) {
	return
}

type LoginProfile struct {
	LoginMethod string
	LoginName   string
	Password    string
}

func (l *LoginProfile) Load(profile db_models.LoginProfile) {
	l.LoginMethod = profile.LoginMethod
	l.LoginName = profile.LoginName
	l.Password = profile.Password
}
