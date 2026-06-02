package domain

type Province struct {
	ID           int    `gorm:"primaryKey;autoIncrement:false" json:"id"`
	ProvinceName string `gorm:"type:varchar(255);not null" json:"province_name"`

	Cities      []City        `gorm:"foreignKey:ProvinceID" json:"cities,omitempty"`
	Stores      []Store       `gorm:"foreignKey:ProvinceID" json:"stores,omitempty"`
	UserAddress []UserAddress `gorm:"foreignKey:ProvinceID" json:"user_addresses,omitempty"`
}

func (Province) TableName() string {
	return "province"
}

type City struct {
	ID         int    `gorm:"primaryKey;autoIncrement:false" json:"id"`
	ProvinceID int    `gorm:"not null;index" json:"province_id"`
	CityName   string `gorm:"type:varchar(255);not null" json:"city_name"`
	Type       string `gorm:"type:varchar(100);not null" json:"type"`

	Province Province      `gorm:"foreignKey:ProvinceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"province"`
	Stores   []Store       `gorm:"foreignKey:CityID" json:"stores,omitempty"`
	Address  []UserAddress `gorm:"foreignKey:CityID" json:"addresses,omitempty"`
}

func (City) TableName() string {
	return "city"
}
