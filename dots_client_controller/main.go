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

	log "github.com/sirupsen/logrus"
	common "github.com/nttdots/go-dots/dots_common"
)

var (
	requestName   string
	requestMethod string
	cuid          string
	mid           string
	sid           string
	jsonFilePath  string
	socket        string
)

func init() {
	defaultValue := ""

	flag.StringVar(&requestName, "request", defaultValue, "Request Name")
	flag.StringVar(&requestMethod, "method", defaultValue, "Request Method(Get/Post/Put/Delete)")
	flag.StringVar(&cuid, "cuid", defaultValue, "Client Unique Identifier on Uri-Path. mandatory in Put/Get/Delete")
	// TODO: cdid for gateway
	flag.StringVar(&mid, "mid", defaultValue, "Identifier for the mitigation request on Uri-Path. mandatory in Put/Delete")
	flag.StringVar(&sid, "sid", defaultValue, "Session Identifier is an identifier for the DOTS signal channel session configuration data represented as an integer.")
	flag.StringVar(&jsonFilePath, "json", defaultValue, "Request Json file")
	flag.StringVar(&socket, "socket", common.DEFAULT_CLIENT_SOCKET_FILE, "dots client socket")
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
	if cuid != "" {
		if mid == "" {
			u.Path = path.Join(u.Path, "server", requestName) + "/cuid=" + cuid
		} else {
			u.Path = path.Join(u.Path, "server", requestName) + "/cuid=" + cuid + "/mid=" + mid
		}
	} else if sid != "" {
		u.Path = path.Join(u.Path, "server", requestName) + "/sid=" + sid
	} else {
		u.Path = path.Join(u.Path, "server", requestName)
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
	resp, err := client.Do(request)

	if err != nil {
		fmt.Printf("dots_client_controller.main -- %s : %s \n", requestMethod, err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()
	buff := new(bytes.Buffer)
	buff.ReadFrom(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("dots_client return error")
	}

	//dump received data.
	log.Infof("dots_client_controller.main -- dots_client response :%s\n", string(buff.String()))
}
