package tracking

import (
    "context"
    "fmt"
    "time"

    "delivery/internal/orders"

    "github.com/go-redis/redis/v8"
)

type Tracker struct {
    orders *orders.Service
    rdb *redis.Client
}

func NewTracker(s *orders.Service, rdb *redis.Client) *Tracker {
    return &Tracker{orders: s, rdb: rdb}
}

func (t *Tracker) Run(ctx context.Context, tick time.Duration) {
    ticker := time.NewTicker(tick)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done(): return
        case <-ticker.C:
            // fetch in-progress orders and progress them
            var inProg []orders.Order
            t.orders.db.Where("cancelled = false AND status != 'delivered'").Find(&inProg)
            for _, o := range inProg {
                go func(id uint) {
                    if err := t.orders.ProgressStatus(ctx, id); err != nil {
                        fmt.Println("progress error:", err)
                    }
                }(o.ID)
            }
        }
    }
}
