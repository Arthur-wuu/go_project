package quote

import (
	"flag"
	"fmt"
	"github.com/BastionPay/bas-export/db"
	"github.com/BastionPay/bas-export/config"
	"strconv"
	"sync"
)

type DataMgr struct {
	mSqlDb      db.DbMgr
	mExitCh    chan bool
	mRunFlag   bool
	sync.WaitGroup
	sync.Mutex
}

func (this *DataMgr) Init() (string) {

	this.mExitCh = make(chan bool)
	//err = this.mSqlDb.Init(&db.DbOptions{
	//	Host:        config.GConfig.Db.Host,
	//	Port:        config.GConfig.Db.Port,
	//	User:        config.GConfig.Db.User,
	//	Pass:        config.GConfig.Db.Pwd,
	//	DbName:      config.GConfig.Db.Quote_db,
	//	MaxIdleConn: config.GConfig.Db.Max_idle_conn,
	//	MaxOpenConn: config.GConfig.Db.Max_open_conn,
	//})

	intport,err := strconv.Atoi(config.GConfig.Db.Port)
	if err != nil {
		return "string to int error"
	}
	port := flag.Int("port", intport, "the port for mysql,default:3306")
	addr := flag.String("addr", config.GConfig.Db.Host, "the address for mysql,default:127.0.0.1")
	user := flag.String("user", config.GConfig.Db.User, "the username for login mysql,default:dbuser")
	pwd := flag.String("pwd", config.GConfig.Db.Pwd, "the password for login mysql by the username,default:Admin@123")
	db := flag.String("db", config.GConfig.Db.Db, "the port for me to listen on,default:auditlogdb")
	//tabs := flag.String("tables", "im", "the tables will export data, multi tables separator by comma, default:op_log,sc_log,sys_log")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", *user, *pwd, *addr, *port, *db)

	return dataSourceName
}


func (this *DataMgr) Stop() error {

	this.mSqlDb.Close()
	return nil
}


/***********************内部接口分割线*************************/

