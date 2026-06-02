package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
	RoleSuper Role = "super"
)

type User struct {
	ID           string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email        string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	IsVerified   bool   `gorm:"default:false" json:"is_verified"`
	IsActive     bool   `gorm:"default:true" json:"is_active"`
	Role         Role   `gorm:"type:varchar(20);default:'user'" json:"role"`
	ReferralCode string `gorm:"type:varchar(50);uniqueIndex;not null" json:"referral_code"`

	Password  *string `gorm:"type:varchar(255)" json:"-"`
	FirstName *string `gorm:"type:varchar(100)" json:"first_name"`
	LastName  *string `gorm:"type:varchar(100)" json:"last_name"`
	Avatar    *string `gorm:"type:text" json:"avatar"`
	AvatarID  *string `gorm:"type:varchar(100)" json:"avatar_id"`

	ReferrerID    *string `gorm:"type:uuid;index" json:"referrer_id"`
	Referrer      *User   `gorm:"foreignKey:ReferrerID" json:"referrer,omitempty"`
	ReferredUsers []User  `gorm:"foreignKey:ReferrerID" json:"referred_users,omitempty"`

	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID" json:"refresh_tokens,omitempty"`
	Addresses     []UserAddress  `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
	// Carts         []Cart         `gorm:"foreignKey:UserID" json:"carts,omitempty"`
	// Orders        []Order        `gorm:"foreignKey:UserID" json:"orders,omitempty"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

type Token struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Token     string    `gorm:"type:text;uniqueIndex;not null" json:"token"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
}

func (Token) TableName() string {
	return "token"
}

type RefreshToken struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"type:text;uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`

	User *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

type UserAddress struct {
	ID            string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID        string  `gorm:"type:uuid;not null;index" json:"user_id"`
	FirstName     string  `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName      string  `gorm:"type:varchar(100);not null" json:"last_name"`
	Address       string  `gorm:"type:text;not null" json:"address"`
	ProvinceID    int     `gorm:"not null" json:"province_id"`
	CityID        int     `gorm:"not null" json:"city_id"`
	PostalCode    *string `gorm:"type:varchar(20)" json:"postal_code"`
	Latitude      float64 `gorm:"not null" json:"latitude"`
	Longitude     float64 `gorm:"not null" json:"longitude"`
	IsMainAddress bool    `gorm:"default:false" json:"is_main_address"`
	IsActive      bool    `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	User     *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	Province *Province `gorm:"foreignKey:ProvinceID" json:"province,omitempty"`
	UserCity *City     `gorm:"foreignKey:CityID" json:"user_city,omitempty"`
}

func (UserAddress) TableName() string {
	return "user_address"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return
}
