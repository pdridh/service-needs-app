package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {

	config := struct {
		Host string
		Port string
	}{
		Host: "localhost",
		Port: "8080",
	} //TODO move this to some kind of config loading thingy pkg

	srv := http.NewServeMux()
	srv.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	httpServer := &http.Server{
		Addr:         net.JoinHostPort(config.Host, config.Port),
		Handler:      srv,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listening on %s\n", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
}
