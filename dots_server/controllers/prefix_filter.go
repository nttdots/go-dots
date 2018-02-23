package controllers

import (
    "strings"
    "github.com/nttdots/go-dots/dots_common"
    "github.com/nttdots/go-dots/dots_server/models"
)

type PrefixFilter struct {
    Controller
    prefix     []string
    controller ControllerInterface
}

func NewPrefixFilter(path string, controller ControllerInterface) ControllerInterface {
    prefix := strings.Split(path, "/")
    if 0 < len(prefix) && prefix[0] == "" {
        prefix = prefix[1:]
    }
    if 0 < len(prefix) && prefix[len(prefix) - 1] == "" {
        prefix = prefix[:len(prefix) - 1]
    }

    return &PrefixFilter{ Controller{}, prefix, controller }
}

func (p *PrefixFilter) forward(req Request, customer *models.Customer, method ServiceMethod) (res Response, err error) {
    startsWith := func(u []string, p []string) bool {
        if len(u) < len(p) { return false }

        for i, s := range p {
            if u[i] != s { return false }
        }
        return true
    }

    if startsWith(req.Uri, p.prefix) {
        return method(req, customer)
    } else {
        return Response{
            Code: dots_common.NotFound,
            Type: dots_common.NonConfirmable,
        }, nil
    }
}

func (p *PrefixFilter) HandleGet(req Request, customer *models.Customer) (res Response, err error) {
    return p.forward(req, customer, p.controller.HandleGet)
}

func (p *PrefixFilter) HandlePost(req Request, customer *models.Customer) (res Response, err error) {
    return p.forward(req, customer, p.controller.HandlePost)
}

func (p *PrefixFilter) HandleDelete(req Request, customer *models.Customer) (res Response, err error) {
    return p.forward(req, customer, p.controller.HandleDelete)
}

func (p *PrefixFilter) HandlePut(req Request, customer *models.Customer) (res Response, err error) {
    return p.forward(req, customer, p.controller.HandlePut)
}
