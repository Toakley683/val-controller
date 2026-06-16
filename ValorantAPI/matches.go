package valorantapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

type internalValorantCurrentMatch struct {
	Subject string `json:"Subject"`
	MatchID string `json:"MatchID"`
}

type ValorantCurrentMatch struct {
	Subject        string
	PregameMatchID string
	GameMatchID    string
}

/* [ - Get Pregame MatchID - ] */
/* [ URL: 'https://glz-{Region}-1.{Shard}.a.pvp.net/pregame/v1/players/{PUUID}' ] */
func (context ValorantAPIContext) getPregameGamePlayer(playerContext *ValorantPlayerContext) (*internalValorantCurrentMatch, error) {

	APICall := newAPICall()

	APICall.URL = "https://glz-" + context.localRegionData.Region + "-1." + context.localRegionData.Shard + ".a.pvp.net/pregame/v1/players/" + playerContext.UUID
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	APICall.SetBearerAuth(context)

	err = APICall.Call()
	if err != nil && err.Error() != VALORANT_NO_VALID_RESPONSE {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &internalValorantCurrentMatch{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [ - Get Game MatchID - ] */
/* [ URL: 'https://glz-{Region}-1.{Shard}.a.pvp.net/core-game/v1/players/{PUUID}' ] */
func (context ValorantAPIContext) getMainGamePlayer(playerContext *ValorantPlayerContext) (*internalValorantCurrentMatch, error) {

	APICall := newAPICall()

	APICall.URL = "https://glz-" + context.localRegionData.Region + "-1." + context.localRegionData.Shard + ".a.pvp.net/core-game/v1/players/" + playerContext.UUID
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	APICall.SetBearerAuth(context)

	err = APICall.Call()
	if err != nil && err.Error() != VALORANT_NO_VALID_RESPONSE {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &internalValorantCurrentMatch{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [ - Name Struct - ] */

type ValorantNameStruct struct {
	DisplayName string `json:"DisplayName"`
	Subject     string `json:"Subject"`
	GameName    string `json:"GameName"`
	TagLine     string `json:"TagLine"`
}

/* [ - Get Valorant names from PlayerUUID - ] */
func (context ValorantAPIContext) SubjectsToNames(Subjects []string) (map[string]ValorantNameStruct, error) {

	APICall := newAPICall()

	APICall.URL = "https://pd." + context.localRegionData.Shard + ".a.pvp.net/name-service/v2/players"

	array, err := json.Marshal(Subjects)
	if err != nil {
		return nil, err
	}

	APICall.Body = strings.NewReader(string(array))
	APICall.Method = "PUT"

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

	Response := &[]ValorantNameStruct{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	FinalData := map[string]ValorantNameStruct{}

	for _, v := range *Response {

		FinalData[v.Subject] = v

	}

	return FinalData, nil

}

/* [ - Get Current Pregame & Match ID - ] */
func (context ValorantAPIContext) GetCurrentGamePlayer(playerContext *ValorantPlayerContext) (*ValorantCurrentMatch, error) {

	pregameInfo, err := context.getPregameGamePlayer(playerContext)
	if err != nil {
		return nil, err
	}

	curgameInfo, err := context.getMainGamePlayer(playerContext)
	if err != nil {
		return nil, err
	}

	FinalMatch := ValorantCurrentMatch{}

	FinalMatch.Subject = playerContext.UUID
	FinalMatch.PregameMatchID = pregameInfo.MatchID
	FinalMatch.GameMatchID = curgameInfo.MatchID

	return &FinalMatch, nil

}

type PlayerIdentity struct {
	Subject                string `json:"Subject"`
	PlayerCardID           string `json:"PlayerCardID"`
	PlayerTitleID          string `json:"PlayerTitleID"`
	AccountLevel           int    `json:"AccountLevel"`
	PreferredLevelBorderID string `json:"PreferredLevelBorderID"`
	Incognito              bool   `json:"Incognito"`
	HideAccountLevel       bool   `json:"HideAccountLevel"`
	GameName               string
	TagLine                string
}

/* [ - Current Pregame & Match Data - ] */
type valorantMainGameData struct {
	MatchID          string `json:"MatchID"`
	Version          int64  `json:"Version"`
	State            string `json:"State"`
	MapID            string `json:"MapID"`
	ModeID           string `json:"ModeID"`
	ProvisioningFlow string `json:"ProvisioningFlow"`
	GamePodID        string `json:"GamePodID"`
	AllMUCName       string `json:"AllMUCName"`
	TeamMUCName      string `json:"TeamMUCName"`
	TeamVoiceID      string `json:"TeamVoiceID"`
	TeamMatchToken   string `json:"TeamMatchToken"`
	IsReconnectable  bool   `json:"IsReconnectable"`
	Players          []struct {
		PlayerID        string         `json:"Subject"`
		TeamID          string         `json:"TeamID"`
		CharacterID     string         `json:"CharacterID"`
		PlayerIdentity  PlayerIdentity `json:"PlayerIdentity"`
		IsCoach         bool           `json:"IsCoach"`
		IsAssociated    bool           `json:"IsAssociated"`
		PlatformType    string         `json:"PlatformType"`
		PremierPrestige struct {
		} `json:"PremierPrestige"`
	} `json:"Players"`
	MatchmakingData interface{} `json:"MatchmakingData"`
}

/* [ - Get Game Data - ] */
/* [ URL: 'https://glz-{Region}-1.{Shard}.a.pvp.net/core-game/v1/matches/{MATCH ID}' ] */
func (context ValorantAPIContext) getMainGameData(MatchID string) (*valorantMainGameData, error) {

	APICall := newAPICall()

	APICall.URL = "https://glz-" + context.localRegionData.Region + "-1." + context.localRegionData.Shard + ".a.pvp.net/core-game/v1/matches/" + MatchID
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

	Response := &valorantMainGameData{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

type valorantItemLoadout struct {
	DisplayIcon string
	DisplayName string
}

type valorantMatchTeamPlayer struct {
	Subject                 string  `json:"Subject"`
	CharacterID             AgentID `json:"CharacterID"`
	CharacterName           string
	CharacterDisplayIcon    string
	CharacterSelectionState string `json:"CharacterSelectionState"`
	PeakRank                string
	PeakRankDisplayIcon     string
	CurrentRank             string
	CurrentRankDisplayIcon  string
	LastMatchPartyID        string
	PregamePlayerState      string         `json:"PregamePlayerState"`
	CompetitiveTier         int            `json:"CompetitiveTier"`
	PlayerIdentity          PlayerIdentity `json:"PlayerIdentity"`
	IsCaptain               bool           `json:"IsCaptain"`
	PlatformType            string         `json:"PlatformType"`
	//Expressions             map[string]valorantExpression
	Items map[WeaponID]valorantItemLoadout `json:"Items"`
}

type ValorantMatchTeam struct {
	TeamID     string                    `json:"TeamID"`
	TeamNumber int                       `json:"TeamNumber"`
	Players    []valorantMatchTeamPlayer `json:"Players"`
}

type valorantPreGameData struct {
	ID      string `json:"ID"`
	Version int64  `json:"Version"`
	Teams   []struct {
		TeamID     string `json:"TeamID"`
		TeamNumber int    `json:"TeamNumber"`
	}
	AllyTeam           ValorantMatchTeam `json:"AllyTeam"`
	EnemyTeam          ValorantMatchTeam `json:"EnemyTeam"`
	ObserverSubjects   []interface{}     `json:"ObserverSubjects"`
	MatchCoaches       []interface{}     `json:"MatchCoaches"`
	EnemyTeamSize      int               `json:"EnemyTeamSize"`
	EnemyTeamLockCount int               `json:"EnemyTeamLockCount"`
	PregameState       string            `json:"PregameState"`
	MapID              string            `json:"MapID"`
	MapSelectPool      []interface{}     `json:"MapSelectPool"`
	BannedMapIDs       []interface{}     `json:"BannedMapIDs"`
	CastedVotes        struct {
	} `json:"CastedVotes"`
	MapSelectSteps        []interface{} `json:"MapSelectSteps"`
	MapSelectStep         int           `json:"MapSelectStep"`
	Team1                 string        `json:"Team1"`
	GamePodID             string        `json:"GamePodID"`
	Mode                  string        `json:"Mode"`
	LocalAccessibleAgents []interface{} `json:"LocalAccessibleAgents"`
	VoiceSessionID        string        `json:"VoiceSessionID"`
	MUCName               string        `json:"MUCName"`
	TeamMatchToken        string        `json:"TeamMatchToken"`
	QueueID               string        `json:"QueueID"`
	ProvisioningFlowID    string        `json:"ProvisioningFlowID"`
	IsRanked              bool          `json:"IsRanked"`
	PhaseTimeRemainingNS  int64         `json:"PhaseTimeRemainingNS"`
	StepTimeRemainingNS   int           `json:"StepTimeRemainingNS"`
	AltModesFlagADA       bool          `json:"altModesFlagADA"`
	TournamentMetadata    interface{}   `json:"TournamentMetadata"`
	RosterMetadata        interface{}   `json:"RosterMetadata"`
	GameRules             struct {
		AccoladesEnabled   string `json:"AccoladesEnabled"`
		AllowGameModifiers string `json:"AllowGameModifiers"`
	} `json:"GameRules"`
}

/* [ - Get Pregame Data - ] */
/* [ URL: 'https://glz-{Region}-1.{Shard}.a.pvp.net/pregame/v1/matches/{MATCH ID}' ] */
func (context ValorantAPIContext) GetPreGameData(MatchID string) (*valorantPreGameData, error) {

	APICall := newAPICall()

	APICall.URL = "https://glz-" + context.localRegionData.Region + "-1." + context.localRegionData.Shard + ".a.pvp.net/pregame/v1/matches/" + MatchID
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

	Response := &valorantPreGameData{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [[ - Setup Map ID - ]] */

type MapID string

func (id MapID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

func (id *MapID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*id = MapID(s)
	return nil
}

/* [[ - Setup Character ID - ]] */

type AgentID string

func (id AgentID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

func (id *AgentID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*id = AgentID(s)
	return nil
}

type AgentData struct {
	Status int `json:"status"`
	Data   struct {
		UUID                     string   `json:"uuid"`
		DisplayName              string   `json:"displayName"`
		Description              string   `json:"description"`
		DeveloperName            string   `json:"developerName"`
		DisplayIcon              string   `json:"displayIcon"`
		DisplayIconSmall         string   `json:"displayIconSmall"`
		BustPortrait             string   `json:"bustPortrait"`
		FullPortrait             string   `json:"fullPortrait"`
		FullPortraitV2           string   `json:"fullPortraitV2"`
		KillfeedPortrait         string   `json:"killfeedPortrait"`
		MinimapPortrait          string   `json:"minimapPortrait"`
		Background               string   `json:"background"`
		BackgroundGradientColors []string `json:"backgroundGradientColors"`
	} `json:"data"`
}

/* [[ - Setup Character Data from ID - ]] */
func (id AgentID) GetData() (*AgentData, error) {

	APICall := newAPICall()

	APICall.URL = "https://valorant-api.com/v1/agents/" + string(id)
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

	Response := &AgentData{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	if Response.Status != 200 {
		return nil, errors.New(VALORANT_NO_VALID_RESPONSE)
	}

	return Response, nil

}

/* [[ - Match Data Type - ]] */
type MatchData struct {
	MatchID   string
	IsPregame bool
	AllyTeam  ValorantMatchTeam
	EnemyTeam ValorantMatchTeam
	MapID     MapID
}

type valorantRankStruct struct {
	UUID        string
	CurrentRank int
	PeakRank    int
}

var lastMatchData *MatchData = nil
var lastPreMatchData *MatchData = nil

/* [[ - Get Match Data - ]] */
func (context ValorantAPIContext) GetGameData(useCache bool) (*MatchData, error) {

	LocalPlayer, err := context.GetLocalPlayerContext()
	if err != nil {
		return nil, err
	}

	CurrentMatch, err := context.GetCurrentGamePlayer(LocalPlayer)
	if err != nil {
		return nil, err
	}

	if useCache {

		if lastPreMatchData != nil {

			if lastPreMatchData.MatchID == CurrentMatch.PregameMatchID && CurrentMatch.PregameMatchID != "" {
				return lastPreMatchData, nil
			}

		}

		if lastMatchData != nil {

			if lastMatchData.MatchID == CurrentMatch.GameMatchID && CurrentMatch.GameMatchID != "" {
				return lastMatchData, nil
			}

		}

	}

	CurrentMatchLoadouts, err := context.ValorantGetMatchLoadouts(CurrentMatch, LocalPlayer)
	if err != nil {
		if err.Error() != VALORANT_NO_VALID_RESPONSE {
			return nil, err
		}
	}

	// Check if in Coregame or Pregame

	FinalMatchData := &MatchData{}

	if CurrentMatch.PregameMatchID != "" && CurrentMatch.GameMatchID == "" {

		var err error
		FinalMatchData, err = context.constructPregameData(CurrentMatch, CurrentMatchLoadouts)
		if err != nil {
			return nil, err
		}

		FinalMatchData.IsPregame = true

		lastPreMatchData = FinalMatchData

	}

	if CurrentMatch.GameMatchID != "" {

		var err error
		FinalMatchData, err = context.constructMainGameData(CurrentMatch, LocalPlayer, CurrentMatchLoadouts)
		if err != nil {
			return nil, err
		}

		FinalMatchData.IsPregame = false

		lastMatchData = FinalMatchData

	}

	return FinalMatchData, nil

}

func (context ValorantAPIContext) constructPregameData(CurrentMatch *ValorantCurrentMatch, CurrentMatchLoadout *ValorantMatchLoadouts) (*MatchData, error) {

	//GetValorantPlayer()

	data, err := context.GetPreGameData(CurrentMatch.PregameMatchID)
	if err != nil {
		return nil, err
	}

	ValorantNames := []string{}

	rankChan := make(chan valorantRankStruct)
	var wg sync.WaitGroup

	for _, v := range append(data.AllyTeam.Players, data.EnemyTeam.Players...) {

		ValorantNames = append(ValorantNames, v.Subject)

		wg.Add(1)

		go func() {

			defer wg.Done()

			for attempt := 0; attempt < 5; attempt++ {

				if attempt >= 5 {
					break
				}

				CurrentRank, PeakRank, err := context.GetCurrentAndPeakRank(GetValorantPlayer(v.Subject))
				if err != nil {
					fmt.Println("MMR:", err)

					if err.Error() == VALORANT_API_TIMEOUT {

						time.Sleep(time.Millisecond * 500)

					}
					continue
				}

				rankStruct := valorantRankStruct{
					UUID:        v.Subject,
					CurrentRank: CurrentRank,
					PeakRank:    PeakRank,
				}

				rankChan <- rankStruct
				break

			}

		}()

	}

	go func() {
		wg.Wait()
		close(rankChan)
	}()

	RankLookupTable := map[string]valorantRankStruct{}
	for rankData := range rankChan {
		RankLookupTable[rankData.UUID] = rankData
	}

	TeammateLookup, err := context.EstimatePlayerTeammates(ValorantNames, "pregame"+CurrentMatch.PregameMatchID)
	if err != nil {
		return nil, err
	}

	names, err := context.SubjectsToNames(ValorantNames)
	if err != nil {
		return nil, err
	}

	for index, v := range data.AllyTeam.Players {

		rankData := RankLookupTable[v.Subject]

		Items, err := context.constructItemData(CurrentMatchLoadout, v.Subject)
		if err != nil {
			return nil, err
		}

		data.AllyTeam.Players[index].PlayerIdentity.GameName = names[v.Subject].GameName
		data.AllyTeam.Players[index].PlayerIdentity.TagLine = names[v.Subject].TagLine

		if v.PlayerIdentity.Incognito == true {

			data.AllyTeam.Players[index].PlayerIdentity.GameName = "Player " + strconv.Itoa(index+1)
			data.AllyTeam.Players[index].PlayerIdentity.TagLine = ""

		}

		data.AllyTeam.Players[index].CharacterDisplayIcon = "https://media.valorant-api.com/playercards/" + v.PlayerIdentity.PlayerCardID + "/smallart.png"

		data.AllyTeam.Players[index].PeakRank = CompetitiveTierLookup[rankData.PeakRank]
		data.AllyTeam.Players[index].LastMatchPartyID = TeammateLookup[v.Subject]
		data.AllyTeam.Players[index].PeakRankDisplayIcon = fmt.Sprintf(`https://media.valorant-api.com/competitivetiers/564d8e28-c226-3180-6285-e48a390db8b1/%v/largeicon.png`, rankData.PeakRank)
		data.AllyTeam.Players[index].CurrentRank = CompetitiveTierLookup[rankData.CurrentRank]
		data.AllyTeam.Players[index].CurrentRankDisplayIcon = fmt.Sprintf(`https://media.valorant-api.com/competitivetiers/564d8e28-c226-3180-6285-e48a390db8b1/%v/largeicon.png`, rankData.CurrentRank)
		data.AllyTeam.Players[index].Items = Items

	}

	for index, v := range data.EnemyTeam.Players {

		rankData := RankLookupTable[v.Subject]

		Items, err := context.constructItemData(CurrentMatchLoadout, v.Subject)
		if err != nil {
			return nil, err
		}

		if v.PlayerIdentity.Incognito == true {

			data.AllyTeam.Players[index].PlayerIdentity.GameName = "Player " + strconv.Itoa(index+1)
			data.AllyTeam.Players[index].PlayerIdentity.TagLine = ""

		}

		data.EnemyTeam.Players[index].CharacterDisplayIcon = "https://media.valorant-api.com/playercards/" + v.PlayerIdentity.PlayerCardID + "/smallart.png"

		data.EnemyTeam.Players[index].PeakRank = CompetitiveTierLookup[rankData.PeakRank]
		data.EnemyTeam.Players[index].LastMatchPartyID = TeammateLookup[v.Subject]
		data.EnemyTeam.Players[index].PeakRankDisplayIcon = fmt.Sprintf(`https://media.valorant-api.com/competitivetiers/564d8e28-c226-3180-6285-e48a390db8b1/%v/largeicon.png`, rankData.PeakRank)
		data.EnemyTeam.Players[index].CurrentRank = CompetitiveTierLookup[rankData.CurrentRank]
		data.EnemyTeam.Players[index].CurrentRankDisplayIcon = fmt.Sprintf(`https://media.valorant-api.com/competitivetiers/564d8e28-c226-3180-6285-e48a390db8b1/%v/largeicon.png`, rankData.CurrentRank)
		data.EnemyTeam.Players[index].Items = Items

	}

	FinalMatchData := MatchData{}

	FinalMatchData.MatchID = data.ID
	FinalMatchData.MapID = MapID(data.MapID)
	FinalMatchData.AllyTeam = data.AllyTeam
	FinalMatchData.EnemyTeam = data.EnemyTeam

	return &FinalMatchData, nil

}

func (context ValorantAPIContext) constructMainGameData(CurrentMatch *ValorantCurrentMatch, LocalPlayer *ValorantPlayerContext, CurrentMatchLoadout *ValorantMatchLoadouts) (*MatchData, error) {

	data, err := context.getMainGameData(CurrentMatch.GameMatchID)
	if err != nil {
		return nil, err
	}

	TeamLookup := map[string]string{}
	ValorantNames := []string{}

	rankChan := make(chan valorantRankStruct)
	var wg sync.WaitGroup

	for _, v := range data.Players {

		TeamLookup[v.PlayerIdentity.Subject] = v.TeamID
		ValorantNames = append(ValorantNames, v.PlayerID)

		wg.Add(1)

		go func() {

			defer wg.Done()

			for attempt := 0; attempt < 5; attempt++ {

				if attempt >= 5 {
					break
				}

				CurrentRank, PeakRank, err := context.GetCurrentAndPeakRank(GetValorantPlayer(v.PlayerID))
				if err != nil {
					fmt.Println(err)

					if err.Error() == VALORANT_API_TIMEOUT {

						time.Sleep(time.Millisecond * 500)

					}
					continue
				}

				rankStruct := valorantRankStruct{
					UUID:        v.PlayerID,
					CurrentRank: CurrentRank,
					PeakRank:    PeakRank,
				}

				rankChan <- rankStruct
				break

			}

		}()

	}

	go func() {
		wg.Wait()
		close(rankChan)
	}()

	RankLookupTable := map[string]valorantRankStruct{}
	for rankData := range rankChan {
		RankLookupTable[rankData.UUID] = rankData
	}

	TeammateLookup, err := context.EstimatePlayerTeammates(ValorantNames, "coregame"+CurrentMatch.GameMatchID)
	if err != nil {
		return nil, err
	}

	names, err := context.SubjectsToNames(ValorantNames)
	if err != nil {
		return nil, err
	}

	AllyTeam := ValorantMatchTeam{}
	EnemyTeam := ValorantMatchTeam{}

	AllyTeamIndex := 0
	EnemyTeamIndex := 0

	for _, v := range data.Players {

		TeamOfPlayer := TeamLookup[v.PlayerIdentity.Subject]

		v.PlayerIdentity.GameName = names[v.PlayerID].GameName
		v.PlayerIdentity.TagLine = names[v.PlayerID].TagLine

		if TeamOfPlayer != "" {

			Items, err := context.constructItemData(CurrentMatchLoadout, v.PlayerID)
			if err != nil {
				continue
			}

			agentData, err := AgentID(v.CharacterID).GetData()
			if err != nil {
				continue
			}

			rankData := RankLookupTable[v.PlayerID]

			NewPlayerData := valorantMatchTeamPlayer{
				Subject:                 v.PlayerID,
				CharacterID:             AgentID(v.CharacterID),
				CharacterName:           agentData.Data.DisplayName,
				CharacterDisplayIcon:    agentData.Data.DisplayIcon,
				CharacterSelectionState: "selected",
				PregamePlayerState:      "joined",
				LastMatchPartyID:        TeammateLookup[v.PlayerID],
				PlayerIdentity:          v.PlayerIdentity,
				IsCaptain:               false,
				PlatformType:            v.PlatformType,
				PeakRank:                CompetitiveTierLookup[rankData.PeakRank],
				PeakRankDisplayIcon:     fmt.Sprintf(`https://media.valorant-api.com/competitivetiers/564d8e28-c226-3180-6285-e48a390db8b1/%v/largeicon.png`, rankData.PeakRank),
				CurrentRank:             CompetitiveTierLookup[rankData.CurrentRank],
				CurrentRankDisplayIcon:  fmt.Sprintf(`https://media.valorant-api.com/competitivetiers/564d8e28-c226-3180-6285-e48a390db8b1/%v/largeicon.png`, rankData.CurrentRank),
				Items:                   Items,
			}

			if TeamOfPlayer == TeamLookup[LocalPlayer.UUID] {

				AllyTeamIndex++

				if v.PlayerIdentity.Incognito == true {

					NewPlayerData.PlayerIdentity.GameName = "Player " + strconv.Itoa(AllyTeamIndex)
					NewPlayerData.PlayerIdentity.TagLine = ""

				}

				AllyTeam.Players = append(AllyTeam.Players, NewPlayerData)
				AllyTeam.TeamID = v.TeamID

			} else {

				EnemyTeamIndex++

				if v.PlayerIdentity.Incognito == true {

					NewPlayerData.PlayerIdentity.GameName = "Player " + strconv.Itoa(EnemyTeamIndex)
					NewPlayerData.PlayerIdentity.TagLine = ""

				}

				EnemyTeam.Players = append(EnemyTeam.Players, NewPlayerData)
				EnemyTeam.TeamID = v.TeamID

			}

		}

	}

	FinalMatchData := &MatchData{}

	FinalMatchData.MatchID = data.MatchID
	FinalMatchData.MapID = MapID(data.MapID)
	FinalMatchData.AllyTeam = AllyTeam
	FinalMatchData.EnemyTeam = EnemyTeam

	return FinalMatchData, nil

}

func (context ValorantAPIContext) constructItemData(CurrentMatchLoadout *ValorantMatchLoadouts, playerID string) (map[WeaponID]valorantItemLoadout, error) {

	Items := map[WeaponID]valorantItemLoadout{}

	for _, v := range CurrentMatchLoadout.Loadouts[playerID].Loadout.Items {

		displayName := ""
		displayIcon := ""

		data1, err := WeaponSkinID(v.Sockets[TypeID("bcef87d6-209b-46c6-8b19-fbe40bd95abc")].Item.ID).GetInformation()
		if err != nil {
			return nil, err
		}

		displayName = data1.Data.DisplayName
		displayIcon = data1.Data.DisplayIcon

		data2, err := WeaponChromaID(v.Sockets[TypeID("3ad1b2b2-acdb-4524-852f-954a76ddae0a")].Item.ID).GetInformation() // Get Skin Chroma
		if err != nil {
			return nil, err
		}

		displayIcon = data2.Data.FullRender
		if displayName == "" {
			displayName = data2.Data.DisplayName
		}

		Items[v.ID] = valorantItemLoadout{
			DisplayName: displayName,
			DisplayIcon: displayIcon,
		}

	}

	return Items, nil

}

/* [[ - Match Loadout Data - ]] */

type ValorantPlayerMatchLoadoutItems struct {
	ID      WeaponID `json:"ID"`
	TypeID  TypeID   `json:"TypeID"`
	Sockets map[TypeID]struct {
		Item struct {
			ID     string `json:"ID"`
			TypeID TypeID `json:"TypeID"`
		} `json:"Item"`
	} `json:"Sockets"`
}

type internalValorantLoadout struct {
	Subject     string `json:"Subject"`
	Expressions struct {
		AESSelections []struct {
			AssetID string `json:"AssetID"`
			TypeID  string `json:"TypeID"`
		} `json:"AESSelections"`
	} `json:"Expressions"`
	Items map[string]ValorantPlayerMatchLoadoutItems `json:"Items"`
}

type ValorantPlayerMatchLoadout struct {
	Loadout internalValorantLoadout `json:"Loadout"`
}

type internalValorantPregamePlayerMatchLoadout struct {
	Subject     string `json:"Subject"`
	Expressions struct {
		AESSelections []struct {
			AssetID string `json:"AssetID"`
			TypeID  string `json:"TypeID"`
		} `json:"AESSelections"`
	} `json:"Expressions"`
	Items map[string]ValorantPlayerMatchLoadoutItems `json:"Items"`
}

type ValorantMatchLoadouts struct {
	Loadouts map[string]ValorantPlayerMatchLoadout `json:"Loadouts"`
}
type internalValorantMatchLoadouts struct {
	Loadouts []ValorantPlayerMatchLoadout `json:"Loadouts"`
}
type internalValorantPregameMatchLoadouts struct {
	Loadouts []internalValorantPregamePlayerMatchLoadout `json:"Loadouts"`
}

var LoadoutCache = map[string]*ValorantMatchLoadouts{}

/* [[ - Get Match Loadout Data - ]] */
func (context ValorantAPIContext) ValorantGetMatchLoadouts(CurrentMatch *ValorantCurrentMatch, LocalPlayer *ValorantPlayerContext) (*ValorantMatchLoadouts, error) {

	usedID := ""

	if CurrentMatch.PregameMatchID != "" && CurrentMatch.GameMatchID == "" {
		usedID = CurrentMatch.PregameMatchID
	} else {
		usedID = CurrentMatch.GameMatchID
	}

	if LoadoutCache[usedID] != nil {
		return LoadoutCache[usedID], nil
	}

	Response := &internalValorantMatchLoadouts{}

	if CurrentMatch.PregameMatchID != "" && CurrentMatch.GameMatchID == "" {

		pregameLoadout, err := context.valorantGetPreGameMatchLoadouts(CurrentMatch)
		if err != nil {
			return nil, err
		}

		Response = pregameLoadout

	} else {

		coregameLoadout, err := context.valorantGetCoreGameMatchLoadouts(CurrentMatch)
		if err != nil {
			return nil, err
		}

		Response = coregameLoadout

	}

	Final := &ValorantMatchLoadouts{}
	Final.Loadouts = map[string]ValorantPlayerMatchLoadout{}

	for _, v := range Response.Loadouts {
		Final.Loadouts[v.Loadout.Subject] = v
	}

	for k := range LoadoutCache {
		delete(LoadoutCache, k)
	}

	LoadoutCache[usedID] = Final

	return Final, nil

}

/* [[ - Get PreGame Match Loadout Data - ]] */
func (context ValorantAPIContext) valorantGetPreGameMatchLoadouts(CurrentMatch *ValorantCurrentMatch) (*internalValorantMatchLoadouts, error) {

	APICall := newAPICall()

	APICall.URL = "https://glz-" + context.localRegionData.Region + "-1." + context.localRegionData.Shard + ".a.pvp.net/pregame/v1/matches/" + CurrentMatch.PregameMatchID + "/loadouts"

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

	GeneralResponse := &internalValorantMatchLoadouts{}
	PregameResponse := &internalValorantPregameMatchLoadouts{}

	err = json.Unmarshal(FullBody, &PregameResponse)
	if err != nil {
		return nil, err
	}

	GeneralResponse.Loadouts = make([]ValorantPlayerMatchLoadout, len(PregameResponse.Loadouts))

	for i, v := range PregameResponse.Loadouts {

		internalElement := ValorantPlayerMatchLoadout{
			Loadout: internalValorantLoadout{
				Subject:     v.Subject,
				Expressions: v.Expressions,
				Items:       v.Items,
			},
		}

		GeneralResponse.Loadouts[i] = internalElement

	}

	return GeneralResponse, nil

}

/* [[ - Get CoreGame Match Loadout Data - ]] */
func (context ValorantAPIContext) valorantGetCoreGameMatchLoadouts(CurrentMatch *ValorantCurrentMatch) (*internalValorantMatchLoadouts, error) {

	APICall := newAPICall()

	APICall.URL = "https://glz-" + context.localRegionData.Region + "-1." + context.localRegionData.Shard + ".a.pvp.net/core-game/v1/matches/" + CurrentMatch.GameMatchID + "/loadouts"

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

	Response := &internalValorantMatchLoadouts{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [[ - Select Agent - ]] */
func (context ValorantAPIContext) SelectAgent(CurrentPregameID string, agentID AgentID) (*MatchData, error) {

	APICall := newAPICall()

	APICall.URL = "https://glz-" + context.localRegionData.Region + "-1." + context.localRegionData.Shard + ".a.pvp.net/pregame/v1/matches/" + CurrentPregameID + "/select/" + string(agentID)

	APICall.Body = nil
	APICall.Method = "POST"

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

	return nil, nil

}
