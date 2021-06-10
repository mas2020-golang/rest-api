package utils

import "github.com/google/uuid"

var Server = ServerT{}

type ServerT struct {
	TokenPwd string // password for the algo to sign the token
}

// GeneratePwd creates a random password to use for the jwt signature
func (s *ServerT) GeneratePwd() {
	s.TokenPwd = uuid.NewString()
}
