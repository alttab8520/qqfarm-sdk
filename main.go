package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/alttab8520/qqfarm-sdk/internal/api"
	"github.com/alttab8520/qqfarm-sdk/internal/version"
)

func main() {
	host := flag.String("host", "127.0.0.1", "listen host")
	port := flag.String("port", "8765", "listen port")
	flag.Parse()

	addr := *host + ":" + *port
	log.Printf("qqfarm-sdk %s http://%s/docs", version.Version, addr)
	if err := http.ListenAndServe(addr, api.NewMux(api.NewHub(nil))); err != nil {
		log.Fatal(err)
	}
}
