package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func handleAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	addr := flag.String("addr", ":80", "Address to bind to")
	webPath := flag.String("web-path", "./web/dist", "Path to website assets")
	flag.Parse()

	http.HandleFunc("/api", handleAPI)
	http.Handle("/", http.FileServer(http.Dir(*webPath)))
	log.Printf("Serving requests on %s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
