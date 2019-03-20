package controller

import (
	"net/http"
)

func Initialize() {
	registerRoutes()

	registerFileServers()
}

func registerRoutes() {
}

func registerFileServers() {
	http.Handle("/public/",
		http.FileServer(http.Dir("assets")))
	http.Handle("/public/lib/",
		http.StripPrefix("/public/lib/",
			http.FileServer(http.Dir("node_modules"))))
}
