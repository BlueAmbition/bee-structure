package datetime

import (
	"fmt"
	"time"
)

// 格式化时间
func Format(timeStr string, formal string) string {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeLayout, timeStr, loc)
	ts := tmp.Unix() //转化为时间戳 类型是int64
	needTime := time.Unix(ts, 0).Format(formal)
	return needTime
}

//通过当前时间获取距离现在的时间
//duration 间隔数字
//unit Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"
func GetTimeByNow(duration int, unit string) time.Time {
	now := time.Now() //质押时间
	tag := fmt.Sprintf("%v%v", duration, unit)
	tagDuration, _ := time.ParseDuration(tag)
	return now.Add(tagDuration)
}
