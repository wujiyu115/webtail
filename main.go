package main

import (
	"net/http"

	"github.com/Akagi201/light"
	log "github.com/Sirupsen/logrus"
	// "golang.org/x/net/websocket"
)

func main() {
	root := light.New()
	hub := newHub()

	go hub.run()
	go errLogTail(hub)

	root.Get("/", handleHome)
	root.Get("/ws_err", func(w http.ResponseWriter, r *http.Request) {
		handleErrWebSocket(hub, w, r)
	})
	root.Get("/realtime_err", handleRealTimeErr)
	root.Get("/err_page", handleErrPage)
	root.Get("/clean_err", handleCleanErr)
	root.Get("/report_err", handleReportErr)

	log.Printf("HTTP listening at: %v", opts.ListenAddr)
	http.ListenAndServe(opts.ListenAddr, root)

}
