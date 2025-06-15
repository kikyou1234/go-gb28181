package model

type TimeItem struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type DateItem struct {
	Date  string     `json:"date"`
	Items []TimeItem `json:"items"`
}

type DataContent struct {
	DayNum  int        `json:"daynum"`
	TimeNum int        `json:"timenum"`
	List    []DateItem `json:"list"`
}
