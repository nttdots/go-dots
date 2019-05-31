package models

import (
	"time"
)

/*
 * Struct for Fresh Session Configuration.
 */
 type ActiveSessionConfiguration struct {
	SessionId       int
	MaxAge	        uint
	LastRefresh     time.Time
}

var ascMap map[int]ActiveSessionConfiguration = make(map[int]ActiveSessionConfiguration)

func GetFreshSessionMap() map[int]ActiveSessionConfiguration {
	return ascMap
}

/*
 * Remove by customer id an active session configuration out of map when it is deleted
 *  or expired max-age and reset to default value
 */
func RemoveActiveSessionConfiguration(customerId int) {
	_, isPresent := ascMap[customerId]
	if isPresent {
		delete(ascMap, customerId)
	}
}

/*
 * Refresh or update an active session configuration when client request Get with sid
 *  or request Put with sid to update session configuration
 */
func RefreshActiveSessionConfiguration(customerId int, sid int, maxAge uint) {
	asc, isPresent := ascMap[customerId]
	if isPresent {
		asc.LastRefresh = time.Now()
		asc.MaxAge = maxAge
		asc.SessionId = sid
		ascMap[customerId] = asc
	} else {
		asc = ActiveSessionConfiguration {
			sid,
			maxAge,
			time.Now(),
		}
		ascMap[customerId] = asc
	}
}