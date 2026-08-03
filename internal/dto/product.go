package dto

type ProductQueryRequest struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
	Search   string `query:"search"`
	Category string `query:"category"`
	StoreID  string `query:"storeId"`
}
