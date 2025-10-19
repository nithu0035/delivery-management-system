package orders

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
)

type Service struct {
    db *gorm.DB
    rdb *redis.Client
}

func NewService(db *gorm.DB, rdb *redis.Client) *Service {
    return &Service{db: db, rdb: rdb}
}

func (s *Service) Create(ctx context.Context, customerID uint, items interface{}) (*Order, error) {
    b, _ := json.Marshal(items)
    o := &Order{CustomerID: customerID, Items: string(b), Status: "created"}
    if err := s.db.WithContext(ctx).Create(o).Error; err != nil {
        return nil, err
    }
    // cache status
    s.rdb.Set(ctx, fmt.Sprintf("order:%d:status", o.ID), o.Status, 0)
    return o, nil
}

func (s *Service) GetByID(ctx context.Context, id uint) (*Order, error) {
    var o Order
    if err := s.db.WithContext(ctx).First(&o, id).Error; err != nil {
        return nil, err
    }
    return &o, nil
}

func (s *Service) ListForUser(ctx context.Context, userID uint) ([]Order, error) {
    var os []Order
    if err := s.db.WithContext(ctx).Where("customer_id = ?", userID).Find(&os).Error; err != nil {
        return nil, err
    }
    return os, nil
}

func (s *Service) Cancel(ctx context.Context, id uint, userID uint, isAdmin bool) error {
    var o Order
    if err := s.db.WithContext(ctx).First(&o, id).Error; err != nil { return err }
    if o.Cancelled { return errors.New("already cancelled") }
    if !isAdmin && o.CustomerID != userID { return errors.New("not allowed") }
    o.Cancelled = true
    o.Status = "cancelled"
    if err := s.db.WithContext(ctx).Save(&o).Error; err != nil { return err }
    s.rdb.Set(ctx, fmt.Sprintf("order:%d:status", o.ID), o.Status, 0)
    // publish update
    s.rdb.Publish(ctx, "orders:updates", fmt.Sprintf("%d|%s", o.ID, o.Status))
    return nil
}

func (s *Service) ProgressStatus(ctx context.Context, id uint) error {
    var o Order
    if err := s.db.WithContext(ctx).First(&o, id).Error; err != nil { return err }
    if o.Cancelled || o.Status == "delivered" || o.Status == "cancelled" { return nil }
    orderStates := []string{"created", "dispatched", "in_transit", "delivered"}
    var idx int
    for i, st := range orderStates { if st == o.Status { idx = i; break } }
    if idx+1 < len(orderStates) {
        o.Status = orderStates[idx+1]
        o.UpdatedAt = time.Now()
        if err := s.db.WithContext(ctx).Save(&o).Error; err != nil { return err }
        s.rdb.Set(ctx, fmt.Sprintf("order:%d:status", o.ID), o.Status, 0)
        s.rdb.Publish(ctx, "orders:updates", fmt.Sprintf("%d|%s", o.ID, o.Status))
    }
    return nil
}
