package valorantapi

import (
	"encoding/json"
)

/* [[ - Setup Type ID - ]] */

type TypeID string

func (id TypeID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

func (id *TypeID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*id = TypeID(s)
	return nil
}

func (id *TypeID) GetTypeName() string {

	switch string(*id) {
	case "03a572de-4234-31ed-d344-ababa488f981":
		return "Flex"
	case "d5f120f8-ff8c-4aac-92ea-f2b5acbe9475":
		return "Spray"

	case "51c9eb99-3e6b-4658-801f-a5a7fd64bb9d":
		return "Weapon"
	case "bcef87d6-209b-46c6-8b19-fbe40bd95abc":
		return "Skin"
	case "3ad1b2b2-acdb-4524-852f-954a76ddae0a":
		return "Skin Chroma"
	case "e7c63390-eda7-46e0-bb7a-a6abdacd2433":
		return "Skin Level"

	default:
		return "Could not identify"
	}

}

func UseTypeFunction[T any](typeID TypeID, id string) (*T, error) {

	switch string(typeID) {
	case "03a572de-4234-31ed-d344-ababa488f981":
		return any("Flex").(*T), nil // Flex
	case "d5f120f8-ff8c-4aac-92ea-f2b5acbe9475":
		return any("Spray").(*T), nil // Spray

	case "51c9eb99-3e6b-4658-801f-a5a7fd64bb9d":
		weapon, err := WeaponID(id).GetInformation()
		return any(weapon).(*T), err // Returns ( *ValorantWeaponIDSingle, error )
	case "bcef87d6-209b-46c6-8b19-fbe40bd95abc":
		skin, err := WeaponSkinID(id).GetInformation()
		return any(skin).(*T), err // Returns ( *ValorantWeaponSkinIDSingle, error )
	case "3ad1b2b2-acdb-4524-852f-954a76ddae0a":
		chroma, err := WeaponChromaID(id).GetInformation()
		return any(chroma).(*T), err // Returns ( *ValorantWeaponChromaSingle, error )
	case "e7c63390-eda7-46e0-bb7a-a6abdacd2433":
		skin_level, err := WeaponSkinLevelID(id).GetInformation()
		return any(skin_level).(*T), err // Returns ( *ValorantWeaponLevelSingle, error )

	default:
		return any("Could not identify").(*T), nil
	}

}
