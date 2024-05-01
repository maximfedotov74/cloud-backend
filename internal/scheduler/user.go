package scheduler

import (
	"context"

	"github.com/go-co-op/gocron"
	"github.com/maximfedotov74/cloud-api/internal/shared/db"
)

type UserScheduler struct {
	cron *gocron.Scheduler
	db   db.PostgresClient
}

func NewUserScheduler(cron *gocron.Scheduler, db db.PostgresClient) *UserScheduler {
	return &UserScheduler{cron: cron, db: db}
}

func (s *UserScheduler) Start() {
	ctx := context.Background()
	go s.a(ctx)
}

// `select a.action_id, am.product_model_id from action as a
//     inner join action_model as am on a.action_id = am.action_id
//     where current_timestamp >= end_date and is_activated = true;`

func (s *UserScheduler) a(ctx context.Context) {}
