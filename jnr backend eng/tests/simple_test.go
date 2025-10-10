package tests

import (
    "context"
    "os"
    "testing"

    "delivery/internal/users"
    "delivery/internal/orders"

    "github.com/go-redis/redis/v8"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setup() (*users.Service, *orders.Service, func()) {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&users.User{}, &orders.Order{})
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    userSvc := users.NewService(db)
    orderSvc := orders.NewService(db, rdb)
    return userSvc, orderSvc, func() { rdb.Close(); os.RemoveAll("/tmp/testdb") }
}

func TestCreateAndProgress(t *testing.T) {
    userSvc, orderSvc, teardown := setup()
    defer teardown()
    ctx := context.Background()
    u, err := userSvc.Create(ctx, "a@b.com", "pass", "customer")
    if err != nil { t.Fatal(err) }
    o, err := orderSvc.Create(ctx, u.ID, map[string]interface{}{"item":"book"})
    if err != nil { t.Fatal(err) }
    // progress once
    if err := orderSvc.ProgressStatus(ctx, o.ID); err != nil { t.Fatal(err) }
    got, _ := orderSvc.GetByID(ctx, o.ID)
    if got.Status == "created" { t.Fatalf("status not progressed: %s", got.Status) }
}
