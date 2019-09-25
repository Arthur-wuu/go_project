package jobs

import (
	"BastionPay/merchant-workattence-api/models"
	"BastionPay/merchant-workattence-api/modules"
	"github.com/robfig/cron"
)

var Schd Scheduler

type Scheduler struct {
	Cron *cron.Cron
	//Tasks []Task
}

type Task struct {
	TimeStr string
	F       func()
}

func (this *Scheduler) RunCron() {
	awdModels := new(models.AwardRecord)
	rcModels := new(models.RubbishClassify)
	smModels := new(models.StaffMotivation)
	dTalkM := modules.New()
	dTalkRxM := modules.NewRuoXi()
	stallM := new(modules.Staff)
	var Tasks = []Task{
		//{"*/2 * * * * *", func(){fmt.Printf("cron test")}},
		//{"@hourly", awdModels.ResendRecord},
		{"0 */20 * * * *", awdModels.ResendRecord},
		{"10 */20 * * * *", rcModels.ResendRecord},
		{"20 */20 * * * *", smModels.ResendRecord},
		{"@every 1h55m", dTalkM.SetAuthAccess},
		{"*/10 * 0-1,10-12,22-23 * * *", dTalkM.AttenRewardSend},
		{"*/30 * 2-9,13-21 * * *", dTalkM.AttenRewardSend},
		{"*/20,59 * 12-15 * * *", dTalkM.OvertimeAwardSend},
		{"@every 1h55m", dTalkRxM.SetAuthAccess},
		{"*/10 * 0-1,10-12,22-23 * * *", dTalkRxM.AttenRewardSend},
		{"*/30 * 2-9,13-21 * * *", dTalkRxM.AttenRewardSend},
		{"*/20,59 * 12-15 * * *", dTalkRxM.OvertimeAwardSend},
		{"0 0 */6 * * *", stallM.SyncStaffInfo},
	}
	this.Cron = cron.New()

	for _, task := range Tasks {
		this.Cron.AddFunc(task.TimeStr, task.F)
	}

	this.Cron.Start()

	//this.Tasks = Tasks
}
