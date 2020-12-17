//Package main retrieves the gameIds from the NHL stats api from provided date input.
package main

import (
	"time"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"os"
	"fmt"
)

type getGameIds struct {
	Copyright    string `json:"copyright"`
	TotalItems   int    `json:"totalItems"`
	TotalEvents  int    `json:"totalEvents"`
	TotalGames   int    `json:"totalGames"`
	TotalMatches int    `json:"totalMatches"`
	Wait         int    `json:"wait"`
	Dates        []struct {
		Date         string `json:"date"`
		TotalItems   int    `json:"totalItems"`
		TotalEvents  int    `json:"totalEvents"`
		TotalGames   int    `json:"totalGames"`
		TotalMatches int    `json:"totalMatches"`
		Games        []struct {
			GamePk   int       `json:"gamePk"`
			Link     string    `json:"link"`
			GameType string    `json:"gameType"`
			Season   string    `json:"season"`
			GameDate time.Time `json:"gameDate"`
			Status   struct {
				AbstractGameState string `json:"abstractGameState"`
				CodedGameState    string `json:"codedGameState"`
				DetailedState     string `json:"detailedState"`
				StatusCode        string `json:"statusCode"`
				StartTimeTBD      bool   `json:"startTimeTBD"`
			} `json:"status"`
			Teams struct {
				Away struct {
					LeagueRecord struct {
						Wins   int    `json:"wins"`
						Losses int    `json:"losses"`
						Ot     int    `json:"ot"`
						Type   string `json:"type"`
					} `json:"leagueRecord"`
					Score int `json:"score"`
					Team  struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
						Link string `json:"link"`
					} `json:"team"`
				} `json:"away"`
				Home struct {
					LeagueRecord struct {
						Wins   int    `json:"wins"`
						Losses int    `json:"losses"`
						Ot     int    `json:"ot"`
						Type   string `json:"type"`
					} `json:"leagueRecord"`
					Score int `json:"score"`
					Team  struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
						Link string `json:"link"`
					} `json:"team"`
				} `json:"home"`
			} `json:"teams"`
			Venue struct {
				Name string `json:"name"`
				Link string `json:"link"`
			} `json:"venue,omitempty"`
			Content struct {
				Link string `json:"link"`
			} `json:"content"`
			Venue2 struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Link string `json:"link"`
			} `json:"venue,omitempty"`
		} `json:"games"`
		Events  []interface{} `json:"events"`
		Matches []interface{} `json:"matches"`
	} `json:"dates"`
}

func main() {

	// use 20200901 as an example
	// https://statsapi.web.nhl.com/api/v1/schedule?startDate=2020-09-01&endDate=2020-09-01

	if len(os.Args) < 2 {
		fmt.Println("Enter a date in yyyymmdd format")
		os.Exit(1)
	}

	dateArg := os.Args[1]

	if len(dateArg) != 8 {
		fmt.Println("Enter a date in yyyymmdd format")
		os.Exit(1)
	}

	yearArg := dateArg[0:4]
	monthArg := dateArg[4:6]
	dayArg := dateArg[6:]

	ya, err := strconv.Atoi(yearArg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ma, err := strconv.Atoi(monthArg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	da, err := strconv.Atoi(dayArg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	if ya > 2030 || ya < 2010 || ma < 1 || ma > 12 || da < 1 || da > 31 {
		fmt.Println("Enter a date in yyymmdd format")
		os.Exit(1)
	}
	// fmt.Println(ya)
	// fmt.Println(ma)
	// fmt.Println(da)
	// fmt.Printf("yearArg is: %T\n", ya)
	// fmt.Printf("monthArg is: %T\n", ma)
	// fmt.Printf("dayArg is: %T\n", da)

	requestStr := yearArg + "-" + monthArg + "-" + dayArg
	jsonLoc := "https://statsapi.web.nhl.com/api/v1/schedule?startDate=" + requestStr + "&endDate=" + requestStr

	fmt.Println(jsonLoc, "\n")

	// get json

	// get response of request, should be http 200 OK, or error.
	response, err := http.Get(jsonLoc)
	if err != nil {
		panic(err.Error())
	
	}
	//fmt.Print(response)

	// get body of http request (not the response/code)
	// body returns a slice of bytes []byte or []uint8
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
    	panic(err.Error())
	}

	// prints body, a byte slice ([]byte) as a string
	//fmt.Printf("%s\n", body)	
	//fmt.Printf("%T\n", body)

	// slice of byte conversion to string
	// s := string(body)
	// fmt.Println(s)

	var data getGameIds
	err = json.Unmarshal(body, &data)

	if err != nil {
    	panic(err.Error())
	}

	//fmt.Println(data.Dates[0].Games[0].GamePk)
	//fmt.Println(data.Dates[0].Games[1].GamePk)
	//fmt.Println(data.Dates)
	//fmt.Printf("%T\n", data.Dates)
	
	// var gamePk []int

	// for i := 0; i < len(data.Dates[0].Games); i++ {
	// 	item := data.Dates[0].Games[i].GamePk
	// 	//fmt.Println(item)
	// 	gamePk = append(gamePk, item) 
		
	// }

	// //fmt.Println(gamePk)
	// // loop over slice
	// for _, v := range gamePk {
	// 	fmt.Println(v)
	// }

	//var outGames int[]

	//dates := data.Dates
	//fmt.Printf("%T\n", dates)

	//for _, v := range dates {
	//	fmt.Println(k)
	//}
	//fmt.Println(data.Dates)

	//fmt.Println(data.Dates[0].Games[0].Season)
	//fmt.Println(data.Dates[0].Games[0].Status.DetailedState)

	var season []string
	var gamePk []string
	var gameState []string

	for i := 0; i < len(data.Dates[0].Games); i++ {
		s := data.Dates[0].Games[i].Season

		season = append(season, s)

		gpk := strconv.Itoa(data.Dates[0].Games[i].GamePk)
		gpk = gpk[5:]
		gamePk = append(gamePk, gpk)

		ds := data.Dates[0].Games[i].Status.DetailedState

		gameState = append(gameState, ds) 

	}

	//fmt.Println(season)

	//var gamePk []string

	// for i := 0; i < len(data.Dates[0].Games); i++ {
	// 	item := strconv.Itoa(data.Dates[0].Games[i].GamePk)
	// 	//fmt.Println(item)
	// 	item = item[5:]
	// 	//fmt.Println(item)
	// 	gamePk = append(gamePk, item)
	// }

	//fmt.Println(gamePk)

	//var gameState []string

	// for i := 0; i < len(data.Dates[0].Games); i++ {
	// 	item := data.Dates[0].Games[i].Status.DetailedState

	// 	gameState = append(gameState, item)
	// }

	//fmt.Println(gameState)

	//fmt.Println("Length of data.Dates[0].Gameslen", len(data.Dates[0].Games))
	fmt.Println(data.Dates[0].Date, "\n")
	fmt.Println("SEASON      GAMEID   STATE")

	if len(gamePk) == len(gameState) {
        for i := range season {
			fmt.Println(season[i] + "    " + gamePk[i]  + "    " +  gameState[i])
        }
    }
}