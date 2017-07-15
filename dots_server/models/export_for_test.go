package models

import (
	"github.com/nttdots/go-dots/dots_server/db_models"
)

func BlockerSelection(d *LoadBaseBlockerSelection, scope *MitigationScope) (b Blocker, err error) {
	return d.selection(scope)
}

func NewLoadBaseBlockerSelection() *LoadBaseBlockerSelection {
	return &LoadBaseBlockerSelection{}
}

func NewBlockerBase(id int64, capacity, load int, connector DeviceConnector, loginInfo *LoginProfile) BlockerBase {
	return BlockerBase{
		id:        id,
		capacity:  capacity,
		load:      load,
		connector: connector,
		loginInfo: loginInfo,
	}
}

func ToBlocker(blocker db_models.Blocker, loginProfile db_models.LoginProfile, blockerParameters []db_models.BlockerParameter) Blocker {

	var b Blocker

	base := BlockerBase{
		id:        blocker.Id,
		capacity:  blocker.Capacity,
		load:      blocker.Load,
		loginInfo: new(LoginProfile),
	}
	base.loginInfo.Load(loginProfile)

	paramMap := BlockerParametersToMap(blockerParameters)

	switch blocker.Type {
	case BLOCKER_TYPE_GoBGP_RTBH:
		b = NewGoBgpRtbhReceiver(base, paramMap)
	}

	return b
}
