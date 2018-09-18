package main

import (
    "bufio"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/nttdots/go-dots/dots_server/models"
	"github.com/nttdots/go-dots/libcoap"
	"github.com/nttdots/go-dots/dots_common/messages"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
)

type TableName string
const (
	MITIGATION_SCOPE      TableName = "mitigation_scope"
	SESSION_CONFIGURATION TableName = "signal_session_configuration"
)

/*
 * Listen for notification from DB
 */
func listenDB (context *libcoap.Context) {
	config := dots_config.GetServerSystemConfig()
	port := config.Network.DBNotificationPort
	listen, err := net.Listen("tcp4", ":" + strconv.Itoa(port))
	if err != nil {
		log.Debugf("[MySQL-Notification]:Socket listening on port %+v failed,%+v", port, err)
		os.Exit(1)
	}
	log.Debugf("[MySQL-Notification]:Begin listening on port: %+v", port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Debugf("[MySQL-Notification]:Error : %+v", err)
			continue
		}
		go handler(conn, context)
	}

}

/*
 * Handle notifcation from DB
 */
func handler(conn net.Conn, context *libcoap.Context) {

	defer conn.Close()

	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
	)

ILOOP:
	for {
		n, err := r.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			break ILOOP
		case nil:
			log.Debugf("[MySQL-Notification]: Received data changed notification from DB trigger for: %+v", data)
			if isTransportOver(data) {
				break ILOOP
			}

			// Parse json data from notification to parameters
			var mapData map[string]interface{}
			err := json.Unmarshal([]byte (data), &mapData)
			if err != nil {
				log.Errorf("[MySQL-Notification]: Failed to encode json message to map data.")
				return
			}
			if mapData["table_trigger"].(string) == string(MITIGATION_SCOPE) {
				id, iErr := strconv.ParseInt(mapData["id"].(string), 10, 64)
				cid, cErr := strconv.Atoi(mapData["cid"].(string))
				cuid := mapData["cuid"].(string)
				mid, mErr := strconv.Atoi(mapData["mid"].(string))
				status, sErr := strconv.Atoi(mapData["status"].(string))
				if iErr != nil || mErr != nil || sErr != nil || cErr != nil {
					log.Debugf("[MySQL-Notification]:Failed to parse string to integer")
					return
				}
				uriPath := messages.MessageTypes[messages.MITIGATION_REQUEST].Path
				query := uriPath + "/cuid=" + cuid + "/mid=" + strconv.Itoa(mid)

				// Check duplicate mitigation when PUT a new mitigation before delete an expired mitigation
				mids, err := models.GetMitigationIds(cid, cuid)
				if err != nil {
					log.Errorf("[MySQL-Notification]: Error: %+v", err)
					return
				}
				dup := isDuplicateMitigation(mids, mid)
				if dup && status == models.Terminated {
					// Skip notify, just delete the expired mitigation
					log.Debugf("[MySQL-Notification]: Skip Notification for this mitigation (mid=%+v, id=%+v) due to duplicate with another existing active mitigation", mid, id)
				} else {
					// Notify status changed to those clients who are observing this mitigation request
					log.Debug("[MySQL-Notification]: Send notification if obsevers exists")
					context.NotifyOnce(query)
				}

				// If mitigation status was changed to 6: (attack mitigation is now terminated), delete this mitigation after notifying
				if status == models.Terminated {
					log.Debugf("[MySQL-Notification]: Mitigation was terminated. Delete corresponding sub-resource and mitigation request.")
					err = models.DeleteMitigationScope(cid, cuid, mid, id)
					if err != nil {
						log.Errorf("[MySQL-Notification]: Delete mitigation scope error: %+v", err)
						return
					}
					// Keep resource when there is a duplication
					if !dup {
						context.DeleteResourceByQuery(query)
					}
				}
			} else if mapData["table_trigger"].(string) == string(SESSION_CONFIGURATION) {

				// Notify session configuration changed to those clients who are observing this mitigation request
				log.Debug("[MySQL-Notification]: Send notification if obsevers exists")
				cid := mapData["cid"].(string)
				uriPath := messages.MessageTypes[messages.SESSION_CONFIGURATION].Path
				query := uriPath + "/customerId=" + cid
				context.NotifyOnce(query)
			}
		default:
			log.Debugf("[MySQL-Notification]: Failed to receive data:%+v", err)
			return
		}
	}
}

/*
 * Check if nofified data has been transported completely
 */
func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\r\n\r\n")
	return
}

/*
 * Check if there is a duplication mitigation in DB
 */
func isDuplicateMitigation(mids []int, mid int) bool {
	count := 0
	for _, id := range mids {
		if id == mid {
			count++
		}
	}
	if count > 1 {
		return true
	} else {
		return false
	}
}