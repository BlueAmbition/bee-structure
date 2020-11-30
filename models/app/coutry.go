package app

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type Country struct {
	Id          int64     `orm:"column(id)" json:"id"`
	Country     string    `orm:"column(country)" json:"country"`
	CountryEn   string    `orm:"column(country_en)" json:"country_en"`
	MobileCode  string    `orm:"column(mobile_code)" json:"mobile_code"`
	Description string    `orm:"column(description)" json:"description"`
	AllowReg    uint      `orm:"column(allow_reg)" json:"allow_reg"`
	Status      uint      `orm:"column(status)" json:"status"`
	CreatedAt   time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt   time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (m *Country) TableName() string {
	return "country"
}

func init() {
	orm.RegisterModel(new(Country))
}

//获取允许注册的国家信息
func GetAllowCountry() []Country {
	o := orm.NewOrm()
	var list []Country
	o.Raw("SELECT id,country,country_en,mobile_code FROM `country` WHERE `status`=1 AND allow_reg=1").
		QueryRows(&list)
	return list
}

//获取国家信息
func GetCountryById(countryId int64) Country {
	var country Country
	o := orm.NewOrm()
	sql := fmt.Sprintf("SELECT id,country,country_en,mobile_code FROM `country` WHERE id=%v Limit 1;", countryId)
	o.Raw(sql).QueryRow(&country)
	return country
}
