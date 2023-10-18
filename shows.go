// watchlist - a web server for movie and show lists
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

type Shows struct {
	sync.Mutex
	Shows map[string]*Show `json:"shows"`
}

type Show struct {
	Id      string `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Year    int    `json:"year,omitempty"`
	KindOf  string `json:"filmOrTv,omitempty"`
	Genre   string `json:"genre,omitempty"`
	ImdbUrl string `json:"imdbUrl,omitempty"`
}

func NewShowsStore(path string) (*Shows, error) {
	s := &Shows{Shows: make(map[string]*Show)}
	return s, s.Load(path)
}

func (s *Shows) DeleteById(id string) bool {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.Shows[id]; !ok {
		return false
	}
	delete(s.Shows, id)
	return true
}

func (s *Shows) FetchAll() ([]Show, bool) {
	s.Lock()
	defer s.Unlock()
	var shows []Show
	for _, show := range s.Shows {
		shows = append(shows, *show)
	}
	sort.Slice(shows, func(i, j int) bool {
		return shows[i].Title < shows[j].Title
	})
	return shows, true
}

func (s *Shows) FetchById(id string) (Show, bool) {
	s.Lock()
	defer s.Unlock()
	show, ok := s.Shows[id]
	if !ok {
		return Show{}, false
	}
	return *show, true
}

func (s *Shows) Load(path string) error {
	s.Lock()
	defer s.Unlock()
	if data, err := os.ReadFile(path + ".json"); err != nil {
		return err
	} else if err = json.Unmarshal(data, s); err != nil {
		return err
	}
	return nil
}

func (s *Shows) Save(path string) error {
	s.Lock()
	defer s.Unlock()
	bak := time.Now().UTC().Format("2006.01.02.15.04.05")
	if data, err := json.MarshalIndent(s, "", "  "); err != nil {
		return err
	} else if err = os.WriteFile(path+"."+bak+".json", data, 0644); err != nil {
		return fmt.Errorf("backup: %w", err)
	} else if err = os.WriteFile(path+".json", data, 0644); err != nil {
		return fmt.Errorf("actual: %w", err)
	}
	return nil
}
