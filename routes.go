// watchlist - a web server for movie and show lists
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"github.com/mdhender/watchlist/way"
	"net/http"
)

func (s *Server) routes() http.Handler {
	r := way.NewRouter()

	// create routes
	r.HandleFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		user := s.currentUser(r)
		if !user.IsAuthenticated() {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, "/shows", http.StatusSeeOther)
	})
	r.HandleFunc("GET", "/help", func(w http.ResponseWriter, r *http.Request) {
		user := s.currentUser(r)
		if !user.IsAuthenticated() {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintf(w, "<p>to be done later</p>")
	})
	r.HandleFunc("GET", "/reload", func(w http.ResponseWriter, r *http.Request) {
		user := s.currentUser(r)
		if !user.IsAuthenticated() {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		s.reloadAll()
		http.Redirect(w, r, "/shows", http.StatusSeeOther)
	})
	r.HandleFunc("GET", "/settings", func(w http.ResponseWriter, r *http.Request) {
		user := s.currentUser(r)
		if !user.IsAuthenticated() {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintf(w, "<p>to be done later</p>")
	})
	r.HandleFunc("GET", "/search", s.getSearch())
	r.HandleFunc("GET", "/shows", s.getShows())
	r.HandleFunc("DELETE", "/shows/:id", s.deleteShowsById())
	r.HandleFunc("GET", "/shows/:id", s.getShowsById())
	r.HandleFunc("GET", "/signin/:id", s.signin())
	r.HandleFunc("GET", "/signout", s.signout())

	// all other routes to be served as files requiring authentication
	fs := http.FileServer(http.Dir("../web"))
	r.HandleFunc("GET", "/...", func(w http.ResponseWriter, r *http.Request) {
		user := s.currentUser(r)
		if !user.IsAuthenticated() {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fs.ServeHTTP(w, r)
	})

	return r
}
