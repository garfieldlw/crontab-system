package main

import (
	"context"
	"github.com/garfieldlw/crontab-system/cron/master/service/health"
	"github.com/garfieldlw/crontab-system/cron/master/service/init-cron"
	"github.com/garfieldlw/crontab-system/cron/master/service/schedule"
)

func main() {
	ctx := context.Background()
	go init_cron.Register(ctx)
	init_cron.Flows(ctx)
	health.CheckHealth(ctx)
	schedule.DoFlow(ctx)
}
