package valorantapi

import (
	"encoding/json"
	"fmt"
	"io"
)

type OwnedItemsAgents struct {
	ItemTypeID   string `json:"ItemTypeID"`
	Entitlements []struct {
		TypeID string      `json:"TypeID"`
		ItemID string      `json:"ItemID"`
		Tiers  interface{} `json:"Tiers"`
	} `json:"Entitlements"`
}

var OwnedAgentsCache map[string]AgentID

func (context ValorantAPIContext) GetOwnedAgentData() (map[string]AgentID, error) {

	if OwnedAgentsCache != nil {
		return OwnedAgentsCache, nil
	}

	LocalPlayer, err := context.GetLocalPlayerContext()
	if err != nil {
		return nil, err
	}

	APICall := newAPICall()

	APICall.URL = "https://pd." + context.localRegionData.Shard + ".a.pvp.net/store/v1/entitlements/" + LocalPlayer.UUID + "/01bb38e1-da47-4e6a-9b3d-945fe4655707"

	APICall.Body = nil
	APICall.Method = "GET"

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

	Response := &OwnedItemsAgents{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	OwnedAgentLookup := map[string]AgentID{}

	for _, v := range Response.Entitlements {

		OwnedAgentLookup[v.ItemID] = AgentID(v.ItemID)

	}

	fmt.Println(OwnedAgentLookup)

	OwnedAgentsCache = OwnedAgentLookup

	return OwnedAgentLookup, nil

}
