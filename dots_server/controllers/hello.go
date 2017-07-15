package controllers

import (
	"fmt"

	"github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/models"
)

/*
 * Hello API controller returns response messages to the client.
 */
type Hello struct {
	Controller
}

/*
 * Handles Hello POST requests.
 */
func (h *Hello) Post(req interface{}, customer *models.Customer) (res Response, err error) {

	hr := req.(*messages.HelloRequest)
	rr := messages.HelloResponse{
		Message: fmt.Sprintf("hello, \"%s\"!", hr.Message),
	}

	res = Response{
		Code: dots_common.Valid,
		Type: dots_common.Acknowledgement,
		Body: rr,
	}

	return
}
