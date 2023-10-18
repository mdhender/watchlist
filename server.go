// watchlist - a web server for movie and show lists
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	http.Server
	stores struct {
		sync.Mutex
		accounts *Accounts
		shows    *Shows
	}
}

func (s *Server) Serve() error {
	if s.Addr == ":" {
		return fmt.Errorf("missing port")
	} else if s.Handler == nil {
		return fmt.Errorf("missing handler")
	}

	// set up stuff so that we can gracefully shut down the server and application
	serverCh := make(chan struct{})
	go func() {
		log.Printf("[server] serving %q\n", s.Addr)
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[server] exited with: %v", err)
		}
		close(serverCh)
	}()

	// create a catch for signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt
	<-signalCh

	// use the context to shut down the application
	log.Printf("[server] received interrupt, shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("[server] failed to shutdown server: %s", err)
	}

	// If we got this far, it was an interrupt, so don't exit cleanly
	return fmt.Errorf("interrupted and stopped")
}
