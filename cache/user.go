package cache

import (
	"bee-structure/functions/redis"
	"bee-structure/models/user"
	"encoding/json"
	"strconv"
)

//获取用户缓存信息
func GetUserInfo(userId int64) user.UserInfo {
	var (
		value  string
		object user.UserInfo
	)
	value, _ = redis.GetString(0, "user_info:"+strconv.FormatInt(userId, 10))
	if value == "" {
		SetUserInfo(userId)
		value, _ = redis.GetString(0, "user_info:"+strconv.FormatInt(userId, 10))
	}
	json.Unmarshal([]byte(value), &object)
	return object
}

//设置用户缓存信息
func SetUserInfo(userId int64) {
	userInfo := user.GetUserInfo(userId)
	if userInfo.Id > 0 {
		value, _ := json.Marshal(userInfo)
		redis.SetString(0, "user_info:"+strconv.FormatInt(userId, 10), value, 0)
	}
}

//删除用户缓存
func DelUserInfo(userId int64) {
	redis.DelKey(0, "user_info:"+strconv.FormatInt(userId, 10))
}
