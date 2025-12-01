package cron

import "github.com/robfig/cron/v3"

type Scheduler struct {
	// tripService service.TripService
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) CronStart() {
	c := cron.New(cron.WithSeconds())

	// 每秒鐘執行一次
	executeJob(c, Second, []func(){})

	// 每分鐘執行一次
	executeJob(c, Minute, []func(){
		// func(){
		// 	for i := 0; i <= 10; i++{
		// 		println(i)
		// 		time.Sleep(1000 * time.Millisecond)
		// 	}
		// },
	})
	// 每5分鐘執行一次
	executeJob(c, Five, []func(){
		// persist.SaveGpsFmRdsToMongo,
	})

	// 每10分鐘執行一次
	executeJob(c, Ten, []func(){
		// persist.SaveTripFmMongoToMaria,
	})

	// 每15分鐘執行一次
	executeJob(c, Quarter, []func(){})

	// 每半小時執行一次
	executeJob(c, HalfHour, []func(){})

	// 每小時執行一次
	executeJob(c, Hour, []func(){
		// data.GetOnlineDevice,
	})

	// 每半天執行一次（每日00:00, 12:00）
	executeJob(c, HalfDay, []func(){})

	// 每天執行一次（每日00:00）
	executeJob(c, Day, []func(){
		// logafa.StartRotateFile,
	})

	c.Start()
}
