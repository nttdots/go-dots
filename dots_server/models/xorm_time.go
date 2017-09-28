package models

import "time"

const time_layout = "2006-01-02 15:04:05"

// MySQL用に指定時刻の文字列を返却
func GetMySqlTime(targetTime time.Time) string {
    return targetTime.Format(time_layout)
}

// MySQL用の時刻(文字列)からtime.Timeに変換し返却
// 変換出来ない文字列が指定された場合はtime.Now()を返却
func GetSysTime(targetTime string) time.Time {
    chgTime, err := time.ParseInLocation(time_layout, targetTime, time.Local)
    if err != nil {
        return time.Now()
    }
    return chgTime
}

// 指定されたtime.Timeに指定された秒数を加算し返却
func AddSecond(targetTime time.Time, addTime int64) time.Time {
    return targetTime.Add(time.Duration(addTime) * time.Second)
}

// 指定されたtime.Timeに指定された分数を加算し返却
func AddMinute(targetTime time.Time, addTime int64) time.Time {
    return targetTime.Add(time.Duration(addTime) * time.Minute)
}

// 指定されたtime.Timeに指定された時間数を加算し返却
func AddHour(targetTime time.Time, addTime int64) time.Time {
    return targetTime.Add(time.Duration(addTime) * time.Hour)
}
