package main

import (
	"database/sql"
	"sync"
)

type App struct {
	AppRepo
	*sql.DB
	sessions *sync.Map
}

/*
The session store is very simple and niave it.
It will potentially grow unbounded over time.
Like much of this project an OTS solution should be used.
*/

// Create a new session for the user and return.
func (a App) SessionCreate(username string) *Session {
	sesh := SessionNew(username)
	a.sessions.Store(sesh.Id, sesh)
	return &sesh
}

func (a App) SessionDelete(id string) {
	a.sessions.Delete(id)
}

// Get session by Id and return, or return nil if not found or expired.
func (a App) SessionGetById(id string) *Session {
	entry, ok := a.sessions.Load(id)
	if !ok {
		return nil
	}

	sesh, _ := entry.(Session)

	if sesh.isExpired() {
		a.sessions.Delete(id)
		return nil
	} else {
		return &sesh
	}
}
