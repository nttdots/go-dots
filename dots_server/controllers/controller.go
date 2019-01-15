/*
 * Package "controllers" provides Dots API controllers.
 */
package controllers

import (
	"github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_server/models"
	"github.com/nttdots/go-dots/libcoap"
)

/*
 * Controller API interface.
 * It provides function interfaces correspondent to CoAP methods.
 */
type ControllerInterface interface {
	// Service methods
	HandleGet(Request, *models.Customer) (Response, error)
	HandlePost(Request, *models.Customer) (Response, error)
	HandleDelete(Request, *models.Customer) (Response, error)
	HandlePut(Request, *models.Customer) (Response, error)
}

type ServiceMethod func(req Request, customer *models.Customer) (Response, error)

/*
 * Controller super class.
 */
type Controller struct {
}

/*
 * Handles CoAP Get requests.
 */
func (c *Controller) HandleGet(req Request, customer *models.Customer) (Response, error) {
	return Response{
		Code: dots_common.MethodNotAllowed,
		Type: dots_common.NonConfirmable,
	}, nil
}

/*
 * Handles CoAP Post requests.
 */
func (c *Controller) HandlePost(req Request, customer *models.Customer) (Response, error) {
	return Response{
		Code: dots_common.MethodNotAllowed,
		Type: dots_common.NonConfirmable,
	}, nil
}

/*
 * Handles CoAP Delete requests.
 */
func (c *Controller) HandleDelete(req Request, customer *models.Customer) (Response, error) {
	return Response{
		Code: dots_common.MethodNotAllowed,
		Type: dots_common.NonConfirmable,
	}, nil
}

/*
 * Handles CoAP Put requests.
 */
func (c *Controller) HandlePut(req Request, customer *models.Customer) (Response, error) {
	return Response{
		Code: dots_common.MethodNotAllowed,
		Type: dots_common.NonConfirmable,
	}, nil
}

/*
 * Regular API request
 */
type Request struct {
	Code     libcoap.Code
	Type     libcoap.Type
	Uri      []string     // Full request URI
	PathInfo []string     // URI-Paths after with prefix
	Queries  []string     // Uri-Queries
	Body     interface{}
	Options  []libcoap.Option
}

/*
 * Regular API response
 */
type Response struct {
	Code dots_common.Code
	Type dots_common.Type // not used
	Options []libcoap.Option
	Body interface{}
}

/*
 * Error API response
 */
type Error struct {
	Code dots_common.Code
	Type dots_common.Type
	Body interface{}
}

func (e Error) Error() string {
	return e.Code.String()
}
