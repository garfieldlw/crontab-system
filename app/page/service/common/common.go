package common

type PaginationInputModel struct {
	SortValue string `json:"sort_value"` //排序方式，id desc
	Offset    int64  `json:"offset"`     //跳过的数据条数
	Limit     int64  `json:"limit"`      //单页条数
}

type PaginationOutputModel struct {
	Total  int64 `json:"total"`  //总条数
	Offset int64 `json:"offset"` //跳过的数据条数，同传入数据
	Limit  int64 `json:"limit"`  //单页条数，同传入数据
}
