package job

import (
	eventsvc "booking-svc/internal/service/event"

	"github.com/robfig/cron"
	"go.uber.org/zap"
)

type cancelBookingJob struct {
	scheduler *cron.Cron
	service   *eventsvc.EventService
	logger    *zap.Logger
}

func NewCancelBookingJob(svc *eventsvc.EventService, logger *zap.Logger) *cancelBookingJob {
	return &cancelBookingJob{service: svc, logger: logger}
}

func (j *cancelBookingJob) Run() {
	j.scheduler = cron.New()
	j.scheduler.AddFunc("0 * * * *", func() {
		err := j.service.CancelBooking()
		if err != nil {
			j.logger.Error("cancel booking job failed", zap.Error(err))
		}
	})
	j.scheduler.Start()
}
