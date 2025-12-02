package initial

import "PetTrack/core/cron"

func InitCron(s *Services) {
	cron.NewScheduler(
		s.Trip,
	).CronStart()
}
