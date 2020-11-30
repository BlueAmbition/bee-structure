package cache

import (
	"bee-structure/functions/redis"
	"bee-structure/models/app"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

//获取词组失败
func GetCacheError(lang string) string {
	tips := ""
	switch lang {
	case "zh_cn":
		tips = "获取词组失败。"
		break
	case "zh_hk":
		tips = "獲取片語失敗。"
		break
	case "ja_jp":
		tips = "フレーズの取得に失敗しました"
		break
	case "ko_kr":
		tips = "구문 가 져 오 는 데 실 패 했 습 니 다."
		break
	default:
		tips = "Failed to get phrase."
	}
	return tips
}

//删除语言信息，API相关缓存
func DelLang() bool {
	return redis.DelKey(0, "language_list")
}

//删除语言词组信息，API相关缓存
func DelLangWords() {
	redis.DelKey(0, "app_lang")
	list := redis.KeysList(0, "lang:*")
	if list != nil {
		for _, v := range list {
			redis.DelKey(0, v)
		}
	}
}

//设置APP词组刷新Key
func SetWordsRefreshTs() {
	ts := time.Now().Unix()
	redis.SetString(0, "app_words_refresh_ts", ts, 0)
}

//获取APP词组刷新Key
func GetWordsRefreshTs() int64 {
	value, _ := redis.GetString(0, "app_words_refresh_ts")
	if value == "" {
		SetWordsRefreshTs()
		value, _ = redis.GetString(0, "app_words_refresh_ts")
	}
	ts, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return ts
}

//获取多国语言词提示信息
func GetTipsWord(key string, lang string) (string, error) {
	flag := redis.KeyExists(0, "lang:"+key)
	if !flag {
		SetTipsLangWords()
	}
	hash := redis.GetHash(0, "lang:"+key, lang)
	if hash[0] == nil {
		errorTips := GetCacheError(lang)
		return "", errors.New(errorTips)
	}
	value := hash[0].([]uint8)
	return string(value), nil
}

//设置多国语言提示缓存
func SetTipsLangWords() {
	list := app.GetLangWordsByModular("api")
	for _, v := range list {
		redis.SetHash(0, "lang:"+v.LanguageKey, "en_us", v.EnUs, "zh_cn", v.ZhCN, "zh_hk", v.ZhHk)
	}
}

//设置多国语言App缓存
func SetAppLangWords() {
	list := app.GetLangWordsByModular("app")
	for _, v := range list {
		value, _ := json.Marshal(v)
		redis.SetSortSet(0, "app_lang", v.Id, string(value))
	}
}

//获取APP多国语言Json
func GetAppLangWords() map[string]app.LanguageWordSimple {
	var (
		err      error
		flag     bool
		list     []interface{}
		JsonByte []byte
		object   app.LanguageWordSimple
		wordsMap = make(map[string]app.LanguageWordSimple)
	)
	flag = redis.KeyExists(0, "app_lang")
	if !flag {
		SetAppLangWords()
	}
	list = redis.GetSortSet(0, "app_lang", 0, -1, false)
	if list != nil {
		for _, v := range list {
			JsonByte = v.([]uint8)
			err = json.Unmarshal(JsonByte, &object)
			if err == nil {
				wordsMap[object.LanguageKey] = object
			}
		}
	}
	return wordsMap
}

//语言列表
func SetLangList() {
	list, _ := app.LanguageList()
	for _, v := range list {
		value, _ := json.Marshal(v)
		redis.SetSortSet(0, "language_list", v.Id, string(value))
	}
}

//获取APP多国语言Json
func GetLangList() []app.LanguageExt {
	var (
		err      error
		flag     bool
		list     []interface{}
		JsonByte []byte
		object   app.LanguageExt
		langList = make([]app.LanguageExt, 0)
	)
	flag = redis.KeyExists(0, "language_list")
	if !flag {
		SetLangList()
	}
	list = redis.GetSortSet(0, "language_list", 0, -1, false)
	if list != nil {
		for _, v := range list {
			JsonByte = v.([]uint8)
			err = json.Unmarshal(JsonByte, &object)
			if err == nil {
				langList = append(langList, object)
			}
		}
	}
	return langList
}

//获取语言信息
func GetLangByLocalKey(localKey string) app.LanguageExt {
	var item app.LanguageExt
	list := GetLangList()
	for _, v := range list {
		if v.LocalCode == localKey {
			item = v
			break
		}
	}
	return item
}
