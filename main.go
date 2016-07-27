package main

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/golang/geo/s2"
)

var latitude = 45.4159064
var longitude = -75.6934399
var altitude = 72

type Uint64Slice []uint64

func (s Uint64Slice) Len() int           { return len(s) }
func (s Uint64Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s Uint64Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type responseRtmStart struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	Url   string       `json:"url"`
	Self  responseSelf `json:"self"`
}

type responseSelf struct {
	Id string `json:"id"`
}

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

type LoginJar struct {
	Jar map[string][]*http.Cookie
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
	// session.Global.Close()
	// session.Global = session.NewCookieManagerOptions(session.NewInMemStore(), &session.CookieMngrOptions{AllowHTTP: true})

	// dat, err := ioutil.ReadFile("pokemon.json")
	// if err != nil {
	// 	log.Error(err)
	// }

	// pokemon := []Pokemon{}
	// json.Unmarshal(dat, &pokemon)

	// for _, v := range pokemon {
	// 	fmt.Println(v.Name)
	// }

	// username := "PhireSlack"
	// password := "mainstreet"

	// // login(username, password)

	// provider, err := auth.NewProvider("ptc", username, password)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// // Set the coordinates from where you're connecting
	// location := &api.Location{
	// 	Lon: longitude,
	// 	Lat: latitude,
	// 	Alt: 72.0,
	// }

	// // Start new session and connect
	// session := api.NewSession(provider, location, false)
	// session.Init()

	// // Start querying the API
	// player, err := session.GetPlayer()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// out, err := json.Marshal(player)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(string(out))

	// slack part
	// token := "xoxb-63465485520-BOhuvOoBsWw7309OgJiAvppM"
	// url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	// resp, err := http.Get(url)
	// //...
	// body, err := ioutil.ReadAll(resp.Body)
	// //...
	// var respObj responseRtmStart
	// err = json.Unmarshal(body, &respObj)

	// if err != nil {
	// 	log.Error(err)
	// }

	// fmt.Println(respObj)

	getCellIDs(latitude, longitude, 10)
}

func getCellIDs(lat, long float64, rad int) {
	var walk []uint64

	ll := s2.LatLngFromDegrees(lat, long)
	origin := s2.CellIDFromLatLng(ll).Parent(15)
	walk = append(walk, origin.Pos())
	right := origin.Next()
	left := origin.Prev()

	for i := 0; i < rad; i++ {
		walk = append(walk, right.Pos())
		walk = append(walk, left.Pos())
		right = right.Next()
		left = left.Prev()
	}

	sort.Sort(Uint64Slice(walk))

	fmt.Println(walk)
}

// func (s Uint64Slice) Sort() {
// 	sort.Sort(s)
// }
