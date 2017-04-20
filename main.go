package main

import (
	"net/http"

	"github.com/Akagi201/light"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
)

func main() {
	go mainLogInsert()
	root := light.New()

	root.Get("/", handleHome)
	root.Get("/follow", websocket.Handler(handleFollow).ServeHTTP)
	root.Get("/realtime_err", handleRealTimeErr)
	root.Get("/err_page", handleErrPage)
	root.Get("/clean_err", handleCleanErr)

	log.Printf("HTTP listening at: %v", opts.ListenAddr)
	http.ListenAndServe(opts.ListenAddr, root)

}
