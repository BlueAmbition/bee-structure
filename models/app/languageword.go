package app

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type LanguageWord struct {
	Id          int64     `orm:"column(id)"`
	LanguageKey string    `orm:"column(language_key)"`
	ZhCN        string    `orm:"column(zh_cn)"`  // 中文简体
	ZhHk        string    `orm:"column(zh_hk)"`  // 中文繁体
	EnUs        string    `orm:"column(en_us)"`  //英文
	KoKr        string    `orm:"column(ko_kr)"`  //韩文
	JaJp        string    `orm:"column(ja_jp)"`  //日文
	Status      int       `orm:"column(status)"` // 状态
	Module      string    `orm:"column(module)"`
	Remark      string    `orm:"column(remark)"` // 备注
	CreatedAt   time.Time `orm:"column(created_at);auto_now_add;type(datetime)"`
	UpdatedAt   time.Time `orm:"column(updated_at);auto_now;type(datetime)"`
}

func (m *LanguageWord) TableName() string {
	return "language_word"
}

func init() {
	orm.RegisterModel(new(LanguageWord))
}

//通过模块获取躲过语言词组
func GetLangWordsByModular(modular string) []LanguageWordExt {
	o := orm.NewOrm()
	var list []LanguageWordExt
	sql := "SELECT id,language_key,zh_cn,zh_hk,en_us,ko_kr,ja_jp FROM language_word WHERE `module`=? AND `status`=1;"
	o.Raw(sql, modular).QueryRows(&list)
	return list
}
