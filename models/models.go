package models

//while dealing with json, for input and output, it will look like stockid, but for golang backend usage, it will look like StockId
type Stock struct {
	StockId int64  `json:"stockid"`
	Name    string `json:"name"`
	Price   int64  `json:"price"`
	Company string `json:"company"`
}
