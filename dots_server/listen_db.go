package main

import (
    "bufio"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/nttdots/go-dots/dots_server/models"
	"github.com/nttdots/go-dots/libcoap"
	"github.com/nttdots/go-dots/dots_server/controllers"
	"github.com/nttdots/go-dots/dots_common/messages"
)

var port = 9999

/*
 * Listen for notification from DB
 */
func listenDB (context *libcoap.Context) {
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
			log.Debugf("[MySQL-Notification]: Received mitigation status changed notification from DB for :", data)
			if isTransportOver(data) {
				break ILOOP
			}
			// Notify status changed to those clients who are observing this mitigation request
			log.Debug("[MySQL-Notification]: Send notification if obsevers exists")
			uriPath := messages.MessageTypes[messages.MITIGATION_REQUEST].Path
			id, cuid, mid, status, query := context.NotifyOnce(data, uriPath)

			idValue, iErr := strconv.ParseInt(id, 10, 64)
			midValue, mErr := strconv.Atoi(mid)
			statusValue, sErr := strconv.Atoi(status)
			if iErr != nil || mErr != nil || sErr != nil {
				log.Debugf("[MySQL-Notification]:Failed to parse string to integer")
				return
			}
			// If mitigation status was changed to 6: (attack mitigation is now terminated), delete this mitigation after notifying
			if statusValue == models.Terminated {
				log.Debugf("[MySQL-Notification]: Mitigation was terminated. Delete corresponding sub-resource and mitigation request.", models.Terminated)
				context.DeleteResourceByQuery(query)
				controllers.DeleteMitigation(0, cuid, midValue, idValue)
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
