package valorantapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

type valorantCompetitiveTierLookup struct {
	Status int `json:"status"`
	Data   []struct {
		UUID            string `json:"uuid"`
		AssetObjectName string `json:"assetObjectName"`
		Tiers           []struct {
			Tier                 int         `json:"tier"`
			TierName             string      `json:"tierName"`
			Division             string      `json:"division"`
			DivisionName         string      `json:"divisionName"`
			Color                string      `json:"color"`
			BackgroundColor      string      `json:"backgroundColor"`
			SmallIcon            string      `json:"smallIcon"`
			LargeIcon            string      `json:"largeIcon"`
			RankTriangleDownIcon interface{} `json:"rankTriangleDownIcon"`
			RankTriangleUpIcon   interface{} `json:"rankTriangleUpIcon"`
		} `json:"tiers"`
		AssetPath string `json:"assetPath"`
	} `json:"data"`
}

var CompetitiveTierLookup = map[int]string{}

func init() {

	APICall := newAPICall()

	APICall.URL = "https://valorant-api.com/v1/competitivetiers"
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return
	}

	err = APICall.Call()
	if err != nil {
		return
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return
	}

	Response := &valorantCompetitiveTierLookup{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return
	}

	for _, v := range Response.Data[0].Tiers {

		CompetitiveTierLookup[v.Tier] = v.TierName

	}

}

type valorantSeasonData struct {
	SeasonID                   string         `json:"SeasonID"`
	NumberOfWins               int            `json:"NumberOfWins"`
	NumberOfWinsWithPlacements int            `json:"NumberOfWinsWithPlacements"`
	NumberOfGames              int            `json:"NumberOfGames"`
	Rank                       int            `json:"Rank"`
	CapstoneWins               int            `json:"CapstoneWins"`
	LeaderboardRank            int            `json:"LeaderboardRank"`
	CompetitiveTier            int            `json:"CompetitiveTier"`
	RankedRating               int            `json:"RankedRating"`
	WinsByTier                 map[string]int `json:"WinsByTier"`
	GamesNeededForRating       int            `json:"GamesNeededForRating"`
	TotalWinsNeededForRank     int            `json:"TotalWinsNeededForRank"`
}

/* [ - MMR Data - ] */
type valorantMMRData struct {
	Version                 int64  `json:"Version"`
	Subject                 string `json:"Subject"`
	LatestCompetitiveUpdate struct {
		MatchID                        string `json:"MatchID"`
		MapID                          string `json:"MapID"`
		QueueID                        string `json:"QueueID"`
		SeasonID                       string `json:"SeasonID"`
		MatchStartTime                 int64  `json:"MatchStartTime"`
		MatchLength                    int    `json:"MatchLength"`
		TierAfterUpdate                int    `json:"TierAfterUpdate"`
		TierBeforeUpdate               int    `json:"TierBeforeUpdate"`
		RankedRatingAfterUpdate        int    `json:"RankedRatingAfterUpdate"`
		RankedRatingBeforeUpdate       int    `json:"RankedRatingBeforeUpdate"`
		RankedRatingEarned             int    `json:"RankedRatingEarned"`
		RankedRatingPerformanceBonus   int    `json:"RankedRatingPerformanceBonus"`
		RankedRatingRefundApplied      int    `json:"RankedRatingRefundApplied"`
		NewMapIncentiveRRForgiven      int    `json:"NewMapIncentiveRRForgiven"`
		CompetitiveMovement            string `json:"CompetitiveMovement"`
		AFKPenalty                     int    `json:"AFKPenalty"`
		WasDerankProtected             bool   `json:"WasDerankProtected"`
		WasDerankProtectionReplenished bool   `json:"WasDerankProtectionReplenished"`
	} `json:"LatestCompetitiveUpdate"`
	NewPlayerExperienceFinished   bool   `json:"NewPlayerExperienceFinished"`
	IsActRankBadgeHidden          bool   `json:"IsActRankBadgeHidden"`
	IsLeaderboardAnonymized       bool   `json:"IsLeaderboardAnonymized"`
	OnboardingFlowV2Enabled       bool   `json:"OnboardingFlowV2Enabled"`
	OnboardingStatus              string `json:"OnboardingStatus"`
	IsAtDerankProtectedTier       bool   `json:"IsAtDerankProtectedTier"`
	DerankProtectedGamesRemaining int    `json:"DerankProtectedGamesRemaining"`
	DerankProtectedStatus         string `json:"DerankProtectedStatus"`
	QueueSkills                   struct {
		Competitive struct {
			TotalGamesNeededForRating         int                           `json:"TotalGamesNeededForRating"`
			TotalGamesNeededForLeaderboard    int                           `json:"TotalGamesNeededForLeaderboard"`
			CurrentSeasonGamesNeededForRating int                           `json:"CurrentSeasonGamesNeededForRating"`
			SeasonalInfoBySeasonID            map[string]valorantSeasonData `json:"SeasonalInfoBySeasonID"`
		} `json:"competitive"`
		Newmap struct {
			TotalGamesNeededForRating         int                           `json:"TotalGamesNeededForRating"`
			TotalGamesNeededForLeaderboard    int                           `json:"TotalGamesNeededForLeaderboard"`
			CurrentSeasonGamesNeededForRating int                           `json:"CurrentSeasonGamesNeededForRating"`
			SeasonalInfoBySeasonID            map[string]valorantSeasonData `json:"SeasonalInfoBySeasonID"`
		} `json:"newmap"`
		Seeding struct {
			TotalGamesNeededForRating         int                           `json:"TotalGamesNeededForRating"`
			TotalGamesNeededForLeaderboard    int                           `json:"TotalGamesNeededForLeaderboard"`
			CurrentSeasonGamesNeededForRating int                           `json:"CurrentSeasonGamesNeededForRating"`
			SeasonalInfoBySeasonID            map[string]valorantSeasonData `json:"SeasonalInfoBySeasonID"`
		} `json:"seeding"`
		Skirmishascension1V1 struct {
			TotalGamesNeededForRating         int                           `json:"TotalGamesNeededForRating"`
			TotalGamesNeededForLeaderboard    int                           `json:"TotalGamesNeededForLeaderboard"`
			CurrentSeasonGamesNeededForRating int                           `json:"CurrentSeasonGamesNeededForRating"`
			SeasonalInfoBySeasonID            map[string]valorantSeasonData `json:"SeasonalInfoBySeasonID"`
		} `json:"skirmishascension1v1"`
		Skirmishascension2V2 struct {
			TotalGamesNeededForRating         int                           `json:"TotalGamesNeededForRating"`
			TotalGamesNeededForLeaderboard    int                           `json:"TotalGamesNeededForLeaderboard"`
			CurrentSeasonGamesNeededForRating int                           `json:"CurrentSeasonGamesNeededForRating"`
			SeasonalInfoBySeasonID            map[string]valorantSeasonData `json:"SeasonalInfoBySeasonID"`
		} `json:"skirmishascension2v2"`
		Swiftplay struct {
			TotalGamesNeededForRating         int                           `json:"TotalGamesNeededForRating"`
			TotalGamesNeededForLeaderboard    int                           `json:"TotalGamesNeededForLeaderboard"`
			CurrentSeasonGamesNeededForRating int                           `json:"CurrentSeasonGamesNeededForRating"`
			SeasonalInfoBySeasonID            map[string]valorantSeasonData `json:"SeasonalInfoBySeasonID"`
		} `json:"swiftplay"`
		Unrated struct {
			TotalGamesNeededForRating         int                           `json:"TotalGamesNeededForRating"`
			TotalGamesNeededForLeaderboard    int                           `json:"TotalGamesNeededForLeaderboard"`
			CurrentSeasonGamesNeededForRating int                           `json:"CurrentSeasonGamesNeededForRating"`
			SeasonalInfoBySeasonID            map[string]valorantSeasonData `json:"SeasonalInfoBySeasonID"`
		} `json:"unrated"`
	} `json:"QueueSkills"`
}

/* [ - Get MMR Data - ] */
/* [ URL: 'https://glz-{Region}-1.{Shard}.a.pvp.net/pregame/v1/matches/{MATCH ID}' ] */
func (context ValorantAPIContext) getMMRData(player *ValorantPlayerContext) (*valorantMMRData, error) {

	APICall := newAPICall()

	APICall.URL = "https://pd." + context.localRegionData.Shard + ".a.pvp.net/mmr/v1/players/" + player.UUID
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	APICall.SetBearerAuth(context)

	err = APICall.Call()
	if err != nil {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &valorantMMRData{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	fmt.Println(Response)

	return Response, nil

}

type RankCacheItem struct {
	CurrentRank    int
	CurrentSeason  string
	PeakRank       int
	PeakRankSeason string
	LastChecked    time.Time
}

var RankCache = map[string]*RankCacheItem{}

func (context ValorantAPIContext) GetCurrentAndPeakRank(player *ValorantPlayerContext) (rankData *RankCacheItem, err error) {

	if RankCache[player.UUID] != nil {

		// If cache is younger than 30 minutes

		if time.Since(RankCache[player.UUID].LastChecked) <= time.Minute*30 {

			return RankCache[player.UUID], nil

		}

	}

	var Data *valorantMMRData

	Data, err = context.getMMRData(player)
	if err != nil {
		return nil, err
	}

	fmt.Println(Data)

	currentRank := Data.LatestCompetitiveUpdate.TierAfterUpdate
	currentSeason := Data.LatestCompetitiveUpdate.SeasonID

	if currentSeason == "" {
		currentSeason = "564d8e28-c226-3180-6285-e48a390db8b1"
	}

	highestRank := 0
	highestSeasonID := ""

	for s, v := range Data.QueueSkills.Competitive.SeasonalInfoBySeasonID {

		for rank := range v.WinsByTier {
			numRank, err := strconv.Atoi(rank)

			if err != nil {
				numRank = v.CompetitiveTier
			}

			if highestRank < numRank {
				highestRank = numRank
				highestSeasonID = s
			}

		}

	}

	if highestSeasonID == "" {
		highestSeasonID = currentSeason
	}

	RankCache[player.UUID] = &RankCacheItem{
		CurrentRank:    currentRank,
		CurrentSeason:  currentSeason,
		PeakRank:       highestRank,
		PeakRankSeason: highestSeasonID,
		LastChecked:    time.Now(),
	}

	return RankCache[player.UUID], nil

}

type ValorantCompetitive struct {
	Status int `json:"status"`
	Data   []struct {
		UUID                 string    `json:"uuid"`
		StartTime            time.Time `json:"startTime"`
		EndTime              time.Time `json:"endTime"`
		SeasonUUID           string    `json:"seasonUuid"`
		CompetitiveTiersUUID string    `json:"competitiveTiersUuid"`
		Borders              []struct {
			UUID         string      `json:"uuid"`
			Level        int         `json:"level"`
			WinsRequired int         `json:"winsRequired"`
			DisplayIcon  string      `json:"displayIcon"`
			SmallIcon    interface{} `json:"smallIcon"`
			AssetPath    string      `json:"assetPath"`
		} `json:"borders"`
		AssetPath string `json:"assetPath"`
	} `json:"data"`
}

var SeasonLookupCache map[string]string

/* Gets all Seasons with their UUID as lookup to get the competitive rank uuid */
func (context ValorantAPIContext) GetSeasonLookup() (map[string]string, error) {

	if SeasonLookupCache != nil {
		return SeasonLookupCache, nil
	}

	APICall := newAPICall()

	APICall.URL = "https://valorant-api.com/v1/seasons/competitive"
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	err = APICall.Call()
	if err != nil {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &ValorantCompetitive{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	if Response.Status != 200 {
		return nil, errors.New(VALORANT_NO_VALID_RESPONSE)
	}

	Output := map[string]string{}

	for _, v := range Response.Data {

		Output[v.SeasonUUID] = v.CompetitiveTiersUUID

	}

	SeasonLookupCache = Output

	return SeasonLookupCache, nil

}
