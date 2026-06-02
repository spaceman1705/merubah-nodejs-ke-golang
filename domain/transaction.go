package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusWaitingPayment      OrderStatus = "WAITING_PAYMENT"
	OrderStatusWaitingConfirmation OrderStatus = "WAITING_CONFIRMATION"
	OrderStatusConfirmed           OrderStatus = "CONFIRMED"
	OrderStatusCancelled           OrderStatus = "CANCELLED"
	OrderStatusPrescribed          OrderStatus = "PRESCRIBED"
	OrderStatusShipped             OrderStatus = "SHIPPED"
	OrderStatusDelivered           OrderStatus = "DELIVERED"
)

type StockAction string

const (
	StockActionIn  StockAction = "IN"
	StockActionOut StockAction = "OUT"
)

type Cart struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ProductID string `gorm:"type:uuid;not null" json:"product_id"`
	StoreID   string `gorm:"type:uuid;not null" json:"store_id"`
	UserID    string `gorm:"type:uuid;not null" json:"user_id"`
	Quantity  int    `gorm:"default:1" json:"quantity"`

	Product *Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product,omitempty"`
	Store   *Store   `gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"store,omitempty"`
	User    *User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

func (Cart) TableName() string { return "carts" }
func (c *Cart) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return
}

type Order struct {
	ID                string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrderNumber       string      `gorm:"type:varchar(100);uniqueIndex;not null" json:"order_number"`
	Subtotal          float64     `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	DiscountAmount    float64     `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	ShippingCost      float64     `gorm:"type:decimal(10,2);not null" json:"shipping_cost"`
	ShippingDiscount  float64     `gorm:"type:decimal(10,2);default:0" json:"shipping_discount"`
	TotalAmount       float64     `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	EstimatedDelivery *string     `gorm:"type:varchar(100)" json:"estimated_delivery"`
	VoucherCodeUsed   *string     `gorm:"type:varchar(100)" json:"voucher_code_used"`
	Status            OrderStatus `gorm:"type:varchar(50);not null" json:"status"`
	PaymentProof      *string     `gorm:"type:text" json:"payment_proof"`

	PaymentConfirmedAt *time.Time `json:"payment_confirmed_at"`
	ShippedAt          *time.Time `json:"shipped_at"`
	ConfirmedAt        *time.Time `json:"confirmed_at"`
	CancelledAt        *time.Time `json:"cancelled_at"`
	CancellationReason *string    `gorm:"type:text" json:"cancellation_reason"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	UserID        string `gorm:"type:uuid;not null" json:"user_id"`
	UserAddressID string `gorm:"type:uuid;not null" json:"user_address_id"`
	StoreID       string `gorm:"type:uuid;not null;index" json:"store_id"`

	User        *User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user,omitempty"`
	UserAddress *UserAddress `gorm:"foreignKey:UserAddressID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user_address,omitempty"`
	Store       *Store       `gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"store,omitempty"`
	OrderItems  []OrderItem  `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

func (Order) TableName() string { return "orders" }
func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == "" {
		o.ID = uuid.NewString()
	}
	return
}

type OrderItem struct {
	ID             string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrderID        string  `gorm:"type:uuid;not null" json:"order_id"`
	ProductID      string  `gorm:"type:uuid;not null" json:"product_id"`
	Quantity       int     `gorm:"not null" json:"quantity"`
	Price          float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	DiscountAmount float64 `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	Subtotal       float64 `gorm:"type:decimal(10,2);not null" json:"subtotal"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Order   *Order   `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"product,omitempty"`
}

func (OrderItem) TableName() string { return "order_items" }
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	if oi.ID == "" {
		oi.ID = uuid.NewString()
	}
	return
}

type StockJournal struct {
	ID             string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ProductStockID string      `gorm:"type:uuid;not null;index" json:"product_stock_id"`
	Action         StockAction `gorm:"type:varchar(20);not null" json:"action"`
	Quantity       int         `gorm:"not null" json:"quantity"`
	Note           *string     `gorm:"type:text" json:"note"`

	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`

	ProductStock *ProductStock `gorm:"foreignKey:ProductStockID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product_stock,omitempty"`
}

func (StockJournal) TableName() string { return "stock_journals" }
func (sj *StockJournal) BeforeCreate(tx *gorm.DB) (err error) {
	if sj.ID == "" {
		sj.ID = uuid.NewString()
	}
	return
}
