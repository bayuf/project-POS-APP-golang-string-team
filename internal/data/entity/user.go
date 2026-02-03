package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPermissions map[string]bool

// Implementasi interface GORM agar bisa Scan/Value JSON
func (p UserPermissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *UserPermissions) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB")
	}
	return json.Unmarshal(b, &p)
}

const (
	PermissionDashboard = "dashboard"
	PermissionReports   = "reports"
	PermissionInventory = "inventory"
	PermissionOrders    = "orders"
	PermissionCustomers = "customers"
	PermissionSettings  = "settings"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Name         string          `gorm:"size:100;not null"`
	Email        string          `gorm:"size:100;not null;unique"`
	PasswordHash string          `gorm:"type:text;not null;default:''"`
	Phone        string          `gorm:"size:20;not null;default:''"`
	Role         string          `gorm:"size:20;not null"` // superadmin, admin, staff
	Permissions  UserPermissions `gorm:"type:jsonb;default:'{\"dashboard\":true}'"`
	AvatarURL    string          `gorm:"size:255;not null;default:'public/img/user/default.jpg'"`
	BirthDate    time.Time       `gorm:"type:date"`
	Salary       int64           `gorm:"not null;default:0"`

	ShiftStart string `gorm:"type:time"`
	ShiftEnd   string `gorm:"type:time"`

	Address          string  `gorm:"text;not null;default:''"`
	AdditionalDetail *string `gorm:"text;"`

	IsActive  bool           `gorm:"default:true"`
	CreatedAt time.Time      `gorm:"type:timestamptz"`
	UpdatedAt time.Time      `gorm:"type:timestamptz"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index"`

	OTPRequests   []OTPRequest
	Notifications []Notification
	Session       []Session `gorm:"foreignKey:UserID"`
}
