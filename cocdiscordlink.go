package cocdiscordlink

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

type LinkResponse struct {
	PlayerTag string `json:"playerTag"`
	DiscordId string `json:"discordId"`
}

const (
	BaseUrl  = "https://cocdiscordlink.azurewebsites.net/api"
	LoginUrl = BaseUrl + "/login"
	LinksUrl = BaseUrl + "/links/"
)

var (
	ErrNoResults      = errors.New("no results")
	ErrTooManyResults = errors.New("too many results")

	client = &http.Client{}
)

func New(u, p string) (s *Session, err error) {
	s = &Session{
		username: u,
		password: p,
	}
	err = s.authorize()
	return
}

func (s *Session) GetLinkFromPlayerTag(tag string) (result string, err error) {
	body, err := s.getSingle(tag)
	var r []LinkResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return
	}
	if len(r) < 1 {
		err = ErrNoResults
	} else if len(r) > 1 {
		err = ErrTooManyResults
	}
	if err != nil {
		return
	}
	result = r[0].DiscordId
	return
}

func (s *Session) GetLinksFromDiscordId(id string) (results []string, err error) {
	body, err := s.getSingle(id)
	var r []LinkResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return
	}
	if len(r) < 1 {
		err = ErrNoResults
		return
	}
	results = make([]string, len(r))
	for i, v := range r {
		results[i] = v.PlayerTag
	}
	return
}

func (s *Session) getSingle(q string) (result []byte, err error) {
	req, err := http.NewRequest("GET", LinksUrl + url.QueryEscape(q), nil)
	req.Header.Add("Authorization", "Bearer " + s.Token)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	result, err = io.ReadAll(resp.Body)
	return
}

func (s *Session) authorize() (err error) {
	resp, err := client.Post(
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

