package valorantapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

/* [ - Get Valorant Content - ] */
/* [ URL: 'https://shared.{Shard}.a.pvp.net/content-service/v3/content' ] */

type ValorantFetchedSeason struct {
	ID        string    `json:"ID"`
	Name      string    `json:"Name"`
	Type      string    `json:"Type"`
	StartTime time.Time `json:"StartTime"`
	EndTime   time.Time `json:"EndTime"`
	IsActive  bool      `json:"IsActive"`
}

type ValorantFetchedEvent struct {
	ID        string    `json:"ID"`
	Name      string    `json:"Name"`
	StartTime time.Time `json:"StartTime"`
	EndTime   time.Time `json:"EndTime"`
	IsActive  bool      `json:"IsActive"`
}

type ValorantFetchedContent struct {
	DisabledIDs []string                `json:"DisabledIDs"`
	Seasons     []ValorantFetchedSeason `json:"Seasons"`
	Events      []ValorantFetchedEvent  `json:"Events"`
}

func (context ValorantAPIContext) GetFetchContent() (*ValorantFetchedContent, error) {

	APICall := newAPICall()

	fmt.Println("Shard:", context.localRegionData.Shard)

	APICall.URL = "https://shared." + context.localRegionData.Shard + ".a.pvp.net/content-service/v3/content"
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

	Response := &ValorantFetchedContent{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [ - Get Account XP - ] */
/* [ URL: 'https://pd.{Shard}.a.pvp.net/account-xp/v1/players/{PUUID}' ] */

type ValorantAccountXPProgress struct {
	Level int `json:"Level"`
	XP    int `json:"XP"`
}

type ValorantAccountXPHistory struct {
	ID            string    `json:"ID"`
	MatchStart    time.Time `json:"MatchStart"`
	StartProgress struct {
		Level int `json:"Level"`
		XP    int `json:"XP"`
	} `json:"StartProgress"`
	EndProgress struct {
		Level int `json:"Level"`
		XP    int `json:"XP"`
	} `json:"EndProgress"`
	XPDelta   int `json:"XPDelta"`
	XPSources []struct {
		ID     string `json:"ID"`
		Amount int    `json:"Amount"`
	} `json:"XPSources"`
	XPMultipliers []interface{} `json:"XPMultipliers"`
}

type ValorantAccountXP struct {
	Version                   int                        `json:"Version"`
	UUID                      string                     `json:"Subject"`
	Progress                  ValorantAccountXPProgress  `json:"Progress"`
	History                   []ValorantAccountXPHistory `json:"History"`
	LastTimeGrantedFirstWin   time.Time                  `json:"LastTimeGrantedFirstWin"`
	NextTimeFirstWinAvailable time.Time                  `json:"NextTimeFirstWinAvailable"`
}

func (context ValorantAPIContext) GetAccountXP(playerContext ValorantPlayerContext) (*ValorantAccountXP, error) {

	APICall := newAPICall()

	APICall.URL = "https://pd." + context.localRegionData.Shard + ".a.pvp.net/account-xp/v1/players/" + playerContext.UUID
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

	Response := &ValorantAccountXP{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [ - Get Local Loadout - ] */
/* [ URL: 'https://pd.{Shard}.a.pvp.net/personalization/v2/players/{PUUID}/playerloadout' ] */

type ValorantLocalLoadoutGuns struct {
	SkinName        string
	ID              WeaponID          `json:"ID"`
	SkinID          WeaponSkinID      `json:"SkinID"`
	SkinLevelID     WeaponSkinLevelID `json:"SkinLevelID"`
	ChromaID        WeaponChromaID    `json:"ChromaID"`
	CharmInstanceID string            `json:"CharmInstanceID,omitempty"`
	CharmID         string            `json:"CharmID,omitempty"`
	CharmLevelID    string            `json:"CharmLevelID,omitempty"`
	Attachments     []interface{}     `json:"Attachments"`
}

type ValorantLocalExpression struct {
	TypeID  string `json:"TypeID"`
	AssetID string `json:"AssetID"`
}

type ValorantLocalLoadoutIdentity struct {
	PlayerCardID           string  `json:"PlayerCardID"`
	PlayerTitleID          TitleID `json:"PlayerTitleID"`
	AccountLevel           int     `json:"AccountLevel"`
	PreferredLevelBorderID string  `json:"PreferredLevelBorderID"`
	Incognito              bool    `json:"Incognito"`
	HideAccountLevel       bool    `json:"HideAccountLevel"`
}

type ValorantLocalLoadout struct {
	Subject  string                       `json:"Subject"`
	Version  int                          `json:"Version"`
	Guns     []ValorantLocalLoadoutGuns   `json:"Guns"`
	Sprays   []ValorantLocalExpression    `json:"ActiveExpressions"`
	Identity ValorantLocalLoadoutIdentity `json:"Identity"`
}

/* Gets currently active loadout (ONLY WORKS ON LOCAL PLAYER) */
func (context ValorantAPIContext) GetAccountLoadout(playerContext *ValorantPlayerContext) (*ValorantLocalLoadout, error) {

	APICall := newAPICall()

	APICall.URL = "https://pd." + context.localRegionData.Shard + ".a.pvp.net/personalization/v3/players/" + playerContext.UUID + "/playerloadout"
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

	Response := &ValorantLocalLoadout{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* Sets currently active loadout (ONLY WORKS ON LOCAL PLAYER) + (DOESN'T SHOW CHANGES, BUT CHANGES STILL MADE) */
func (context ValorantAPIContext) SetAccountLoadout(loadout ValorantLocalLoadout, playerContext *ValorantPlayerContext) error {

	jsonBytes, err := json.MarshalIndent(loadout, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	APICall := newAPICall()

	APICall.URL = "https://pd." + context.localRegionData.Shard + ".a.pvp.net/personalization/v3/players/" + playerContext.UUID + "/playerloadout"
	APICall.Body = strings.NewReader(string(jsonBytes))
	APICall.Method = "PUT"

	err = APICall.SetRequest()
	if err != nil {
		return err
	}

	APICall.SetBearerAuth(context)

	err = APICall.Call()
	if err != nil {
		return err
	}

	return nil

}
