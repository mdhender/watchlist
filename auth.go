// watchlist - a web server for movie and show lists
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/mdhender/watchlist/way"
	"net/http"
	"time"
)

type User struct {
	Id string
}

func (u User) IsAuthenticated() bool {
	return u.Id != ""
}

func (s *Server) currentUser(r *http.Request) User {
	cookie, err := r.Cookie("watchlist")
	if err != nil {
		return User{}
	} else if cookie.Value == "" {
		return User{}
	}
	user, ok := s.fetchUserById(cookie.Value)
	if !ok {
		return User{}
	}
	return user
}

func (s *Server) signin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// delete the session cookie, just in case
		http.SetCookie(w, &http.Cookie{
			Name:     "watchlist",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
		})

		id := way.Param(r.Context(), "id")
		user, ok := s.fetchUserById(id)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// create a new session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "watchlist",
			Path:     "/",
			Value:    user.Id,
			Expires:  time.Now().Add(2 * 7 * 24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
		})

		http.Redirect(w, r, "/shows", http.StatusSeeOther)
	}
}

func (s *Server) signout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "watchlist",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
