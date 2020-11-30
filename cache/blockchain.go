package cache

import (
	"bee-structure/functions/redis"
	"bee-structure/models/blockchain"
	"encoding/json"
)

//设置区块链配置表缓存信息
func SetBlockChainConfigs() {
	list := blockchain.GetBlockChainConfigs()
	for _, v := range list {
		value, _ := json.Marshal(v)
		redis.SetSortSet(0, "block_chain_configs", v.Id, string(value))
	}
}

//获取区块链配置表缓存信息
func GetBlockChainConfigs() []blockchain.BlockChainConfig {
	var (
		err            error
		flag           bool
		list           []interface{}
		JsonByte       []byte
		object         blockchain.BlockChainConfig
		blockChainList = make([]blockchain.BlockChainConfig, 0)
	)
	flag = redis.KeyExists(0, "block_chain_configs")
	if !flag {
		SetBlockChainConfigs()
	}
	list = redis.GetSortSet(0, "block_chain_configs", 0, -1, false)
	if list != nil {
		for _, v := range list {
			JsonByte = v.([]uint8)
			err = json.Unmarshal(JsonByte, &object)
			if err == nil {
				blockChainList = append(blockChainList, object)
			}
		}
	}
	return blockChainList
}

//获取单个区块链配置
func GetBlockChainConfig(key string) blockchain.BlockChainConfig {
	list := GetBlockChainConfigs()
	for _, v := range list {
		if v.Key == key {
			return v
		}
	}
	return blockchain.BlockChainConfig{}
}
