package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	valorantapi "val-controller/ValorantAPI"
)

type ClientUpdate struct {
	IsLoaded   bool   `json:"isLoaded"`
	IsInMatch  bool   `json:"isInMatch"`
	IsGameOpen bool   `json:"gameOpen"`
	TokenTest  string `json:"tokenTest"`
}

// App struct
type App struct {
	ctx                context.Context
	blockMainThread    chan struct{}
	clientUpdate       ClientUpdate
	valorantAPIContext *valorantapi.ValorantAPIContext
}

// NewApp creates a new App application struct
func NewApp() *App {
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

var loadoutStore *JSONStore
var OwnedAgentLookup = map[string]valorantapi.AgentID{}

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

	go func() {

		for {

			a.valorantAPIContext, err = valorantapi.GetLocalValorantAPIContext()
			if err != nil {
				fmt.Println(err)
				time.Sleep(1 * time.Second)
				continue
			}

			_, err := a.valorantAPIContext.GetEntitlementToken()
			if err != nil {
				fmt.Println(err)

				a.clientUpdate.IsGameOpen = false
				a.clientUpdate.IsLoaded = false

				a.updateClient(a.clientUpdate)
				time.Sleep(1 * time.Second)
				continue

			}

			a.clientUpdate.IsGameOpen = true
			a.clientUpdate.IsLoaded = true
			a.updateClient(a.clientUpdate)

			OwnedAgentLookup, err = a.valorantAPIContext.GetOwnedAgentData()

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

type SavedLoadout struct {
	LoadoutData valorantapi.ValorantLocalLoadout `json:"LoadoutData"`
	NameLookup  map[string]string                `json:"NameLookup"`
}

func (a *App) GetLoadouts() error {

	return a.sendLoadout()

}

func (a *App) sendLoadout() error {

	writeData := map[string]SavedLoadout{}
	err := loadoutStore.readJSON(&writeData)
	if err != nil {
		return err
	}

	runtime.EventsEmit(a.ctx, "on_loadout_update", writeData)

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
		NameLookup[string(v.SkinID)] = data.Data.DisplayName
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

	err = a.GetLoadouts()
	if err != nil {
		return err
	}

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

	err = a.GetLoadouts()
	if err != nil {
		return err
	}

	return nil

}

func (a *App) LoadSavedLoadout(name string) error {

	fmt.Println("Loading loadout.. '" + name + "'")

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

	data := writeData[name]

	data.LoadoutData.Identity.HideAccountLevel = currentLoadout.Identity.HideAccountLevel
	data.LoadoutData.Identity.Incognito = currentLoadout.Identity.Incognito
	data.LoadoutData.Identity.PreferredLevelBorderID = currentLoadout.Identity.PreferredLevelBorderID

	err = a.valorantAPIContext.SetAccountLoadout(data.LoadoutData, localPlayer)
	if err != nil {
		return err
	}

	return nil

}

type MatchData struct {
	IsLoaded   bool   `json:"isLoaded"`
	IsInMatch  bool   `json:"isInMatch"`
	IsGameOpen bool   `json:"gameOpen"`
	TokenTest  string `json:"tokenTest"`
}

func (a *App) GetMatch() *valorantapi.MatchData {

	data := a.updateMatchData()

	return data

}

func (a *App) updateMatchData() *valorantapi.MatchData {

	// Get Match Data

	if a.clientUpdate.IsGameOpen {

		data, err := a.valorantAPIContext.GetGameData(true)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		fmt.Println("Sent")

		runtime.EventsEmit(a.ctx, "update_match", data)

		return data

	}

	return nil

}

func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func (a *App) SelectRandomAgent() error {

	if a.clientUpdate.IsGameOpen {

		LocalPlayer, err := a.valorantAPIContext.GetLocalPlayerContext()
		if err != nil {
			fmt.Println(err)
			return err
		}

		CurrentMatch, err := a.valorantAPIContext.GetCurrentGamePlayer(LocalPlayer)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if CurrentMatch.PregameMatchID != "" && CurrentMatch.GameMatchID == "" {

			data, err := a.valorantAPIContext.GetPreGameData(CurrentMatch.PregameMatchID)
			if err != nil {
				fmt.Println(err)
				return err
			}

			if data.ID == CurrentMatch.PregameMatchID {

				AgentLookup := map[string]bool{}

				for _, v := range data.AllyTeam.Players {
					fmt.Println(v.CharacterID)
					AgentLookup[string(v.CharacterID)] = true
				}

				CopyTable := []valorantapi.AgentID{}
				for i, v := range OwnedAgentLookup {
					if !AgentLookup[i] {
						CopyTable = append(CopyTable, v)
						fmt.Println(i)
					}
				}

				randID := randRange(0, len(CopyTable))

				a.valorantAPIContext.SelectAgent(CurrentMatch.PregameMatchID, CopyTable[randID])

			}

		}

	}

	return nil

}
