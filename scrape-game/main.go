package main

import (
	"time"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strings"
	"strconv"
	"fmt"
	"os"
)


func main() {
	//i, err := toSecs("05:02")
	//if err == nil {
	//	fmt.Println(i)
	//}
	
	// Get user arguments
	seasonArg := os.Args[1]					// Specify 20142015 for the 2014-2015 season
	shortSeasonArg :=  seasonArg[0:4]		// The starting year of the season
	gameArg := os.Args[2]					// Specify a gameId 20100, or a range 20100-20105
	//gameIds = []							// List of gameIds to scrape

	if len(os.Args) < 3 {
		fmt.Println("Enter a date in yyyyyyyy format and gameId")
		os.Exit(1)
	}

	if len(seasonArg) != 8 {
		fmt.Println("Enter a date in yyyyyyyy format, ie: 20192020")
		os.Exit(1)
	}

	if len(gameArg) < 5 {
		fmt.Println("Enter a gameId or gameId range, ie: 20100 or 20100-20105")
		os.Exit(1)
	}
	
	// List of season+gameIds that won't use the json pbp
	//fallbackGameIds = ["20152016-20823"]	
	//fallbackGameIds = []		

	var gameIDs []int

	if !strings.Contains(gameArg, "-") {
		//gameIDs := gameArg

		startID := gameArg[0:5]
		istartID, err := strconv.Atoi(startID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		gameIDs = append(gameIDs, istartID)

		//fmt.Println(gameIDs)
	} else {
		startID := gameArg[0:5]
		endID := gameArg[6:]
		istartID, err := strconv.Atoi(startID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		iendID, err := strconv.Atoi(endID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//fmt.Println(startID)
		//fmt.Println(endID)
		for i := istartID; i <= iendID; i++ {
			//fmt.Println(i)
			gameIDs = append(gameIDs, i)
		}
		//fmt.Println(gameIDs)
	}

	// Scrape data for each game

 	// Converts full team names used in json (e.g., the event team) to json abbreviations (e.g., sjs)
	teamAbbrevs := map[string]string{}

	teamAbbrevs["carolina hurricanes"] = "car"
	teamAbbrevs["columbus blue jackets"] = "cbj"
	teamAbbrevs["new jersey devils"] = "njd"
	teamAbbrevs["new york islanders"] = "nyi"
	teamAbbrevs["new york rangers"] = "nyr"
	teamAbbrevs["philadelphia flyers"] = "phi"
	teamAbbrevs["pittsburgh penguins"] = "pit"
	teamAbbrevs["washington capitals"] = "wsh"
	teamAbbrevs["boston bruins"] = "bos"
	teamAbbrevs["buffalo sabres"] = "buf"
	teamAbbrevs["detroit red wings"] = "det"
	teamAbbrevs["florida panthers"] = "fla"
	teamAbbrevs["montreal canadiens"] = "mtl"
	teamAbbrevs["ottawa senators"] = "ott"
	teamAbbrevs["tampa bay lightning"] = "tbl"
	teamAbbrevs["toronto maple leafs"] = "tor"
	teamAbbrevs["chicago blackhawks"] = "chi"
	teamAbbrevs["colorado avalanche"] = "col"
	teamAbbrevs["dallas stars"] = "dal"
	teamAbbrevs["minnesota wild"] = "min"
	teamAbbrevs["nashville predators"] = "nsh"
	teamAbbrevs["st. louis blues"] = "stl"
	teamAbbrevs["winnipeg jets"] = "wpg"
	teamAbbrevs["anaheim ducks"] = "ana"
	teamAbbrevs["arizona coyotes"] = "ari"
	teamAbbrevs["calgary flames"] = "cgy"
	teamAbbrevs["edmonton oilers"] = "edm"
	teamAbbrevs["los angeles kings"] = "lak"
	teamAbbrevs["san jose sharks"] = "sjs"
	teamAbbrevs["vancouver canucks"] = "van"
	teamAbbrevs["vegas golden knights"] = "vgk"

	//fmt.Println(teamAbbrevs)

	// Situations and stats to record

	var scoreSits = []int{-3, -2, -1, 0, 1, 2, 3}
	var strengthSits = []string{"ownGPulled", "oppGPulled", "sh45", "sh35", "sh34", "pp54", "pp53", "pp43", "ev5", "ev4", "ev3", "other"}
	var teamStats = []string{"toi", "gf", "ga", "sf", "sa", "bsf", "bsa", "msf", "msa", "foWon", "foLost", "ofo", "dfo", "nfo", "penTaken", "penDrawn"}
	//var playerStats = []string{"toi", "ig", "is", "ibs", "ims", "ia1", "ia2", "blocked", "gf", "ga", "sf", "sa", "bsf", "bsa", "msf", "msa", "foWon", "foLost", "ofo", "dfo", "nfo", "penTaken", "penDrawn"}

	// foWon: team won face-offs, individually won face-offs
	// foLost: team lost face-offs, individually lost face-offs
	// ig, is, ibs, ims, ia1, ia2: individual goals, shots, blocked shots, missed shots, primary assists, secondary assists
	// blocked: shots blocked by the individual
	// penDrawn and penTaken don't account for delayed penalties - this is apparent in the team stats:
	//	If teamA draws a penalty but gets a chance to pull their goalie and play for a few seconds before the penalty is actually called, the penDrawn stat will be counted towards the "ownGPulled" situation
	//	Similarly, teamB's penTaken stat will be counted towards the "oppGPulled" situation

	for _, gameID := range gameIDs {
		//fmt.Println(g)

		if gameID < 20000 || gameID >= 40000 {
			fmt.Println("Invalid gameId: ", gameID)
			//continue
			os.Exit(1)
		} else {
			fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
			fmt.Println("Processing game: ", gameID)

		}

		// Download input files

		inDir := "nhl-data/"								// Where the input files are stored
		//outDir := "data-for-db/"							// Where the output files (to be written to database) are stored

		// Input file urls
		// Old shiftcharts API URL
		//shiftJsonUrl = "http://www.nhl.com/stats/rest/shiftcharts?cayenneExp=gameId=" + str(shortSeasonArg) + "0" + str(gameId)
		// New shiftcharts API URL: https://api.nhle.com/stats/rest/en/shiftcharts?cayenneExp=gameId=2019021081
		shiftJsonUrl := "https://api.nhle.com/stats/rest/en/shiftcharts?cayenneExp=gameId=" + shortSeasonArg + "0" + strconv.Itoa(gameID)
		pbpJsonUrl := "https://statsapi.web.nhl.com/api/v1/game/" + shortSeasonArg + "0" + strconv.Itoa(gameID) + "/feed/live"
		fmt.Println("Shift JSON: ", shiftJsonUrl)
		fmt.Println("PBP JSON: ",  pbpJsonUrl)

		// Downloaded input file names
		shiftJson := seasonArg + "-" + strconv.Itoa(gameID) + "-shifts.json"
		pbpJson := seasonArg + "-" + strconv.Itoa(gameID) + "-events.json"

		//fmt.Println(shiftJson, pbpJson)

		var fileNames []string
		var fileUrls []string

		fileNames = append(fileNames, shiftJson, pbpJson)
		fileUrls = append(fileUrls, shiftJsonUrl, pbpJsonUrl)

		//fmt.Println(fileNames, fileUrls)

		for i, f := range fileNames {

			_, err := os.Stat(inDir + f)
			//fmt.Println(res, err)
			if os.IsNotExist(err) {
				fmt.Println("Downloading", f)

				// get json

				// get response of request, should be http 200 OK, or error.
				response, err := http.Get(fileUrls[i])
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
				//_, err = ioutil.WriteFile(fileNames[i], []byte(body), 0644)
				err = os.MkdirAll(inDir, os.ModePerm)
				if err != nil {
					panic(err.Error())
				}

				err = ioutil.WriteFile(inDir + f, []byte(body), 0644)
    			if err != nil {
					panic(err.Error())
    			}
				//fmt.Printf("%s\n", body)
				//fmt.Printf("%T\n", body)

				// slice of byte conversion to string
				// s := string(body)
				// fmt.Println(s)
			} else {
				fmt.Println(f, "already exists.")
			}
			//fmt.Println(i, v)
		}
		fmt.Println("- - - - -")

		// Parse pbpJson

		//var pbpData pbp
		//var result map[string]interface{}
		var pbpData Data

		file, err := ioutil.ReadFile(inDir + pbpJson)
		
		if err != nil {
        	panic(err.Error())
    	}
    	err = json.Unmarshal(file, &pbpData)
    	if err != nil {
        	panic(err.Error())
		}

		//fmt.Println(pbpData)
		
		gameDateTime := pbpData.GameData.DateTime.DateTime
		//gameEndDateTime := pbpData.GameData.DateTime.EndDateTime
		sGameDate := gameDateTime.Format("20060102150405")
		iGameDate, err := strconv.Atoi(sGameDate)
		fmt.Println(iGameDate)
		//fmt.Println(gameDateTime)
		//fmt.Println(gameEndDateTime)
		//fmt.Printf("%T", iGameDate)

		//game := pbpData.GameData.Game.Pk
		//fmt.Println(game)

		//teams := pbpData.GameData.Teams.Home.Name
		//fmt.Println(teams)

		//players := pbpData.GameData.Players
		//fmt.Println(players)
		//fmt.Printf("%T", players)

		// linescore := pbpData.LiveData.Linescore
		// fmt.Println(linescore)

		// rosters := pbpData.LiveData.Boxscore.Teams
		// fmt.Println(rosters)

		//fmt.Println(len(pbpData.GameData.Players[0]))

		// for i := 0; i < len(pbpData.GameData.Players[0]); i++ {
		//  	s := pbpData.GameData.Players[0]

		// }

		// for k,v := range players {
		// 	fmt.Println(k,v)
		// }
		//fmt.Println(players)
		
		// value prints in curly braces because it's a struct
		players := pbpData.GameData.Players

		//for k,_ := range pbpData.GameData.Players {
		
		var fnames []string
		for k,v := range players {
			//fmt.Println(v)
			//fmt.Printf("%v", v)
			//fmt.Println(k[2:])
			k = k[2:]
			//fmt.Println(k,v.FullName)
			fnames = append(fnames, v.FullName)
		}
		//fmt.Printf("%T", players)

		//fmt.Println(fnames)

		teams := pbpData.GameData.Teams

		// Prepare team output

		for _, v := range teams { 
			outTeamsIceSit := strings.ToLower(v.Abbreviation) // iceSit = 'home' or 'away
			fmt.Println(outTeamsIceSit)
			
			for _, strSit := range strengthSits {
				outTeamsIceSitStrSit := strSit
				fmt.Println(outTeamsIceSitStrSit)
				
				for _, scSit := range scoreSits {
					outTeamsIceSitScSit := scSit
					fmt.Println(outTeamsIceSitScSit)
					
					for _, stat := range teamStats {
						outTeamsiceSitStrSitScSitStat := stat
						outTeamsiceSitStrSitScSitStat = "0"
						fmt.Println("--", outTeamsiceSitStrSitScSitStat)
						

					}
					

				}
				
			}
		}

		// Prepare players output
		for _, pId := range players {
			outPlayersPidPrimaryPos := strings.ToLower(pId.PrimaryPosition.Abbreviation)
			fmt.Println(outPlayersPidPrimaryPos)
			outPlayersPidFirstName := pId.FirstName
			fmt.Println(outPlayersPidFirstName)
			outPlayersPidLastName := pId.LastName
			fmt.Println(outPlayersPidLastName)

		}

		//fmt.Println("--", pbpData.LiveData.Plays.PenaltyPlays)
		//fmt.Println("--", pbpData.LiveData.Plays.PlaysByPeriod)
		//fmt.Println("--", pbpData.LiveData.Plays.CurrentPlay)
		//fmt.Println("---", pbpData.LiveData.Plays.AllPlays)
		//fmt.Println("----", pbpData.LiveData.Plays.AllPlays)
		//for k,v := range pbpData.LiveData.Plays {
		// 	fmt.Println(k,v)
		//}



		//fmt.Println(pbpData.GameData.Players.FullName)
		//fmt.Println(v)
		// l := len(pbpData.GameData.Players)
		// fmt.Println(l)
		// res := pbpData.GameData.Players["ID8474884"].FullName
		// fmt.Println(res)

		// allplays := pbpData.LiveData.AllPlays
		// //fmt.Println(allplays)
		// fmt.Printf("%T", allplays)

		// for k,v := range allplays {
		// 	fmt.Println(k,v)
		// }

		// for i := 0; i < len(pbpData.LiveData.AllPlays.scoringPlays); i++ {
		// 	fmt.Println(i)
		// }

		// linescore := pbpData.LiveData.LineScore
		// fmt.Println(linescore)

		// for k,v := range linescore {
		// 	fmt.Println(k,v)
		// }
		fmt.Println(pbpData.LiveData.Linescore)

		//fmt.Println("--", pbpData.LiveData.Boxscore.BoxscoreTeams.BoxscoreTeamsAway.BoxscoreTeamsAwayPlayers)
		fmt.Println("--", pbpData.LiveData.Boxscore)

		fmt.Println("---", pbpData.LiveData.Decisions)

	
	}

}

type Data struct {
	Copyright string   `json:"copyright"`
	GamePk    int      `json:"gamePk"`
    Link      string   `json:"link"`
	GameData  GameData `json:"gameData"`
	LiveData  LiveData `json:"liveData"`
}

type GameData struct {
	Players map[string]Player `json:"players"`
	DateTime Datetime `json:"dateTime"`
	Game Game `json:"game"`
	Status Status `json:"status"`
	Teams map[string]Team `json:"teams"`

}

type Game struct {
	Pk     int    `json:"pk"`
	Season string `json:"season"`
	Type   string `json:"type"`
}

type Datetime struct {
	DateTime    time.Time `json:"dateTime"`
	EndDateTime time.Time `json:"endDateTime"`
}

type Status struct {
	AbstractGameState string `json:"abstractGameState"`
	CodedGameState    string `json:"codedGameState"`
	DetailedState     string `json:"detailedState"`
	StatusCode        string `json:"statusCode"`
	StartTimeTBD      bool   `json:"startTimeTBD"`
}

type Team struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Link  string `json:"link"`
	Abbreviation    string `json:"abbreviation"`
	TriCode         string `json:"triCode"`
	TeamName        string `json:"teamName"`
	LocationName    string `json:"locationName"`
	FirstYearOfPlay string `json:"firstYearOfPlay"`
	ShortName       string `json:"shortName"`
	OfficialSiteURL string `json:"officialSiteUrl"`
	FranchiseID     int    `json:"franchiseId"`
	Active          bool   `json:"active"`
	Venue struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Link     string `json:"link"`
		City     string `json:"city"`
		TimeZone struct {
			ID     string `json:"id"`
			Offset int    `json:"offset"`
			Tz     string `json:"tz"`
		} `json:"timeZone"`
	} `json:"venue"`
	Division        struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Link string `json:"link"`
	} `json:"division"`
	Conference struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Link string `json:"link"`
	} `json:"conference"`
	Franchise struct {
		FranchiseID int    `json:"franchiseId"`
		TeamName    string `json:"teamName"`
		Link        string `json:"link"`
	} `json:"franchise"`
}

type Player struct {
	ID                 int    `json:"id"`
    FullName           string `json:"fullName"`
    Link               string `json:"link"`
    FirstName          string `json:"firstName"`
    LastName           string `json:"lastName"`
    PrimaryNumber      string `json:"primaryNumber"`
    BirthDate          string `json:"birthDate"`
    CurrentAge         int    `json:"currentAge"`
    BirthCity          string `json:"birthCity"`
    BirthStateProvince string `json:"birthStateProvince"`
    BirthCountry       string `json:"birthCountry"`
    Nationality        string `json:"nationality"`
    Height             string `json:"height"`
    Weight             int    `json:"weight"`
    Active             bool   `json:"active"`
    AlternateCaptain   bool   `json:"alternateCaptain"`
    Captain            bool   `json:"captain"`
    Rookie             bool   `json:"rookie"`
    ShootsCatches      string `json:"shootsCatches"`
    RosterStatus       string `json:"rosterStatus"`
	CurrentTeam		   CurrentTeam `json:"currentTeam"`
    PrimaryPosition PrimaryPosition `json:"primaryPosition"`
}

type CurrentTeam struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Link    string `json:"link"`
	TriCode string `json:"triCode"`
}

type PrimaryPosition struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Abbreviation string `json:"abbreviation"`
}

type LiveData struct {
	//Plays map[string]interface{} `json:"plays"`
	//Plays map[string]Plays `json:"plays"`
	Plays        Plays    `json:"plays"`
	AllPlays     AllPlays `json:"allPlays"`
	Linescore Linescore `json:"linescore"`
	Boxscore  Boxscore `json:"boxscore"`
	Decisions Decisions `json:"decisions"`
	//Boxscore map[string]Boxscore `json:"boxscore"`
	//Livedata map[string]Livedata `json:'livedata"`

}

type Decisions struct {
	LiveDataDecisionsWinner LiveDataDecisionsWinner `json:"winner"`
	LiveDataDecisionsLoser LiveDataDecisionsLoser `json:"loser"`

}

type LiveDataDecisionsWinner struct {
	ID       int    `json:"id"`
	FullName string `json:"fullName"`
	Link     string `json:"link"`
}

type LiveDataDecisionsLoser struct {
	ID       int    `json:"id"`
	FullName string `json:"fullName"`
	Link     string `json:"link"`
}


type Plays struct {
	AllPlays AllPlays `json:"allPlays"`
	ScoringPlays  []int `json:"scoringPlays"`
	PenaltyPlays  []int `json:"penaltyPlays"`
	PlaysByPeriod PlaysByPeriod `json:"playsByPeriod"`
	CurrentPlay CurrentPlay `json:"currentPlay"`
}

type PlaysByPeriod []struct {
	StartIndex int   `json:"startIndex"`
	Plays      []int `json:"plays"`
	EndIndex   int   `json:"endIndex"`
}

type CurrentPlay struct {
	//CurrentPlay CurrentPlay `json:"currentPlay"`
	CurrentPlayResult CurrentPlayResult `json:"result"`
	CurrentPlayAbout CurrentPlayAbout `json:"about"`
	CurrentPlayAboutGoals CurrentPlayAboutGoals `json:"goals"`
	CurrentPlayCoordinates CurrentPlayCoordinates `json:"coordinates"`
}

type CurrentPlayResult struct {
	Event       string `json:"event"`
	EventCode   string `json:"eventCode"`
	EventTypeID string `json:"eventTypeId"`
	Description string `json:"description"`
}

type CurrentPlayAbout struct {
	EventIdx            int       `json:"eventIdx"`
	EventID             int       `json:"eventId"`
	Period              int       `json:"period"`
	PeriodType          string    `json:"periodType"`
	OrdinalNum          string    `json:"ordinalNum"`
	PeriodTime          string    `json:"periodTime"`
	PeriodTimeRemaining string    `json:"periodTimeRemaining"`
	DateTime            time.Time `json:"dateTime"`
}

type CurrentPlayAboutGoals struct {
	Away int `json:"away"`
	Home int `json:"home"`
}

type CurrentPlayCoordinates struct {

}

type AllPlays []struct {
	Allplays       AllPlays       `json:"allPlays"`
	AllPlaysResult AllPlaysResult `json:"result"`
	AllPlaysAbout  AllPlaysAbout  `json:"about"`
	AllPlaysAboutGoals  AllPlaysAboutGoals  `json:"goals"`
	AllPlaysCoordinates AllPlaysCoordinates `json:"coordinates,omitempty"`
	AllPlaysTeam AllPlaysTeam `json:"team,omitempty"`
}

type AllPlaysResult struct {
	Event       string `json:"event"`
	EventCode   string `json:"eventCode"`
	EventTypeID string `json:"eventTypeId"`
	Description string `json:"description"`
	SecondaryType string `json:"secondaryType"`
	PenaltySeverity string `json:"penaltySeverity"`
	PenaltyMinutes  int    `json:"penaltyMinutes"`
	GameWinningGoal bool `json:"gameWinningGoal"`
	EmptyNet        bool `json:"emptyNet"`
	AllPlaysResultStrength AllPlaysResultStrength `json:"strength"`
}

type AllPlaysResultStrength struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type AllPlaysAbout struct {
	EventIdx            int       `json:"eventIdx"`
	EventID             int       `json:"eventId"`
	Period              int       `json:"period"`
	PeriodType          string    `json:"periodType"`
	OrdinalNum          string    `json:"ordinalNum"`
	PeriodTime          string    `json:"periodTime"`
	PeriodTimeRemaining string    `json:"periodTimeRemaining"`
	DateTime            time.Time `json:"dateTime"`
}

type AllPlaysAboutGoals struct {
	Away int `json:"away"`
	Home int `json:"home"`
}

type AllPlaysCoordinates struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type AllPlaysPlayers []struct {
	AllPlaysPlayers AllPlaysPlayers `json:"players,omitempty"`
	AllPlaysPlayer AllPlaysPlayer `json:"players,omitempty"`
	PlayerType string `json:"playerType"`
}

type AllPlaysPlayer struct {
	ID       int    `json:"id"`
	FullName string `json:"fullName"`
	Link     string `json:"link"`
}

type AllPlaysTeam struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Link    string `json:"link"`
	TriCode string `json:"triCode"`
}

type Linescore struct {
	CurrentPeriod              int    `json:"currentPeriod"`
	CurrentPeriodOrdinal       string `json:"currentPeriodOrdinal"`
	CurrentPeriodTimeRemaining string `json:"currentPeriodTimeRemaining"`
	LinescorePeriods LinescorePeriods `json:"periods"`
	LinescoreShootoutInfo LinescoreShootoutInfo `json:"shootoutInfo"`
	LinescoreTeams LinescoreTeams `json:"teams"`
	PowerPlayStrength string `json:"powerPlayStrength"`
	HasShootout       bool   `json:"hasShootout"`
	LinescoreIntermissionInfo LinescoreIntermissionInfo `json:"intermissionInfo"`
	LinescorePowerPlayInfo LinescorePowerPlayInfo `json:"powerPlayInfo"`
	
}

type LinescorePeriods []struct {
	PeriodType string    `json:"periodType"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Num        int       `json:"num"`
	OrdinalNum string    `json:"ordinalNum"`
	LinescorePeriodsHome LinescorePeriodsHome `json:"home"`
	LinescorePeriodsAway LinescorePeriodsAway `json:"away"`
}

type LinescorePeriodsHome struct {
	Goals       int    `json:"goals"`
	ShotsOnGoal int    `json:"shotsOnGoal"`
	RinkSide    string `json:"rinkSide"`
}

type LinescorePeriodsAway struct {
	Goals       int    `json:"goals"`
	ShotsOnGoal int    `json:"shotsOnGoal"`
	RinkSide    string `json:"rinkSide"`
}

type LinescoreShootoutInfo struct {
	LinescoreShootoutInfoAway LinescoreShootoutInfoAway `json:"away"`
	LinescoreShootoutInfoHome LinescoreShootoutInfoHome `json:"home"`
}

type LinescoreShootoutInfoAway struct {
	Scores   int `json:"scores"`
	Attempts int `json:"attempts"`
}

type LinescoreShootoutInfoHome struct {
	Scores   int `json:"scores"`
	Attempts int `json:"attempts"`
}

type LinescoreTeams struct {
	LinescoreTeamsHome LinescoreTeamsHome `json:"home"`
	LinescoreTeamsAway LinescoreTeamsAway `json:"away"`
}

type LinescoreTeamsHome struct {
	LinescoreTeamsHomeTeam LinescoreTeamsHomeTeam `json:"team"`
	Goals        int  `json:"goals"`
	ShotsOnGoal  int  `json:"shotsOnGoal"`
	GoaliePulled bool `json:"goaliePulled"`
	NumSkaters   int  `json:"numSkaters"`
	PowerPlay    bool `json:"powerPlay"`
}

type LinescoreTeamsHomeTeam struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Link         string `json:"link"`
	Abbreviation string `json:"abbreviation"`
	TriCode      string `json:"triCode"`
}

type LinescoreTeamsAway struct {
	LinescoreTeamsAwayTeam LinescoreTeamsAwayTeam `json:"team"`
	Goals        int  `json:"goals"`
	ShotsOnGoal  int  `json:"shotsOnGoal"`
	GoaliePulled bool `json:"goaliePulled"`
	NumSkaters   int  `json:"numSkaters"`
	PowerPlay    bool `json:"powerPlay"`
}

type LinescoreTeamsAwayTeam struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Link         string `json:"link"`
	Abbreviation string `json:"abbreviation"`
	TriCode      string `json:"triCode"`
}

type LinescoreIntermissionInfo struct {
	IntermissionTimeRemaining int  `json:"intermissionTimeRemaining"`
	IntermissionTimeElapsed   int  `json:"intermissionTimeElapsed"`
	InIntermission            bool `json:"inIntermission"`
}

type LinescorePowerPlayInfo struct {
	SituationTimeRemaining int  `json:"situationTimeRemaining"`
	SituationTimeElapsed   int  `json:"situationTimeElapsed"`
	InSituation            bool `json:"inSituation"`
}

type Boxscore struct {
	BoxscoreTeams BoxscoreTeams `json:"teams"`
	BoxscoreOfficials BoxscoreOfficials `json:"officials"`
	//BoxscoreTeamsAwayPlayers map[string]BoxscoreTeamsAwayPlayers `json:"players"`
	//BoxscoreTeamsAwayPlayers BoxscoreTeamsAwayPlayers `json:"players"`
}

type BoxscoreOfficials []struct {
	BoxscoreOfficialsOfficial BoxscoreOfficialsOfficial `json:"official"`
	OfficialType string `json:"officialType"`
}

type BoxscoreOfficialsOfficial struct {
	ID       int    `json:"id"`
	FullName string `json:"fullName"`
	Link     string `json:"link"`
}

type BoxscoreTeams struct {
	BoxscoreTeamsAway BoxscoreTeamsAway `json:"away"`
	BoxscoreTeamsHome BoxscoreTeamsHome `json:"home"`
}

type BoxscoreTeamsAway struct {
	BoxscoreTeamsAwayTeam BoxscoreTeamsAwayTeam `json:"team"`
	BoxscoreTeamsAwayTeamStats BoxscoreTeamsAwayTeamStats `json:"teamStats"`
	BoxscoreTeamsAwayPlayers map[string]BoxscoreTeamsAwayPlayer `json:"players"`
	Goalies   []int `json:"goalies"`
	Skaters   []int `json:"skaters"`
	OnIce     []int `json:"onIce"`
	BoxscoreTeamsAwayOnIcePlus BoxscoreTeamsAwayOnIcePlus `json:"onIcePlus"`
	Scratches  []int         `json:"scratches"`
	PenaltyBox []interface{} `json:"penaltyBox"`
	BoxscoreTeamsAwayCoaches BoxscoreTeamsAwayCoaches `json:"coaches"`

}

type BoxscoreTeamsHome struct {
	BoxscoreTeamsHomeTeam BoxscoreTeamsHomeTeam `json:"team"`
	BoxscoreTeamsHomeTeamStats BoxscoreTeamsHomeTeamStats `json:"teamStats"`
	BoxscoreTeamsHomePlayers map[string]BoxscoreTeamsHomePlayer `json:"players"`
	Goalies   []int `json:"goalies"`
	Skaters   []int `json:"skaters"`
	OnIce     []int `json:"onIce"`
	BoxscoreTeamsHomeOnIcePlus BoxscoreTeamsHomeOnIcePlus `json:"onIcePlus"`
	Scratches  []int         `json:"scratches"`
	PenaltyBox []interface{} `json:"penaltyBox"`
	BoxscoreTeamsHomeCoaches BoxscoreTeamsHomeCoaches `json:"coaches"`
}

type BoxscoreTeamsAwayOnIcePlus []struct {
	PlayerID      int `json:"playerId"`
	ShiftDuration int `json:"shiftDuration"`
	Stamina       int `json:"stamina"`
}

type BoxscoreTeamsHomeOnIcePlus []struct {
	PlayerID      int `json:"playerId"`
	ShiftDuration int `json:"shiftDuration"`
	Stamina       int `json:"stamina"`
}

type BoxscoreTeamsAwayCoaches []struct {
	BoxscoreTeamsAwayCoachesPerson BoxscoreTeamsAwayCoachesPerson `json:"person"`
	BoxscoreTeamsAwayCoachesPosition BoxscoreTeamsAwayCoachesPosition `json:"position"`
}

type BoxscoreTeamsHomeCoaches []struct {
	BoxscoreTeamsHomeCoachesPerson BoxscoreTeamsHomeCoachesPerson `json:"person"`
	BoxscoreTeamsHomeCoachesPosition BoxscoreTeamsHomeCoachesPosition `json:"position"`
}

type BoxscoreTeamsAwayCoachesPerson struct {
	FullName string `json:"fullName"`
	Link     string `json:"link"`
}

type BoxscoreTeamsHomeCoachesPerson struct {
	FullName string `json:"fullName"`
	Link     string `json:"link"`
}

type BoxscoreTeamsAwayCoachesPosition struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Abbreviation string `json:"abbreviation"`
}

type BoxscoreTeamsHomeCoachesPosition struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Abbreviation string `json:"abbreviation"`
}

type BoxscoreTeamsAwayTeam struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Link         string `json:"link"`
	Abbreviation string `json:"abbreviation"`
	TriCode      string `json:"triCode"`
}

type BoxscoreTeamsHomeTeam struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Link         string `json:"link"`
	Abbreviation string `json:"abbreviation"`
	TriCode      string `json:"triCode"`
}

type BoxscoreTeamsAwayTeamStats struct {
	BoxscoreTeamAwayTeamStatsTeamSkaterStats BoxscoreTeamAwayTeamStatsTeamSkaterStats `json:"teamSkaterStats"`
}

type BoxscoreTeamsHomeTeamStats struct {
	BoxscoreTeamHomeTeamStatsTeamSkaterStats BoxscoreTeamHomeTeamStatsTeamSkaterStats `json:"teamSkaterStats"`
}

type BoxscoreTeamAwayTeamStatsTeamSkaterStats struct {
	Goals                  int     `json:"goals"`
	Pim                    int     `json:"pim"`
	Shots                  int     `json:"shots"`
	PowerPlayPercentage    string  `json:"powerPlayPercentage"`
	PowerPlayGoals         float64 `json:"powerPlayGoals"`
	PowerPlayOpportunities float64 `json:"powerPlayOpportunities"`
	FaceOffWinPercentage   string  `json:"faceOffWinPercentage"`
	Blocked                int     `json:"blocked"`
	Takeaways              int     `json:"takeaways"`
	Giveaways              int     `json:"giveaways"`
	Hits                   int     `json:"hits"`
}

type BoxscoreTeamHomeTeamStatsTeamSkaterStats struct {
	Goals                  int     `json:"goals"`
	Pim                    int     `json:"pim"`
	Shots                  int     `json:"shots"`
	PowerPlayPercentage    string  `json:"powerPlayPercentage"`
	PowerPlayGoals         float64 `json:"powerPlayGoals"`
	PowerPlayOpportunities float64 `json:"powerPlayOpportunities"`
	FaceOffWinPercentage   string  `json:"faceOffWinPercentage"`
	Blocked                int     `json:"blocked"`
	Takeaways              int     `json:"takeaways"`
	Giveaways              int     `json:"giveaways"`
	Hits                   int     `json:"hits"`
}

type BoxscoreTeamsAwayPlayers struct {
	BoxscoreTeamsAwayPlayer BoxscoreTeamsAwayPlayer `json:"players"`

}

type BoxscoreTeamsHomePlayers struct {
	BoxscoreTeamsHomePlayer BoxscoreTeamsHomePlayer `json:"players"`

}

type BoxscoreTeamsAwayPlayer struct {
	BoxscoreTeamsAwayPlayerPerson BoxscoreTeamsAwayPlayerPerson `json:"person"`
	JerseyNumber string `json:"jerseyNumber"`
	BoxscoreTeamsAwayPlayerPosition BoxscoreTeamsAwayPlayerPosition `json:"position"`
	BoxscoreTeamsAwayPlayerStats BoxscoreTeamsAwayPlayerStats `json:"stats"`

}

type BoxscoreTeamsHomePlayer struct {
	BoxscoreTeamsHomePlayerPerson BoxscoreTeamsHomePlayerPerson `json:"person"`
	JerseyNumber string `json:"jerseyNumber"`
	BoxscoreTeamsHomePlayerPosition BoxscoreTeamsHomePlayerPosition `json:"position"`
	BoxscoreTeamsHomePlayerStats BoxscoreTeamsHomePlayerStats `json:"stats"`

}

type BoxscoreTeamsAwayPlayerPerson struct {
	ID            int    `json:"id"`
	FullName      string `json:"fullName"`
	Link          string `json:"link"`
	ShootsCatches string `json:"shootsCatches"`
	RosterStatus  string `json:"rosterStatus"`
}

type BoxscoreTeamsHomePlayerPerson struct {
	ID            int    `json:"id"`
	FullName      string `json:"fullName"`
	Link          string `json:"link"`
	ShootsCatches string `json:"shootsCatches"`
	RosterStatus  string `json:"rosterStatus"`
}

type BoxscoreTeamsAwayPlayerPosition struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Abbreviation string `json:"abbreviation"`
}

type BoxscoreTeamsHomePlayerPosition struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Abbreviation string `json:"abbreviation"`
}

type BoxscoreTeamsAwayPlayerStats struct {
	BoxscoreTeamsAwayPlayersStatsSkaterStats BoxscoreTeamsAwayPlayersStatsSkaterStats `json:"skaterStats"`
	BoxscoreTeamsAwayPlayersStatsGoalieStats BoxscoreTeamsAwayPlayersStatsGoalieStats `json:"goalieStats"`
}

type BoxscoreTeamsHomePlayerStats struct {
	BoxscoreTeamsHomePlayersStatsSkaterStats BoxscoreTeamsHomePlayersStatsSkaterStats `json:"skaterStats"`
	BoxscoreTeamsHomePlayersStatsGoalieStats BoxscoreTeamsHomePlayersStatsGoalieStats `json:"goalieStats"`
}

type BoxscoreTeamsAwayPlayersStatsSkaterStats struct {
	TimeOnIce            string `json:"timeOnIce"`
	Assists              int    `json:"assists"`
	Goals                int    `json:"goals"`
	Shots                int    `json:"shots"`
	Hits                 int    `json:"hits"`
	PowerPlayGoals       int    `json:"powerPlayGoals"`
	PowerPlayAssists     int    `json:"powerPlayAssists"`
	PenaltyMinutes       int    `json:"penaltyMinutes"`
	FaceOffWins          int    `json:"faceOffWins"`
	FaceoffTaken         int    `json:"faceoffTaken"`
	Takeaways            int    `json:"takeaways"`
	Giveaways            int    `json:"giveaways"`
	ShortHandedGoals     int    `json:"shortHandedGoals"`
	ShortHandedAssists   int    `json:"shortHandedAssists"`
	Blocked              int    `json:"blocked"`
	PlusMinus            int    `json:"plusMinus"`
	EvenTimeOnIce        string `json:"evenTimeOnIce"`
	PowerPlayTimeOnIce   string `json:"powerPlayTimeOnIce"`
	ShortHandedTimeOnIce string `json:"shortHandedTimeOnIce"`
}

type BoxscoreTeamsHomePlayersStatsSkaterStats struct {
	TimeOnIce            string `json:"timeOnIce"`
	Assists              int    `json:"assists"`
	Goals                int    `json:"goals"`
	Shots                int    `json:"shots"`
	Hits                 int    `json:"hits"`
	PowerPlayGoals       int    `json:"powerPlayGoals"`
	PowerPlayAssists     int    `json:"powerPlayAssists"`
	PenaltyMinutes       int    `json:"penaltyMinutes"`
	FaceOffWins          int    `json:"faceOffWins"`
	FaceoffTaken         int    `json:"faceoffTaken"`
	Takeaways            int    `json:"takeaways"`
	Giveaways            int    `json:"giveaways"`
	ShortHandedGoals     int    `json:"shortHandedGoals"`
	ShortHandedAssists   int    `json:"shortHandedAssists"`
	Blocked              int    `json:"blocked"`
	PlusMinus            int    `json:"plusMinus"`
	EvenTimeOnIce        string `json:"evenTimeOnIce"`
	PowerPlayTimeOnIce   string `json:"powerPlayTimeOnIce"`
	ShortHandedTimeOnIce string `json:"shortHandedTimeOnIce"`
}

type BoxscoreTeamsAwayPlayersStatsGoalieStats struct {
	TimeOnIce                  string  `json:"timeOnIce"`
	Assists                    int     `json:"assists"`
	Goals                      int     `json:"goals"`
	Pim                        int     `json:"pim"`
	Shots                      int     `json:"shots"`
	Saves                      int     `json:"saves"`
	PowerPlaySaves             int     `json:"powerPlaySaves"`
	ShortHandedSaves           int     `json:"shortHandedSaves"`
	EvenSaves                  int     `json:"evenSaves"`
	ShortHandedShotsAgainst    int     `json:"shortHandedShotsAgainst"`
	EvenShotsAgainst           int     `json:"evenShotsAgainst"`
	PowerPlayShotsAgainst      int     `json:"powerPlayShotsAgainst"`
	Decision                   string  `json:"decision"`
	SavePercentage             float64 `json:"savePercentage"`
	PowerPlaySavePercentage    float64 `json:"powerPlaySavePercentage"`
	EvenStrengthSavePercentage float64 `json:"evenStrengthSavePercentage"`
}

type BoxscoreTeamsHomePlayersStatsGoalieStats struct {
	TimeOnIce                  string  `json:"timeOnIce"`
	Assists                    int     `json:"assists"`
	Goals                      int     `json:"goals"`
	Pim                        int     `json:"pim"`
	Shots                      int     `json:"shots"`
	Saves                      int     `json:"saves"`
	PowerPlaySaves             int     `json:"powerPlaySaves"`
	ShortHandedSaves           int     `json:"shortHandedSaves"`
	EvenSaves                  int     `json:"evenSaves"`
	ShortHandedShotsAgainst    int     `json:"shortHandedShotsAgainst"`
	EvenShotsAgainst           int     `json:"evenShotsAgainst"`
	PowerPlayShotsAgainst      int     `json:"powerPlayShotsAgainst"`
	Decision                   string  `json:"decision"`
	SavePercentage             float64 `json:"savePercentage"`
	PowerPlaySavePercentage    float64 `json:"powerPlaySavePercentage"`
	EvenStrengthSavePercentage float64 `json:"evenStrengthSavePercentage"`
}


