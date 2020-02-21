package messages

import (
	"fmt"
	"reflect"
	"io/ioutil"
	"os"
	"encoding/json"
	"github.com/ugorji/go/codec"
	"github.com/nttdots/go-dots/libcoap"
	"github.com/nttdots/go-dots/dots_common"
)


var hbRequestPath = ".well-known/dots/hb"

/* 
 * Create heartbeat json data
 */
func CreateHeartBeatJsonData(jsonFileName string, hbValue bool) {
	data := HeartBeatRequest{
		HeartBeat: HeartBeat{
			PeerHbStatus: &hbValue,
		},
	}
	file, _ := json.MarshalIndent(data, "", " ")
	
	_ = ioutil.WriteFile(jsonFileName, file, 0644)
}

// validate heartbeat mechanism
func ValidateHeartBeatMechanism(request *libcoap.Pdu) (body *HeartBeatRequest, errMsg string) {
    bodyUnmarshal, err := UnmarshalCbor(request, reflect.TypeOf(HeartBeatRequest{}))
    if err != nil {
        errMsg := fmt.Sprintf("Unmarshal Cbor failed. Error: %+v", err)
        return nil, errMsg
    }
    // peer-hb-status is mandatory attribute
    body = bodyUnmarshal.(*HeartBeatRequest)
    if body.HeartBeat.PeerHbStatus == nil {
        errMsg := fmt.Sprintf("Missing the mandatory attribute")
        return body, errMsg
    }
    cdid, cuid, mid, err := ParseURIPath(request.Path())
    if cuid != "" || cdid != "" || mid != nil {
        errMsg := fmt.Sprintf("The DOTS heartbeat MUST NOT have a 'cuid', 'cdid,' or 'mid' Uri-Path.")
        return body, errMsg
    }
    return body, ""
}

// create new heartbeat message
func NewHeartBeatMessage(session libcoap.Session,jsonFileName string, hbValue bool) (*libcoap.Pdu, error) {
    CreateHeartBeatJsonData(jsonFileName, hbValue)
    jsonData, err := readJsonFile(jsonFileName)
    if err != nil {
        return nil, err
    }
    data, err := dumpCbor(jsonData)
    if err != nil {
        return nil, err
    }
    pdu := &libcoap.Pdu{}
    pdu.Type = libcoap.TypeNon
    pdu.Code = libcoap.RequestPut
    pdu.Data = data
    pdu.Token = dots_common.RandStringBytes(8)
	pdu.MessageID = session.NewMessageID()
	pdu.SetOption(libcoap.OptionContentFormat, uint16(libcoap.AppCbor))
    pdu.SetPathString(hbRequestPath)
    return pdu, nil
}

/*
 * convert this Request into the Cbor format.
 */
 func dumpCbor(jsonData []byte) ([]byte, error) {
    m := HeartBeatRequest{}
	err := json.Unmarshal(jsonData, &m)
	if err != nil {
		return nil, fmt.Errorf("Can't Convert Json to Message Object: %v\n", err)

    }
    
	var buf []byte
	e := codec.NewEncoderBytes(&buf, dots_common.NewCborHandle())
    err = e.Encode(m)
	if err != nil {
        return nil, err
	}
	return buf, nil
}

/*
 * readJsonFile is a function that loads a JSON file and returns []byte.
 */
 func readJsonFile(path string) (jsonData []byte, err error) {
	jsonData = nil
	_, err = os.Stat(path)
	if err != nil {
		if len(path) == 0 {
			//log.log.With Println("Need request json file")
		} else {
			fmt.Printf("Not Found %s \n", path)
		}
		return
	}

	jsonData, err = ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("dots_client_controller.readJsonFile -- File error: %v\n", err)
		return
	}
	return
}