// watchlist - a web server for movie and show lists
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

// Package main implements a web server for shared movie and show lists.
package main

import (
	"log"
	"net"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)

	host, port := "", "8080"

	var err error

	// create a new http server with good values for timeouts and transports
	s := &Server{}
	s.Addr = net.JoinHostPort(host, port)
	s.IdleTimeout = 10 * time.Second
	s.ReadTimeout = 2 * time.Second
	s.WriteTimeout = 2 * time.Second

	if s.stores.accounts, err = NewAccountsStore("accounts"); err != nil {
		log.Fatal(err)
	} else if s.stores.shows, err = NewShowsStore("shows"); err != nil {
		log.Fatal(err)
	}
	s.Handler = s.routes()

	defer func(started time.Time) {
		log.Printf("[main] elapsed time %v\n", time.Now().Sub(started))
	}(time.Now())

	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
