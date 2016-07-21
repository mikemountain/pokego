package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/icza/session"
	"golang.org/x/net/publicsuffix"
)

var latitude = 45.4159064
var longitude = -75.6934399
var altitude = 72

type LoginGET struct {
	LT         string `json:"lt"`
	Execuction string `json:"execution"`
}

type LoginPOST struct {
	LoginGET
	EventID  string `json:"_eventId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Pokemon struct {
	Number         string   `json:"Number"`
	Name           string   `json:"Name"`
	Classification string   `json:"Classification"`
	TypeI          []string `json:"Type I"`
	TypeII         []string `json:"Type II"`
	Weaknesses     []string `json:"Weaknesses"`
	FastAttacks    []string `json:"Fast Attack(s)"`
	Weight         string   `json:"Weight"`
	Height         string   `json:"Height"`
	NextEvoReq
	NextEvos []NextEvo
}

type NextEvoReq struct {
	Amount int    `json:"Amount"`
	Name   string `json:"Name"`
}

type NextEvo struct {
	Number string `json:"Number"`
	Name   string `json:"Name"`
}

func main() {
	session.Global.Close()
	session.Global = session.NewCookieManagerOptions(session.NewInMemStore(), &session.CookieMngrOptions{AllowHTTP: true})

	dat, err := ioutil.ReadFile("pokemon.json")
	if err != nil {
		log.Error(err)
	}

	pokemon := []Pokemon{}
	json.Unmarshal(dat, &pokemon)

	username := "PhireSlack"
	password := "mainstreet"

	login(username, password)
}

func login(user string, pass string) {
	loginStr := "https://sso.pokemon.com/sso/login?service=https%3A%2F%2Fsso.pokemon.com%2Fsso%2Foauth2.0%2FcallbackAuthorize"

	loginURL, _ := url.Parse(loginStr)

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", loginStr, nil)
	req.Header.Set("User-Agent", "niantic")
	if err != nil {
		log.Error(err)
	}

	sess := session.Get(req)
	sess = session.NewSession()
	// jar.SetCookies(loginURL, req.Header.Get("Set-Cookie"))
	// session.Add(sess, w?????)

	if sess == nil {
		fmt.Printf("else: %v\n", sess)
		log.Fatal("i don't know")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Jar:       jar,
	}

	response, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer response.Body.Close()

	client.Jar.SetCookies(loginURL, response.Cookies())

	jdata := LoginGET{}

	dat, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
	}

	json.Unmarshal(dat, &jdata)
	// fmt.Println(jdata.LT)

	loginPost := &LoginPOST{
		LoginGET: jdata,
		EventID:  "submit",
		Username: user,
		Password: pass,
	}

	resLoginPost, _ := json.Marshal(loginPost)
	// fmt.Println(string(resLoginPost))

	req, err = http.NewRequest("POST", loginStr, bytes.NewBuffer(resLoginPost))
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("User-Agent", "niantic")
	req.AddCookie(response.Cookies()[0])
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error("Server return non-200 status: %v\n", resp.Status)
	}

	// ticket := resp.Header.Get("Location")
	fmt.Printf("%+v\n", resp)

	// loc, loc2 := resp.Location()

	// fmt.Printf("%+v\n", loc)
	// fmt.Printf("%+v\n", loc2)

}
