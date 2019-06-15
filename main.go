package main

import (
	"net/http"
	"os"

	"github.com/diegohce/logger"
)

var (
	log *logger.Logger
	srv *http.Server
)

func main() {

	bindAddr := os.Getenv("CHAOSPROXY_BINDADDR")
	if bindAddr == "" {
		bindAddr = ":6667"
	}

	log = logger.New("chaosproxy::")

	if err := loadConfig(); err != nil {
		log.Error().Fatalln(err, "loading config file")
	}

	proxies := createProxies()

	log.Info().Println("Starting chaos-proxy on", bindAddr)

	srv = &http.Server{Addr: bindAddr, Handler: proxies}

	log.Error().Println(srv.ListenAndServe())
	//log.Error().Fatal(http.ListenAndServe(bindAddr, proxies))
}

