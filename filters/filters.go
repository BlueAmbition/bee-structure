package filters

import (
	"bee-structure/cache"
	"bee-structure/functions/app"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strings"
)

//版本是否需要更新
func AppUpgrade() {
	var App = func(ctx *context.Context) {
		v1 := ctx.Input.Header("Version")
		v1 = strings.ToLower(v1)
		v1Num := app.GetVersion(v1)
		version := cache.GetAppVersion()
		v2 := strings.ToLower(version.Version)
		v2Num := app.GetVersion(v2)
		diffVersion := v2Num - v1Num
		if diffVersion > 1 || diffVersion < 0 || v1Num == 0 {
			version.Force = 1
		}
		mapReturn := make(map[string]interface{})
		if diffVersion != 0 && version.Force == 1 {
			lang := ctx.Input.Header("Show-Language")
			if lang == "" {
				lang = "en-US"
			}
			lang = app.GetAppLang(lang)
			version.NeedUpgrade = true
			title, _ := cache.GetTipsWord("app_title", lang)
			version.Title = title
			content, _ := cache.GetTipsWord("app_content", lang)
			version.Content = content
			msg, _ := cache.GetTipsWord("force_upgrade", lang)
			ctx.Output.Header("content-type", "application/json")
			mapReturn["code"] = 301
			mapReturn["msg"] = msg
			mapReturn["status"] = "success"
			mapReturn["data"] = version
			jsonData, _ := json.Marshal(mapReturn)
			ctx.WriteString(string(jsonData))
		}
	}
	beego.InsertFilter("/v1/*", beego.BeforeRouter, App)
}
