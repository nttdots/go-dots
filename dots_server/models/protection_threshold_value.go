package models

import (
    "time"
)

type ProtectionThresholdValue struct {
	Id               int64
	ProtectionId     int64
	ThresholdPackets int
	ThresholdBytes   int64
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
    }
    return
}

func CreateProtectionThresholdValueModel(p Protection, thresholdPackets int, thresholdBytes int64) (ProtectionThresholdValue) {
    retProtectionThresholdValue := ProtectionThresholdValue{
        ProtectionId: p.Id(),
        ThresholdPackets: thresholdPackets,
        ThresholdBytes: thresholdBytes,
    }

    return retProtectionThresholdValue
}

func (p *ProtectionThresholdValue) SetId(id int64) {
    if p != nil {
        p.Id = id
    }
}
