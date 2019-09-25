package db

import (
	"BastionPay/bas-base/config"
	"database/sql"
	"errors"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var (
	Url       = "" //"root@tcp(127.0.0.1:3306)/wallet"
	database  string
	usertable = "user_property"
	db        *sql.DB

	q = map[string]string{}

	accountQ = map[string]string{
		"readUserLevel": "SELECT public_key, source_ip, user_class, level, is_frozen, audite_status,user_name, user_mobile, user_email, country_code, language from %s.%s where user_key = ? limit ? offset ?",
	}

	st = map[string]*sql.Stmt{}
)

// 用户权限信息
type UserLevel struct {
	PublicKey    string `json:"public_key"`
	SourceIP     string `json:"source_ip"`
	UserClass    int    `json:"user_class"`
	Level        int    `json:"level"`
	IsFrozen     rune   `json:"is_frozen"`
	AuditeStatus int    `json:"audite_status"`
	UserName     string
	CountryCode  string
	UserMobile  string
	UserEmail   string
	Language    string
}

func Init(configPath string) {
	var d *sql.DB
	var err error

	err = config.LoadJsonNode(configPath, "db", &Url)
	if err != nil {
		l4g.Crashf("", err)
	}

	parts := strings.Split(Url, "/")
	if len(parts) != 2 {
		l4g.Crashf("Invalid database url")
	}

	if len(parts[1]) == 0 {
		l4g.Crashf("Invalid database name")
	}

	//url := parts[0]
	database = parts[1]

	//if d, err = sql.Open("mysql", url+"/"); err != nil {
	//	l4g.Crashf(err)
	//}
	//if _, err := d.Exec("CREATE DATABASE IF NOT EXISTS " + database); err != nil {
	//	l4g.Crashf(err)
	//}
	//d.Close()
	if d, err = sql.Open("mysql", Url); err != nil {
		l4g.Crashf("", err)
	}
	// http://www.01happy.com/golang-go-sql-drive-mysql-connection-pooling/
	d.SetMaxOpenConns(2000)
	d.SetMaxIdleConns(1000)
	d.Ping()
	//if _, err = d.Exec(accountdb.UsersSchema); err != nil {
	//	l4g.Crash(err)
	//}

	db = d

	for query, statement := range accountQ {
		prepared, err := db.Prepare(fmt.Sprintf(statement, database, usertable))
		if err != nil {
			l4g.Crashf("", err)
		}
		st[query] = prepared
	}
}

func ReadUserLevel(userKey string) (*UserLevel, error) {
	var r *sql.Rows
	var err error

	r, err = st["readUserLevel"].Query(userKey, 1, 0)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if !r.Next() {
		return nil, errors.New("row no next")
	}

	ul := &UserLevel{}
	if err := r.Scan(&ul.PublicKey, &ul.SourceIP, &ul.UserClass, &ul.Level, &ul.IsFrozen, &ul.AuditeStatus, &ul.UserName, &ul.UserMobile, &ul.UserEmail, &ul.CountryCode, &ul.Language); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no rows")
		}
		return nil, err
	}
	if r.Err() != nil {
		return nil, err
	}

	return ul, nil
}
