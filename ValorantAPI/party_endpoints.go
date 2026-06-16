package valorantapi

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

type ClientPlatform struct {
	PlatformType      string `json:"platformType"`
	PlatformOS        string `json:"platformOS"`
	PlatformOSVersion string `json:"platformOSVersion"`
	PlatformChipset   string `json:"platformChipset"`
	PlatformDevice    string `json:"platformDevice"`
}

type ValorantPartyContext struct {
	Subject            string         `json:"Subject"`
	Version            int64          `json:"Version"`
	PartyID            string         `json:"CurrentPartyID"`
	Invites            interface{}    `json:"Invites"`
	Requests           []interface{}  `json:"Requests"`
	PlatformInfo       ClientPlatform `json:"PlatformInfo"`
	IsCrossPlayEnabled bool           `json:"IsCrossPlayEnabled"`
}

/* Gets party context of player (ONLY WORKS ON LOCAL PLAYER) */
func (context ValorantAPIContext) GetPartyContext(playerContext ValorantPlayerContext) (*ValorantPartyContext, error) {

	APICall := newAPICall()

	APICall.URL = "https://glz-" + context.localRegionData.Region + "-1." + context.localRegionData.Shard + ".a.pvp.net/parties/v1/players/" + playerContext.UUID
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

	Response := &ValorantPartyContext{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

type ServerPing struct {
	Ping      int    `json:"Ping"`
	GamePodID string `json:"GamePodID"`
}

type ValorantPartyMember struct {
	PlayerUUID                          string                       `json:"Subject"`
	CompetitiveTier                     int                          `json:"CompetitiveTier"`
	PlayerIdentity                      ValorantLocalLoadoutIdentity `json:"PlayerIdentity"`
	SeasonalBadgeInfo                   interface{}                  `json:"SeasonalBadgeInfo"`
	IsOwner                             bool                         `json:"IsOwner"`
	QueueEligibleRemainingAccountLevels int                          `json:"QueueEligibleRemainingAccountLevels"`
	Pings                               []ServerPing                 `json:"Pings"`
	IsReady                             bool                         `json:"IsReady"`
	IsModerator                         bool                         `json:"IsModerator"`
	UseBroadcastHUD                     bool                         `json:"UseBroadcastHUD"`
	PlatformType                        string                       `json:"PlatformType"`
	PreferredAgentInfo                  interface{}                  `json:"PreferredAgentInfo"`
	PremierPrestige                     struct{}                     `json:"PremierPrestige"`
}

type Membership struct {
	TeamOne        interface{} `json:"teamOne"`
	TeamTwo        interface{} `json:"teamTwo"`
	TeamSpectate   interface{} `json:"teamSpectate"`
	TeamOneCoaches interface{} `json:"teamOneCoaches"`
	TeamTwoCoaches interface{} `json:"teamTwoCoaches"`
}

type CustomGameSettingsData struct {
	Map       string      `json:"Map"`
	Mode      string      `json:"Mode"`
	UseBots   bool        `json:"UseBots"`
	GamePod   string      `json:"GamePod"`
	GameRules interface{} `json:"GameRules"`
}

type CustomGameData struct {
	Settings              CustomGameSettingsData `json:"Settings"`
	Membership            Membership             `json:"Membership"`
	MaxPartySize          int                    `json:"MaxPartySize"`
	AutoBalanceEnabled    bool                   `json:"AutobalanceEnabled"`
	AutoBalanceMinPlayers int                    `json:"AutobalanceMinPlayers"`
	HasRecoveryData       bool                   `json:"HasRecoveryData"`
}

type ValorantParty struct {
	ID                    string                `json:"ID"`
	MUCName               string                `json:"MUCName"`
	VoiceRoomID           string                `json:"VoiceRoomID"`
	Version               int64                 `json:"Version"`
	ClientVersion         string                `json:"ClientVersion"`
	Members               []ValorantPartyMember `json:"Members"`
	State                 string                `json:"State"`
	PreviousState         string                `json:"PreviousState"`
	StateTransitionReason string                `json:"StateTransitionReason"`
	Accessibility         string                `json:"Accessibility"`
	CustomGameData        CustomGameData        `json:"CustomGameData"`
	MatchmakingData       struct {
		QueueID                 string   `json:"QueueID"`
		PreferredGamePods       []string `json:"PreferredGamePods"`
		SkillDisparityRRPenalty int      `json:"SkillDisparityRRPenalty"`
	} `json:"MatchmakingData"`
	Invites           interface{}   `json:"Invites"`
	Requests          []interface{} `json:"Requests"`
	ExternalInvites   interface{}   `json:"ExternalInvites"`
	QueueEntryTime    time.Time     `json:"QueueEntryTime"`
	ErrorNotification struct {
		ErrorType      string      `json:"ErrorType"`
		ErroredPlayers interface{} `json:"ErroredPlayers"`
	} `json:"ErrorNotification"`
	RestrictedSeconds    int           `json:"RestrictedSeconds"`
	EligibleQueues       []string      `json:"EligibleQueues"`
	QueueIneligibilities []interface{} `json:"QueueIneligibilities"`
	CheatData            struct {
		GamePodOverride         string `json:"GamePodOverride"`
		ForcePostGameProcessing bool   `json:"ForcePostGameProcessing"`
	} `json:"CheatData"`
	XPBonuses  []interface{} `json:"XPBonuses"`
	InviteCode string        `json:"InviteCode"`
}

/* Gets current party from context (ONLY WORKS ON LOCAL PLAYER) */
func (context ValorantAPIContext) GetParty(partyContext ValorantPartyContext) (*ValorantParty, error) {

	APICall := newAPICall()

	APICall.URL = "https://glz-" + context.localRegionData.Region + "-1." + context.localRegionData.Shard + ".a.pvp.net/parties/v1/parties/" + partyContext.PartyID
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

	Response := &ValorantParty{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

type teammateChanItem struct {
	PlayerUUID string
	PartyID    string
}

var teammateCache = map[string]map[string]string{}

/* Estimates the partyID of the player */
func (context ValorantAPIContext) EstimatePlayerTeammates(playerIDs []string, matchID string) (map[string]string, error) {

	if teammateCache[matchID] != nil {
		return teammateCache[matchID], nil
	}

	Lookup := map[string]string{}

	outputChan := make(chan teammateChanItem)
	var wg sync.WaitGroup

	for _, v := range playerIDs {

		wg.Add(1)

		go func() {

			defer wg.Done()

			LastMatch, err := context.ValorantLastMatch(v, "0", "1")
			if err != nil {
				fmt.Println(err)
				return
			}

			if len(LastMatch.History) <= 0 {
				return
			}

			ValorantPartyDetails, err := context.ValorantMatchPartyDetails(LastMatch.History[0].MatchID)
			if err != nil {
				fmt.Println(err)
				return
			}

			ValorantPartyLookup := map[string]string{}

			for _, w := range ValorantPartyDetails.Players {

				ValorantPartyLookup[w.Subject] = w.PartyID
			}

			rankStruct := teammateChanItem{
				PlayerUUID: v,
				PartyID:    ValorantPartyLookup[v],
			}

			outputChan <- rankStruct

		}()

	}

	go func() {
		wg.Wait()
		close(outputChan)
	}()

	for teammateLookup := range outputChan {
		Lookup[teammateLookup.PlayerUUID] = teammateLookup.PartyID
	}

	for k := range teammateCache {
		delete(teammateCache, k)
	}

	teammateCache[matchID] = Lookup

	return Lookup, nil

}

type ValorantMatchHistory struct {
	Subject    string `json:"Subject"`
	BeginIndex int    `json:"BeginIndex"`
	EndIndex   int    `json:"EndIndex"`
	Total      int    `json:"Total"`
	History    []struct {
		MatchID       string `json:"MatchID"`
		GameStartTime int64  `json:"GameStartTime"`
		QueueID       string `json:"QueueID"`
	} `json:"History"`
}

func (context ValorantAPIContext) ValorantLastMatch(playerUUID string, startIndex string, endIndex string) (*ValorantMatchHistory, error) {

	APICall := newAPICall()

	APICall.URL = "https://pd." + context.localRegionData.Shard + ".a.pvp.net/match-history/v1/history/" + playerUUID + "?startIndex=" + startIndex + "&endIndex=" + endIndex
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

	Response := &ValorantMatchHistory{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

type ValorantMatchPartyDetails struct {
	Players []struct {
		Subject string `json:"subject"`
		TeamID  string `json:"teamId"`
		PartyID string `json:"partyId"`
	} `json:"players"`
}

func (context ValorantAPIContext) ValorantMatchPartyDetails(matchID string) (*ValorantMatchPartyDetails, error) {

	APICall := newAPICall()

	APICall.URL = "https://pd." + context.localRegionData.Shard + ".a.pvp.net/match-details/v1/matches/" + matchID
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

	Response := &ValorantMatchPartyDetails{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}
