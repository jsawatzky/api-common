package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jsawatzky/go-common/log"
	"github.com/jsawatzky/go-common/metrics"
)

var (
	logger = log.GetLogger("http")
)

func runServer(server *http.Server, serverErrors chan<- error) {
	logger.Info("Starting server listening on %s", server.Addr)
	switch err := server.ListenAndServe(); err {
	case http.ErrServerClosed:
		logger.Info("Stopped server listening on %s", server.Addr)
		return
	default:
		logger.Error("Error in server listening on %s", server.Addr)
		serverErrors <- err
	}
}

func StartWithMetrics(server *http.Server) error {
	return Start(server, metrics.Server())
}

func Start(servers ...*http.Server) error {
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, syscall.SIGTERM, syscall.SIGINT)

	serverErrors := make(chan error, len(servers))

	var wg sync.WaitGroup
	wg.Add(len(servers))

	for _, srv := range servers {
		go func(wg *sync.WaitGroup, s *http.Server) {
			defer wg.Done()
			runServer(s, serverErrors)
		}(&wg, srv)
	}

	select {
	case err := <-serverErrors:
		return err
	case <-interrupts:
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	shutdownErrors := make(chan error, len(servers))
	wg.Add(len(servers))

	for _, srv := range servers {
		go func(wg *sync.WaitGroup, s *http.Server) {
			defer wg.Done()
			if err := s.Shutdown(ctx); err != nil {
				shutdownErrors <- err
			}
		}(&wg, srv)
	}

	wg.Wait()

	if len(shutdownErrors) > 0 {
		return <-shutdownErrors
	}

	return nil
}
