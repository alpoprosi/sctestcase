package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sctestcase/cache"
	"sctestcase/server"
	"syscall"
	"time"
)

func main() {
	counter, err := cache.NewCounter("data/")
	if err != nil {
		log.Fatalf("Initializing counter: %v", err)
	}

	srv := server.NewServer(counter)

	httpServer := newHttpServer(srv)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	fmt.Println("[INFO] Server Started...")

	<-done
	fmt.Println("[INFO] Server Stopped...")

	shutdownServer(srv)
}

func newHttpServer(s *server.Server) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.Handle)

	addr := "127.0.0.1:8888"

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	s.AddServer(srv)

	return srv
}

func shutdownServer(s *server.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer func() {
		cancel()
	}()

	if err := s.Shutdown(ctx); err != nil {
		fmt.Println("[ERROR] Shutdown: %w", err)
	}

	fmt.Println("[INFO] Server exited gracefully!")
}
