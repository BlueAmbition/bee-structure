package app

import (
	"strconv"
	"strings"
)

//获取语言
func GetAppLang(headerLang string) string {
	//转换成language_word表中的相应的语言字段
	lang := strings.ToLower(headerLang)
	//过滤语言设置默认
	if lang != "zh-cn" && lang != "en-us" && lang != "zh-hk" {
		lang = "en-us"
	}
	lang = strings.Replace(lang, "-", "_", -1)
	return lang
}

/**
 * 通用返回Map对象
 *@param version 版本号。
 *@return 版本转换的整数
 */
func GetVersion(version string) int64 {
	versionStr := strings.Replace(version, "v", "", 1)
	versionStr = strings.Replace(versionStr, ".", "", -1)
	versionNum, _ := strconv.ParseInt(versionStr, 10, 64)
	return versionNum
}