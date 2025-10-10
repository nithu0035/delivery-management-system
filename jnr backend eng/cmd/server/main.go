package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "delivery/internal/auth"
    "delivery/internal/orders"
    "delivery/internal/tracking"
    "delivery/internal/users"

    "github.com/gorilla/mux"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "github.com/go-redis/redis/v8"
)

func main() {
    dsn := os.Getenv("DATABASE_DSN")
    if dsn == "" {
        dsn = "host=localhost user=demo password=demo dbname=deliverydb port=5432 sslmode=disable"
    }
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect db:", err)
    }

    // Auto migrate (simple for assignment)
    if err := db.AutoMigrate(&users.User{}, &orders.Order{}); err != nil {
        log.Fatal(err)
    }

    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "localhost:6379"
    }
    rdb := redis.NewClient(&redis.Options{Addr: redisAddr})

    // Create services
    userSvc := users.NewService(db)
    orderSvc := orders.NewService(db, rdb)
    authSvc := auth.NewService(os.Getenv("JWT_SECRET"), userSvc)

    // Start background tracker
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    tracker := tracking.NewTracker(orderSvc, rdb)
    go tracker.Run(ctx, 5*time.Second) // progress every 5s

    // Router
    r := mux.NewRouter()
    auth.RegisterRoutes(r, authSvc, userSvc, orderSvc)
    orders.RegisterRoutes(r, authSvc, orderSvc, userSvc)

    srv := &http.Server{
        Addr: ":8080",
        Handler: r,
        ReadTimeout: 15 * time.Second,
        WriteTimeout: 15 * time.Second,
    }

    fmt.Println("Server listening on :8080")
    log.Fatal(srv.ListenAndServe())
}
