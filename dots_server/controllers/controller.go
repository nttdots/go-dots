/*
 * Package "controllers" provides Dots API controllers.
 */
package controllers

import (
	"github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_server/models"
)

/*
 * Controller API interface.
 * It provides function interfaces correspondent to CoAP methods.
 */
type ControllerInterface interface {
	Get(interface{}, *models.Customer) (Response, error)
	Post(interface{}, *models.Customer) (Response, error)
	Delete(interface{}, *models.Customer) (Response, error)
	Put(interface{}, *models.Customer) (Response, error)
}

/*
 * Controller super class.
 */
type Controller struct {
}

/*
 * Handles CoAP Get requests.
 */
func (c *Controller) Get(interface{}, *models.Customer) (Response, error) {
	return Response{
		Code: dots_common.MethodNotAllowed,
		Type: dots_common.NonConfirmable,
	}, nil
}

/*
 * Handles CoAP Post requests.
 */
func (c *Controller) Post(interface{}, *models.Customer) (Response, error) {
	return Response{
		Code: dots_common.MethodNotAllowed,
		Type: dots_common.NonConfirmable,
	}, nil
}

/*
 * Handles CoAP Delete requests.
 */
func (c *Controller) Delete(interface{}, *models.Customer) (Response, error) {
	return Response{
		Code: dots_common.MethodNotAllowed,
		Type: dots_common.NonConfirmable,
	}, nil
}

/*
 * Handles CoAP Put requests.
 */
func (c *Controller) Put(interface{}, *models.Customer) (Response, error) {
	return Response{
		Code: dots_common.MethodNotAllowed,
		Type: dots_common.NonConfirmable,
	}, nil
}

/*
 * Regular API response
 */
type Response struct {
	Code dots_common.Code
	Type dots_common.Type
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
