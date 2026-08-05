package dto

type CreateAddressRequest struct {
	FirstName     string  `json:"firstName" validate:"required"`
	LastName      string  `json:"lastName" validate:"required"`
	Address       string  `json:"address" validate:"required"`
	ProvinceID    int     `json:"provinceId" validate:"required"`
	CityID        int     `json:"cityId" validate:"required"`
	PostalCode    *string `json:"postalCode"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	IsMainAddress bool    `json:"isMainAddress"`
}
