package tool

import "github.com/astaxie/beego/orm"

type Pager struct {
	TotalCount int64        `json:"total_count"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	List       []orm.Params `json:"list"`
}

type Row struct {
	RowCount int64 `json:"row_count"`
}
