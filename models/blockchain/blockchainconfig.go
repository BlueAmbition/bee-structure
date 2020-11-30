package blockchain

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type BlockChainConfig struct {
	Id          int64     `orm:"column(id);"`
	Key         string    `orm:"column(key);"`
	Value       string    `orm:"column(value);"`
	Description string    `orm:"column(description);"`
	Status      int       `orm:"column(status);"`
	CreatedAt   time.Time `orm:"column(created_at);auto_now_add;type(datetime)"`
	UpdatedAt   time.Time `orm:"column(updated_at);auto_now;type(datetime)"`
}

func (m *BlockChainConfig) TableName() string {
	return "block_chain_config"
}

func init() {
	orm.RegisterModel(new(BlockChainConfig))
}

//获取允许显示的币种
func GetBlockChainConfigs() []BlockChainConfig {
	var (
		configList []BlockChainConfig
	)
	o := orm.NewOrm()
	sql := "SELECT * FROM `block_chain_config` WHERE `status`=1;"
	o.Raw(sql).QueryRows(&configList)
	return configList
}
