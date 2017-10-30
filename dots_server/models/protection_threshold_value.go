package models

import "time"


type ProtectionThresholdValue struct {
	Id               int64
	ProtectionId     int64
	ThresholdPackets int
	ThresholdBytes   int64
	ExaminationStart time.Time
	ExaminationEnd   time.Time
	Created          time.Time
	Updated          time.Time
}

func NewProtectionThresholdValue() (s *ProtectionThresholdValue) {
    s = &ProtectionThresholdValue{
        Id: 0,
        ProtectionId: 0,
        ThresholdPackets: 0,
        ThresholdBytes: 0,
        ExaminationStart: time.Unix(0, 0),
        ExaminationEnd: time.Unix(0, 0),
    }
    return
}

func CreateProtectionThresholdValueModel(pId int64, thresholdPackets int, thresholdBytes int64, examinationStart time.Time, examinationEnd time.Time) (ProtectionThresholdValue) {
    retProtectionThresholdValue := ProtectionThresholdValue{
        ProtectionId: pId,
        ThresholdPackets: thresholdPackets,
        ThresholdBytes: thresholdBytes,
        ExaminationStart: examinationStart,
        ExaminationEnd: examinationEnd,
    }

    return retProtectionThresholdValue
}

func (p *ProtectionThresholdValue) SetId(id int64) {
    if p != nil {
        p.Id = id
    }
}
