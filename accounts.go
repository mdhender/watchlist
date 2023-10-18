// watchlist - a web server for movie and show lists
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Accounts struct {
	sync.Mutex
	Accounts map[string]*Account `json:"accounts"`
}

type Account struct {
	Id     string `json:"id"`
	Handle string `json:"handle"`
}

func NewAccountsStore(path string) (*Accounts, error) {
	a := &Accounts{Accounts: make(map[string]*Account)}
	return a, a.Load(path)
}

func (a *Accounts) FetchUser(id string) (User, bool) {
	a.Lock()
	defer a.Unlock()
	acct, ok := a.Accounts[id]
	if !ok {
		return User{}, false
	}
	return User{
		Id: acct.Id,
	}, true
}

func (a *Accounts) Load(path string) error {
	a.Lock()
	defer a.Unlock()
	if data, err := os.ReadFile(path + ".json"); err != nil {
		return err
	} else if err = json.Unmarshal(data, a); err != nil {
		return err
	}
	return nil
}

func (a *Accounts) Save(path string) error {
	a.Lock()
	defer a.Unlock()
	bak := time.Now().UTC().Format("2006.01.02.15.04.05")
	if data, err := json.MarshalIndent(a, "", "  "); err != nil {
		return err
	} else if err = os.WriteFile(path+"."+bak+".json", data, 0644); err != nil {
		return fmt.Errorf("backup: %w", err)
	} else if err = os.WriteFile(path+".json", data, 0644); err != nil {
		return fmt.Errorf("actual: %w", err)
	}
	return nil
}
