// watchlist - a web server for movie and show lists
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"github.com/mdhender/watchlist/way"
	"net/http"
	"strings"
)

type PagePayload struct {
	Title   string
	Flash   string
	Content any
}

type ShowsPayload struct {
	Search string
	Shows  []ShowPayload
}

type ShowPayload struct {
	Search  string
	Id      string
	Title   string
	Year    string
	KindOf  string
	Genre   string
	ImdbUrl string // URL of IMDB entry
}

func (s *Server) deleteShowsById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := way.Param(r.Context(), "id")
		user := s.currentUser(r)
		if !user.IsAuthenticated() {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if ok := s.deleteShowById(id); !ok {
			s.addFlash(w, "Show not found!")
		} else {
			s.addFlash(w, "Show deleted!")
		}
		http.Redirect(w, r, "/shows", http.StatusSeeOther)
	}
}

func (s *Server) getSearch() http.HandlerFunc {
	t, err := s.newTemplate("table-shows")
	if err != nil {
		panic(fmt.Sprintf("[server] getSearch: %v", err))
	}
	t.headers = append(t.headers, []string{"Content-Type", "text/html"})

	return func(w http.ResponseWriter, r *http.Request) {
		user := s.currentUser(r)
		if !user.IsAuthenticated() {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		var content []ShowPayload
		shows, _ := s.searchAllShows(strings.TrimSpace(r.URL.Query().Get("q")))
		for _, show := range shows {
			var yyyy string
			if show.Year != 0 {
				yyyy = fmt.Sprintf("%04d", show.Year)
			}
			content = append(content, ShowPayload{
				Id:      show.Id,
				Title:   show.Title,
				Year:    yyyy,
				KindOf:  show.KindOf,
				Genre:   show.Genre,
				ImdbUrl: show.ImdbUrl,
			})
		}
		t.renderTemplate(w, r, "table-shows", content)
	}
}

func (s *Server) getShows() http.HandlerFunc {
	t, err := s.newTemplate("layout", "search-shows", "table-shows", "shows")
	if err != nil {
		panic(fmt.Sprintf("[server] getShows: %v", err))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		user := s.currentUser(r)
		if !user.IsAuthenticated() {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		payload := PagePayload{Flash: s.getFlash(w, r)}
		var content ShowsPayload
		shows, _ := s.fetchAllShows()
		for _, show := range shows {
			var yyyy string
			if show.Year != 0 {
				yyyy = fmt.Sprintf("%04d", show.Year)
			}
			content.Shows = append(content.Shows, ShowPayload{
				Id:      show.Id,
				Title:   show.Title,
				Year:    yyyy,
				KindOf:  show.KindOf,
				Genre:   show.Genre,
				ImdbUrl: show.ImdbUrl,
			})
		}
		payload.Content = content

		t.render(w, r, payload)
	}
}

func (s *Server) getShowsById() http.HandlerFunc {
	t, err := s.newTemplate("layout", "shows_by_id")
	if err != nil {
		panic(fmt.Sprintf("[server] getShowsById: %v", err))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := way.Param(r.Context(), "id")
		user := s.currentUser(r)
		if !user.IsAuthenticated() {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var payload PagePayload
		show, ok := s.fetchShowById(id)
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var yyyy string
		if show.Year != 0 {
			yyyy = fmt.Sprintf("%04d", show.Year)
		}
		payload.Content = ShowPayload{
			Id:      show.Id,
			Title:   show.Title,
			Year:    yyyy,
			KindOf:  show.KindOf,
			Genre:   show.Genre,
			ImdbUrl: show.ImdbUrl,
		}

		t.render(w, r, payload)
	}
}
