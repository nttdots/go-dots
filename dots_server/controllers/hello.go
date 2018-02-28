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
func (h *Hello) HandlePost(newReq Request, customer *models.Customer) (res Response, err error) {

	req := newReq.Body

	if req == nil {
		res = Response {
			Type: dots_common.NonConfirmable,
			Code: dots_common.BadRequest,
			Body: nil,
		}
		return
	}

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
