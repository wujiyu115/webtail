package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/hpcloud/tail"
	"golang.org/x/net/websocket"
)

func handleErrPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	size := 5
	totals := dbTotalCount()

	curPages := r.URL.Query()["curPage"]
	curPageStr := "1"
	if len(curPages) != 0 {
		curPageStr = curPages[0]
	}
	curPage, _ := strconv.Atoi(curPageStr)
	res := Paginator(curPage, size, totals)

	data := make(map[string]interface{})
	data["paginator"] = res
	data["errors"] = dbQuery(curPage, size)
	// fmt.Println(data)

	t := template.Must(template.New("base").Parse(string(MustAsset("data/template/err_page.html"))))
	if err := t.Execute(w, &data); err != nil {
		log.Printf("Template execute failed, err: %v", err)
		return
	}
}

func handleRealTimeErr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.Must(template.New("base").Parse(string(MustAsset("data/template/realtime_err.html"))))
	v := struct {
		Host string
		Log  string
	}{
		r.Host,
		opts.Log,
	}
	if err := t.Execute(w, &v); err != nil {
		log.Printf("Template execute failed, err: %v", err)
		return
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.Must(template.New("base").Parse(string(MustAsset(opts.Template))))
	v := struct {
		Host string
		Log  string
	}{
		r.Host,
		opts.Log,
	}
	if err := t.Execute(w, &v); err != nil {
		log.Printf("Template execute failed, err: %v", err)
		return
	}
}

func sendWebSocket(ws *websocket.Conn, data string) {
	// log.Println(data)
	ws.Write([]byte(data))

}

func handleFollow(ws *websocket.Conn) {
	t, err := tail.TailFile(opts.Log, tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}})
	if err != nil {
		log.Printf("tail file failed, err: %v", err)
		return
	}
	for line := range t.Lines {
		sendWebSocket(ws, line.Text)
	}
}

func handleCleanErr(w http.ResponseWriter, r *http.Request) {
	row := dbDelete()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("成功清除" + fmt.Sprintf("%d", row) + "条记录"))
}
