package lib

/*
https://godoc.org/github.com/go-sql-driver/mysql
https://github.com/go-sql-driver/mysql#usage
https://github.com/go-sql-driver/mysql/wiki/Examples

user@unix(/path/to/socket)/dbname?charset=utf8
user:password@tcp(localhost:5555)/dbname?charset=utf8
user:password@/dbname
user:password@tcp([de:ad:be:ef::ca:fe]:80)/dbname
*/

import (
	"fmt"
	"log"
)

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	. "github.com/ubrabbit/go-server/common"
)

var (
	g_MysqlDB *sql.DB = nil
)

func checkDBConn(conn *sql.DB) {
	if conn == nil {
		log.Fatalf("DB Conn %v is not inited!!!!!", conn)
	}
}

func InitMysql(host string, port int, dbname string, username string, password string) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", username, password, host, port, dbname)
	fmt.Println("InitMysql ", dataSourceName)
	db, err := sql.Open("mysql", dataSourceName)
	CheckFatal(err)

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	err = db.Ping()
	CheckFatal(err)

	g_MysqlDB = db
	fmt.Println(fmt.Sprintf("Connect Mysql %s Succ", host))
	return db, err
}

func CloseMysql() {
	if g_MysqlDB != nil {
		g_MysqlDB.Close()
	}
}

func MysqlQuery(sql_stmt interface{}, arg ...interface{}) []map[string]interface{} {
	checkDBConn(g_MysqlDB)

	var rows *sql.Rows = nil
	var err error = nil

	switch st := sql_stmt.(type) {
	case *sql.Stmt:
		rows, err = st.Query(arg...)
	case string:
		rows, err = g_MysqlDB.Query(st, arg...)
	default:
		log.Fatalf("MysqlQuery error stmt: %v", sql_stmt)
	}
	CheckFatal(err)

	columns, err := rows.Columns()
	CheckFatal(err)

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	data := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		CheckFatal(err)

		record := make(map[string]interface{})
		for i, element := range values {
			switch value := element.(type) {
			case []byte:
				record[columns[i]] = string(value)
			default:
				record[columns[i]] = value
			}
		}
		fmt.Println(record)
		data = append(data, record)
	}

	fmt.Println("data is  ", data)
	return data
}

func execStmt(sql_stmt interface{}, arg ...interface{}) (sql.Result, error) {
	var result sql.Result = nil
	var err error = nil
	switch st := sql_stmt.(type) {
	case *sql.Stmt:
		result, err = st.Exec(arg...)
	case string:
		result, err = g_MysqlDB.Exec(st, arg...)
	default:
		log.Fatalf("execStmt err sql_stmt: %v", sql_stmt)
	}

	return result, err
}

func MysqlUpdate(sql_stmt interface{}, arg ...interface{}) int64 {
	checkDBConn(g_MysqlDB)

	result, err := execStmt(sql_stmt, arg...)
	CheckFatal(err)

	num, err := result.RowsAffected()
	CheckFatal(err)

	return num
}

func MysqlInsert(sql_stmt interface{}, arg ...interface{}) int64 {
	checkDBConn(g_MysqlDB)

	result, err := execStmt(sql_stmt, arg...)
	CheckFatal(err)

	lastid, err := result.LastInsertId()
	CheckFatal(err)
	return lastid
}

func MysqlDelete(sql_stmt interface{}, arg ...interface{}) int64 {
	checkDBConn(g_MysqlDB)

	result, err := execStmt(sql_stmt, arg...)
	CheckFatal(err)

	num, err := result.RowsAffected()
	CheckFatal(err)
	return num
}

func MysqlTransaction() *sql.Tx {
	checkDBConn(g_MysqlDB)

	tranx, err := g_MysqlDB.Begin()
	CheckFatal(err)

	return tranx
}

func MysqlPrepare(sql string) *sql.Stmt {
	checkDBConn(g_MysqlDB)

	stmt, err := g_MysqlDB.Prepare(sql)
	CheckFatal(err)

	return stmt
}

func MysqlSeekDB(sql_stmt string) chan map[string]interface{} {
	checkDBConn(g_MysqlDB)

	sql_stmt = fmt.Sprintf("%s LIMIT ?,?", sql_stmt)
	fmt.Println("sql_stmt is ", sql_stmt)

	stmt, err := g_MysqlDB.Prepare(sql_stmt)
	CheckFatal(err)

	ch := make(chan map[string]interface{}, 100)

	var res_list []map[string]interface{}
	go func() {
		start, seek_cnt := 0, 100
		for {
			res_list = MysqlQuery(stmt, start, seek_cnt)
			if len(res_list) <= 0 {
				break
			}
			for _, record := range res_list {
				ch <- record
			}
			start += seek_cnt
		}
		ch <- nil
	}()

	return ch
}
