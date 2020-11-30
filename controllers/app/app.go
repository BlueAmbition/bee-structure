package app

import (
	"bee-structure/cache"
	"bee-structure/controllers/base"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

type AppController struct {
	base.BaseController
}

// @Title 获取APP多国语言词组
// @Description 获取APP多国语言词组
// @Success 200 {"code":200,"data":[{"id":1,"language_key":"string 语言key","zh_cn":"string 简体中文释义","zh_hk":"string 繁体中文释义","en_us":"string 英文释义","ko_kr":"string 韩文释义","ja_jp":"string 日文释义","module":"string 模块释义"}]}
// @router /app-lang [get]
func (c *AppController) AppLang() {
	list := cache.GetAppLangWords()
	msg := c.ReturnMsg("success_get_data")
	c.ResJson(base.SuccessCode, msg, list)
}

// @Title 获取选择语言列表
// @Description 获取选择语言列表
// @Success 200 {"code":200,"data":[{"id":1,"language":"string 语言","language_en":"string 语言英文名","local_code":"string 地区码zh-CN","icon":"string 图标地址哦"}],"msg":"Successful data acquisition","status":"success"}
// @router /lang-list [get]
func (c *AppController) LangList() {
	list := cache.GetLangList()
	msg := c.ReturnMsg("success_get_data")
	c.ResJson(base.SuccessCode, msg, list)
}

/**
 * 通用返回Map对象
 *@param version 版本号。
 *@return 版本转换的整数
 */
func getVersionNum(version string) int64 {
	versionStr := strings.Replace(version, "v", "", 1)
	versionStr = strings.Replace(versionStr, ".", "", -1)
	versionNum, _ := strconv.ParseInt(versionStr, 10, 64)
	return versionNum
}

// @Title 检查更新
// @Description 检查App更新
// @Success 200 {"code":200,"data":{"id":"int64 主键","app_name":"string app名","version":"string 版本号","intro":"string 简介","force":"int 1:强制升级","android_url":"string 安卓包下载地址","ios_url":"string ios包下载地址","need_upgrade":"bool 是否需要升级","title": "string 升级标题","content": "string 升级内容","created_at":"2019-09-24T14:43:38+08:00","updated_at":"2020-01-08T15:07:02+08:00"},"msg":"Successful data acquisition","status":"success"}
// @router /upgrade [get]
func (c *AppController) Upgrade() {
	v1 := c.Ctx.Input.Header("Version")
	v1 = strings.ToLower(v1)
	v1Num := getVersionNum(v1)
	version := cache.GetAppVersion()
	v2 := strings.ToLower(version.Version)
	v2Num := getVersionNum(v2)
	diffVersion := v2Num - v1Num
	if diffVersion > 0 || diffVersion < 0 {
		version.NeedUpgrade = true
	}
	if diffVersion > 1 || diffVersion < 0 || v1Num == 0 {
		version.Force = 1
	}
	version.Title = c.ReturnMsg("app_title")
	version.Content = c.ReturnMsg("app_content")
	msg := c.ReturnMsg("success_get_data")
	c.ResJson(base.SuccessCode, msg, version)
}

// @Title APP初始化配置信息
// @Description APP初始化配置信息
// @Success 200 {"code":200,"data":{"im_app_key":"string IM AppKey","im_secret":"string IM Secret","news_link":"分享地址","img_domain":"string 图片域名","agreement_dir":"string 用户协议在线目录","words_refresh_ts":"多语言刷新时间戳"},"msg":"Successful data acquisition","status":"success"}
// @router /init-configs [get]
func (c *AppController) InitConfigs() {
	mapReturn := make(map[string]interface{})
	domain := c.GetImgDomain()
	imAppKey := beego.AppConfig.String("im::app_key")
	imSecret := beego.AppConfig.String("im::secret")
	newsLink := beego.AppConfig.String("url::news_link")
	agreementDir := beego.AppConfig.String("url::agreement")
	mapReturn["news_link"] = newsLink
	mapReturn["img_domain"] = domain
	mapReturn["agreement_dir"] = agreementDir
	mapReturn["im_app_key"] = imAppKey
	mapReturn["im_secret"] = imSecret
	mapReturn["words_refresh_ts"] = cache.GetWordsRefreshTs()
	msg := c.ReturnMsg("success_get_data")
	c.ResJson(base.SuccessCode, msg, mapReturn)
}
