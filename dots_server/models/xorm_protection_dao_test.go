package models_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/nttdots/go-dots/dots_server/db_models"
	"github.com/nttdots/go-dots/dots_server/models"
	log "github.com/sirupsen/logrus"
)

// implements Protection
type testProtectionStruct struct {
	models.ProtectionBase

	TestParam1 int
	TestParam2 string
}

var testProtectionBase models.ProtectionBase
var testForwardedDataInfo, testBlockedDataInfo *models.ProtectionStatus
var testForwardedPeakThroughput, testForwardedAverageThroughput *models.ThroughputData
var testBlockedPeakThroughput, testBlockedAverageThroughput *models.ThroughputData
var testParam1 int
var testParam2 string
var testUpdProtectionBase models.ProtectionBase
var testUpdForwardedDataInfo, testUpdBlockedDataInfo *models.ProtectionStatus
var testUpdForwardedPeakThroughput, testUpdForwardedAverageThroughput *models.ThroughputData
var testUpdBlockedPeakThroughput, testUpdBlockedAverageThroughput *models.ThroughputData
var testUpdParam1 int
var testUpdParam2 string

func protectionSampleDataCreate() {
	loc, _ := time.LoadLocation("Asia/Tokyo")

	// create a test protection instance.
	testForwardedPeakThroughput = models.NewThroughputData(0, 1000, 1200)
	testForwardedAverageThroughput = models.NewThroughputData(0, 500, 600)
	testBlockedPeakThroughput = models.NewThroughputData(0, 11000, 11200)
	testBlockedAverageThroughput = models.NewThroughputData(0, 1500, 1600)
	testForwardedDataInfo = models.NewProtectionStatus(0, 10000, 20000, testForwardedPeakThroughput, testForwardedAverageThroughput)
	testBlockedDataInfo = models.NewProtectionStatus(0, 20000, 30000, testBlockedPeakThroughput, testBlockedAverageThroughput)

	testProtectionBase = models.NewProtectionBase(
		0,
		111222,
		true,
		time.Date(2017, 1, 2, 10, 00, 00, 0, loc),
		time.Date(2017, 2, 1, 9, 59, 59, 0, loc),
		time.Date(2017, 1, 5, 11, 22, 33, 44, loc),
		nil,
		testForwardedDataInfo,
		testBlockedDataInfo)

	testParam1 = 99999
	testParam2 = "TestValue2"

	// create a test protection instance to update.
	testUpdForwardedPeakThroughput = models.NewThroughputData(0, 1009, 1209)
	testUpdForwardedAverageThroughput = models.NewThroughputData(0, 509, 609)
	testUpdBlockedPeakThroughput = models.NewThroughputData(0, 11009, 112009)
	testUpdBlockedAverageThroughput = models.NewThroughputData(0, 1509, 1609)
	testUpdForwardedDataInfo = models.NewProtectionStatus(0, 10009, 20009, testUpdForwardedPeakThroughput, testUpdForwardedAverageThroughput)
	testUpdBlockedDataInfo = models.NewProtectionStatus(0, 20009, 30009, testUpdBlockedPeakThroughput, testUpdBlockedAverageThroughput)
	testUpdProtectionBase = models.NewProtectionBase(
		0,
		111222,
		false,
		time.Date(2010, 11, 12, 10, 00, 00, 0, loc),
		time.Date(2010, 12, 11, 9, 59, 59, 0, loc),
		time.Date(2010, 10, 15, 11, 22, 33, 44, loc),
		nil,
		testUpdForwardedDataInfo,
		testUpdBlockedDataInfo)
	testUpdParam1 = 123123
	testUpdParam2 = "TestValue999"
}

func TestCreateProtection2(t *testing.T) {
	rtbhProtection := models.NewRTBHProtection(
		testProtectionBase,
		map[string][]string{
			models.RTBH_PROTECTION_CUSTOMER_ID: {strconv.Itoa(testParam1)},
			models.RTBH_PROTECTION_TARGET:      {testParam2},
		},
	)

	db_np, err := models.CreateProtection2(rtbhProtection)
	if err != nil {
		t.Errorf("CreateProtection2 err: %s", err)
		return
	}
	np, _ := models.GetProtectionById(db_np.Id)
	log.WithFields(log.Fields{
		"p": np,
	}).Debug("get_protection")

	testProtectionBase = models.NewProtectionBase(
		np.Id(),
		np.MitigationId(),
		np.IsEnabled(),
		np.StartedAt(),
		np.FinishedAt(),
		np.RecordTime(),
		nil,
		np.ForwardedDataInfo(),
		np.BlockedDataInfo())
}

func TestDeleteProtectionById(t *testing.T) {

	// preparing for the test
	engine, _ := models.ConnectDB()
	var p db_models.Protection
	ok, _ := engine.Where("id=?", 100).Get(&p)
	if !ok {
		t.Errorf("protection id error: %d", 100)
		return
	}

	var params []db_models.ProtectionParameter
	engine.Where("protection_id = ?", 100).Find(&params)
	if len(params) != 2 {
		t.Errorf("protection_parameters error: %d, %d", 100, len(params))
		return
	}

	var status []db_models.ProtectionStatus
	engine.In("id", p.ForwardedDataInfoId, p.BlockedDataInfoId).Find(&status)
	if len(status) != 2 {
		t.Errorf("protection_status error: %d, %d", 100, len(status))
		return
	}

	var throuputData []db_models.ThroughputData
	engine.In("id", status[0].AverageThroughputId, status[0].PeakThroughputId, status[1].AverageThroughputId, status[1].PeakThroughputId).Find(&throuputData)
	if len(throuputData) != 4 {
		t.Errorf("throuputData error: %d, %d", 100, len(throuputData))
		return
	}

	models.DeleteProtectionById(100)

	throuputData = make([]db_models.ThroughputData, 0)
	engine.In("id", status[0].AverageThroughputId, status[0].PeakThroughputId, status[1].AverageThroughputId, status[1].PeakThroughputId).Find(&throuputData)
	if len(throuputData) != 0 {
		t.Errorf("throuputData delete error: %d, %d", 100, len(throuputData))
		return
	}

	status = make([]db_models.ProtectionStatus, 0)
	engine.In("id", p.ForwardedDataInfoId, p.BlockedDataInfoId).Find(&status)
	if len(status) != 0 {
		t.Errorf("protection_status delete error: %d, %d", 100, len(status))
		return
	}

	params = make([]db_models.ProtectionParameter, 0)
	engine.Where("protection_id = ?", 100).Find(&params)
	if len(params) != 0 {
		t.Errorf("protection_parameters delete error: %d, %d", 100, len(params))
		return
	}

	ok, _ = engine.ID(100).Get(&p)
	if ok {
		t.Errorf("protection delete error: %d", 100)
		return
	}
}

func TestGetProtectionBase(t *testing.T) {
	protection, err := models.GetProtectionBase(testProtectionBase.MitigationId())
	if err != nil {
		t.Errorf("get protection err: %s", err)
	}

	if protection.MitigationId() != testProtectionBase.MitigationId() {
		t.Errorf("got %s, want %s", protection.MitigationId(), testProtectionBase.MitigationId())
	}

	if protection.IsEnabled() != testProtectionBase.IsEnabled() {
		t.Errorf("got %b, want %b", protection.IsEnabled(), testProtectionBase.IsEnabled())
	}

	if protection.StartedAt().Unix() != testProtectionBase.StartedAt().Unix() {
		t.Errorf("got %d, want %d", protection.StartedAt().Unix(), testProtectionBase.StartedAt().Unix())
	}

	if protection.FinishedAt().Unix() != testProtectionBase.FinishedAt().Unix() {
		t.Errorf("got %d, want %d", protection.FinishedAt().Unix(), testProtectionBase.FinishedAt().Unix())
	}

	if protection.RecordTime().Unix() != testProtectionBase.RecordTime().Unix() {
		t.Errorf("got %d, want %d", protection.RecordTime().Unix(), testProtectionBase.RecordTime().Unix())
	}

	/* temporarily commenting out these lines.
	if protection.ForwardedDataInfo().TotalPackets() != testForwardedDataInfo.TotalPackets() {
		t.Errorf("got %d, want %d", protection.ForwardedDataInfo().TotalPackets(), testForwardedDataInfo.TotalPackets())
	}

	if protection.ForwardedDataInfo().TotalBits() != testForwardedDataInfo.TotalBits() {
		t.Errorf("got %d, want %d", protection.ForwardedDataInfo().TotalBits(), testForwardedDataInfo.TotalBits())
	}

	if protection.BlockedDataInfo().TotalPackets() != testBlockedDataInfo.TotalPackets() {
		t.Errorf("got %d, want %d", protection.BlockedDataInfo().TotalPackets(), testBlockedDataInfo.TotalPackets())
	}

	if protection.BlockedDataInfo().TotalBits() != testBlockedDataInfo.TotalBits() {
		t.Errorf("got %d, want %d", protection.BlockedDataInfo().TotalBits(), testBlockedDataInfo.TotalBits)
	}

	if protection.ForwardedDataInfo().PeakThroughput().Pps() != testForwardedPeakThroughput.Pps() {
		t.Errorf("got %d, want %d", protection.ForwardedDataInfo().PeakThroughput().Pps(), testForwardedPeakThroughput.Pps)
	}

	if protection.ForwardedDataInfo().PeakThroughput().Bps() != testForwardedPeakThroughput.Bps() {
		t.Errorf("got %d, want %d", protection.ForwardedDataInfo().PeakThroughput().Bps(), testForwardedPeakThroughput.Bps)
	}

	if protection.ForwardedDataInfo().AverageThroughput().Pps() != testForwardedAverageThroughput.Pps() {
		t.Errorf("got %d, want %d", protection.ForwardedDataInfo().AverageThroughput().Pps(), testForwardedAverageThroughput.Pps)
	}

	if protection.ForwardedDataInfo().AverageThroughput().Bps() != testForwardedAverageThroughput.Bps() {
		t.Errorf("got %d, want %d", protection.ForwardedDataInfo().AverageThroughput().Bps(), testForwardedAverageThroughput.Bps)
	}

	if protection.BlockedDataInfo().PeakThroughput().Pps() != testBlockedPeakThroughput.Pps() {
		t.Errorf("got %d, want %d", protection.BlockedDataInfo().PeakThroughput().Pps(), testBlockedPeakThroughput.Pps)
	}

	if protection.BlockedDataInfo().PeakThroughput().Bps() != testBlockedPeakThroughput.Bps() {
		t.Errorf("got %d, want %d", protection.BlockedDataInfo().PeakThroughput().Bps(), testBlockedPeakThroughput.Bps)
	}

	if protection.BlockedDataInfo().AverageThroughput().Pps() != testBlockedAverageThroughput.Pps() {
		t.Errorf("got %d, want %d", protection.BlockedDataInfo().AverageThroughput().Pps(), testBlockedAverageThroughput.Pps)
	}

	if protection.BlockedDataInfo().AverageThroughput().Bps() != testBlockedAverageThroughput.Bps() {
		t.Errorf("got %d, want %d", protection.BlockedDataInfo().AverageThroughput().Bps(), testBlockedAverageThroughput.Bps)
	}
	*/

}

func TestGetProtectionParameters(t *testing.T) {
	// data check
	engine, err := models.ConnectDB()
	if err != nil {
		t.Errorf("database connect error: %s", err)
		return
	}

	protection := db_models.Protection{}
	_, err = engine.Where("mitigation_id = ?", testProtectionBase.MitigationId()).Get(&protection)
	if err != nil {
		t.Errorf("get protection err: %s", err)
	}

	protectionParameters, err := models.GetProtectionParameters(protection.Id)
	if err != nil {
		t.Errorf("get protectionParameters err: %s", err)
	}
	if len(protectionParameters) != 2 {
		t.Errorf("protectionParameters record count got %d, want 1", len(protectionParameters))
	}

	if protectionParameters[0].Key != "customerId" {
		t.Errorf("got %s, want %s", protectionParameters[0].Key, "customerId")
	}

	if protectionParameters[0].Value != strconv.Itoa(testParam1) {
		t.Errorf("got %s, want %s", protectionParameters[0].Value, strconv.Itoa(testParam1))
	}

	if protectionParameters[1].Key != "target" {
		t.Errorf("got %s, want %s", protectionParameters[1].Key, "target")
	}

	if protectionParameters[1].Value != testParam2 {
		t.Errorf("got %s, want %s", protectionParameters[1].Value, testParam2)
	}

}
func TestUpdateProtection(t *testing.T) {
	// CreateData id setting
	p, err := models.GetProtectionBase(testProtectionBase.MitigationId())
	if err != nil {
		t.Errorf("get protection err: %s", err)
		return
	}
	testUpdForwardedDataInfo.SetId(p.ForwardedDataInfo().Id())
	testUpdForwardedPeakThroughput.SetId(p.ForwardedDataInfo().PeakThroughput().Id())
	testUpdForwardedAverageThroughput.SetId(p.ForwardedDataInfo().AverageThroughput().Id())
	testUpdBlockedDataInfo.SetId(p.BlockedDataInfo().Id())
	testUpdBlockedPeakThroughput.SetId(p.BlockedDataInfo().PeakThroughput().Id())
	testUpdBlockedAverageThroughput.SetId(p.BlockedDataInfo().AverageThroughput().Id())

	rtbhProtection := models.NewRTBHProtection(
		models.NewProtectionBase(
			testProtectionBase.Id(),
			testUpdProtectionBase.MitigationId(),
			testUpdProtectionBase.IsEnabled(),
			testUpdProtectionBase.StartedAt(),
			testUpdProtectionBase.FinishedAt(),
			testUpdProtectionBase.RecordTime(),
			nil,
			testUpdForwardedDataInfo,
			testUpdBlockedDataInfo,
		),
		map[string][]string{
			models.RTBH_PROTECTION_CUSTOMER_ID: {strconv.Itoa(testUpdParam1)},
			models.RTBH_PROTECTION_TARGET:      {testUpdParam2},
		},
	)

	err = models.UpdateProtection(rtbhProtection)
	if err != nil {
		t.Errorf("CreateProtection2 err: %s", err)
		return
	}

	protection, err := models.GetProtectionBase(testUpdProtectionBase.MitigationId())
	if err != nil {
		t.Errorf("get protection err: %s", err)
		return
	}

	if protection.MitigationId() != testUpdProtectionBase.MitigationId() {
		t.Errorf("got %s, want %s", protection.MitigationId(), testUpdProtectionBase.MitigationId())
		return
	}

	if protection.IsEnabled() != testUpdProtectionBase.IsEnabled() {
		t.Errorf("got %b, want %b", protection.IsEnabled(), testUpdProtectionBase.IsEnabled())
		return
	}

	if protection.StartedAt().Unix() != testUpdProtectionBase.StartedAt().Unix() {
		t.Errorf("got %d, want %d", protection.StartedAt().Unix(), testUpdProtectionBase.StartedAt().Unix())
	}

	if protection.FinishedAt().Unix() != testUpdProtectionBase.FinishedAt().Unix() {
		t.Errorf("got %d, want %d", protection.FinishedAt().Unix(), testUpdProtectionBase.FinishedAt().Unix())
	}

	if protection.RecordTime().Unix() != testUpdProtectionBase.RecordTime().Unix() {
		t.Errorf("got %d, want %d", protection.RecordTime().Unix(), testUpdProtectionBase.RecordTime().Unix())
	}

	protectionParameters, err := models.GetProtectionParameters(protection.Id())
	if err != nil {
		t.Errorf("get protectionParameters err: %s", err)
	}
	if len(protectionParameters) != 2 {
		t.Errorf("protectionParameters record count got %d, want 1", len(protectionParameters))
	}

	if protectionParameters[0].Key != "customerId" {
		t.Errorf("got %s, want %s", protectionParameters[0].Key, "customerId")
	}

	if protectionParameters[0].Value != strconv.Itoa(testUpdParam1) {
		t.Errorf("got %s, want %s", protectionParameters[0].Value, strconv.Itoa(testUpdParam1))
	}

	if protectionParameters[1].Key != "target" {
		t.Errorf("got %s, want %s", protectionParameters[1].Key, "target")
	}

	if protectionParameters[1].Value != testUpdParam2 {
		t.Errorf("got %s, want %s", protectionParameters[1].Value, testUpdParam2)
	}

	/* temporarily commenting out these lines
	if protection.ForwardedDataInfo().TotalPackets() != testUpdForwardedDataInfo.TotalPackets() {
		t.Errorf("got %d, want %d", protection.ForwardedDataInfo().TotalPackets(), testUpdForwardedDataInfo.TotalPackets())
	}

	if protection.ForwardedDataInfo().TotalBits() != testUpdForwardedDataInfo.TotalBits() {
		t.Errorf("got %d, want %d", protection.ForwardedDataInfo().TotalBits(), testUpdForwardedDataInfo.TotalBits())
	}

	if protection.BlockedDataInfo().TotalPackets() != testUpdBlockedDataInfo.TotalPackets() {
		t.Errorf("got %d, want %d", protection.BlockedDataInfo().TotalPackets(), testUpdBlockedDataInfo.TotalPackets())
	}

	if protection.BlockedDataInfo().TotalBits() != testUpdBlockedDataInfo.TotalBits() {
		t.Errorf("got %d, want %d", protection.BlockedDataInfo().TotalBits(), testUpdBlockedDataInfo.TotalBits())
	}

	if protection.ForwardedDataInfo().PeakThroughput().Pps() != testUpdForwardedPeakThroughput.Pps() {
		t.Errorf("got %d, want %d", protection.ForwardedDataInfo().PeakThroughput().Pps(), testUpdForwardedPeakThroughput.Pps())
	}

	if protection.ForwardedDataInfo().PeakThroughput().Bps() != testUpdForwardedPeakThroughput.Bps() {
		t.Errorf("got %d, want %d", protection.ForwardedDataInfo().PeakThroughput().Bps(), testUpdForwardedPeakThroughput.Bps())
	}

	if protection.ForwardedDataInfo().AverageThroughput().Pps() != testUpdForwardedAverageThroughput.Pps() {
		t.Errorf("got %d, want %d", protection.ForwardedDataInfo().AverageThroughput().Pps(), testUpdForwardedAverageThroughput.Pps())
	}

	if protection.ForwardedDataInfo().AverageThroughput().Bps() != testUpdForwardedAverageThroughput.Bps() {
		t.Errorf("got %d, want %d", protection.ForwardedDataInfo().AverageThroughput().Bps(), testUpdForwardedAverageThroughput.Bps())
	}

	if protection.BlockedDataInfo().PeakThroughput().Pps() != testUpdBlockedPeakThroughput.Pps() {
		t.Errorf("got %d, want %d", protection.BlockedDataInfo().PeakThroughput().Pps(), testUpdBlockedPeakThroughput.Pps)
	}

	if protection.BlockedDataInfo().PeakThroughput().Bps() != testUpdBlockedPeakThroughput.Bps() {
		t.Errorf("got %d, want %d", protection.BlockedDataInfo().PeakThroughput().Bps(), testUpdBlockedPeakThroughput.Bps)
	}

	if protection.BlockedDataInfo().AverageThroughput().Pps() != testUpdBlockedAverageThroughput.Pps() {
		t.Errorf("got %d, want %d", protection.BlockedDataInfo().AverageThroughput().Pps(), testUpdBlockedAverageThroughput.Pps)
	}

	if protection.BlockedDataInfo().AverageThroughput().Bps() != testUpdBlockedAverageThroughput.Bps() {
		t.Errorf("got %d, want %d", protection.BlockedDataInfo().AverageThroughput().Bps(), testUpdBlockedAverageThroughput.Bps)
	}
	*/

}

func TestDeleteProtection(t *testing.T) {
	err := models.DeleteProtection(testProtectionBase.MitigationId())
	if err != nil {
		t.Errorf("delete protection err: %s", err)
	}

	engine, err := models.ConnectDB()
	if err != nil {
		t.Errorf("database connect error: %s", err)
	}

	// Protection
	chkProtection := db_models.Protection{}
	_, err = engine.Where("id=?", testProtectionBase.Id()).Get(&chkProtection)
	if err != nil {
		t.Errorf("select protection err: %s", err)
	}
	if chkProtection.Id > 0 {
		t.Errorf("not delete protection %d", testProtectionBase.MitigationId())
	}

	// ForwardedDataInfo
	chkProtectionStatus := db_models.ProtectionStatus{}
	_, err = engine.Where("id=?", testProtectionBase.ForwardedDataInfo().Id()).Get(&chkProtectionStatus)
	if err != nil {
		t.Errorf("select ForwardedDataInfo err: %s", err)
	}
	if chkProtectionStatus.Id > 0 {
		t.Error("not delete ForwardedDataInfo")
	}

	// BlockedDataInfo
	chkProtectionStatus = db_models.ProtectionStatus{}
	_, err = engine.Where("id=?", testProtectionBase.BlockedDataInfo().Id()).Get(&chkProtectionStatus)
	if err != nil {
		t.Errorf("select BlockedDataInfo err: %s", err)
	}
	if chkProtectionStatus.Id > 0 {
		t.Error("not delete BlockedDataInfo")
	}

	// ForwardedDataInfo -> PeakThroughputData
	chkThroughputData := db_models.ThroughputData{}
	_, err = engine.Where("id=?", testProtectionBase.ForwardedDataInfo().PeakThroughput().Id()).Get(&chkThroughputData)
	if err != nil {
		t.Errorf("select ForwardedDataInfo().PeakThroughputData err: %s", err)
	}
	if chkProtectionStatus.Id > 0 {
		t.Error("not delete ForwardedDataInfo().PeakThroughputData")
	}

	// ForwardedDataInfo -> AverageThroughputData
	chkThroughputData = db_models.ThroughputData{}
	_, err = engine.Where("id=?", testProtectionBase.ForwardedDataInfo().AverageThroughput().Id()).Get(&chkThroughputData)
	if err != nil {
		t.Errorf("select ForwardedDataInfo().AverageThroughputData err: %s", err)
	}
	if chkProtectionStatus.Id > 0 {
		t.Error("not delete ForwardedDataInfo().AverageThroughputData")
	}

	// BlockedDataInfo -> PeakThroughputData
	chkThroughputData = db_models.ThroughputData{}
	_, err = engine.Where("id=?", testProtectionBase.BlockedDataInfo().PeakThroughput().Id()).Get(&chkThroughputData)
	if err != nil {
		t.Errorf("select BlockedDataInfo.PeakThroughputData err: %s", err)
	}
	if chkProtectionStatus.Id > 0 {
		t.Error("not delete BlockedDataInfo.PeakThroughputData")
	}

	// BlockedDataInfo -> AverageThroughputData
	chkThroughputData = db_models.ThroughputData{}
	_, err = engine.Where("id=?", testProtectionBase.BlockedDataInfo().AverageThroughput().Id()).Get(&chkThroughputData)
	if err != nil {
		t.Errorf("select BlockedDataInfo.AverageThroughputData err: %s", err)
	}
	if chkProtectionStatus.Id > 0 {
		t.Error("not delete BlockedDataInfo.AverageThroughputData")
	}
}

func TestCreateProtectionThresholdValue(t *testing.T) {
	nowTime := time.Now()
	testProtectionThresholdValue := models.ProtectionThresholdValue{
		ProtectionId:     testProtectionBase.Id(),
		ThresholdPackets: 5000,
		ThresholdBytes:   20000,
		ExaminationStart: nowTime,
		ExaminationEnd: nowTime.Add(2 * time.Minute),
	}
	err := models.CreateProtectionThresholdValue(&testProtectionThresholdValue)
	if err != nil {
		t.Errorf("create protection_threshold_value err: %s", err)
	}

	engine, err := models.ConnectDB()
	if err != nil {
		t.Errorf("database connect error: %s", err)
	}

	// ProtectionThresholdValue
	chkProtectionThresholdValue := db_models.ProtectionThresholdValue{}
	_, err = engine.Where("id=?", testProtectionThresholdValue.Id).Get(&chkProtectionThresholdValue)
	if err != nil {
		t.Errorf("select protection_threshold_value err: %s", err)
	}

	if chkProtectionThresholdValue.ProtectionId != testProtectionThresholdValue.ProtectionId {
		t.Errorf("ProtectionId got %d, want %d", chkProtectionThresholdValue.ProtectionId, testProtectionThresholdValue.ProtectionId)
	}

	if chkProtectionThresholdValue.ThresholdPackets != testProtectionThresholdValue.ThresholdPackets {
		t.Errorf("ThresholdPackets got %d, want %d", chkProtectionThresholdValue.ThresholdPackets, testProtectionThresholdValue.ThresholdPackets)
	}

	if chkProtectionThresholdValue.ThresholdBytes != testProtectionThresholdValue.ThresholdBytes {
		t.Errorf("ThresholdBytes got %d, want %d", chkProtectionThresholdValue.ThresholdBytes, testProtectionThresholdValue.ThresholdBytes)
	}

	if chkProtectionThresholdValue.ExaminationStart != models.GetSysTime(models.GetMySqlTime(testProtectionThresholdValue.ExaminationStart)) {
		t.Errorf("ExaminationStart got %s, want %s", chkProtectionThresholdValue.ExaminationStart, models.GetSysTime(models.GetMySqlTime(testProtectionThresholdValue.ExaminationStart)))
	}

	if chkProtectionThresholdValue.ExaminationEnd != models.GetSysTime(models.GetMySqlTime(testProtectionThresholdValue.ExaminationEnd)) {
		t.Errorf("ExaminationEnd got %s, want %s", chkProtectionThresholdValue.ExaminationEnd, models.GetSysTime(models.GetMySqlTime(testProtectionThresholdValue.ExaminationEnd)))
	}

	testProtectionThresholdValue.ThresholdPackets = testProtectionThresholdValue.ThresholdPackets + 12345
	testProtectionThresholdValue.ThresholdBytes = testProtectionThresholdValue.ThresholdBytes + 987654
	testProtectionThresholdValue.ExaminationStart = nowTime.Add(3 * time.Minute)
	testProtectionThresholdValue.ExaminationEnd = nowTime.Add(5 * time.Minute)
	err = models.CreateProtectionThresholdValue(&testProtectionThresholdValue)
	if err != nil {
		t.Errorf("update protection_threshold_value err: %s", err)
	}

	// ProtectionThresholdValue
	chkProtectionThresholdValue = db_models.ProtectionThresholdValue{}
	_, err = engine.Where("id=?", testProtectionThresholdValue.Id).Get(&chkProtectionThresholdValue)
	if err != nil {
		t.Errorf("select protection_threshold_value err: %s", err)
	}

	if chkProtectionThresholdValue.ProtectionId != testProtectionThresholdValue.ProtectionId {
		t.Errorf("ProtectionId got %d, want %d", chkProtectionThresholdValue.ProtectionId, testProtectionThresholdValue.ProtectionId)
	}

	if chkProtectionThresholdValue.ThresholdPackets != testProtectionThresholdValue.ThresholdPackets {
		t.Errorf("ThresholdPackets got %d, want %d", chkProtectionThresholdValue.ThresholdPackets, testProtectionThresholdValue.ThresholdPackets)
	}

	if chkProtectionThresholdValue.ThresholdBytes != testProtectionThresholdValue.ThresholdBytes {
		t.Errorf("ThresholdBytes got %d, want %d", chkProtectionThresholdValue.ThresholdBytes, testProtectionThresholdValue.ThresholdBytes)
	}

	if chkProtectionThresholdValue.ExaminationStart != models.GetSysTime(models.GetMySqlTime(testProtectionThresholdValue.ExaminationStart)) {
		t.Errorf("ExaminationStart got %s, want %s", chkProtectionThresholdValue.ExaminationStart, models.GetSysTime(models.GetMySqlTime(testProtectionThresholdValue.ExaminationStart)))
	}

	if chkProtectionThresholdValue.ExaminationEnd != models.GetSysTime(models.GetMySqlTime(testProtectionThresholdValue.ExaminationEnd)) {
		t.Errorf("ExaminationEnd got %s, want %s", chkProtectionThresholdValue.ExaminationEnd, models.GetSysTime(models.GetMySqlTime(testProtectionThresholdValue.ExaminationEnd)))
	}

}
