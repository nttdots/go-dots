package controllers

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	common "github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/models"
)

/*
 * Controller for the session_configuration API.
 */
type SessionConfiguration struct {
	Controller
}

/*
 * Handles session_configuration POST requests and start the mitigation.
 *  1. Validate the received session configuration requests.
 *  2. return the validation results.
 *
 * parameter:
 *  request request message
 *  customer request source Customer
 * return:
 *  res response message
 *  err error
 */
func (m *SessionConfiguration) Post(request interface{}, customer *models.Customer) (res Response, err error) {

	payload := request.(*messages.SignalConfig)
	sessionConfigurationPayloadDisplay(payload)

	// validate
	signalSessionConfiguration := models.NewSignalSessionConfiguration(payload.SessionId, payload.HeartbeatInterval,
		payload.MissingHbAllowed, payload.MaxRetransmit, payload.AckTimeout, payload.AckRandomFactor)
	v := models.SignalConfigurationValidator{}
	validateResult := v.Validate(signalSessionConfiguration, *customer)
	if !validateResult {
		goto ResponseNG
	} else {
		// Register SignalConfigurationParameter
		_, err = models.CreateSignalSessionConfiguration(*signalSessionConfiguration, *customer)
		if err != nil {
			goto ResponseNG
		}

		goto ResponseOK
	}

ResponseNG:
// on validation error
	res = Response{
		Type: common.NonConfirmable,
		Code: common.BadRequest,
		Body: nil,
	}
	return
ResponseOK:
// on validation success
	res = Response{
		Type: common.NonConfirmable,
		Code: common.Created,
		Body: nil,
	}
	return
}

/*
 * Parse the request body and display the contents of the messages to stdout.
*/
func sessionConfigurationPayloadDisplay(data *messages.SignalConfig) {

	var result string = "\n"
	result += fmt.Sprintf("   \"%s\": %d\n", "session-id", data.SessionId)
	result += fmt.Sprintf("   \"%s\": %d\n", "heartbeat-interval", data.HeartbeatInterval)
	result += fmt.Sprintf("   \"%s\": %d\n", "missing-hb-allowed", data.MissingHbAllowed)
	result += fmt.Sprintf("   \"%s\": %d\n", "max-retransmit", data.MaxRetransmit)
	result += fmt.Sprintf("   \"%s\": %d\n", "ack-timeout", data.AckTimeout)
	result += fmt.Sprintf("   \"%s\": %f\n", "ack-random-factor", data.AckRandomFactor)
	log.Infoln(result)
}
