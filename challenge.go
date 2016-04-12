// Program for Curbside Job Challenge posted at https://shopcurbside.com/jobs/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	challengeURL = "http://challenge.shopcurbside.com/"
	sessionPath  = "get-session"
	startID      = "startx"
)

type CurbsideChallengeRsp struct {
	Depth    int         `json:"depth"`
	ID       string      `json:"id"`
	Message  string      `json:"message"`
	Secret   string      `json:"secret"`
	ErrorMsg string      `json:"error"`
	Next     interface{} `json:"next"` // A []string or string
}

var cachedSessionID string

func main() {
	cacheSessionID()
	challenge(startID)
}

func challenge(id string) {
	s, err := queryChallengeServer(cachedSessionID, id)
	if err != nil {
		log.Fatal("Query Challenge Server failed:", err)
	}
	switch {
	case s.Secret != "":
		fmt.Print(s.Secret)
	case s.Next != nil:
		switch v := s.Next.(type) {
		case string:
			challenge(v)
		case []interface{}:
			for _, id := range v {
				challenge(id.(string))
			}
		}
	case s.ErrorMsg != "":
		fmt.Printf("%#v\n", s)
		log.Fatal(s.ErrorMsg)
		cacheSessionID()
		challenge(id)
	}
}

func queryChallengeServer(session, id string) (*CurbsideChallengeRsp, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", challengeURL+id, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Session", session)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	rsp := new(CurbsideChallengeRsp)
	if err := json.Unmarshal(body, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func cacheSessionID() {
	resp, err := http.Get(challengeURL + sessionPath)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	cachedSessionID = string(body)
}
