package main

import (
    "bufio"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"encoding/json"

	"github.com/nttdots/go-dots/dots_server/models"
	"github.com/nttdots/go-dots/dots_server/models/data"
	"github.com/nttdots/go-dots/libcoap"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/controllers"
	"github.com/nttdots/go-dots/dots_common/types/data"
	"github.com/nttdots/go-dots/dots_server/db_models"
	log "github.com/sirupsen/logrus"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
	data_controllers "github.com/nttdots/go-dots/dots_server/controllers/data"
)

type TableName string
const (
	MITIGATION_SCOPE         TableName = "mitigation_scope"
	SESSION_CONFIGURATION    TableName = "signal_session_configuration"
	PREFIX_ADDRESS_RANGE     TableName = "prefix"
	DATA_ACLS                TableName = "data_acls"
	TELEMETRY_PRE_MITIGATION TableName = "telemetry_pre_mitigation"

	TELEMETRY_TRAFFIC                 TableName = "telemetry_traffic"
	TELEMETRY_TOTAL_ATTACK_CONNECTION TableName = "telemetry_total_attack_connection"
	TELEMETRY_ATTACK_DETAIL           TableName = "telemetry_attack_detail"
	TELEMETRY_SOURCE_COUNT            TableName = "telemetry_source_count"
	TELEMETRY_TOP_TALKER              TableName = "telemetry_top_talker"
	TELEMETRY_SOURCE_PREFIX           TableName = "telemetry_source_prefix"
	TELEMETRY_SOURCE_PORT_RANGE       TableName = "telemetry_source_port_range"
	TELEMETRY_SOURCE_ICMP_TYPE_RANGE  TableName = "telemetry_source_icmp_type_range"

	URI_FILTERING_TRAFFIC                      TableName = "uri_filtering_traffic"
	URI_FILTERING_TRAFFIC_PROTOCOL             TableName = "uri_filtering_traffic_per_protocol"
	URI_FILTERING_TRAFFIC_PORT                 TableName = "uri_filtering_traffic_per_port"
	URI_FILTERING_TOTAL_ATTACK_CONNECTION      TableName = "uri_filtering_total_attack_connection"
	URI_FILTERING_TOTAL_ATTACK_CONNECTION_PORT TableName = "uri_filtering_total_attack_connection_port"
	URI_FILTERING_ATTACK_DETAIL                TableName = "uri_filtering_attack_detail"
	URI_FILTERING_SOURCE_COUNT                 TableName = "uri_filtering_source_count"
	URI_FILTERING_TOP_TALKER                   TableName = "uri_filtering_top_talker"
	URI_FILTERING_SOURCE_PREFIX                TableName = "uri_filtering_source_prefix"
	URI_FILTERING_SOURCE_PORT_RANGE            TableName = "uri_filtering_source_port_range"
	URI_FILTERING_ICMP_TYPE_RANGE              TableName = "uri_filtering_icmp_type_range"
)

/*
 * Listen for notification from DB
 */
func listenDB (context *libcoap.Context) {
	config := dots_config.GetServerSystemConfig()
	port := config.Network.DBNotificationPort
	listen, err := net.Listen("tcp4", ":" + strconv.Itoa(port))
	if err != nil {
		log.Errorf("[MySQL-Notification]:Socket listening on port %+v failed,%+v", port, err)
		os.Exit(1)
	}
	log.Debugf("[MySQL-Notification]:Begin listening on port: %+v", port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Errorf("[MySQL-Notification]:Error : %+v", err)
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
				cdid := mapData["cdid"].(string)
				mid, mErr := strconv.Atoi(mapData["mid"].(string))
				status, sErr := strconv.Atoi(mapData["status"].(string))
				if iErr != nil || mErr != nil || sErr != nil || cErr != nil {
					log.Errorf("[MySQL-Notification]:Failed to parse string to integer")
					return
				}
				handleNotifyMitigation(id, cid, cuid, cdid, mid, status, context, false)
			} else if mapData["table_trigger"].(string) == string(SESSION_CONFIGURATION) {

				// Notify session configuration changed to those clients who are observing this mitigation request
				log.Debug("[MySQL-Notification]: Send notification if obsevers exists")
				cid := mapData["cid"].(string)
				uriPath := messages.MessageTypes[messages.SESSION_CONFIGURATION].Path
				query := uriPath + "/customerId=" + cid
				context.EnableResourceDirty(query)
			} else if mapData["table_trigger"].(string) == string(PREFIX_ADDRESS_RANGE) {

				// re-check ip address range for each mitigation request, acl that are inactive
				log.Debug("[MySQL-Notification]: Re-check ip-address range for mitigations and acls")
				cid, err := strconv.Atoi(mapData["cid"].(string))
				if err != nil {
					log.Errorf("[MySQL-Notification]: Failed to parse string to integer")
					return
				}
				// Get customer from customer id
				customer, err := models.GetCustomer(cid)
				if err != nil {
					return
				}
				log.Printf("[MySQL-Notification]: new addressrange: %+v", customer.CustomerNetworkInformation.AddressRange)

				// Re-check ip address range for mitigations
				err = controllers.RecheckIpRangeForMitigations(&customer)
				if err != nil {
					log.Errorf("[MySQL-Notification]: Re-check ip range for mitigations failed. Error: %+v", err)
					return
				}

				// Re-check ip address range for acls
				err = data_controllers.RecheckIpRangeForAcls(&customer)
				if err != nil {
					log.Errorf("[MySQL-Notification]: Re-check ip range for acls failed. Error: %+v", err)
					return
				}
			} else if mapData["table_trigger"].(string) == string(DATA_ACLS) {
				handleNotifyACL(mapData["aclId"].(string), context)
			} else if mapData["table_trigger"].(string) == string(TELEMETRY_PRE_MITIGATION) {
				handleNotifyTelemetryPreMitigation(mapData, context)
			} else if mapData["table_trigger"].(string) == string(TELEMETRY_TRAFFIC) || mapData["table_trigger"].(string) == string(TELEMETRY_TOTAL_ATTACK_CONNECTION) ||
				mapData["table_trigger"].(string) == string(TELEMETRY_ATTACK_DETAIL) || mapData["table_trigger"].(string) == string(TELEMETRY_SOURCE_COUNT) ||
				mapData["table_trigger"].(string) == string(TELEMETRY_TOP_TALKER) || mapData["table_trigger"].(string) == string(TELEMETRY_SOURCE_PREFIX) ||
				mapData["table_trigger"].(string) == string(TELEMETRY_SOURCE_PORT_RANGE) || mapData["table_trigger"].(string) == string(TELEMETRY_SOURCE_ICMP_TYPE_RANGE) {
				handleNotifyTelemetryMitigationStatusUpdate(mapData, context)
			} else if mapData["table_trigger"].(string) == string(URI_FILTERING_TRAFFIC) || mapData["table_trigger"].(string) == string(URI_FILTERING_TRAFFIC_PROTOCOL) ||
				mapData["table_trigger"].(string) == string(URI_FILTERING_TRAFFIC_PORT) || mapData["table_trigger"].(string) == string(URI_FILTERING_TOTAL_ATTACK_CONNECTION) ||
				mapData["table_trigger"].(string) == string(URI_FILTERING_TOTAL_ATTACK_CONNECTION_PORT) || mapData["table_trigger"].(string) == string(URI_FILTERING_ATTACK_DETAIL) ||
				mapData["table_trigger"].(string) == string(URI_FILTERING_SOURCE_COUNT) || mapData["table_trigger"].(string) == string(URI_FILTERING_TOP_TALKER) ||
				mapData["table_trigger"].(string) == string(URI_FILTERING_SOURCE_PREFIX) || mapData["table_trigger"].(string) == string(URI_FILTERING_SOURCE_PORT_RANGE) ||
				mapData["table_trigger"].(string) == string(URI_FILTERING_ICMP_TYPE_RANGE) {
				// Handle notification telemetry pre-mitigation
				handleNotifyUriFilteringTelemetryPreMitigation(mapData, context)
			}
		default:
			log.Errorf("[MySQL-Notification]: Failed to receive data:%+v", err)
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

// Hanlde notify mitigation
func handleNotifyMitigation(id int64, cid int, cuid string, cdid string, mid int, status int, context *libcoap.Context, isNotifyTelemetryAttributes bool) {
	uriPath := messages.MessageTypes[messages.MITIGATION_REQUEST].Path
	var queryPaths []string
	uriQueryPaths := libcoap.GetUriFilterMitigationByValue(mid)
	if len(uriQueryPaths) > 0 {
		// Enable uri-query resource
		for _, v := range uriQueryPaths {
			uriQueryPath := strings.Split(v, "?")
			queries := strings.Split(uriQueryPath[1], "&")
			s, err := models.GetMitigationScope(0, "", 0, id, queries)
			if err != nil {
				log.Errorf("[MySQL-Notification]: Failed to get mitigation scope")
				return
			}
			if s != nil {
				queryPaths = append(queryPaths, v)
			}
		}
	} else {
		if cdid == "" {
			queryPaths = append(queryPaths, uriPath + "/cuid=" + cuid + "/mid=" + strconv.Itoa(mid))
		} else {
			queryPaths = append(queryPaths, uriPath + "/cdid="+ cdid + "/cuid=" + cuid + "/mid=" + strconv.Itoa(mid))
		}
	}

	// Check duplicate mitigation when PUT a new mitigation before delete an expired mitigation
	mids, err := models.GetMitigationIds(cid, cuid)
	if err != nil {
		log.Errorf("[MySQL-Notification]: Error: %+v", err)
		return
	}
	dup := isDuplicateMitigation(mids, mid)

	// Check observer resource and handle expired mitigation
	if dup && status == models.Terminated {
		// Skip notify, just delete the expired mitigation
		log.Debugf("[MySQL-Notification]: Skip Notification for this mitigation (mid=%+v, id=%+v) due to duplicate with another existing active mitigation", mid, id)
		controllers.DeleteMitigation(cid, cuid, mid, id)
	} else {
		// Set mitigation map for notification mitigation telemetry attributes
		if isNotifyTelemetryAttributes {
			isServerOriginatedTelemetry, err := getServerOriginatedTelemetry(cid, cuid)
			if err != nil {
				log.Error("[MySQL-Notification]: Failed to get server-originated-telemetry")
				return
			}
			// If the 'server-originated-telemetry' set to false, DOTS server will not notify
			if !isServerOriginatedTelemetry {
				log.Debug("[MySQL-Notification]: DOTS server will not notify to DOST client due to 'server-originated-telemetry' set to false")
				return
			}
		}
		// Notify status changed to those clients who are observing this mitigation request
		log.Debug("[MySQL-Notification]: Send notification if obsevers exists")
		for _, query := range queryPaths {
			resource := context.EnableResourceDirty(query)
			// If mitigation status was changed to Terminated and resource is not being observed => set resource status to removable
			var isObserved bool
			if resource != nil {
				isObserved = resource.IsObserved()
			} else {
				log.Warnf("[MySQL-Notification]: Not found any resource with query: %+v", query)
			}

			if status == models.Terminated && !isObserved {
				controllers.DeleteMitigation(cid, cuid, mid, id)
				// Keep resource when there is a duplication
				if !dup && resource != nil {
					resource.ToRemovableResource()
				}
				// If the last mitigation is expired, the server will remove the resource all
				if len(mids) == 1 && mids[0] == mid && resource != nil {
					uriPathSplit := strings.Split(resource.UriPath(), "/mid")
					resourceAll := context.GetResourceByQuery(&uriPathSplit[0])
					if resourceAll != nil {
						resourceAll.ToRemovableResource()
					}
				}
				break
			}
		}
	}
}

/*
 * Handle notify Acl when the acl's activation-type is updated
 */
func handleNotifyACL(aclIDString string, context *libcoap.Context) {
	aclID, err := strconv.ParseInt(aclIDString, 10, 64)
	if err != nil {
		log.Errorf("[MySQL-Notification]: Failed to parse string to integer")
		return
	}

	// Get acl by id
	acl, err := data_models.FindACLByID(aclID)
	if err != nil {
		log.Errorf("[MySQL-Notification]: Failed to get Acl from DB")
	}

	// Get data client by id
	client, err := data_models.FindClientByID(acl.ClientId)
	if err != nil {
		log.Errorf("[MySQL-Notification]: Failed to Get data_client")
		return
	}

	// Get control filtering by acl name
	controlFilteringList, err := models.GetControlFilteringByACLName(acl.Name)
	// If the acl's activation-type is not-type(the acl is deleted or expired) and the control filtering doesn't exist, remove acl from DB
	if len(controlFilteringList) == 0 && *acl.ACL.ActivationType == data_types.ActivationType_NotType {
		log.Debug("[MySQL-Notification]: Remove ACL from DB")
		err = models.RemoveACLByID(aclID, acl)
		if err != nil {
			log.Errorf("Failed to remove Acl from DB")
		}
	} else {
		uriPath := messages.MessageTypes[messages.MITIGATION_REQUEST].Path
		for _, ctrlFiltering := range controlFilteringList {
			// get mitigation scope by mitigation scope id
			mitigation, err := models.GetMitigationScope(0, "", 0, ctrlFiltering.MitigationScopeId, nil)
			if err != nil || mitigation == nil {
				log.Error("Failed to get mitigation scope")
				return
			}
			if mitigation.Customer.Id == client.CustomerId && mitigation.ClientIdentifier == client.Cuid {
				query := uriPath + "/cuid=" + mitigation.ClientIdentifier + "/mid=" + strconv.Itoa(mitigation.MitigationId)
				log.Debug("[MySQL-Notification]: Send notification if obsevers exists")
				context.EnableResourceDirty(query)
			}
			// If the acl's activation-type is not-type(the acl is deleted or expired), remove acl and control filtering
			if *acl.ACL.ActivationType == data_types.ActivationType_NotType {
				err = models.RemoveACLByID(aclID, acl)
				if err != nil {
					log.Errorf("Failed to remove Acl from DB")
				}

				err = models.RemoveControlFilteringByID(ctrlFiltering.Id, ctrlFiltering)
				if err != nil {
					log.Errorf("Failed to remove control filtering from DB")
				}
			}
		}
	}
}

// Handle notify telemetry pre-mitigation
func handleNotifyTelemetryPreMitigation(mapData map[string]interface{}, context *libcoap.Context) {
	id, err := strconv.Atoi(mapData["id"].(string))
	if err != nil {
		log.Errorf("[MySQL-Notification]: Failed to parse string to integer")
		return
	}
	// Get telemetry pre-mitigation
	preMitigation, err := models.GetTelemetryPreMitigationById(int64(id))
	if err != nil {
		log.Errorf("[MySQL-Notification]: Failed to Get telemetry pre-mitigation. Error: %+v", err)
		return
	}
	isServerOriginatedTelemetry, err := getServerOriginatedTelemetry(preMitigation.CustomerId, preMitigation.Cuid)
	if err != nil {
		log.Error("[MySQL-Notification]: Failed to get server-originated-telemetry")
		return
	}
	// If the 'server-originated-telemetry' set to false, DOTS server will not notify
	if !isServerOriginatedTelemetry {
		log.Debug("[MySQL-Notification]: DOTS server will not notify to DOST client due to 'server-originated-telemetry' set to false")
		return
	}
	uriPath := messages.MessageTypes[messages.TELEMETRY_PRE_MITIGATION_REQUEST].Path
	var query string
	if preMitigation.Cdid == "" {
		query = uriPath + "/cuid=" + preMitigation.Cuid + "/tmid=" + strconv.Itoa(preMitigation.Tmid)
	} else {
		query = uriPath + "/cuid=" + preMitigation.Cuid + "cdid=" + preMitigation.Cdid + "/tmid=" + strconv.Itoa(preMitigation.Tmid)
	}
	log.Debug("[MySQL-Notification]: Send notification if obsevers exists")
	context.EnableResourceDirty(query)
}

// Handle notification telemetry pre-mitigation aggregated by server
func handleNotifyUriFilteringTelemetryPreMitigation(mapData map[string]interface{}, context *libcoap.Context) {
	var preMitigation db_models.UriFilteringTelemetryPreMitigation
	if mapData["table_trigger"].(string) == string(URI_FILTERING_TRAFFIC) || mapData["table_trigger"].(string) == string(URI_FILTERING_TOTAL_ATTACK_CONNECTION) {
		prefixType := mapData["prefix_type"].(string)
		prefixTypeId, err := strconv.Atoi(mapData["prefix_type_id"].(string))
		if err != nil {
			log.Errorf("[MySQL-Notification]: Failed to parse string to integer")
			return
		}
		if prefixType == string(models.TARGET_PREFIX) {
			preMitigation, err = getTelemetryPreMitigation(int64(prefixTypeId), 0, 0)
			if err != nil {
				return
			}
		} else {
			preMitigation, err = getTelemetryPreMitigation(0, 0, int64(prefixTypeId))
			if err != nil {
				return
			}
		}
	} else if mapData["table_trigger"].(string) == string(URI_FILTERING_TRAFFIC_PROTOCOL) || mapData["table_trigger"].(string) == string(URI_FILTERING_TRAFFIC_PORT) ||
		mapData["table_trigger"].(string) == string(URI_FILTERING_TOTAL_ATTACK_CONNECTION_PORT) || mapData["table_trigger"].(string) == string(URI_FILTERING_ATTACK_DETAIL) {
		telePreMitigationId, err := strconv.Atoi(mapData["tele_pre_mitigation_id"].(string))
		if err != nil {
			log.Errorf("[MySQL-Notification]: Failed to parse string to integer")
			return
		}
		preMitigation, err = getTelemetryPreMitigation(int64(telePreMitigationId), 0, 0)
		if err != nil {
			return
		}
	} else if mapData["table_trigger"].(string) == string(URI_FILTERING_SOURCE_COUNT) || mapData["table_trigger"].(string) == string(URI_FILTERING_TOP_TALKER) {
		attackDetailId, err := strconv.Atoi(mapData["tele_attack_detail_id"].(string))
		if err != nil {
			log.Errorf("[MySQL-Notification]: Failed to parse string to integer")
			return
		}
		preMitigation, err = getTelemetryPreMitigation(0, int64(attackDetailId), 0)
		if err != nil {
			return
		}
	} else {
		topTalkerId, err := strconv.Atoi(mapData["tele_top_talker_id"].(string))
		if err != nil {
			log.Errorf("[MySQL-Notification]: Failed to parse string to integer")
			return
		}
		preMitigation, err = getTelemetryPreMitigation(0, 0, int64(topTalkerId))
		if err != nil {
			return
		}
	}
	isServerOriginatedTelemetry, err := getServerOriginatedTelemetry(preMitigation.CustomerId, preMitigation.Cuid)
	if err != nil {
		log.Error("[MySQL-Notification]: Failed to get server-originated-telemetry")
		return
	}
	// If the 'server-originated-telemetry' set to false, DOTS server will not notify
	if !isServerOriginatedTelemetry {
		log.Debug("[MySQL-Notification]: DOTS server will not notify to DOST client due to 'server-originated-telemetry' set to false")
		return
	}
	uriQueryPaths := libcoap.GetUriFilterByValue(preMitigation.Tmid)
	log.Debug("[MySQL-Notification]: Send notification if obsevers exists")
	if len(uriQueryPaths) > 0 {
		// Enable uri-query resource
		for _, v := range uriQueryPaths {
			uriQueryPath := strings.Split(v, "?")
			queries := strings.Split(uriQueryPath[1], "&")
			isExist, err := models.IsExistTelemetryPreMitigationValueByQueries(queries, preMitigation.CustomerId, preMitigation.Cuid, preMitigation)
			if err != nil {
				return
			}
			if isExist {
				context.EnableResourceDirty(v)
			}
		}
	} else {
		// Enable resource
		uriPath := messages.MessageTypes[messages.TELEMETRY_PRE_MITIGATION_REQUEST].Path
		var query string
		if preMitigation.Cdid == "" {
			query = uriPath + "/cuid=" + preMitigation.Cuid + "/tmid=" + strconv.Itoa(preMitigation.Tmid)
		} else {
			query = uriPath + "/cuid=" + preMitigation.Cuid + "cdid=" + preMitigation.Cdid + "/tmid=" + strconv.Itoa(preMitigation.Tmid)
		}
		context.EnableResourceDirty(query)
	}
}

// Handle notify telemetry mitigation status update
func handleNotifyTelemetryMitigationStatusUpdate(mapData map[string]interface{}, context *libcoap.Context) {
	var mitigation db_models.MitigationScope
	if mapData["table_trigger"].(string) == string(TELEMETRY_TRAFFIC) || mapData["table_trigger"].(string) == string(TELEMETRY_TOTAL_ATTACK_CONNECTION) {
		prefixType := mapData["prefix_type"].(string)
		prefixTypeId, err := strconv.Atoi(mapData["prefix_type_id"].(string))
		if err != nil {
			log.Errorf("[MySQL-Notification]: Failed to parse string to integer")
			return
		}
		if prefixType == string(models.TARGET_PREFIX) {
			mitigation, err = getMitigationFromTelemetryAttributes(int64(prefixTypeId), 0, 0)
			if err != nil {
				return
			}
		} else {
			mitigation, err = getMitigationFromTelemetryAttributes(0, 0, int64(prefixTypeId))
			if err != nil {
				return
			}
		}
	} else if mapData["table_trigger"].(string) == string(TELEMETRY_ATTACK_DETAIL) {
		mitigationId, err := strconv.Atoi(mapData["mitigation_scope_id"].(string))
		if err != nil {
			log.Errorf("[MySQL-Notification]: Failed to parse string to integer")
			return
		}
		mitigation, err = getMitigationFromTelemetryAttributes(int64(mitigationId), 0, 0)
		if err != nil {
			return
		}
	} else if mapData["table_trigger"].(string) == string(TELEMETRY_SOURCE_COUNT) || mapData["table_trigger"].(string) == string(TELEMETRY_TOP_TALKER) {
		attackDetailId, err := strconv.Atoi(mapData["tele_attack_detail_id"].(string))
		if err != nil {
			log.Errorf("[MySQL-Notification]: Failed to parse string to integer")
			return
		}
		mitigation, err = getMitigationFromTelemetryAttributes(0, int64(attackDetailId), 0)
		if err != nil {
			return
		}
	} else {
		topTalkerId, err := strconv.Atoi(mapData["tele_top_talker_id"].(string))
		if err != nil {
			log.Errorf("[MySQL-Notification]: Failed to parse string to integer")
			return
		}
		mitigation, err = getMitigationFromTelemetryAttributes(0, 0, int64(topTalkerId))
		if err != nil {
			return
		}
	}
	handleNotifyMitigation(mitigation.Id, mitigation.CustomerId, mitigation.ClientIdentifier, mitigation.ClientDomainIdentifier,
		mitigation.MitigationId, mitigation.Status, context, true)
}

// Get server originated telemetry from telemetry configuration
func getServerOriginatedTelemetry(customerId int, cuid string) (bool, error) {
	isServerOriginatedTelemetry := false
	// Get telemetry setup
	teleSetupList, err := models.GetTelemetrySetupByCuidAndSetupType(customerId, cuid, string(models.TELEMETRY_CONFIGURATION))
	if err != nil {
		log.Errorf("[MySQL-Notification]: Failed to Get telemetry setup configuration. Error: %+v", err)
		return isServerOriginatedTelemetry, err
	}
	// Get telemetry configuration
	if len(teleSetupList) > 0 {
		teleConfig, err := models.GetTelemetryConfiguration(teleSetupList[0].Id)
		if err != nil {
			log.Errorf("[MySQL-Notification]: Failed to Get telemetry configuration. Error: %+v", err)
			return isServerOriginatedTelemetry, err
		}
		isServerOriginatedTelemetry = teleConfig.ServerOriginatedTelemetry
	} else {
		defaultValue := dots_config.GetServerSystemConfig().DefaultTelemetryConfiguration
		isServerOriginatedTelemetry = defaultValue.ServerOriginatedTelemetry
	}
	return isServerOriginatedTelemetry, nil
}

// Get telemetry pre-mitigation
func getTelemetryPreMitigation(telemetryPreMitigationId int64, attackDetailId int64, topTalkerId int64) (db_models.UriFilteringTelemetryPreMitigation, error) {
	preMitigation := db_models.UriFilteringTelemetryPreMitigation{}
	if topTalkerId > 0 {
		// Get uri_filtering_top_talker
		topTalker, err := models.GetUriFilteringTopTalkerById(int64(topTalkerId))
		if err != nil {
			return preMitigation, err
		}
		attackDetailId = topTalker.TeleAttackDetailId
	}
	if attackDetailId > 0 {
		// Get uri_filtering_attack-detail
		attackDetail, err := models.GetUriFilteringAttackDetailById(attackDetailId)
		if err != nil {
			return preMitigation, err
		}
		telemetryPreMitigationId = attackDetail.TelePreMitigationId
	}
	if telemetryPreMitigationId > 0 {
		// Get uri_filtering_telemetry_pre_mitigation
		tmpPreMitigation, err := models.GetUriFilteringTelemetryPreMitigationById(telemetryPreMitigationId)
		if err != nil {
			return preMitigation, err
		}
		preMitigation = tmpPreMitigation
	}
	return preMitigation, nil
}

// Get mitigation from telemetry attributes
func getMitigationFromTelemetryAttributes(mitigationId int64, attackDetailId int64, topTalkerId int64) (db_models.MitigationScope, error) {
	mitigation := db_models.MitigationScope{}
	if topTalkerId > 0 {
		// Get uri_filtering_top_talker
		topTalker, err := models.GetTelemetryTopTalkerById(int64(topTalkerId))
		if err != nil {
			return mitigation, err
		}
		attackDetailId = topTalker.TeleAttackDetailId
	}
	if attackDetailId > 0 {
		// Get uri_filtering_attack-detail
		attackDetail, err := models.GetTelemetryAttackDetailById(attackDetailId)
		if err != nil {
			return mitigation, err
		}
		mitigationId = attackDetail.MitigationScopeId
	}
	if mitigationId > 0 {
		// Get uri_filtering_telemetry_pre_mitigation
		tmpMitigation, err := models.GetMitigationScopeById(mitigationId)
		if err != nil {
			return mitigation, err
		}
		mitigation = *tmpMitigation
	}
	return mitigation, nil
}