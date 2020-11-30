package tool

type PagerInterface struct {
	TotalCount int64         `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	List       []interface{} `json:"list"`
}
