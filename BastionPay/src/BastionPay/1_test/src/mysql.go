package src

// 从Mysql中导出数据到CSV文件。

import (
"database/sql"
"encoding/csv"
"fmt"
"os"
_ "github.com/go-sql-driver/mysql"
"flag"
"strings"
)

var (
	tables         = make([]string, 0)
	dataSourceName = ""
)

const (
	driverNameMysql = "mysql"

	helpInfo = `Usage of mysqldataexport:
  -port int
        the port for mysql,default:3306
  -addr string
        the address for mysql,default:127.0.0.1
  -user string
        the username for login mysql,default:root

  -pwd  string
        the password for login mysql by the username,default:root
  -db   string
        the port for me to listen on,default:mydb
  -tables string
        the tables will export data, multi tables separator by comma, default:test_table
    `
)

func init() {

	port := flag.Int("port", 3306, "the port for mysql,default:32085")
	addr := flag.String("addr", "127.0.0.1", "the address for mysql,default:10.146.145.67")
	user := flag.String("user", "root", "the username for login mysql,default:dbuser")
	pwd := flag.String("pwd", "wsy123456008ik,>LO(", "the password for login mysql by the username,default:Admin@123")
	db := flag.String("db", "mydb", "the port for me to listen on,default:auditlogdb")
	tabs := flag.String("tables", "im", "the tables will export data, multi tables separator by comma, default:op_log,sc_log,sys_log")

	//flag.Usage = usage

	flag.Parse()

	tables = append(tables, strings.Split(*tabs, ",")...)

	dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", *user, *pwd, *addr, *port, *db)
}

func main() {

	count := len(tables)
	ch := make(chan bool, count)

	db, err := sql.Open(driverNameMysql, dataSourceName)
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	for _, table := range tables {
		go querySQL(db, table, ch)
	}

	for i := 0; i < count; i++ {
		<-ch
	}
	fmt.Println("Done!")
}

func querySQL(db *sql.DB, table string, ch chan bool) {
	fmt.Println("开始处理：", table)
	rows, err := db.Query(fmt.Sprintf("SELECT * from %s", table))

	if err != nil {
		panic(err)
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	//values：一行的所有值,把每一行的各个字段放到values中，values长度==列数
	values := make([]sql.RawBytes, len(columns))
	// print(len(values))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	//存所有行的内容totalValues
	totalValues := make([][]string, 0)
	for rows.Next() {

		//存每一行的内容
		var s []string

		//把每行的内容添加到scanArgs，也添加到了values
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		for _, v := range values {
			s = append(s, string(v))
			// print(len(s))
		}
		totalValues = append(totalValues, s)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error())
	}
	writeToCSV(table+".csv", columns, totalValues)
	ch <- true
}

//writeToCSV
func writeToCSV(file string, columns []string, totalValues [][]string) {
	f, err := os.Create(file)
	// fmt.Println(columns)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)
	for i, row := range totalValues {
		//第一次写列名+第一行数据
		if i == 0 {
			w.Write(columns)
			w.Write(row)
		} else {
			w.Write(row)
		}
		fmt.Println("处理完毕totalValues：", totalValues)
	}
	w.Flush()
	fmt.Println("处理完毕：", file)
}

//func usage() {
//	fmt.Fprint(os.Stderr, helpInfo)
//	flag.PrintDefaults()
//}

