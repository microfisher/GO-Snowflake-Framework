package example

import "math/big"

type ResultObject struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type ResultPaging struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	PageIndex int         `json:"page_index"`
	PageSize  int         `json:"page_size"`
	PageCount int         `json:"page_count"`
}

type Airdrop struct {
	Symbol    string
	Address   string
	Share     *big.Int
	Amount    *big.Int
	Hashcode  string
	Timestamp uint64
}
