package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/nttdots/go-dots/dots_common/messages"
	log "github.com/sirupsen/logrus"
	common "github.com/nttdots/go-dots/dots_common"
)

var (
	requestName   string
	requestMethod string
	cuid          string
	cdid          string
	mid           string
	sid           string
	tsid          string
	tmid          string
	jsonFilePath  string
	socket        string
	observe       string
	ifMatch       string
	targetPrefix   string
	targetPort     string
	targetProtocol string
	targetFqdn     string
	targetUri      string
	aliasName      string
	sourcePrefix   string
	sourcePort     string
	sourceIcmpType string
)

/*
 * Default value in case If-Match option is not specified
 */
 var defaultIfMatchValue = "notIfMatch"

func init() {
	defaultValue := ""

	flag.StringVar(&requestName, "request", defaultValue, "Request Name")
	flag.StringVar(&requestMethod, "method", defaultValue, "Request Method(Get/Post/Put/Delete)")
	flag.StringVar(&cuid, "cuid", defaultValue, "Client Unique Identifier on Uri-Path. mandatory in Put/Get/Delete")
	flag.StringVar(&cdid, "cdid", defaultValue, "Client Domain IDentifier on Uri-Path (only used in debug). optional in Put/Get/Delete")
	flag.StringVar(&mid, "mid", defaultValue, "Identifier for the mitigation request on Uri-Path. mandatory in Put/Delete")
	flag.StringVar(&sid, "sid", defaultValue, "Session Identifier is an identifier for the DOTS signal channel session configuration data represented as an integer.")
	flag.StringVar(&tsid, "tsid", defaultValue, "Telemetry Setup Identifier is an identifier for the DOTS telemetry setup and configuration data represented as an integer.")
	flag.StringVar(&tmid, "tmid", defaultValue, "Telemetry Identifier is an identifier for the DOTS pre-mitigation telemetry data represented as an integer.")
	flag.StringVar(&jsonFilePath, "json", defaultValue, "Request Json file")
	flag.StringVar(&socket, "socket", common.DEFAULT_CLIENT_SOCKET_FILE, "dots client socket")
	flag.StringVar(&observe, "observe", defaultValue, "mitigation request observe")
	flag.StringVar(&ifMatch, "ifMatch", defaultIfMatchValue, "If-Match option")
	flag.StringVar(&targetPrefix, "targetPrefix", defaultValue, "target-prefix parameter")
	flag.StringVar(&targetPort, "targetPort", defaultValue, "target-port parameter")
	flag.StringVar(&targetProtocol, "targetProtocol", defaultValue, "target-protocol parameter")
	flag.StringVar(&targetFqdn, "targetFqdn", defaultValue, "target-fqdn parameter")
	flag.StringVar(&targetUri, "targetUri", defaultValue, "target-uri parameter")
	flag.StringVar(&aliasName, "aliasName", defaultValue, "alias-name parameter")
	flag.StringVar(&sourcePrefix, "sourcePrefix", defaultValue, "source-prefix parameter")
	flag.StringVar(&sourcePort, "sourcePort", defaultValue, "source-port parameter")
	flag.StringVar(&sourceIcmpType, "sourceIcmpType", defaultValue, "source-icmp-type parameter")
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

/*
 * socketExist is a function to check the socket to the server is already opened.
 */
func socketExist(socket string) (err error) {
	file, err := os.Stat(socket)
	if err != nil {
		fmt.Printf("dots_client_controller.socketExist -- NotExist: %s \n", err.Error())
		return
	}

	if (file.Mode() & os.ModeSocket) != os.ModeSocket {
		errMessage := fmt.Sprintf("%s is not a socket", socket)
		fmt.Printf("dots_client_controller.socketExist -- File error: %s  \n", errMessage)
		return errors.New(errMessage)
	}
	return
}

func main() {

	flag.Parse()

	if requestName == "" {
		fmt.Println("  -request option is required")
		os.Exit(1)
	}

	if requestMethod == "" {
		fmt.Println("  -method option is required")
		os.Exit(1)
	}

	common.SetUpLogger()
	log.Infof("method: %s, requestName: %s, cuid: %s, mid: %s", requestMethod, requestName, cuid, mid)

	err := socketExist(socket)
	if err != nil {
		os.Exit(1)
	}

	unixDomainSocketDial := func(network, addr string) (net.Conn, error) {
		return net.Dial("unix", socket)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: unixDomainSocketDial,
		},
	}

	u, err := url.Parse("http://dots_client/server")

	contentType := "application/json"
	u.Path = path.Join(u.Path, "server", requestName)
	if cuid != "" {
		// for mitigation request
		// if exists cdid set cdid first
		if cdid != "" {
			u.Path += "/cdid=" + cdid
		}
		// add cuid
		u.Path += "/cuid=" + cuid
		// add mid if exists
		if mid != "" {
			u.Path += "/mid=" + mid
		}
		// add tsid for telemetry setup request
		if tsid != "" {
			u.Path += "/tsid=" + tsid
		} else if tmid != "" {
			// add tmid for telemetry pre-mitigation request
			u.Path += "/tmid=" + tmid
		}
		if targetPrefix != "" || targetPort != "" || targetProtocol != "" || targetFqdn != "" || targetUri != "" || aliasName != "" ||
	        sourcePrefix != "" || sourcePort != "" || sourceIcmpType != "" {
			var queryPath string
			// target-prefix
			if targetPrefix != "" && queryPath != "" {
				queryPath += "&target-prefix="+targetPrefix
			} else if targetPrefix != "" {
				queryPath += "target-prefix="+targetPrefix
			}
			// target-port
			if targetPort != "" && queryPath != "" {
				queryPath += "&target-port="+targetPort
			} else if targetPort != "" {
				queryPath += "target-port="+targetPort
			}
			// target-protocol
			if targetProtocol != "" && queryPath != "" {
				queryPath += "&target-protocol="+targetProtocol
			} else if targetProtocol != "" {
				queryPath += "target-protocol="+targetProtocol
			}
			// target-fqdn
			if targetFqdn != "" && queryPath != "" {
				queryPath += "&target-fqdn="+targetFqdn
			} else if targetFqdn != "" {
				queryPath += "target-fqdn="+targetFqdn
			}
			// target-uri
			if targetUri != "" && queryPath != "" {
				queryPath += "&target-uri="+targetUri
			} else if targetUri != "" {
				queryPath += "target-uri="+targetUri
			}
			// alias-name
			if aliasName != "" && queryPath != "" {
				queryPath += "&alias-name="+aliasName
			} else if aliasName != "" {
				queryPath += "alias-name="+aliasName
			}
			// source-prefix
			if sourcePrefix != "" && queryPath != "" {
				queryPath += "&source-prefix="+sourcePrefix
			} else if sourcePrefix != "" {
				queryPath += "source-prefix="+sourcePrefix
			}
			// source-port
			if sourcePort != "" && queryPath != "" {
				queryPath += "&source-port="+sourcePort
			} else if sourcePort != "" {
				queryPath += "source-port="+sourcePort
			}
			// source-icmp-type
			if sourceIcmpType != "" && queryPath != "" {
				queryPath += "&source-icmp-type="+sourceIcmpType
			} else if sourceIcmpType != "" {
				queryPath += "source-icmp-type="+sourceIcmpType
			}
			u.Path += "?"
			u.Path += queryPath
		}
	} else if sid != "" {
		// for session configuration
		u.Path += "/sid=" + sid
	}
	var body io.Reader
	if jsonFilePath != "" {
		jsonData, err := readJsonFile(jsonFilePath)
		if err != nil {
			os.Exit(1)
		}
		body = bytes.NewBuffer(jsonData)
	}

	log.Debugf("NewRequest requestMethod=%+v, u=%+v, body=%+v", requestMethod, u, body)
	request, err := http.NewRequest(strings.ToUpper(requestMethod), u.String(), body)
	if err != nil {
		fmt.Printf("request message building error. %s", err.Error())
		os.Exit(1)
	}
	request.Header.Set("Content-Type", contentType)
	if observe != "" {
		request.Header.Set(string(messages.OBSERVE), observe)
	}
	if ifMatch != defaultIfMatchValue {
		request.Header.Set(string(messages.IFMATCH), ifMatch)
	}
	resp, err := client.Do(request)

	if err != nil {
		fmt.Printf("dots_client_controller.main -- %s : %s \n", requestMethod, err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()
	buff := new(bytes.Buffer)
	buff.ReadFrom(resp.Body)

	log.Debug("================***Response***================")
	log.Infof("dots_client_controller.main -- dots_client response code :%s\n", resp.Status)
	// The response code is not 2xx successfully
	if resp.StatusCode >= 300 {
		log.Error("dots_client return failed")
	}

	//dump received data.
	log.Infof("dots_client_controller.main -- dots_client response :%s\n", string(buff.String()))
}
