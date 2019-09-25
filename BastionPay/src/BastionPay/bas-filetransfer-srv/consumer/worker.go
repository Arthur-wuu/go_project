package consumer

import (
	"BastionPay/bas-base/log/zap"
	"BastionPay/bas-filetransfer-srv/db"
	"BastionPay/bas-filetransfer-srv/models"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"reflect"
	"time"
)

var GWorker Worker

type Worker struct {
}

func (this *Worker) Start() {
	go this.run()
}

//目前这种方式 比较慢，应该根据 dbname 多线程执行，这样互不影响
func (this *Worker) run() {
	//workDb := make(map[string] bool)
	for {
		reply, err := db.GRedis.Do("BRPOP", models.EXPORT_List_Key, 60)
		if err != nil {
			log.ZapLog().Error("redis brpop err", zap.Error(err))
			time.Sleep(time.Second * 60)
			continue
		}
		if reply == nil {
			continue
		}
		contents, ok := reply.([]interface{})
		if !ok {
			log.ZapLog().Error("redis type err:" + reflect.TypeOf(reply).String())
			continue
		}
		var content []byte
		log.ZapLog().Info(fmt.Sprintf("%d %s %s", len(contents), contents[0], contents[1]))
		if len(contents) >= 2 {
			content, ok = contents[1].([]byte)
			if !ok {
				log.ZapLog().Error("redis type err:" + reflect.TypeOf(contents[1]).String())
				continue
			}
		}
		task := new(models.TaskExportInfo)
		if err := json.Unmarshal([]byte(content), task); err != nil {
			log.ZapLog().Error("json.Unmarshal err", zap.Error(err), zap.String("taskbody", string(content)))
			continue
		}

		log.ZapLog().Info("start MakeFile", zap.Any("task", *task))
		if err := task.MakeFile(); err != nil {
			log.ZapLog().Error("task MakeFile err", zap.Error(err))
			time.Sleep(time.Second * 10)
			continue
		}
	}
}
