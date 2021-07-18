package utils

import (
	"github.com/google/uuid"
	"github.com/mas2020-golang/goutils/output"
	"os"
)

var Server = ServerT{}

type ServerT struct {
	Logging struct{
		Level int `yaml: level`
	} `yaml: logging`
	TokenPwd string // password for the algo to sign the token
}

// GeneratePwd creates a random password to use for the jwt signature
func (s *ServerT) GeneratePwd() {
	// check for the env variable
	if len(os.Getenv("APP_JWTPWD")) > 0 {
		output.InfoLog("", "use APP_JWTPWD variable as a valid JWT HS256 signing password")
		s.TokenPwd = os.Getenv("APP_JWTPWD")
	} else {
		output.InfoLog("", "use a random PWD as a valid JWT HS256 signing password")
		s.TokenPwd = uuid.NewString()
	}
}
