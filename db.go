package main

import (
	"bytes"
	"database/sql"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/hpcloud/tail"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
	"time"
)

type ErrSt struct {
	Id      int
	LogName string
	Created time.Time
}

var db, err = sql.Open("sqlite3", "./errlog.db")

// dbCheckErr   (conn_err)

func dbInsert(logName string) int64 {
	//插入数据
	stmt, err := db.Prepare("insert into errlog(logname) values(?)")
	dbCheckErr(err)
	defer stmt.Close()
	res, err := stmt.Exec(logName)
	dbCheckErr(err)
	affect, err := res.RowsAffected()
	return affect
}

func dbQuery(page int, size int) []ErrSt {
	//查询数据
	sql := "select * from errlog order by created desc limit " + fmt.Sprintf("%d", size) + " offset " + fmt.Sprintf("%d", (page-1)*size)
	rows, err := db.Query(sql)
	dbCheckErr(err)
	defer rows.Close()

	datas := []ErrSt{}
	for rows.Next() {
		row := ErrSt{}
		var id int
		var logname string
		var created time.Time
		err = rows.Scan(&id, &logname, &created)
		dbCheckErr(err)

		row.Id = id
		row.LogName = logname
		row.Created = created
		datas = append(datas, row)
	}
	return datas
}

func dbTotalCount() int64 {
	sql := "select count(id) from errlog"
	rows, err := db.Query(sql)
	dbCheckErr(err)
	defer rows.Close()

	for rows.Next() {
		var count int64
		err = rows.Scan(&count)
		dbCheckErr(err)
		return count
	}
	return 0
}

func dbDelete() int64 {
	//删除数据
	stmt, err := db.Prepare("delete from errlog where created<?")
	dbCheckErr(err)
	defer stmt.Close()

	t1 := time.Now().Format("2006-01-02")

	res, err := stmt.Exec(t1)
	dbCheckErr(err)

	affect, err := res.RowsAffected()
	dbCheckErr(err)

	return affect
}

func dbClose() {
	db.Close()
}

func dbCheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	errorStr   string = "stack traceback"
	maxErrLine int    = 30
)

var (
	logErrLines int
	logGetErr   bool
	errRecord   bytes.Buffer
)

func logWrite(data string) {
	errRecord.WriteString(data + "\n")
}

func logInsertAndClean() {
	dbInsert(errRecord.String())
	errRecord = bytes.Buffer{}
}

func mainLogInsert() {
	t, err := tail.TailFile(opts.Log, tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}})
	if err != nil {
		log.Printf("tail file failed, err: %v", err)
		return
	}
	for line := range t.Lines {
		if strings.Contains(line.Text, errorStr) {
			if logErrLines != 0 {
				logInsertAndClean()
			}
			logGetErr = true
			logErrLines = 0

			logWrite(line.Text)
			continue
		}

		if !logGetErr {
			continue
		}

		logWrite(line.Text)
		logErrLines = logErrLines + 1

		if logErrLines < maxErrLine {
			continue
		}

		logInsertAndClean()
		logErrLines = 0
		logGetErr = false

	}
}

// func main() {
// 	// dbDelete()
// 	// count := dbTotalCount()
// 	// fmt.Println(count)
// 	// dbInsert("1111111111122222sdf")

// 	datas := dbQuery(1, 5)
// 	// for i := 0; i < len(datas); i++ {
// 	// 	fmt.Println(datas[i].created)
// 	// }
// 	for k, v := range datas {
// 		fmt.Printf("%v -> %v\n", k, v.LogName)
// 	}
// 	dbClose()
// }
