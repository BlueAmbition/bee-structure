package cache

import (
	"bee-structure/functions/redis"
	"bee-structure/models/app"
	"encoding/json"
	"github.com/astaxie/beego"
)

//设置APP最新版缓存
func SetAppVersion() {
	version := app.GetLatestVersion()
	downUrl := beego.AppConfig.String("app::android_link")
	version.AndroidApk = downUrl
	value, _ := json.Marshal(version)
	redis.SetString(0, "app_version", value, 0)
}

//获取App最新版本信息
func GetAppVersion() app.AppVersion {
	value, _ := redis.GetString(0, "app_version")
	if value == "" {
		SetAppVersion()
		value, _ = redis.GetString(0, "app_version")
	}
	var object app.AppVersion
	_ = json.Unmarshal([]byte(value), &object)
	return object
}

