package orders

import (
    "time"
)

type Order struct {
    ID uint `gorm:"primaryKey"`
    CustomerID uint
    Items string `gorm:"type:jsonb"`
    Status string
    Cancelled bool
    CreatedAt time.Time
    UpdatedAt time.Time
}
