package pkg

import (
	"github.com/SinarLog/schedulers/utils"
	"github.com/go-co-op/gocron"
)

func GetScheduler() *gocron.Scheduler {
	scheduler := gocron.NewScheduler(utils.CURRENT_LOC).WaitForSchedule()
	scheduler.TagsUnique()
	scheduler.SingletonModeAll()
	return scheduler
}
