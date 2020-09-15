package main

import (
	"github.com/garfieldlw/crontab-system/cron/worker/service/init-cron"
	"github.com/garfieldlw/crontab-system/cron/worker/service/job"
	"context"
)

func main() {
	ctx := context.Background()
	init_cron.Worker(ctx)
	go init_cron.Register(ctx)
	init_cron.Jobs(ctx)

	job.DoJobs(ctx)
}
