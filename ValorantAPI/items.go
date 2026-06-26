package valorantapi

import (
	"encoding/json"
	"io"
)

var OwnedWeapons map[string]*WeaponChromaID

func (context ValorantAPIContext) GetOwnedWeaponChromas(useCache bool) (map[string]*WeaponChromaID, error) {

	if OwnedWeapons != nil && useCache {
		return OwnedWeapons, nil
	}

	LocalPlayer, err := context.GetLocalPlayerContext()
	if err != nil {
		return nil, err
	}

	APICall := newAPICall()

	APICall.URL = "https://pd." + context.localRegionData.Shard + ".a.pvp.net/store/v1/entitlements/" + LocalPlayer.UUID + "/3ad1b2b2-acdb-4524-852f-954a76ddae0a"

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

	Response := &OwnedItems{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	OwnedWeaponLookup := map[string]*WeaponChromaID{}

	for _, v := range Response.Entitlements {
		chroma := WeaponChromaID(v.ItemID)
		OwnedWeaponLookup[v.ItemID] = &chroma

	}

	OwnedWeapons = OwnedWeaponLookup

	return OwnedWeaponLookup, nil

}

var OwnedWeaponSkins map[string]*WeaponSkinID

func (context ValorantAPIContext) GetOwnedWeaponSkins(useCache bool) (map[string]*WeaponSkinID, error) {

	if OwnedWeaponSkins != nil && useCache {
		return OwnedWeaponSkins, nil
	}

	LocalPlayer, err := context.GetLocalPlayerContext()
	if err != nil {
		return nil, err
	}

	APICall := newAPICall()

	APICall.URL = "https://pd." + context.localRegionData.Shard + ".a.pvp.net/store/v1/entitlements/" + LocalPlayer.UUID + "/e7c63390-eda7-46e0-bb7a-a6abdacd2433"

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

	Response := &OwnedItems{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	OwnedWeaponLookup := map[string]*WeaponSkinID{}

	for _, v := range Response.Entitlements {
		skinID := WeaponSkinID(v.ItemID)
		OwnedWeaponLookup[v.ItemID] = &skinID

	}

	OwnedWeaponSkins = OwnedWeaponLookup

	return OwnedWeaponLookup, nil

}
