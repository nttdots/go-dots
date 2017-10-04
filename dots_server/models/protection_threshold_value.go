package models

import (
    "time"
)

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
        0,
        0,
        0,
        0,
        time.Unix(0, 0),
        time.Unix(0, 0),
        time.Unix(0, 0),
        time.Unix(0, 0),
    }
    return
}

func CreateProtectionThresholdValueModel(p Protection, thresholdPackets int, thresholdBytes int64, examinationStart time.Time, examinationEnd time.Time) (ProtectionThresholdValue) {
    retProtectionThresholdValue := ProtectionThresholdValue{
        ProtectionId: p.Id(),
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
