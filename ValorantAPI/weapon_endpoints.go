package valorantapi

import (
	"encoding/json"
	"errors"
	"io"
)

/* [ - Get All Weapons - ] */
/* [ URL: 'https://valorant-api.com/v1/weapons' ] */

type ValorantAllWeaponsStats struct {
	FireRate            float64     `json:"fireRate"`
	MagazineSize        int         `json:"magazineSize"`
	RunSpeedMultiplier  float64     `json:"runSpeedMultiplier"`
	EquipTimeSeconds    float64     `json:"equipTimeSeconds"`
	ReloadTimeSeconds   float64     `json:"reloadTimeSeconds"`
	FirstBulletAccuracy float64     `json:"firstBulletAccuracy"`
	ShotgunPelletCount  int         `json:"shotgunPelletCount"`
	WallPenetration     string      `json:"wallPenetration"`
	Feature             string      `json:"feature"`
	FireMode            interface{} `json:"fireMode"`
	AltFireType         string      `json:"altFireType"`
	AdsStats            struct {
		ZoomMultiplier      float64 `json:"zoomMultiplier"`
		FireRate            float64 `json:"fireRate"`
		RunSpeedMultiplier  float64 `json:"runSpeedMultiplier"`
		BurstCount          int     `json:"burstCount"`
		FirstBulletAccuracy float64 `json:"firstBulletAccuracy"`
	} `json:"adsStats"`
	AltShotgunStats interface{} `json:"altShotgunStats"`
	AirBurstStats   interface{} `json:"airBurstStats"`
	DamageRanges    []struct {
		RangeStartMeters int     `json:"rangeStartMeters"`
		RangeEndMeters   int     `json:"rangeEndMeters"`
		HeadDamage       float64 `json:"headDamage"`
		BodyDamage       float64 `json:"bodyDamage"`
		LegDamage        float64 `json:"legDamage"`
	} `json:"damageRanges"`
}

type ValorantAllWeaponsShopData struct {
	Cost              int    `json:"cost"`
	Category          string `json:"category"`
	ShopOrderPriority int    `json:"shopOrderPriority"`
	CategoryText      string `json:"categoryText"`
	GridPosition      struct {
		Row    int `json:"row"`
		Column int `json:"column"`
	} `json:"gridPosition"`
	CanBeTrashed bool        `json:"canBeTrashed"`
	Image        interface{} `json:"image"`
	NewImage     string      `json:"newImage"`
	NewImage2    interface{} `json:"newImage2"`
	AssetPath    string      `json:"assetPath"`
}

type WeaponSkinChromas struct {
	UUID          string      `json:"uuid"`
	DisplayName   string      `json:"displayName"`
	DisplayIcon   interface{} `json:"displayIcon"`
	FullRender    string      `json:"fullRender"`
	Swatch        interface{} `json:"swatch"`
	StreamedVideo interface{} `json:"streamedVideo"`
	AssetPath     string      `json:"assetPath"`
}

type WeaponSkinLevels struct {
	UUID          string      `json:"uuid"`
	DisplayName   string      `json:"displayName"`
	LevelItem     interface{} `json:"levelItem"`
	DisplayIcon   string      `json:"displayIcon"`
	StreamedVideo interface{} `json:"streamedVideo"`
	AssetPath     string      `json:"assetPath"`
}

type ValorantAllWeaponsSkins struct {
	UUID            string              `json:"uuid"`
	DisplayName     string              `json:"displayName"`
	ThemeUUID       string              `json:"themeUuid"`
	ContentTierUUID string              `json:"contentTierUuid"`
	DisplayIcon     string              `json:"displayIcon"`
	Wallpaper       interface{}         `json:"wallpaper"`
	AssetPath       string              `json:"assetPath"`
	Chromas         []WeaponSkinChromas `json:"chromas"`
	Levels          []WeaponSkinLevels  `json:"levels"`
}

type ValorantAllWeapon struct {
	UUID            string                     `json:"uuid"`
	DisplayName     string                     `json:"displayName"`
	Category        string                     `json:"category"`
	DefaultSkinUUID string                     `json:"defaultSkinUuid"`
	DisplayIcon     string                     `json:"displayIcon"`
	KillStreamIcon  string                     `json:"killStreamIcon"`
	AssetPath       string                     `json:"assetPath"`
	WeaponStats     ValorantAllWeaponsStats    `json:"weaponStats"`
	ShopData        ValorantAllWeaponsShopData `json:"shopData"`
	Skins           []ValorantAllWeaponsSkins  `json:"skins"`
}

type ValorantAllWeapons struct {
	Status int                 `json:"status"`
	Data   []ValorantAllWeapon `json:"data"`
}

/* Gets every single weapon valorant has ( All skins + All Chromas + All Weapon Stats + All weapon upgrades ) */
func (context ValorantAPIContext) GetAllWeapons() (*ValorantAllWeapons, error) {

	APICall := newAPICall()

	APICall.URL = "https://valorant-api.com/v1/weapons"
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

	Response := &ValorantAllWeapons{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	if Response.Status != 200 {
		return nil, errors.New(VALORANT_NO_VALID_RESPONSE)
	}

	return Response, nil

}

/* [ - Get All Weapon Chromas - ] */
/* [ URL: 'https://valorant-api.com/v1/weapons/skinchromas' ] */

type ValorantAllWeaponChromas struct {
	Status int                 `json:"status"`
	Data   []WeaponSkinChromas `json:"data"`
}

/* Gets every single weapon skin variant */
func (context ValorantAPIContext) GetAllWeaponChromas() (*ValorantAllWeaponChromas, error) {

	APICall := newAPICall()

	APICall.URL = "https://valorant-api.com/v1/weapons/skinchromas"
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

	Response := &ValorantAllWeaponChromas{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	if Response.Status != 200 {
		return nil, errors.New(VALORANT_NO_VALID_RESPONSE)
	}

	return Response, nil

}

/* [ - Get All Weapon Skin Levels - ] */
/* [ URL: 'https://valorant-api.com/v1/weapons/skinlevels' ] */

type ValorantAllWeaponSkinLevels struct {
	Status int                `json:"status"`
	Data   []WeaponSkinLevels `json:"data"`
}

/* Gets every single weapon skin upgrade */
func (context ValorantAPIContext) GetAllWeaponSkinLevels() (*ValorantAllWeaponSkinLevels, error) {

	APICall := newAPICall()

	APICall.URL = "https://valorant-api.com/v1/weapons/skinlevels"
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

	Response := &ValorantAllWeaponSkinLevels{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	if Response.Status != 200 {
		return nil, errors.New(VALORANT_NO_VALID_RESPONSE)
	}

	return Response, nil

}

const (
	VALORANT_WEAPON_ID            = 0
	VALORANT_WEAPON_SKIN_ID       = 1
	VALORANT_WEAPON_CHROMA_ID     = 2
	VALORANT_WEAPON_SKIN_LEVEL_ID = 3
)

/* [[ - Setup Weapon ID - ]] */

type WeaponID string

func (id WeaponID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

func (id *WeaponID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*id = WeaponID(s)
	return nil
}

/* [[ - Setup TitleID - ]] */

type TitleID string

func (id TitleID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

func (id *TitleID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*id = TitleID(s)
	return nil
}

/* [[ - Setup Weapon SkinID - ]] */

type WeaponSkinID string

func (id WeaponSkinID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

func (id *WeaponSkinID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*id = WeaponSkinID(s)
	return nil
}

/* [[ - Setup Weapon ChromaID - ]] */

type WeaponChromaID string

func (id WeaponChromaID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

func (id *WeaponChromaID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*id = WeaponChromaID(s)
	return nil
}

/* [[ - Setup Weapon SkinLevelID - ]] */

type WeaponSkinLevelID string

func (id WeaponSkinLevelID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

func (id *WeaponSkinLevelID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*id = WeaponSkinLevelID(s)
	return nil
}

/* [ - Get Single weapon from ID - ] */
/* [ URL: 'https://valorant-api.com/v1/weapons/{WeaponID}' ] */

type ValorantWeaponIDSingle struct {
	Status int               `json:"status"`
	Data   ValorantAllWeapon `json:"data"`
}

/* Gets single weapon type from ID */
func (weaponID WeaponID) GetInformation() (*ValorantWeaponIDSingle, error) {

	APICall := newAPICall()

	APICall.URL = "https://valorant-api.com/v1/weapons/" + string(weaponID)
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

	Response := &ValorantWeaponIDSingle{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	if Response.Status != 200 {
		return nil, errors.New(VALORANT_NO_VALID_RESPONSE)
	}

	return Response, nil

}

/* [ - Get Single Weapon Skin from ID - ] */
/* [ URL: 'https://valorant-api.com/v1/weapons/skins/{SkinID}' ] */

type ValorantWeaponSkinIDSingle struct {
	Status int                     `json:"status"`
	Data   ValorantAllWeaponsSkins `json:"data"`
}

/* Gets single weapon type from ID */
func (weaponSkinID WeaponSkinID) GetInformation() (*ValorantWeaponSkinIDSingle, error) {

	APICall := newAPICall()

	APICall.URL = "https://valorant-api.com/v1/weapons/skins/" + string(weaponSkinID)
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

	Response := &ValorantWeaponSkinIDSingle{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	if Response.Status != 200 {
		return nil, errors.New(VALORANT_NO_VALID_RESPONSE)
	}

	return Response, nil

}

/* [ - Get Single Chroma from ID - ] */
/* [ URL: 'https://valorant-api.com/v1/weapons/skinchromas/{SkinChromaID}' ] */

type ValorantWeaponChromaSingle struct {
	Status int               `json:"status"`
	Data   WeaponSkinChromas `json:"data"`
}

/* Gets single weapon type from ID */
func (weaponChromaID WeaponChromaID) GetInformation() (*ValorantWeaponChromaSingle, error) {

	APICall := newAPICall()

	APICall.URL = "https://valorant-api.com/v1/weapons/skinchromas/" + string(weaponChromaID)
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

	Response := &ValorantWeaponChromaSingle{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	if Response.Status != 200 {
		return nil, errors.New(VALORANT_NO_VALID_RESPONSE)
	}

	return Response, nil

}

/* [ - Get Single WeaponSkin - ] */
/* [ URL: 'https://valorant-api.com/v1/weapons/skinlevels/{SkinLevelID}' ] */

type ValorantWeaponLevelSingle struct {
	Status int              `json:"status"`
	Data   WeaponSkinLevels `json:"data"`
}

/* Gets single weapon type from ID */
func (weaponSkinLevelID WeaponSkinLevelID) GetInformation() (*ValorantWeaponLevelSingle, error) {

	APICall := newAPICall()

	APICall.URL = "https://valorant-api.com/v1/weapons/skinlevels/" + string(weaponSkinLevelID)
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

	Response := &ValorantWeaponLevelSingle{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	if Response.Status != 200 {
		return nil, errors.New(VALORANT_NO_VALID_RESPONSE)
	}

	return Response, nil

}

type ValorantTitleInformation struct {
	UUID               string `json:"uuid"`
	DisplayName        string `json:"displayName"`
	TitleText          string `json:"titleText"`
	IsHiddenIfNotOwned bool   `json:"isHiddenIfNotOwned"`
	AssetPath          string `json:"assetPath"`
}

/* [ - Get Title Info from ID - ] */
/* [ URL: 'https://valorant-api.com/v1/playertitles/{TitleID}' ] */

type ValorantTitleIDSingle struct {
	Status int                      `json:"status"`
	Data   ValorantTitleInformation `json:"data"`
}

/* Gets title information from ID */
func (titleID TitleID) GetInformation() (*ValorantTitleIDSingle, error) {

	APICall := newAPICall()

	APICall.URL = "https://valorant-api.com/v1/playertitles/" + string(titleID)
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

	Response := &ValorantTitleIDSingle{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	if Response.Status != 200 {
		return nil, errors.New(VALORANT_NO_VALID_RESPONSE)
	}

	return Response, nil

}
