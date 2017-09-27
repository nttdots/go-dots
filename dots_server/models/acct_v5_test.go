package models_test

import (
    "testing"
    "github.com/nttdots/go-dots/dots_server/models"
    "time"
)

func TestTotalPacketsBytesCalc(t *testing.T) {
    var acctV5List []models.AcctV5

    newAcctV5 := models.AcctV5{
        1,
        "A1",
        "10.11.12.13.14.15",
        "20.21.22.23.24.25",
        1,
        "10.100.111.1",
        "192.168.1.0",
        5500,
        5501,
        "12345",
        1,
        10,
        100,
        1,
        time.Unix(0, 0),
        time.Unix(0, 0),
    }
    acctV5List = append(acctV5List, newAcctV5)

    newAcctV5 = models.AcctV5{
        1,
        "A1",
        "10.11.12.13.14.15",
        "20.21.22.23.24.25",
        1,
        "10.100.111.1",
        "192.168.1.0",
        5500,
        5501,
        "12345",
        1,
        20,
        300,
        1,
        time.Unix(0, 0),
        time.Unix(0, 0),
    }
    acctV5List = append(acctV5List, newAcctV5)

    newAcctV5 = models.AcctV5{
        1,
        "A1",
        "10.11.12.13.14.15",
        "20.21.22.23.24.25",
        1,
        "10.100.111.1",
        "192.168.1.0",
        5500,
        5501,
        "12345",
        1,
        30,
        500,
        1,
        time.Unix(0, 0),
        time.Unix(0, 0),
    }
    acctV5List = append(acctV5List, newAcctV5)

    var totalPackets int = 0
    var totalBytes int64 = 0
    for _,acctV5 := range acctV5List {
        totalPackets = totalPackets + acctV5.Packets
        totalBytes = totalBytes + acctV5.Bytes
    }

    packets, bytes := models.TotalPacketsBytesCalc(acctV5List)
    if packets != totalPackets {
        t.Errorf("Packets total got %s, want %s", packets, totalPackets)
    }

    if bytes != totalBytes {
        t.Errorf("Bytes total got %s, want %s", bytes, totalBytes)
    }

}
