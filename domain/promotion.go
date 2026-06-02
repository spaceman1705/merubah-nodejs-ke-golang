package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscountType string

const (
	DiscountTypePercent DiscountType = "PERCENT"
	DiscountTypeNominal DiscountType = "NOMINAL"
)

type DiscountRule string

const (
	DiscountRuleManual      DiscountRule = "MANUAL"
	DiscountRuleMinPurchase DiscountRule = "MIN_PURCHASE"
	DiscountRuleBogo        DiscountRule = "BOGO"
)

type VoucherType string

const (
	VoucherTypeShopping VoucherType = "shopping"
	VoucherTypeShipping VoucherType = "shipping"
)

type Discount struct {
	ID          string       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	StoreID     string       `gorm:"type:uuid;not null;index" json:"store_id"`
	ProductID   string       `gorm:"type:uuid;not null;index" json:"product_id"`
	Type        DiscountType `gorm:"type:varchar(50);not null" json:"type"`
	Rule        DiscountRule `gorm:"type:varchar(50);not null" json:"rule"`
	Value       float64      `gorm:"type:decimal(10,2);not null" json:"value"`
	MinPurchase *float64     `gorm:"type:decimal(10,2)" json:"min_purchase"`
	MaxDiscount *float64     `gorm:"type:decimal(10,2)" json:"max_discount"`

	StartAt time.Time `gorm:"index:idx_start_end_at" json:"start_at"`
	EndAt   time.Time `gorm:"index:idx_start_end_at" json:"end_at"`

	Store   *Store   `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (Discount) TableName() string { return "discounts" }
func (d *Discount) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	return
}

type DiscountUsage struct {
	ID         string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	DiscountID string  `gorm:"type:uuid;not null;uniqueIndex:idx_discount_order" json:"discount_id"`
	OrderID    string  `gorm:"type:uuid;not null;index;uniqueIndex:idx_discount_order" json:"order_id"`
	Amount     float64 `gorm:"type:decimal(10,2);not null" json:"amount"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Discount *Discount `gorm:"foreignKey:DiscountID" json:"discount,omitempty"`
	Order    *Order    `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

func (DiscountUsage) TableName() string { return "discount_usages" }
func (du *DiscountUsage) BeforeCreate(tx *gorm.DB) (err error) {
	if du.ID == "" {
		du.ID = uuid.NewString()
	}
	return
}

// Model Voucher
type Voucher struct {
	ID            string       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	StoreID       *string      `gorm:"type:uuid" json:"store_id"`
	ProductID     *string      `gorm:"type:uuid" json:"product_id"`
	VoucherCode   string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"voucher_code"`
	VoucherType   VoucherType  `gorm:"type:varchar(50);default:'shopping'" json:"voucher_type"`
	DiscountType  DiscountType `gorm:"type:varchar(50);not null" json:"discount_type"`
	DiscountValue float64      `gorm:"type:decimal(10,2);not null" json:"discount_value"`
	MaxDiscount   *float64     `gorm:"type:decimal(10,2)" json:"max_discount"`
	MinPurchase   *float64     `gorm:"type:decimal(10,2)" json:"min_purchase"`
	Stock         *int         `json:"stock"`

	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Store   *Store   `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (Voucher) TableName() string { return "voucher" }
func (v *Voucher) BeforeCreate(tx *gorm.DB) (err error) {
	if v.ID == "" {
		v.ID = uuid.NewString()
	}
	return
}

type UserVoucher struct {
	ID          string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID      string     `gorm:"type:uuid;not null" json:"user_id"`
	VoucherCode string     `gorm:"type:varchar(100);not null" json:"voucher_code"`
	IsUsed      bool       `gorm:"default:false" json:"is_used"`
	UsedAt      *time.Time `json:"used_at"`
	ObtainedAt  time.Time  `json:"obtained_at"`
	ExpiredAt   time.Time  `json:"expired_at"`

	User    *User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user,omitempty"`
	Voucher *Voucher `gorm:"foreignKey:VoucherCode;references:VoucherCode;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"voucher,omitempty"`
}

func (UserVoucher) TableName() string { return "user_voucher" }
func (uv *UserVoucher) BeforeCreate(tx *gorm.DB) (err error) {
	if uv.ID == "" {
		uv.ID = uuid.NewString()
	}
	return
}
