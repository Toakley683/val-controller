package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/go-stack/stack"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	valorantapi "val-controller/ValorantAPI"
)

type ClientUpdate struct {
	IsLoaded   bool   `json:"isLoaded"`
	IsInMatch  bool   `json:"isInMatch"`
	IsGameOpen bool   `json:"gameOpen"`
	TokenTest  string `json:"tokenTest"`
}

type SettingsStoreData struct {
	IsRandomized          bool            `json:"isRandomized"`
	RandomSelectedWeapons map[string]bool `json:"randomSelected"`
}

// App struct
type App struct {
	ctx                context.Context
	blockMainThread    chan struct{}
	clientUpdate       ClientUpdate
	valorantAPIContext *valorantapi.ValorantAPIContext
}

var (
	buildCommit = ""
)

type GithubUpdateStruct struct {
	URL         string    `json:"url"`
	ID          int       `json:"id"`
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Draft       bool      `json:"draft"`
	Immutable   bool      `json:"immutable"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		URL                string `json:"url"`
		ID                 int    `json:"id"`
		NodeID             string `json:"node_id"`
		Name               string `json:"name"`
		Label              string `json:"label"`
		ContentType        string `json:"content_type"`
		State              string `json:"state"`
		Size               int    `json:"size"`
		DownloadCount      int    `json:"download_count"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var (
	client http.Client
)

func init() {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client = http.Client{Transport: tr}

}

func checkForUpdates() error {

	fmt.Println("Current commit:", buildCommit)
	fmt.Println("Latest commit:", buildCommit)

	req, err := http.NewRequest("GET", "https://api.github.com/repos/Toakley683/val-controller/releases/latest", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	FullBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	Response := &GithubUpdateStruct{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return err
	}

	if strconv.Itoa(Response.ID) != buildCommit && len(Response.Assets) > 0 {

		fmt.Println("Update found")

		fmt.Println(Response.Assets[0].BrowserDownloadURL)

	} else {

		fmt.Println("No update found, continuing..")

	}

	return nil

}

// NewApp creates a new App application struct
func NewApp() *App {

	checkForUpdates()

	return &App{}
}

func (a *App) checkForGame() {

	var err error

	for {

		a.clientUpdate.IsLoaded = false

		a.valorantAPIContext, err = valorantapi.GetLocalValorantAPIContext()
		if err != nil {
			fmt.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}

		Token, err := a.valorantAPIContext.GetEntitlementToken()
		if err != nil {
			fmt.Println(err)

			a.clientUpdate.IsGameOpen = false
			a.updateClient(a.clientUpdate)
			time.Sleep(2 * time.Second)
			continue
		}

		a.clientUpdate.IsLoaded = true
		a.clientUpdate.TokenTest = Token.AccessToken
		break

	}

}

type JSONStore struct {
	Directory string
	Filepath  string
	mu        sync.RWMutex
}

func newDataStore(Directory string, fileName string) (*JSONStore, error) {
	if err := os.MkdirAll(Directory, os.ModePerm); err != nil {
		return nil, fmt.Errorf("Failed to get config dir: " + err.Error())
	}

	filePath := filepath.Join(Directory, fileName)

	dataStore := &JSONStore{
		Directory: Directory,
		Filepath:  filePath,
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create empty data structure
		initialData := map[string]string{}
		err = dataStore.writeJSON(initialData)
		if err != nil {
			return nil, err
		}
	}

	return dataStore, nil

}

func (store *JSONStore) readJSON(data interface{}) error {

	if store == nil {
		return fmt.Errorf("Store is invalid")
	}

	store.mu.RLock()
	defer store.mu.RUnlock()

	dataRead, err := os.ReadFile(store.Filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if len(dataRead) <= 0 {
		dataRead = []byte("{}")
	}

	if err := json.Unmarshal(dataRead, data); err != nil {

		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nil

}

func (store *JSONStore) writeJSON(data interface{}) error {

	if store == nil {
		return fmt.Errorf("Store is invalid")
	}

	store.mu.RLock()
	defer store.mu.RUnlock()

	dataMarshal, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(store.Filepath, dataMarshal, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil

}

var settingsStore *JSONStore
var loadoutStore *JSONStore
var lastSeenStore *JSONStore
var OwnedAgentLookup = map[string]valorantapi.AgentID{}
var OwnedWeaponLookup map[string][]valorantapi.ValorantLocalLoadoutGuns
var AllWeapons = &valorantapi.ValorantAllWeapons{}

type LastSeenObject struct {
	Subject    string `json:"Subject"`
	LastSeen   int64  `json:"LastSeen"`
	MatchesAgo int64  `json:"MatchesAgo"`
}

var lastSubject string = ""

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.clientUpdate = ClientUpdate{}

	// Setup jsonFile for saving loadouts

	configDir, err := os.UserConfigDir()
	if err != nil {
		runtime.LogError(ctx, "Failed to get config dir: "+err.Error())
		return
	}

	appDataPath := filepath.Join(configDir, "ValorantController")

	loadoutStore, err = newDataStore(appDataPath, "loadoutStore.json")
	if err != nil {
		runtime.LogError(ctx, "Failed to make data store: "+err.Error())
		return
	}

	lastSeenStore, err = newDataStore(appDataPath, "lastSeenStore.json")
	if err != nil {
		runtime.LogError(ctx, "Failed to make lastSeen store: "+err.Error())
		return
	}

	settingsStore, err = newDataStore(appDataPath, "settings.json")
	if err != nil {
		runtime.LogError(ctx, "Failed to make data store: "+err.Error())
		return
	}

	go func() {

		for {

			a.valorantAPIContext, err = valorantapi.GetLocalValorantAPIContext()
			if err != nil {
				fmt.Println("Get Context Error:", err)
				time.Sleep(1 * time.Second)
				continue
			}

			_, err := a.valorantAPIContext.GetEntitlementToken()
			if err != nil {
				fmt.Println("Get Entitlement Error:", err)

				a.clientUpdate.IsGameOpen = false
				a.clientUpdate.IsLoaded = false

				a.updateClient(a.clientUpdate)
				time.Sleep(1 * time.Second)
				continue

			}

			a.clientUpdate.IsGameOpen = true
			a.clientUpdate.IsLoaded = true
			a.updateClient(a.clientUpdate)

			//doUseCache := true

			Subject, err := a.valorantAPIContext.GetLocalPlayerContext()
			if err != nil {
				fmt.Println("Subject Find Error:", err)
				continue
			}

			if lastSubject != Subject.UUID {
				lastSubject = Subject.UUID

				a.onLoad()
			}

			time.Sleep(1 * time.Second)

		}

	}()

	go func() {

		for {

			a.updateMatchData()

			time.Sleep(1500 * time.Millisecond)

		}

	}()
}

func (a *App) onLoad() {

	fmt.Println("On Load")

	var err error

	OwnedAgentLookup, err = a.valorantAPIContext.GetOwnedAgentData(false)
	if err != nil {
		fmt.Println("Owned Agent Data Error:", err)
		return
	}

	OwnedWeaponChromas, err := a.valorantAPIContext.GetOwnedWeaponChromas(false)
	if err != nil {
		fmt.Println("Get Owned Chroma Error:", err)
		return
	}

	OwnedWeaponSkins, err := a.valorantAPIContext.GetOwnedWeaponSkins(false)
	if err != nil {
		fmt.Println("Get Owned Weapon Error:", err)
		return
	}

	AllWeapons, err = a.valorantAPIContext.GetAllWeapons()

	OwnedWeaponLookup = map[string][]valorantapi.ValorantLocalLoadoutGuns{}

	for _, weapon := range AllWeapons.Data {
		if OwnedWeaponLookup[weapon.UUID] == nil {
			OwnedWeaponLookup[weapon.UUID] = []valorantapi.ValorantLocalLoadoutGuns{}
		}

		for _, skin := range weapon.Skins {

			isDefault := false

			if skin.ThemeUUID == "5a629df4-4765-0214-bd40-fbb96542941f" {

				// Get all standard skins

				isDefault = true

			}

			allChroma := []valorantapi.ValorantLocalLoadoutGuns{}

			highestLevel := 0

			for i, skinLevels := range skin.Levels {

				if OwnedWeaponSkins[skinLevels.UUID] != nil || isDefault {

					if highestLevel < i+1 {
						highestLevel = i + 1
					}

				}

			}

			if highestLevel > 0 {

				ownedChromas := 0

				for _, chroma := range skin.Chromas {

					if OwnedWeaponChromas[chroma.UUID] != nil {

						// Need to add max skin VARIANTS to list of available to use

						allChroma = append(allChroma, valorantapi.ValorantLocalLoadoutGuns{
							ID:          valorantapi.WeaponID(weapon.UUID),
							SkinID:      valorantapi.WeaponSkinID(skin.UUID),
							SkinLevelID: valorantapi.WeaponSkinLevelID(skin.Levels[highestLevel-1].UUID),
							ChromaID:    valorantapi.WeaponChromaID(chroma.UUID),
						})

						ownedChromas++

					}

				}

				// Need to add max skin level to list of available to use

				allChroma = append(allChroma, valorantapi.ValorantLocalLoadoutGuns{
					ID:          valorantapi.WeaponID(weapon.UUID),
					SkinID:      valorantapi.WeaponSkinID(skin.UUID),
					SkinLevelID: valorantapi.WeaponSkinLevelID(skin.Levels[highestLevel-1].UUID),
					ChromaID:    valorantapi.WeaponChromaID(skin.Chromas[0].UUID),
				})

			}

			OwnedWeaponLookup[weapon.UUID] = append(OwnedWeaponLookup[weapon.UUID], allChroma...)
		}
	}

	a.GetLoadouts()

}

func (a *App) CloseWindow() {
	runtime.Quit(a.ctx)
}

func (a *App) MinimizeWindow() {
	runtime.WindowMinimise(a.ctx)
}

func (a *App) SetWindowSize(width, height int) {
	runtime.WindowSetSize(a.ctx, width, height)
}

func (a *App) UpdateCurrentClient() {
	a.updateClient(a.clientUpdate)
}

func (a *App) updateClient(clientUpdate ClientUpdate) {

	runtime.EventsEmit(a.ctx, "updateClient", clientUpdate)

}

func (a *App) GetLoadout() *SavedLoadout {

	writeData := map[string]SavedLoadout{}
	err := loadoutStore.readJSON(&writeData)
	if err != nil {
		return nil
	}

	keys := reflect.ValueOf(writeData).MapKeys()

	return keys[0].Interface().(*SavedLoadout)

}

var RandomLoadout SavedLoadout = SavedLoadout{
	LoadoutData: valorantapi.ValorantLocalLoadout{
		Guns: []valorantapi.ValorantLocalLoadoutGuns{
			{ID: "63e6c2b6-4a8e-869c-3d4c-e38355226584", SkinID: "f454efd1-49cb-372f-7096-d394df615308", ChromaID: "2f93861d-4b2f-2175-af0c-3ba0c736e257"},
			{ID: "55d8a0f4-4274-ca67-fe2c-06ab45efdf58", SkinID: "5305d9c4-4f46-fbf4-9e9a-dea772c263b5", ChromaID: "b33de820-4061-8b85-31ce-808f1a2c58f5"},
			{ID: "9c82e19d-4575-0200-1a81-3eacf00cf872", SkinID: "27f21d97-4c4b-bd1c-1f08-31830ab0be84", ChromaID: "19629ae1-4996-ae98-7742-24a240d41f99"},
			{ID: "ae3de142-4d85-2547-dd26-4e90bed35cf7", SkinID: "724a7f42-4315-eccf-0e76-77bdd3ec2e09", ChromaID: "bf35f404-4a14-6953-ced2-5bafd21639a0"},
			{ID: "ee8e8d15-496b-07ac-e5f6-8fae5d4c7b1a", SkinID: "337cb216-4a6e-d85d-88c2-f29ab317784c", ChromaID: "52221ba2-4e4c-ec76-8c81-3483506d5242"},
			{ID: "ec845bf4-4f79-ddda-a3da-0db3774b2794", SkinID: "acd26127-48ff-8b9e-7ba6-b989af8a4b24", ChromaID: "b71ae8d6-44bb-aa4c-0d2a-dc9ed9e66410"},
			{ID: "910be174-449b-c412-ab22-d0873436b21b", SkinID: "70c97fb2-4d79-d4bb-5173-a1888cd4bfd9", ChromaID: "3d8ffcfe-4786-0180-42d7-e1be18dd1cab"},
			{ID: "44d4e95c-4157-0037-81b2-17841bf2e8e3", SkinID: "f06657f3-48b6-6314-7235-a9a2749df5b9", ChromaID: "dc99ed5a-4d75-87a0-c921-75963ea3c1e1"},
			{ID: "29a0cfab-485b-f5d5-779a-b59f85e204a8", SkinID: "24aee897-4cdc-b0fd-e596-1ba90fa6d1b2", ChromaID: "4b2d5b4f-4955-4208-286c-abadec250cdd"},
			{ID: "1baa85b4-4c70-1284-64bb-6481dfc3bb4e", SkinID: "1c63b43b-43c4-04e4-01c9-7aa1bffa5ac1", ChromaID: "947a28b6-4e0f-61fb-e795-bc9a5e7b7129"},
			{ID: "e336c6b8-418d-9340-d77f-7a9e4cfe0702", SkinID: "1ef6ba68-4dbe-30c7-6bc8-93a6c6f13f04", ChromaID: "5a59bd61-48a9-af61-c00f-4aa21deca9a8"},
			{ID: "42da8ccc-40d5-affc-beec-15aa47b42eda", SkinID: "48ad078a-4dae-2b85-a945-f4b6d1efecbb", ChromaID: "95608504-4c8b-1408-1612-0f8200421c49"},
			{ID: "a03b24d3-4319-996d-0f8c-94bbfba1dfc7", SkinID: "d1f2920f-469a-3431-ad96-96afbd0017f2", ChromaID: "4914f50d-49f9-6424-ca80-9486c45a138d"},
			{ID: "4ade7faa-4cf1-8376-95ef-39884480959b", SkinID: "3bf1e8e0-47e8-f27a-6054-929575f41a54", ChromaID: "0f934388-418a-a9e7-42a7-21b27402e46c"},
			{ID: "c4883e50-4494-202c-3ec3-6b8a9284f00b", SkinID: "fd44b2d5-49ee-77ab-fa56-588f3ac0c268", ChromaID: "1afec971-4170-f29b-1c94-07a0eff270ab"},
			{ID: "462080d1-4035-2937-7c09-27aa2a5c27a7", SkinID: "f01d1307-4299-42f5-2c5e-7dab7e69ab19", ChromaID: "a9aaccca-4cdc-02ea-1d7e-89bbacecc0e2"},
			{ID: "f7e1b454-4ad4-1063-ec0a-159e56b58941", SkinID: "940fb417-4a9c-3004-41f5-3e8f1f4178b2", ChromaID: "31bb2115-4c62-d37c-43c4-11b8fee7f212"},
			{ID: "2f59173c-4bed-b6c3-2191-dea9b58be9c7", SkinID: "12cc9ed2-4430-d2fe-3064-f7a19b1ba7c7", ChromaID: "cac83e5c-47a1-3519-5420-1db1fdbc4892"},
			{ID: "5f0aaf7a-4289-3998-d5ff-eb9a5cf7ef5c", SkinID: "740b9572-44b1-57bf-767e-6aa01811f94d", ChromaID: "66c8d241-4f7c-6652-3aaa-51bafffbd493"},
			{ID: "410b2e0b-4ceb-1321-1727-20858f7f3477", SkinID: "b576b9f1-407d-310a-6009-6287fb6829bc", ChromaID: "ff308f94-427d-e6d3-8fcb-258f5d4c2c29"},
		},
		Sprays: []valorantapi.ValorantLocalExpression{
			{TypeID: "03a572de-4234-31ed-d344-ababa488f981"},
			{TypeID: "d5f120f8-ff8c-4aac-92ea-f2b5acbe9475"},
			{TypeID: "d5f120f8-ff8c-4aac-92ea-f2b5acbe9475"},
			{TypeID: "d5f120f8-ff8c-4aac-92ea-f2b5acbe9475"},
		},
		Identity: valorantapi.ValorantLocalLoadoutIdentity{
			PlayerCardID: "9fb348bc-41a0-91ad-8a3e-818035c4e561",
		},
	},
	NameLookup: map[string]string{
		"1baa85b4-4c70-1284-64bb-6481dfc3bb4e": "Random Ghost",
		"29a0cfab-485b-f5d5-779a-b59f85e204a8": "Random Classic",
		"2f59173c-4bed-b6c3-2191-dea9b58be9c7": "Random Melee",
		"410b2e0b-4ceb-1321-1727-20858f7f3477": "Random Bandit",
		"42da8ccc-40d5-affc-beec-15aa47b42eda": "Random Shorty",
		"44d4e95c-4157-0037-81b2-17841bf2e8e3": "Random Frenzy",
		"462080d1-4035-2937-7c09-27aa2a5c27a7": "Random Spectre",
		"4ade7faa-4cf1-8376-95ef-39884480959b": "Random Guardian",
		"55d8a0f4-4274-ca67-fe2c-06ab45efdf58": "Random Ares",
		"5f0aaf7a-4289-3998-d5ff-eb9a5cf7ef5c": "Random Outlaw",
		"63e6c2b6-4a8e-869c-3d4c-e38355226584": "Random Odin",
		"6b4e1d0c-410e-878b-f151-9b8a8abc83a3": "Random Title",
		"910be174-449b-c412-ab22-d0873436b21b": "Random Bucky",
		"9c82e19d-4575-0200-1a81-3eacf00cf872": "Random Vandal",
		"a03b24d3-4319-996d-0f8c-94bbfba1dfc7": "Random Operator",
		"ae3de142-4d85-2547-dd26-4e90bed35cf7": "Random Bulldog",
		"c4883e50-4494-202c-3ec3-6b8a9284f00b": "Random Marshal",
		"e336c6b8-418d-9340-d77f-7a9e4cfe0702": "Random Sheriff",
		"ec845bf4-4f79-ddda-a3da-0db3774b2794": "Random Judge",
		"ee8e8d15-496b-07ac-e5f6-8fae5d4c7b1a": "Random Phantom",
		"f7e1b454-4ad4-1063-ec0a-159e56b58941": "Random Stinger",
	},
}

type SavedLoadout struct {
	LoadoutData valorantapi.ValorantLocalLoadout `json:"LoadoutData"`
	NameLookup  map[string]string                `json:"NameLookup"`
}
type UpdateLoadoutObj struct {
	Loadouts              map[string]SavedLoadout
	CurrentLoadout        valorantapi.ValorantLocalLoadout
	RandomWeaponsSelected map[string]bool
	IsRandomSelected      bool
}

func (a *App) GetLoadouts() UpdateLoadoutObj {

	err := a.sendLoadout()
	if err != nil {
		fmt.Println("Error getting loadout,", err)
	}

	return UpdateLoadoutObj{}

}

func (a *App) sendLoadout() error {

	writeData := map[string]SavedLoadout{}
	err := loadoutStore.readJSON(&writeData)
	if err != nil {
		return err
	}

	localPlayer, err := a.valorantAPIContext.GetLocalPlayerContext()
	if err != nil {
		return err
	}

	currentLoadout, err := a.valorantAPIContext.GetAccountLoadout(localPlayer)
	if err != nil {
		return err
	}

	writeData["Random"] = RandomLoadout

	storeData := SettingsStoreData{}
	err = settingsStore.readJSON(&storeData)
	if err != nil {
		fmt.Println("Settings Store Error:", err)
		return err
	}

	if storeData.RandomSelectedWeapons == nil {
		storeData.RandomSelectedWeapons = map[string]bool{}
	}

	UploadData := UpdateLoadoutObj{
		Loadouts:              writeData,
		CurrentLoadout:        *currentLoadout,
		RandomWeaponsSelected: storeData.RandomSelectedWeapons,
		IsRandomSelected:      storeData.IsRandomized,
	}

	runtime.EventsEmit(a.ctx, "on_loadout_update", UploadData)

	return nil

}

func (a *App) AddWeaponsToBeRandomized(val map[string]bool) error {

	fmt.Println(val)

	storeData := SettingsStoreData{}
	err := settingsStore.readJSON(&storeData)
	if err != nil {
		fmt.Println("Read Settings Error", err, stack.Trace())
		return err
	}

	storeData.RandomSelectedWeapons = val

	err = settingsStore.writeJSON(storeData)
	if err != nil {
		fmt.Println("Write Settings Error", err, stack.Trace())
		return err
	}

	// Change to courotine that does it in a few seconds to not spam the API

	a.randomizeLoadout()

	return nil

}

func (a *App) SaveCurrentLoadout(name string) error {

	fmt.Println("Saving loadout.. '" + name + "'")

	localPlayer, err := a.valorantAPIContext.GetLocalPlayerContext()
	if err != nil {
		return err
	}

	currentLoadout, err := a.valorantAPIContext.GetAccountLoadout(localPlayer)
	if err != nil {
		return err
	}

	NameLookup := map[string]string{}

	for _, v := range currentLoadout.Guns {
		data, err := v.SkinID.GetInformation()
		if err != nil {
			continue
		}
		NameLookup[string(v.ID)] = data.Data.DisplayName
	}

	data, err := currentLoadout.Identity.PlayerTitleID.GetInformation()
	if err != nil {
		return err
	}

	NameLookup[string(data.Data.UUID)] = data.Data.TitleText

	ToBeSaved := SavedLoadout{
		LoadoutData: *currentLoadout,
		NameLookup:  NameLookup,
	}

	writeData := map[string]SavedLoadout{}
	err = loadoutStore.readJSON(&writeData)
	if err != nil {
		return err
	}

	writeData[name] = ToBeSaved

	loadoutStore.writeJSON(writeData)

	a.GetLoadouts()

	return nil

}

func (a *App) DeleteSavedLoadout(name string) error {

	fmt.Println("Deleting loadout.. '" + name + "'")

	writeData := map[string]SavedLoadout{}
	err := loadoutStore.readJSON(&writeData)
	if err != nil {
		return err
	}

	delete(writeData, name)

	loadoutStore.writeJSON(writeData)

	a.GetLoadouts()

	return nil

}

func (a *App) randomizeLoadout() error {

	storeData := SettingsStoreData{}
	err := settingsStore.readJSON(&storeData)
	if err != nil {
		fmt.Println("Read Settings Error", err, stack.Trace())
		return err
	}

	if storeData.IsRandomized == false {
		return nil
	}

	localPlayer, err := a.valorantAPIContext.GetLocalPlayerContext()
	if err != nil {
		fmt.Println("Get Local Player Context Error:", err, stack.Trace())
		return err
	}

	currentLoadout, err := a.valorantAPIContext.GetAccountLoadout(localPlayer)
	if err != nil {
		fmt.Println("Get Account Loadout Error:", err, stack.Trace())
		return err
	}

	currentLoadoutLookup := map[valorantapi.WeaponID]int{}

	for index, v := range currentLoadout.Guns {

		currentLoadoutLookup[v.ID] = index

	}

	for uuid, isRandomized := range storeData.RandomSelectedWeapons {

		if isRandomized {

			index := currentLoadoutLookup[valorantapi.WeaponID(uuid)]

			randomWeapon := OwnedWeaponLookup[uuid]

			if len(randomWeapon) <= 0 {
				continue
			}

			// Shuffle the array and get the first element == Random

			rand.Shuffle(len(randomWeapon), func(i, j int) {
				randomWeapon[i], randomWeapon[j] = randomWeapon[j], randomWeapon[i]
			})

			randomWeapon[0].CharmInstanceID = currentLoadout.Guns[index].CharmInstanceID
			randomWeapon[0].CharmID = currentLoadout.Guns[index].CharmID
			randomWeapon[0].CharmLevelID = currentLoadout.Guns[index].CharmLevelID
			randomWeapon[0].Attachments = currentLoadout.Guns[index].Attachments

			currentLoadout.Guns[index] = randomWeapon[0]

		}

	}

	err = a.valorantAPIContext.SetAccountLoadout(*currentLoadout, localPlayer)
	if err != nil {
		fmt.Println("Set Account Loadout Error", err, stack.Trace())
		return err
	}
	a.GetLoadouts()

	return nil

}

func (a *App) LoadSavedLoadout(name string, isRandom bool) error {

	fmt.Println("Loading loadout.. '" + name + "'")

	writeData := map[string]SavedLoadout{}
	err := loadoutStore.readJSON(&writeData)
	if err != nil {
		fmt.Println("Read Settings Error", err, stack.Trace())
		return err
	}

	localPlayer, err := a.valorantAPIContext.GetLocalPlayerContext()
	if err != nil {
		fmt.Println("Get Local Player Context Error:", err, stack.Trace())
		return err
	}

	currentLoadout, err := a.valorantAPIContext.GetAccountLoadout(localPlayer)
	if err != nil {
		fmt.Println("Get Account Loadout Error:", err, stack.Trace())
		return err
	}

	data := writeData[name]

	data.LoadoutData.Identity.HideAccountLevel = currentLoadout.Identity.HideAccountLevel
	data.LoadoutData.Identity.Incognito = currentLoadout.Identity.Incognito
	data.LoadoutData.Identity.PreferredLevelBorderID = currentLoadout.Identity.PreferredLevelBorderID

	storeData := SettingsStoreData{}
	err = settingsStore.readJSON(&storeData)
	if err != nil {
		fmt.Println("Read Settings Error", err, stack.Trace())
		return err
	}

	storeData.IsRandomized = isRandom

	err = settingsStore.writeJSON(storeData)
	if err != nil {
		fmt.Println("Write Settings Error", err, stack.Trace())
		return err
	}

	if name != "Random" {

		err = a.valorantAPIContext.SetAccountLoadout(data.LoadoutData, localPlayer)
		if err != nil {
			fmt.Println("Set Account Loadout Error", err, stack.Trace())
			return err
		}

	} else {

		// Set loadout to random weapons based on schema

		if isRandom {

			// Enable random

			a.randomizeLoadout()

		} else {

			// Disable random

		}

	}

	a.GetLoadouts()

	return nil

}

var lastMatchID = ""
var lastMatchData *valorantapi.MatchData

func (a *App) GetMatch() *valorantapi.MatchData {

	data := a.updateMatchData()

	return data

}

var RemoveOldEntryLimit = 204800

// Only keep up to X amount of entries in the lastSeenStore

func (a *App) onMatchEnd(lastMatchData *valorantapi.MatchData) {

	writeData := map[string]LastSeenObject{}
	err := lastSeenStore.readJSON(&writeData)
	if err == nil {

		for _, v := range append(lastMatchData.AllyTeam.Players, lastMatchData.EnemyTeam.Players...) {

			writeData[v.Subject] = LastSeenObject{
				Subject:    v.Subject,
				LastSeen:   time.Now().Unix(),
				MatchesAgo: 0,
			}

		}

		for _, v := range writeData {

			writeData[v.Subject] = LastSeenObject{
				Subject:    v.Subject,
				LastSeen:   v.LastSeen,
				MatchesAgo: v.MatchesAgo + 1,
			}

		}

		if len(writeData) > RemoveOldEntryLimit {

			objects := make([]LastSeenObject, 0, len(writeData))
			for _, obj := range writeData {
				objects = append(objects, obj)
			}

			sort.Slice(objects, func(i, j int) bool {
				return objects[i].LastSeen > objects[j].LastSeen
			})

			for index, obj := range objects {

				if index >= RemoveOldEntryLimit {

					delete(writeData, obj.Subject)

				}
			}

		}

		err = lastSeenStore.writeJSON(writeData)

		if err != nil {

			fmt.Println("Errors with saving lastSeenStore:", err)

		}

		a.randomizeLoadout()

	} else {
		fmt.Println("Read Settings Error", err, stack.Trace())
	}

}

func (a *App) onMatchStart(matchData *valorantapi.MatchData) {

	fmt.Println("Match has started")

}

func (a *App) updateMatchData() *valorantapi.MatchData {

	// Get Match Data

	if a.clientUpdate.IsGameOpen {

		data, err := a.valorantAPIContext.GetGameData(true)
		if err != nil {
			fmt.Println("Get Game Data Error", err, stack.Trace())
			return nil
		}

		if data.MatchID != lastMatchID {

			if lastMatchID == "" {
				lastMatchData = data
				a.onMatchStart(data)
			} else {
				if data.MatchID == "" {
					a.onMatchEnd(lastMatchData)
				}
			}

			lastMatchID = data.MatchID

		}

		readData := map[string]LastSeenObject{}
		err = lastSeenStore.readJSON(&readData)
		if err != nil {
			fmt.Println("Read Last Seen Error", err, stack.Trace())
		}

		for i, v := range data.AllyTeam.Players {
			data.AllyTeam.Players[i].MatchesAgo = readData[v.Subject].MatchesAgo
		}
		for i, v := range data.EnemyTeam.Players {
			data.EnemyTeam.Players[i].MatchesAgo = readData[v.Subject].MatchesAgo
		}

		fmt.Println("Sent")

		runtime.EventsEmit(a.ctx, "update_match", data)

		return data

	}

	return nil

}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func (a *App) ExitCoreGame() error {

	LocalPlayer, err := a.valorantAPIContext.GetLocalPlayerContext()
	if err != nil {
		fmt.Println("Get Local Player Context Error:", err, stack.Trace())
		return err
	}

	CurrentMatch, err := a.valorantAPIContext.GetCurrentGamePlayer(LocalPlayer)
	if err != nil {
		fmt.Println("Get Current Game Player Error:", err, stack.Trace())
		return err
	}

	if CurrentMatch.PregameMatchID == "" && CurrentMatch.GameMatchID != "" {

		a.valorantAPIContext.CoregameQuit(CurrentMatch.GameMatchID, LocalPlayer)

	}

	return nil

}

func (a *App) ExitPregame() error {

	LocalPlayer, err := a.valorantAPIContext.GetLocalPlayerContext()
	if err != nil {
		fmt.Println("Get Local Player Context Error:", err, stack.Trace())
		return err
	}

	CurrentMatch, err := a.valorantAPIContext.GetCurrentGamePlayer(LocalPlayer)
	if err != nil {
		fmt.Println("Get Current Game Player Error:", err, stack.Trace())
		return err
	}

	if CurrentMatch.PregameMatchID != "" && CurrentMatch.GameMatchID == "" {

		a.valorantAPIContext.PregameQuit(CurrentMatch.PregameMatchID)

	}

	return nil

}

func (a *App) SelectRandomAgent() error {

	if a.clientUpdate.IsGameOpen {

		LocalPlayer, err := a.valorantAPIContext.GetLocalPlayerContext()
		if err != nil {
			fmt.Println("Get Local Player Context Error:", err, stack.Trace())
			return err
		}

		CurrentMatch, err := a.valorantAPIContext.GetCurrentGamePlayer(LocalPlayer)
		if err != nil {
			fmt.Println("Get Current Game Player Error:", err, stack.Trace())
			return err
		}

		if CurrentMatch.PregameMatchID != "" && CurrentMatch.GameMatchID == "" {

			data, err := a.valorantAPIContext.GetPreGameData(CurrentMatch.PregameMatchID)
			if err != nil {
				fmt.Println("Get Pregame Data Error:", err, stack.Trace())
				return err
			}

			if data.ID == CurrentMatch.PregameMatchID {

				AgentLookup := map[string]bool{}

				for _, v := range data.AllyTeam.Players {
					AgentLookup[string(v.CharacterID)] = true
				}

				CopyTable := []valorantapi.AgentID{}
				for i, v := range OwnedAgentLookup {
					if !AgentLookup[i] {
						CopyTable = append(CopyTable, v)
					}
				}

				randID := randRange(0, len(CopyTable))

				a.valorantAPIContext.SelectAgent(CurrentMatch.PregameMatchID, CopyTable[randID])

			}

		}

	}

	return nil

}
