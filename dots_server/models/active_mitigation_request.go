package models

import (
	"time"
)

/*
 * Struct for Active Mitigation Request.
 */
 type ActiveMitigationRequest struct {
	MitigationScopeId       int64
	Lifetime				int
	LastModified            time.Time
}

var acmMap map[int64]ActiveMitigationRequest = make(map[int64]ActiveMitigationRequest)

func GetActiveMitigationMap() map[int64]ActiveMitigationRequest{
	return acmMap
}

func AddActiveMitigationRequest(id int64, lifetime int, modified time.Time) {
	acm, isPresent := acmMap[id]
	if isPresent {
		acm.LastModified = modified
		acm.Lifetime = lifetime
		acmMap[id] = acm
	} else {
		acm = ActiveMitigationRequest{
			id,
			lifetime,
			modified,
		}
		acmMap[id] = acm
	}
}

func RemoveActiveMitigationRequest(id int64) {
	_, isPresent := acmMap[id]
	if isPresent {
		delete(acmMap, id)
	}
}

func CreateActiveMitigationRequest(id int64, lifetime int) (acm ActiveMitigationRequest) {
    acm = ActiveMitigationRequest{
        id,
        lifetime,
        time.Now(),
    }
    return
}