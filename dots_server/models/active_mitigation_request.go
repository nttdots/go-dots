package models

import (
	"time"
	"strconv"
	log "github.com/sirupsen/logrus"
)

/*
 * Struct for Active Mitigation Request.
 */
 type ActiveMitigationRequest struct {
	CustomerId       		int
	ClientIdentifier 		string
	ClientDomainIdentifier  string
	MitigationId     		int
	Lifetime				int
	LastModified            time.Time
}

var (
	acmMap map[string]ActiveMitigationRequest = make(map[string]ActiveMitigationRequest)
)

func GetActiveMitigationMap() map[string]ActiveMitigationRequest{
	return acmMap
}

func generateKey(customerId int, cuid string, mid int) (key string) {
	key = strconv.Itoa(customerId) + cuid + strconv.Itoa(mid)
	return
}

func AddActiveMitigationRequest(customerId int, cuid string, cdid string, mid int, lifetime int, modified time.Time) {
	key := generateKey(customerId, cuid, mid)
	acm, isPresent := acmMap[key]
	if isPresent {
		log.Debugf("Mitigation Request lifetime with id: %+v is updated: %+v", acm.MitigationId, lifetime)
		acm.LastModified = modified
		acm.Lifetime = lifetime
		acmMap[key] = acm
	} else {
		acm = ActiveMitigationRequest{
			customerId,
			cuid,
			cdid,
			mid,
			lifetime,
			modified,
		}
		acmMap[key] = acm
	}
}

func RemoveActiveMitigationRequest(customerId int, cuid string, mid int) {
	key := generateKey(customerId, cuid, mid)
	_, isPresent := acmMap[key]
	if isPresent {
		delete(acmMap, key)
	}
}

func CreateActiveMitigationRequest(customerId int, cuid string, cdid string, mid int, lifetime int) (acm ActiveMitigationRequest) {
    acm = ActiveMitigationRequest{
        customerId,
        cuid,
        cdid,
        mid,
        lifetime,
        time.Now(),
    }
    return
}