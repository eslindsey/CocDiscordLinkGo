package cocdiscordlink

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Session struct{
	Token    string
	username string
	password string
}

type TokenResponse struct {
	Token string `json:"token"`
}

const (
	BaseUrl  = "https://cocdiscordlink.azurewebsites.net/api"
	LoginUrl = BaseUrl + "/login"
)

func New(u, p string) (s *Session, err error) {
	s = &Session{
		username: u,
		password: p,
	}
	err = s.authorize()
	return
}

func (s *Session) authorize() (err error) {
	resp, err := http.Post(
		LoginUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf(`{"username":"%s","password":"%s"}`, s.username, s.password)),
	)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var r TokenResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return
	}
	s.Token = r.Token
	return
}

