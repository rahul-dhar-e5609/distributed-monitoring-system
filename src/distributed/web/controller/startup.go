package controller

import (
	"net/http"
)

var ws = newWebsocketController()

// Initialize function is responsible for
// initializing the web app.
func Initialize() {
	registerRoutes()

	registerFileServers()
}

func registerRoutes() {
	http.HandleFunc("/ws", ws.handleMessage)
}

func registerFileServers() {
	http.Handle("/public/",
		http.FileServer(http.Dir("assets")))
	http.Handle("/public/lib/",
		http.StripPrefix("/public/lib/",
			http.FileServer(http.Dir("node_modules"))))
}
