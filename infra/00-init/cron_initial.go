package initial

import "PetTrack/infra/00-core/cron"

func InitCron(s *Services) {
	cron.NewScheduler(
		s.Trip,
	).CronStart()
}
