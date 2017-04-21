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

const (
	errorStr   string = "stack traceback"
	maxErrLine int    = 30
)

var (
	logErrLines int
	logGetErr   bool
	errRecord   bytes.Buffer
)

var db, err = sql.Open("sqlite3", "./errlog.db")

// dbCheckErr   (conn_err)

func errLogInsert(logName string) int64 {
	//插入数据
	stmt, err := db.Prepare("insert into errlog(logname) values(?)")
	dbCheckErr(err)
	defer stmt.Close()
	res, err := stmt.Exec(logName)
	dbCheckErr(err)
	affect, err := res.RowsAffected()
	return affect
}

func errLogQuery(page int, size int) []ErrSt {
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

func errLogTotalCount() int64 {
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

func errLogDelete() int64 {
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

func logWrite(data string) {
	errRecord.WriteString(data + "\n")
}

func logInsertAndClean(hub *Hub) {
	message := errRecord.String()
	errLogInsert(message)
	hub.broadcast <- []byte(message)
	errRecord = bytes.Buffer{}
}

func errLogTail(hub *Hub) {
	t, err := tail.TailFile(opts.Log, tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}})
	if err != nil {
		log.Printf("tail file failed, err: %v", err)
		return
	}
	for line := range t.Lines {
		if strings.Contains(line.Text, errorStr) {
			if logErrLines != 0 {
				logInsertAndClean(hub)
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

		logInsertAndClean(hub)
		logErrLines = 0
		logGetErr = false

	}
}
