package main

import (
	"fmt"
	"strconv"

	"github.com/nttdots/go-dots/coap"
)

/*
 * Output formatted CoAP messages to the stdout.
 */
func CoapHeaderDisplay(request *coap.Message) string {
	var result string = "\n"
	// Version
	result += fmt.Sprint(" 01.. .... = Version: 1\n")
	// Type
	var value, key = getTypeValue(request.Type)
	result += fmt.Sprintf(" ..%s .... = Type: %s(%d)\n", decimalToBinary(value, 2), key, value)
	// Token Length
	var token_length int = len(request.Token)
	result += fmt.Sprintf(" .... %s = Token Length: %d\n", decimalToBinary(token_length, 4), token_length)
	// Code
	value, key = getMethodValue(request.Code)
	result += fmt.Sprintf(" Code: %s (%d)\n", key, value)
	// Message ID
	result += fmt.Sprintf(" Message ID: %d\n", request.MessageID)
	// Token
	result += fmt.Sprintf(" Token: %s\n", request.Token)
	// Options
	var option_count int = 0
	result += getOptionValue(coap.URIHost, request.Options(coap.URIHost), &option_count)
	result += getOptionValue(coap.URIPort, request.Options(coap.URIPort), &option_count)
	result += getOptionValue(coap.MaxAge, request.Options(coap.MaxAge), &option_count)
	result += getOptionValue(coap.ETag, request.Options(coap.ETag), &option_count)
	result += getOptionValue(coap.LocationPath, request.Options(coap.LocationPath), &option_count)
	result += getOptionValue(coap.LocationQuery, request.Options(coap.LocationQuery), &option_count)
	result += getOptionValue(coap.URIPath, request.Options(coap.URIPath), &option_count)
	result += getOptionValue(coap.Observe, request.Options(coap.Observe), &option_count)
	result += getOptionValue(coap.ContentFormat, request.Options(coap.ContentFormat), &option_count)
	result += getOptionValue(coap.IfMatch, request.Options(coap.IfMatch), &option_count)
	result += getOptionValue(coap.IfNoneMatch, request.Options(coap.IfNoneMatch), &option_count)
	result += getOptionValue(coap.ProxyURI, request.Options(coap.ProxyURI), &option_count)
	result += getOptionValue(coap.ProxyScheme, request.Options(coap.ProxyScheme), &option_count)
	result += getOptionValue(coap.Size1, request.Options(coap.Size1), &option_count)
	result += getOptionValue(coap.URIQuery, request.Options(coap.URIQuery), &option_count)
	// End of Options marker
	result += fmt.Sprintf(" End of options marker: %d\n", 255)
	// get content_format string
	var opt_content_format = request.Options(coap.ContentFormat)
	var content_format string = ""
	for _, v := range opt_content_format {
		content_format = GetContentFormatValue(v)
	}
	// Payload
	result += fmt.Sprintf(" Payload: Payload Content-Format: %s, Length: %d\n", content_format, len(request.Payload))
	result += fmt.Sprintf("  Payload Desc: %s\n", content_format)
	result += fmt.Sprintf("  JavaScript Object Notation: %s\n", content_format)
	result += fmt.Sprintf("  Line-based text data: %s\n", content_format)

	return result
}

func getTypeValue(t coap.COAPType) (int, string) {
	switch t {
	case coap.Confirmable:
		return 0, "Confirmable"
	case coap.NonConfirmable:
		return 1, "Non-Confirmable"
	case coap.Acknowledgement:
		return 2, "Acknowledgement"
	case coap.Reset:
		return 3, "Reset"
	default:
		return 0, ""
	}
}

func getMethodValue(t coap.COAPCode) (int, string) {
	switch t {
	case coap.GET:
		return 1, "GET"
	case coap.POST:
		return 2, "POST"
	case coap.PUT:
		return 3, "PUT"
	case coap.DELETE:
		return 4, "DELETE"
	default:
		return 0, ""
	}
}

func getOptionValue(t coap.OptionID, o []interface{}, op_count *int) string {
	var retValue string = ""

	// return nil if the interface equals to nil.
	if o == nil {
		return retValue
	}

	// Convert given OptionIDs to strings.
	for _, v := range o {
		var optName string = ""
		var optValue string = ""
		switch t {
		case coap.URIHost:
			optName = "Uri-Host"
			optValue = fmt.Sprintf("%s", v)
			break
		case coap.URIPort:
			optName = "Uri-Port"
			optValue = fmt.Sprintf("%d", v)
			break
		case coap.MaxAge:
			optName = "Max-Age"
			optValue = fmt.Sprintf("%d", v)
			break
		case coap.ETag:
			optName = "E-Tag"
			optValue = fmt.Sprintf("%s", v)
			break
		case coap.LocationPath:
			optName = "Location-Path"
			optValue = fmt.Sprintf("%s", v)
			break
		case coap.LocationQuery:
			optName = "Location-Query"
			optValue = fmt.Sprintf("%s", v)
			break
		case coap.URIPath:
			optName = "Uri-Path"
			optValue = fmt.Sprintf("%s", v)
			break
		case coap.Observe:
			optName = "Observe"
			optValue = fmt.Sprintf("%d", v)
			break
		case coap.ContentFormat:
			optName = "Content-Format"
			optValue = GetContentFormatValue(v)
			break
		case coap.IfMatch:
			optName = "If-Match"
			optValue = fmt.Sprintf("%s", v)
			break
		case coap.IfNoneMatch:
			optName = "If-None-Match"
			optValue = fmt.Sprintf("%s", v)
			break
		case coap.ProxyURI:
			optName = "Proxy-Uri"
			optValue = fmt.Sprintf("%s", v)
			break
		case coap.ProxyScheme:
			optName = "Proxy-Schema"
			optValue = fmt.Sprintf("%s", v)
			break
		case coap.Size1:
			optName = "Size1"
			optValue = fmt.Sprintf("%d", v)
			break
		case coap.URIQuery:
			optName = "Uri-Query"
			optValue = fmt.Sprintf("%s", v)
			break
		}

		// display the option names if the value is not nil.
		if optName != "" {
			*op_count++
			retValue += fmt.Sprintf(">Opt Name: #%d: %s: %s\n", *op_count, optName, optValue)
		}
	}

	return retValue
}

func GetContentFormatValue(v interface{}) string {
	switch v {
	case coap.TextPlain:
		return "text/plain"
	case coap.AppJSON:
		return "app/json"
	case coap.AppExi:
		return "app/exi"
	case coap.AppLinkFormat:
		return "app/linkformat"
	case coap.AppOctets:
		return "app/octets"
	case coap.AppXML:
		return "app/xml"
	case coap.AppCbor:
		return "app/cbor"
	default:
		return ""
	}
}

func decimalToBinary(dec_value int, digit int) string {
	var bin_value string = ""
	var new_digit int = 0
	// decimal to binary
	for dec_value != 0 {
		bin_value = strconv.Itoa(dec_value%2) + bin_value
		dec_value = dec_value / 2
		new_digit++
	}
	// zero-fill the upper bits
	for new_digit < digit {
		bin_value = "0" + bin_value
		new_digit++
	}

	return bin_value
}
