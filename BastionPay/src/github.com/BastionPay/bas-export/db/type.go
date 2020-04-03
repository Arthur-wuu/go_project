package db

type DbOptions struct {
	Host        string
	Port        string
	User        string
	Pass        string
	DbName      string
	MaxIdleConn int
	MaxOpenConn int
}
