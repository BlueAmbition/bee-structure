package app

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type AppVersion struct {
	Id          int64  `orm:"column(id)" json:"id"`
	AppName     string `orm:"column(app_name)" json:"app_name"`
	Version     string `orm:"column(version)" json:"version"`
	Intro       string `orm:"column(intro)" json:"intro"`
	Force       uint   `orm:"column(force)" json:"force"`
	AndroidUrl  string `orm:"column(android_url)" json:"android_url"`
	IOSUrl      string `orm:"column(ios_url)" json:"ios_url"`
	NeedUpgrade bool   `orm:"-" json:"need_upgrade"` //新增无关数据库字段，是否需要升级
	Title       string `orm:"-" json:"title"`        //新增无关数据库字段，更新标题
	Content     string `orm:"-" json:"content"`      //新增无关数据库字段，更新内容
	//WordsRefreshTs int64     `orm:"-" json:"words_refresh_ts"` //新增无关数据库字段，多国语言刷新时间戳
	AndroidApk string    `orm:"-" json:"android_apk"`
	CreatedAt  time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt  time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (m *AppVersion) TableName() string {
	return "app_version"
}

func init() {
	orm.RegisterModel(new(AppVersion))
}

//获取允许注册的国家信息
func GetLatestVersion() AppVersion {
	o := orm.NewOrm()
	var appVersion AppVersion
	o.Raw("SELECT * FROM `app_version` ORDER BY id desc LIMIT 0,1").
		QueryRow(&appVersion)
	return appVersion
}
