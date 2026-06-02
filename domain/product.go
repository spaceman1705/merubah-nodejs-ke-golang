package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	ID         string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name       string  `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	IsActive   bool    `gorm:"default:false" json:"is_active"`
	Address    string  `gorm:"type:text;not null" json:"address"`
	Latitude   float64 `gorm:"not null" json:"latitude"`
	Longitude  float64 `gorm:"not null" json:"longitude"`
	CityID     int     `gorm:"not null;index" json:"city_id"`
	ProvinceID int     `gorm:"not null;index" json:"province_id"`
	PostalCode *string `gorm:"type:varchar(20)" json:"postal_code"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Admins   *StoreAdmin    `gorm:"foreignKey:StoreID" json:"admins,omitempty"`
	Products []ProductStock `gorm:"foreignKey:StoreID" json:"products,omitempty"`
}

func (Store) TableName() string {
	return "stores"
}

func (s *Store) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	return
}

type StoreAdmin struct {
	ID      string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID  string `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	StoreID string `gorm:"type:uuid;uniqueIndex;not null" json:"store_id"`

	User  *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	Store *Store `gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"store,omitempty"`
}

func (StoreAdmin) TableName() string {
	return "store_admins"
}

type ProductCategory struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name     string `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	IsActive bool   `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

func (ProductCategory) TableName() string {
	return "product_categories"
}

type Product struct {
	ID          string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string  `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Description *string `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	IsActive    bool    `gorm:"default:true" json:"is_active"`

	CategoryID string `gorm:"type:uuid;not null" json:"category_id"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Category *ProductCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Images   []ProductImage   `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	Stocks   []ProductStock   `gorm:"foreignKey:ProductID" json:"stocks,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

type ProductImage struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ProductID string `gorm:"type:uuid;not null" json:"product_id"`
	ImageURL  string `gorm:"type:text;not null" json:"image_url"`
	ImageKey  string `gorm:"type:varchar(255);not null" json:"image_key"`

	Product *Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product,omitempty"`
}

func (ProductImage) TableName() string {
	return "product_images"
}

type ProductStock struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ProductID string `gorm:"type:uuid;not null;index;uniqueIndex:idx_product_store" json:"product_id"`
	StoreID   string `gorm:"type:uuid;not null;index;uniqueIndex:idx_product_store" json:"store_id"`
	Quantity  int    `gorm:"default:0" json:"quantity"`

	// Relasi
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Store   *Store   `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}

func (ProductStock) TableName() string {
	return "product_stocks"
}
