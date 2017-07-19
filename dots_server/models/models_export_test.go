package models

import (
	"time"

	"github.com/nttdots/go-dots/dots_server/db_models"
)

func ToProtectionParameters(obj Protection) []db_models.ProtectionParameter {
	return toProtectionParameters(obj, obj.Id())
}

func NewProtectionBase(id int64, mitigationId int, isEnabled bool, startdAt, finishedAt, recordTime time.Time,
	blocker Blocker, forwardedDataInfo, blockedDataInfo *ProtectionStatus) ProtectionBase {

	return ProtectionBase{
		id,
		mitigationId,
		blocker,
		isEnabled,
		startdAt,
		finishedAt,
		recordTime,
		forwardedDataInfo,
		blockedDataInfo,
	}
}
