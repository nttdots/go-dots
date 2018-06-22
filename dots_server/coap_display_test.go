package main_test

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/nttdots/go-dots/libcoap"
	"github.com/nttdots/go-dots/dots_server"
)

func TestCoapHeaderDisplay(t *testing.T) {
	// test data setting
	token := bytes.NewBuffer(nil)
	token.WriteString("test_token_data")
	payload := bytes.NewBuffer(nil)
	payload.WriteString("test_payload_data")

	// return expected value
	exp_value := make([]string, 40)
	exp_value[0] = ""
	exp_value[1] = " 01.. .... = Version: 1"
	exp_value[2] = " ..00 .... = Type: Confirmable(0)"
	exp_value[3] = " .... 1111 = Token Length: 15"
	exp_value[4] = " Code: POST (2)"
	exp_value[5] = " Message ID: 12"
	exp_value[6] = " Token: test_token_data"
	exp_value[7] = ">Opt Name: #1: Uri-Host: test_host"
	exp_value[8] = ">Opt Name: #2: Uri-Port: 3366"
	exp_value[9] = ">Opt Name: #3: Max-Age: 20"
	exp_value[10] = ">Opt Name: #4: E-Tag: test_e-tag"
	exp_value[11] = ">Opt Name: #5: Location-Path: test_location_path"
	exp_value[12] = ">Opt Name: #6: Location-Query: test_location_query"
	exp_value[13] = ">Opt Name: #7: Uri-Path: .well-known"
	exp_value[14] = ">Opt Name: #8: Uri-Path: v1"
	exp_value[15] = ">Opt Name: #9: Uri-Path: dots-signal"
	exp_value[16] = ">Opt Name: #10: Uri-Path: test"
	exp_value[17] = ">Opt Name: #11: Observe: 10"
	exp_value[18] = ">Opt Name: #12: Content-Format: app/cbor"
	exp_value[19] = ">Opt Name: #13: If-Match: test_if_match"
	exp_value[20] = ">Opt Name: #14: If-None-Match: test_if_none_match"
	exp_value[21] = ">Opt Name: #15: Proxy-Uri: test_proxy_uri"
	exp_value[22] = ">Opt Name: #16: Proxy-Schema: test_proxy_schema"
	exp_value[23] = ">Opt Name: #17: Size1: 22"
	exp_value[24] = ">Opt Name: #18: Uri-Query: test_uri_query"
	exp_value[25] = " End of options marker: 255"
	exp_value[26] = " Payload: Payload Content-Format: app/cbor, Length: 17"
	exp_value[27] = "  Payload Desc: app/cbor"
	exp_value[28] = "  JavaScript Object Notation: app/cbor"
	exp_value[29] = "  Line-based text data: app/cbor"

	var c libcoap.Pdu = libcoap.Pdu{
		Type:      libcoap.TypeCon,
		Code:      libcoap.RequestPost,
		Token:     token.Bytes(),
		MessageID: 12,
		Data:      payload.Bytes(),
	}
	c.SetPathString(".well-known/v1/dots-signal/test")
	c.SetOption(libcoap.OptionUriHost, "test_host")
	c.SetOption(libcoap.OptionUriPort, 3366)
	c.SetOption(libcoap.OptionMaxage, 20)
	c.SetOption(libcoap.OptionEtag, "test_e-tag")
	c.SetOption(libcoap.OptionLocationPath, "test_location_path")
	c.SetOption(libcoap.OptionLocationQuery, "test_location_query")
	c.SetOption(libcoap.OptionObserve, 10)
	c.SetOption(libcoap.OptionContentFormat, libcoap.AppCbor)
	c.SetOption(libcoap.OptionIfMatch, "test_if_match")
	c.SetOption(libcoap.OptionIfNoneMatch, "test_if_none_match")
	c.SetOption(libcoap.OptionProxyUri, "test_proxy_uri")
	c.SetOption(libcoap.OptionProxyScheme, "test_proxy_schema")
	c.SetOption(libcoap.OptionSize1, 22)
	c.SetOption(libcoap.OptionUriQuery, "test_uri_query")

	ret := main.CoapHeaderDisplay(&c)
	// data check
	for i, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(ret, -1) {
		if v != exp_value[i] {
			t.Errorf("line[%d] got %s, want %s", i, v, exp_value[i])
		}
	}
}

func TestGetContentFormatValue(t *testing.T) {
	// test data setting
	impValues := make([]interface{}, 10)
	impValues[0] = libcoap.TextPlain
	impValues[1] = libcoap.AppJSON
	impValues[2] = libcoap.AppExi
	impValues[3] = libcoap.AppLinkFormat
	impValues[4] = libcoap.AppOctets
	impValues[5] = libcoap.AppXML
	impValues[6] = libcoap.AppCbor
	impValues[7] = libcoap.TypeAck

	// return expected value
	expValues := make([]string, 10)
	expValues[0] = "text/plain"
	expValues[1] = "app/json"
	expValues[2] = "app/exi"
	expValues[3] = "app/linkformat"
	expValues[4] = "app/octets"
	expValues[5] = "app/xml"
	expValues[6] = "app/cbor"
	expValues[7] = ""

	for i, v := range impValues {
		cmpValue := main.GetContentFormatValue(v)
		if cmpValue != expValues[i] {
			t.Errorf("line[%d] got %s, want %s", i, cmpValue, expValues[i])
		}
	}
}
