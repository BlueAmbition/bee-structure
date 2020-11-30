package app

type LanguageExt struct {
	Id         int64  `orm:"column(id)" json:"id"`
	Language   string `orm:"column(language)" json:"language"`
	LanguageEn string `orm:"column(language_en)" json:"language_en"`
	LocalCode  string `orm:"column(local_code)" json:"local_code"`
	Icon       string `orm:"column(icon)" json:"icon"`
}

type LanguageWordExt struct {
	Id          int64  `orm:"column(id)" json:"id"`
	LanguageKey string `orm:"column(language_key)" json:"language_key"`
	ZhCN        string `orm:"column(zh_cn)" json:"zh_cn"` // 中文简体
	ZhHk        string `orm:"column(zh_hk)" json:"zh_hk"` // 中文繁体
	EnUs        string `orm:"column(en_us)" json:"en_us"` //英文
	KoKr        string `orm:"column(ko_kr)" json:"ko_kr"` //韩文
	JaJp        string `orm:"column(ja_jp)" json:"ja_jp"` //日文
}

type LanguageWordSimple struct {
	LanguageKey string `json:"language_key"`
	ZhCN        string `json:"zh_cn"` // 中文简体
	ZhHk        string `json:"zh_hk"` // 中文繁体
	EnUs        string `json:"en_us"` //英文
	KoKr        string `json:"ko_kr"` //韩文
	JaJp        string `json:"ja_jp"` //日文
}

type CountryExt struct {
	Id          int64     `json:"id"`
	Country     string    `json:"country"`
	CountryEn   string    `json:"country_en"`
	MobileCode  string    `json:"mobile_code"`
}