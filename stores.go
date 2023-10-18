// watchlist - a web server for movie and show lists
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

func (s *Server) deleteShowById(id string) bool {
	return s.stores.shows.DeleteById(id)
}

func (s *Server) fetchAllShows() ([]Show, bool) {
	return s.stores.shows.FetchAll()
}

func (s *Server) fetchShowById(id string) (Show, bool) {
	return s.stores.shows.FetchById(id)
}

func (s *Server) fetchUserById(id string) (User, bool) {
	return s.stores.accounts.FetchUser(id)
}

func (s *Server) reloadAll() {
	s.stores.Lock()
	defer s.stores.Unlock()
	_ = s.stores.accounts.Load("accounts")
	_ = s.stores.shows.Load("shows")
}

func (s *Server) searchAllShows(q string) ([]Show, bool) {
	return s.fetchAllShows()
}
