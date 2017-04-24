package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func wirteResponse(w http.ResponseWriter, info string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(info))
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

func handleErrPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	size := 5
	totals := errLogTotalCount()

	curPages := r.URL.Query()["curPage"]
	curPageStr := "1"
	if len(curPages) != 0 {
		curPageStr = curPages[0]
	}
	curPage, _ := strconv.Atoi(curPageStr)
	res := Paginator(curPage, size, totals)

	data := make(map[string]interface{})
	data["paginator"] = res
	data["errors"] = errLogQuery(curPage, size)
	// fmt.Println(data)

	t := template.Must(template.New("base").Parse(string(MustAsset("data/template/err_page.html"))))
	if err := t.Execute(w, &data); err != nil {
		log.Printf("Template execute failed, err: %v", err)
		return
	}
}

func handleCleanErr(w http.ResponseWriter, r *http.Request) {
	row := errLogDelete()
	wirteResponse(w, "成功清除"+fmt.Sprintf("%d", row)+"条记录")
}

func handleReportErr(w http.ResponseWriter, r *http.Request) {
	errLogInsert("error")
	wirteResponse(w, "ok")
}

func handleSlaveConn(w http.ResponseWriter, r *http.Request) {
	ip:= strings.Split(r.RemoteAddr, ":")[0]
	salves[ip] = true
	wirteResponse(w, ip)
}