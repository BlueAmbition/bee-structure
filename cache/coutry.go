package cache

import (
	"bee-structure/functions/redis"
	"bee-structure/models/app"
	"encoding/json"
)

//国家列表
func SetCountryList() {
	list := app.GetAllowCountry()
	for _, v := range list {
		value, _ := json.Marshal(v)
		redis.SetSortSet(0, "country_list", v.Id, string(value))
	}
}

//获取国家列表
func GetCountryList() []app.CountryExt {
	var (
		err         error
		flag        bool
		list        []interface{}
		JsonByte    []byte
		object      app.CountryExt
		countryList = make([]app.CountryExt, 0)
	)
	flag = redis.KeyExists(0, "country_list")
	if !flag {
		SetCountryList()
	}
	list = redis.GetSortSet(0, "country_list", 0, -1, false)
	if list != nil {
		for _, v := range list {
			JsonByte = v.([]uint8)
			err = json.Unmarshal(JsonByte, &object)
			if err == nil {
				countryList = append(countryList, object)
			}
		}
	}
	return countryList
}
