package models

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/nttdots/go-dots/dots_server/db_models"
	log "github.com/sirupsen/logrus"
)

/*
 * Stores a Blocker object to the database
 *
 * parameter:
 *  blocker Blocker
 * return:
 *  newBlocker db_models.Blocker
 *  err error
 */
func CreateBlocker(blocker Blocker) (newBlocker db_models.Blocker, err error) {

	var blockerParameters []db_models.BlockerParameter

	engine, err := ConnectDB()
	if err != nil {
		return
	}

	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return
	}

	newBlocker = db_models.Blocker{
		BlockerType: string(blocker.Type()),
		Load:        blocker.Load(),
		Capacity:    blocker.Capacity(),
	}
	_, err = session.Insert(&newBlocker)
	if err != nil {
		goto Rollback
	}

	if blocker.LoginProfile() != nil {
		err = createBlockerLoginProfile(session, *blocker.LoginProfile(), newBlocker.Id)
		if err != nil {
			goto Rollback
		}
	}
	blockerParameters = toBlockerParameters(blocker, newBlocker.Id)

	if len(blockerParameters) > 0 {
		_, err := session.InsertMulti(blockerParameters)
		if err != nil {
			goto Rollback
		}
	}

	session.Commit()
	return
Rollback:
	session.Rollback()
	return
}

/*
 * Stores a LoginInfo object related to a Blocker to the database.
 * Parameter:
 *  session session information
 *  login_info LoginProfile
 *  blocker_id  Blocker ID to which this loginInfo is related.
 * return:
 *  err error
 */
func createBlockerLoginProfile(session *xorm.Session, loginInfo LoginProfile, blockerId int64) (err error) {
	// registered LoginProfile
	newLoginProfile := db_models.LoginProfile{
		LoginMethod: loginInfo.LoginMethod,
		LoginName:   loginInfo.LoginName,
		Password:    loginInfo.Password,
		BlockerId:   blockerId,
	}

	_, err = session.Insert(&newLoginProfile)
	if err != nil {
		session.Rollback()
		log.Infof("LoginProfile insert err: %s", err)
		return
	}

	return
}

/*
 * Stores a Blocker parameters related to a Blocker to the database.
 * Parameter:
 *  session session information
 *  login_info LoginProfile
 *  blocker_id  Blocker ID to which this loginInfo is related.
 * return:
 *  err error
 */
func createBlockerParameters(session *xorm.Session, blockerParameters []db_models.BlockerParameter, blockerId int64) (err error) {
	// registered BlockerParameters
	for _, param := range blockerParameters {
		param.BlockerId = blockerId
		_, err = session.Insert(param)
		if err != nil {
			session.Rollback()
			log.Infof("LoginProfile insert err: %s", err)
			return
		}
	}

	return
}

/*
 * Obtain all the blocker information stored in the database.
 *
 * return:
 *  blockers a list of blockers
 *  error error
 */
func GetBlockers() (blockers []db_models.Blocker, err error) {
	// default value setting
	blockers = []db_models.Blocker{}

	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// Get blocker table data
	err = engine.Find(&blockers)
	if err != nil {
		return
	}

	return
}

/*
 * Obtain the blocker with the least load.
 * If the load values are same, capacity information are used to break the tie.
 */
func GetLowestLoadBlocker(blockerType string) (blocker db_models.Blocker, err error) {
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}


	// ok, err := engine.Where("`load` < `capacity`").OrderBy("`load`, `capacity` desc").Get(&blocker)
	// get first blocker
	ok, err := engine.Where("blocker_type = ? AND `load` < `capacity`", blockerType).OrderBy("`load`, `capacity` desc").Get(&blocker)
	if err != nil {
		return
	}
	if !ok {
		blocker.Id = 0
	}
	return blocker, nil
}

/*
 * Obtain the designated BlockerParameters by the Blocker ID.
 *
 * parameter:
 *  blockerId  Blocker ID
 * return:
 *  BlockerParameter BlockerParameter
 *  error error
 */
func GetBlockerParameters(blockerId int64) (blockerParameters []db_models.BlockerParameter, err error) {
	// default value setting
	blockerParameters = []db_models.BlockerParameter{}

	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// Get blocker table data
	err = engine.Where("blocker_id = ?", blockerId).Find(&blockerParameters)
	if err != nil {
		return
	}

	return

}

/*
 * Obtain the designated BlockerConfigurationParameters by the Blocker COnfiguration ID.
 *
 * parameter:
 *  blockerConfigurationId  the Blocker Configuration ID
 * return:
 *  blockerConfigurationParameters the Blocker Configuration Parameter
 *  error error
 */
func GetBlockerConfigurationParameters(blockerConfigurationId int64) (blockerConfigurationParameters []db_models.BlockerConfigurationParameter, err error) {
	// default value setting
	blockerConfigurationParameters = []db_models.BlockerConfigurationParameter{}

	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// Get blocker table data
	err = engine.Where("blocker_configuration_id = ?", blockerConfigurationId).Find(&blockerConfigurationParameters)
	if err != nil {
		return
	}

	return
}

/*
 * Obtain the designated LoginProfile by the Blocker ID.
 *
 * parameter:
 *  blockerId  Blocker ID
 * return:
 *  LoginProfile LoginProfile
 *  error error
 */
func GetLoginProfile(blockerId int64) (loginProfile db_models.LoginProfile, err error) {
	// default value setting
	loginProfile = db_models.LoginProfile{}

	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// Get blocker table data
	ok, err := engine.Where("blocker_id = ?", blockerId).Get(&loginProfile)
	if err != nil {
		return
	}
	if !ok {
		// No data found
		log.WithField("blocker_id", blockerId).Warn("login_profile data not exist.")
		return
	}

	return

}

/*
 * Delete the Blocker by the Blocker ID.
 *
 * parameter:
 *  blockerId  Blocker ID
 * return:
 *  error error
 */
func DeleteBlockerById(blockerId int64) (err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return
	}

	// Delete blocker_parameters table data
	_, err = session.Delete(db_models.BlockerParameter{BlockerId: blockerId})
	if err != nil {
		log.Errorf("delete blockerParameters error: %s", err)
		goto Rollback
	}

	// Delete login_profile table data
	_, err = session.Delete(db_models.LoginProfile{BlockerId: blockerId})
	if err != nil {
		log.Errorf("delete loginProfile error: %s", err)
		goto Rollback
	}

	// Delete blocker table data
	_, err = session.Id(blockerId).Delete(db_models.Blocker{})
	if err != nil {
		log.Errorf("delete blocker error: %s", err)
		goto Rollback
	}

	session.Commit()
	return
Rollback:
	session.Rollback()
	return
}

func GetBlockerById(blockerId int64) (blocker db_models.Blocker, err error) {
	engine, err := ConnectDB()
	if err != nil {
		return
	}

	ok, err := engine.Id(blockerId).Get(&blocker)
	if !ok {
		return db_models.Blocker{}, nil
	} else {
		return
	}
}

/*
 * Update the load value of a Blocker
 *
 * parameter:
 *  blockerId id
 *  diff diff of the load value
 * return:
 *  error error
 */
func UpdateBlockerLoad(blockerId int64, diff int) (err error) {

	engine, err := ConnectDB()
	if err != nil {
		return
	}
	session := engine.NewSession()
	err = session.Begin()
	if err != nil {
		return
	}
	blocker := db_models.Blocker{}

	ok, err := session.Id(blockerId).Get(&blocker)
	if err != nil {
		goto Rollback
	}

	if ok {
		blocker.Load += diff
		_, err = session.Id(blockerId).Cols("Load").Update(&blocker)
		if err != nil {
			goto Rollback
		}
	}

	session.Commit()
	return
Rollback:
	session.Rollback()
	return
}

func BlockerParametersToMap(params []db_models.BlockerParameter) map[string][]string {
	m := make(map[string][]string)

	for _, p := range params {
		a, ok := m[p.Key]
		if !ok {
			a = make([]string, 0)
		}
		a = append(a, p.Value)
		m[p.Key] = a
	}
	return m
}

func BlockerConfigurationParametersToMap(configParams []db_models.BlockerConfigurationParameter) map[string][]string {
	m := make(map[string][]string)

	for _, p := range configParams {
		a, ok := m[p.Key]
		if !ok {
			a = make([]string, 0)
		}
		a = append(a, p.Value)
		m[p.Key] = a
	}
	return m
}

func toBlocker(blocker db_models.Blocker, blockerConfig *db_models.BlockerConfiguration) (b Blocker, err error) {

	profile, err := GetLoginProfile(blocker.Id)
	if err != nil {
		return nil, err
	}

	base := BlockerBase{
		id:        blocker.Id,
		capacity:  blocker.Capacity,
		load:      blocker.Load,
		loginInfo: new(LoginProfile),
	}
	base.loginInfo.Load(profile)

	// Get blocker parameter and parse it to map
	params, err := GetBlockerParameters(blocker.Id)
	if err != nil {
		return nil, err
	}
	paramMap := BlockerParametersToMap(params)

	// Get blocker configuration parameter and parse it to map
	configParams, err := GetBlockerConfigurationParameters(blockerConfig.Id)
	if err != nil {
		return nil, err
	}
	configParamMap := BlockerConfigurationParametersToMap(configParams)

	switch blocker.BlockerType {
	case BLOCKER_TYPE_GoBGP_RTBH:
		b = NewGoBgpRtbhReceiver(base, paramMap)
	case BLOCKER_TYPE_GoBGP_FLOWSPEC:
		b = NewGoBgpFlowSpecReceiver(base, paramMap, configParamMap)
	case BLOCKER_TYPE_GO_ARISTA:
		b = NewGoAristaReceiver(base, configParamMap)
	default:
		err = errors.New(fmt.Sprintf("invalid blocker type: %s", blocker.BlockerType))
		return nil, err
	}

	log.Debugf("toBlocker result: [%T] %+v", b, b)
	return b, nil
}

/*
 * Get blocker configuration by customerID and targetType
 */
func GetBlockerConfiguration(customerId int, targetType string) (blockerConfiguration *db_models.BlockerConfiguration, err error){
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return nil, err
	}
	blockerConfig := db_models.BlockerConfiguration{}

	_, err = engine.Where("customer_id = ? AND target_type = ?", customerId, targetType).Get(&blockerConfig)
	if err != nil {
		log.Printf("Get blocker_configuration error: %s\n", err)
		return nil, err
	}

    if (blockerConfig.CustomerId != customerId) {
        err = errors.New(fmt.Sprintf("No blocker found. customerId: %d targetType: %s", customerId, targetType))
        return　nil, err
    }

	blockerConfiguration = &blockerConfig

	return blockerConfiguration, nil
}
