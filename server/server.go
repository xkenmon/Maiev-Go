package server

import (
	"github.com/xkenmon/maiev/config"
)

func Server() {
	conf := config.GetConfig()
	r := NewRouter()
	host := conf.GetString("server.host")
	if len(host) == 0 {
		host = "127.0.0.1"
	}
	port := conf.GetString("server.port")
	if len(port) == 0 {
		port = "8080"
	}
	r.Run(host + ":" + port)
}
