package blockchain

//区块链通用返回
type BlockChainRes struct {
	Code   int
	Data   interface{}
	Msg    string
	Status string
}