package app

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Language struct {
	Id         int64     `orm:"column(id)"`
	Language   string    `orm:"column(language)"`
	LanguageEn string    `orm:"column(language_en)"`
	LocalCode  string    `orm:"column(local_code)"`
	Icon       string    `orm:"column(icon)"`
	Status     uint      `orm:"column(status)"`
	CreatedAt  time.Time `orm:"column(created_at);auto_now_add;type(datetime)"`
	UpdatedAt  time.Time `orm:"column(updated_at);auto_now;type(datetime)"`
}

func (m *Language) TableName() string {
	return "language"
}

func init() {
	orm.RegisterModel(new(Language))
}

//语言列表
func LanguageList() ([]LanguageExt, error) {
	o := orm.NewOrm()
	var list []LanguageExt
	sql := "SELECT `id`,`language`,local_code,icon FROM `language` WHERE `status`=1 ORDER BY `sort` DESC,created_at DESC;"
	num, err := o.Raw(sql).QueryRows(&list)
	if err == nil && num > 0 {
		return list, nil
	}
	return list, err
}

func LanguageDetail(code string) Language {
	o := orm.NewOrm()
	language := Language{LocalCode: code}
	err := o.QueryTable(&language).Filter("local_code", language.LocalCode).One(&language)
	if err != nil {
		return language
	}
	return language
}
