package main

import (
	"fmt"
	"strconv"

	"github.com/nttdots/go-dots/libcoap"
)

/*
 * Output formatted CoAP messages to the stdout.
 */
func CoapHeaderDisplay(request *libcoap.Pdu) string {
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
	result += getOptionValue(libcoap.OptionUriHost, request.OptionValues(libcoap.OptionUriHost), &option_count)
	result += getOptionValue(libcoap.OptionUriPort, request.OptionValues(libcoap.OptionUriPort), &option_count)
	result += getOptionValue(libcoap.OptionMaxage, request.OptionValues(libcoap.OptionMaxage), &option_count)
	result += getOptionValue(libcoap.OptionEtag, request.OptionValues(libcoap.OptionEtag), &option_count)
	result += getOptionValue(libcoap.OptionLocationPath, request.OptionValues(libcoap.OptionLocationPath), &option_count)
	result += getOptionValue(libcoap.OptionLocationQuery, request.OptionValues(libcoap.OptionLocationQuery), &option_count)
	result += getOptionValue(libcoap.OptionUriPath, request.OptionValues(libcoap.OptionUriPath), &option_count)
	result += getOptionValue(libcoap.OptionObserve, request.OptionValues(libcoap.OptionObserve), &option_count)
	result += getOptionValue(libcoap.OptionContentFormat, request.OptionValues(libcoap.OptionContentFormat), &option_count)
	result += getOptionValue(libcoap.OptionIfMatch, request.OptionValues(libcoap.OptionIfMatch), &option_count)
	result += getOptionValue(libcoap.OptionIfNoneMatch, request.OptionValues(libcoap.OptionIfNoneMatch), &option_count)
	result += getOptionValue(libcoap.OptionProxyUri, request.OptionValues(libcoap.OptionProxyUri), &option_count)
	result += getOptionValue(libcoap.OptionProxyScheme, request.OptionValues(libcoap.OptionProxyScheme), &option_count)
	result += getOptionValue(libcoap.OptionSize1, request.OptionValues(libcoap.OptionSize1), &option_count)
	result += getOptionValue(libcoap.OptionUriQuery, request.OptionValues(libcoap.OptionUriQuery), &option_count)
	// End of Options marker
	result += fmt.Sprintf(" End of options marker: %d\n", 255)
	// get content_format string
	var opt_content_format = request.OptionValues(libcoap.OptionContentFormat)
	var content_format string = ""
	for _, v := range opt_content_format {
		content_format = GetContentFormatValue(v)
	}
	// Payload
	result += fmt.Sprintf(" Payload: Payload Content-Format: %s, Length: %d\n", content_format, len(request.Data))
	result += fmt.Sprintf("  Payload Desc: %s\n", content_format)
	result += fmt.Sprintf("  JavaScript Object Notation: %s\n", content_format)
	result += fmt.Sprintf("  Line-based text data: %s\n", content_format)

	return result
}

func getTypeValue(t libcoap.Type) (int, string) {
	switch t {
	case libcoap.TypeCon:
		return 0, "Confirmable"
	case libcoap.TypeNon:
		return 1, "Non-Confirmable"
	case libcoap.TypeAck:
		return 2, "Acknowledgement"
	case libcoap.TypeRst:
		return 3, "Reset"
	default:
		return 0, ""
	}
}

func getMethodValue(t libcoap.Code) (int, string) {
	switch t {
	case libcoap.RequestGet:
		return 1, "GET"
	case libcoap.RequestPost:
		return 2, "POST"
	case libcoap.RequestPut:
		return 3, "PUT"
	case libcoap.RequestDelete:
		return 4, "DELETE"
	default:
		return 0, ""
	}
}

func getOptionValue(t libcoap.OptionKey, o []interface{}, op_count *int) string {
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
		case libcoap.OptionUriHost:
			optName = "Uri-Host"
			optValue = fmt.Sprintf("%s", v)
			break
		case libcoap.OptionUriPort:
			optName = "Uri-Port"
			optValue = fmt.Sprintf("%d", v)
			break
		case libcoap.OptionMaxage:
			optName = "Max-Age"
			optValue = fmt.Sprintf("%d", v)
			break
		case libcoap.OptionEtag:
			optName = "E-Tag"
			optValue = fmt.Sprintf("%s", v)
			break
		case libcoap.OptionLocationPath:
			optName = "Location-Path"
			optValue = fmt.Sprintf("%s", v)
			break
		case libcoap.OptionLocationQuery:
			optName = "Location-Query"
			optValue = fmt.Sprintf("%s", v)
			break
		case libcoap.OptionUriPath:
			optName = "Uri-Path"
			optValue = fmt.Sprintf("%s", v)
			break
		case libcoap.OptionObserve:
			optName = "Observe"
			optValue = fmt.Sprintf("%d", v)
			break
		case libcoap.OptionContentFormat:
			optName = "Content-Format"
			optValue = GetContentFormatValue(v)
			break
		case libcoap.OptionIfMatch:
			optName = "If-Match"
			optValue = fmt.Sprintf("%s", v)
			break
		case libcoap.OptionIfNoneMatch:
			optName = "If-None-Match"
			optValue = fmt.Sprintf("%s", v)
			break
		case libcoap.OptionProxyUri:
			optName = "Proxy-Uri"
			optValue = fmt.Sprintf("%s", v)
			break
		case libcoap.OptionProxyScheme:
			optName = "Proxy-Schema"
			optValue = fmt.Sprintf("%s", v)
			break
		case libcoap.OptionSize1:
			optName = "Size1"
			optValue = fmt.Sprintf("%d", v)
			break
		case libcoap.OptionUriQuery:
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
	case libcoap.TextPlain:
		return "text/plain"
	case libcoap.AppJSON:
		return "app/json"
	case libcoap.AppExi:
		return "app/exi"
	case libcoap.AppLinkFormat:
		return "app/linkformat"
	case libcoap.AppOctets:
		return "app/octets"
	case libcoap.AppXML:
		return "app/xml"
	case libcoap.AppCbor:
		return "app/cbor"
	case libcoap.AppDotsCbor:
		return "app/dots+cbor"
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
