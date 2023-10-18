// watchlist - a web server for movie and show lists
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

func (s *Server) newTemplate(files ...string) (*TemplateHandler, error) {
	t := &TemplateHandler{}
	// load requested templates
	for _, file := range files {
		t.files = append(t.files, filepath.Join("../templates", file+".gohtml"))
	}
	_, err := template.ParseFiles(t.files...)
	if err != nil {
		log.Printf("[templates] error %v\n", err)
	}
	return t, err
}

// TemplateHandler implements a handler for loading, compiling, and serving a template.
// from Mat Ryer's Go Programming Blueprints.
type TemplateHandler struct {
	sync.Mutex
	once    sync.Once
	files   []string           // template files to load
	headers [][]string         // response headers
	t       *template.Template // represents a single template
}

func (t *TemplateHandler) render(w http.ResponseWriter, r *http.Request, data any) {
	buf := &bytes.Buffer{}
	var err error
	t.t, err = template.ParseFiles(t.files...)
	if err != nil {
		log.Printf("%s %s: render: parse: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if err = t.t.ExecuteTemplate(buf, "layout", data); err != nil {
		log.Printf("%s %s: render: execute: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for _, kvp := range t.headers {
		key, value := kvp[0], kvp[1]
		w.Header().Set(key, value)
	}
	_, _ = w.Write(buf.Bytes())
}

func (t *TemplateHandler) renderTemplate(w http.ResponseWriter, r *http.Request, name string, data any) {
	buf := &bytes.Buffer{}
	var err error
	t.t, err = template.ParseFiles(t.files...)
	if err != nil {
		log.Printf("%s %s: render: parse: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if err = t.t.ExecuteTemplate(buf, name, data); err != nil {
		log.Printf("%s %s: render: execute: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for _, kvp := range t.headers {
		key, value := kvp[0], kvp[1]
		w.Header().Set(key, value)
	}
	_, _ = w.Write(buf.Bytes())
}
