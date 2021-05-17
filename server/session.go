package main

import (
	"time"
)

const (
	SessionTTL = time.Second * 15
)

type Session struct {
	Id       string
	Username string
	Expiry   time.Time
}

func SessionNew(username string) Session {
	return Session{
		Id:       GetSecureToken(),
		Username: username,
		Expiry:   time.Now().UTC().Add(SessionTTL),
	}
}

func (sesh Session) isExpired() bool {
	return sesh.Expiry.Before(time.Now().UTC())
}
