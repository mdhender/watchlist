// watchlist - a web server for movie and show lists
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"net/http"
	"strings"
)

// create a new session cookie for the flash message
func (s *Server) addFlash(w http.ResponseWriter, msg string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "wl-flash",
		Path:     "/",
		Value:    msg,
		HttpOnly: true,
		Secure:   true,
	})
}

// delete the session cookie for the flash message
func (s *Server) deleteFlash(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "wl-flash",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	})
}

// fetch the flash message from the request and delete from session
func (s *Server) getFlash(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("wl-flash")
	if err != nil {
		return ""
	}
	s.deleteFlash(w)
	return strings.TrimSpace(cookie.Value)
}
