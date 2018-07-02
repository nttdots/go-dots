package data_controllers

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"regexp"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/nttdots/go-dots/dots_server/models"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
)

type ResourceDiscoveryController struct {
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

type XRD struct {
	Xmlns string `xml:"xmlns,attr"`
	Link  []Link
}

const (
	CONTENT_TYPE_YANG_DATA_JSON string = "application/yang-data+json"
	CONTENT_TYPE_XRD_XML        string = "application/xrd+xml"
	XRD_NAMESPACE               string = "http://docs.oasis-open.org/ns/xri/xrd-1.0"
)

func (c *ResourceDiscoveryController) Get(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
	log.Debugf("GET ResourceDiscoveryController")

	contentType := r.Header.Get("Content-Type")

	if contentType == CONTENT_TYPE_XRD_XML || contentType == "" {

    log.Debugf("Content-Type: %+v", CONTENT_TYPE_XRD_XML)
    config := dots_config.GetServerSystemConfig()

		resource := XRD{}
		resource.Xmlns = XRD_NAMESPACE
		resource.Link = []Link{{"restconf", config.Network.HrefOrigin + ":" + strconv.Itoa(config.Network.DataChannelPort) + config.Network.HrefPathname}}

		x, err := xml.MarshalIndent(resource, "", " ")
		if err != nil {
			return ErrorResponse(http.StatusInternalServerError)
		}

        xmlData := string(x)
        reg := regexp.MustCompile("(></Link>)")
        xmlData = reg.ReplaceAllString(xmlData, "/>")

		resp, e := ErrorResponse(http.StatusOK)
		resp.Headers.Set("Content-Type", CONTENT_TYPE_XRD_XML)
		resp.Content = []byte(xmlData)

		return resp, e

	} else if contentType == CONTENT_TYPE_YANG_DATA_JSON {
		log.Debugf("Content-Type: %+v", CONTENT_TYPE_YANG_DATA_JSON)
		return ErrorResponse(http.StatusUnsupportedMediaType)

	}

	return ErrorResponse(http.StatusBadRequest)
}
