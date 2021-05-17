package main

type AuthMethod int

const (
	AuthMethodNone AuthMethod = iota
	AuthMethodBasic
	AuthMethodBearer
)

type AuthContext struct {
	Method AuthMethod
	*Account
	Session *Session
}
