package valorantapi

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

/* [ - Local Lockfile - ] */
/* [ Returns the data in the RiotClient lockfile for later ] */

type LocalLockfile struct {
	Name      string
	ProcessID string
	Port      string
	Password  string
	Protocol  string
}

func (context *ValorantAPIContext) setLocalLockfile() error {

	cacheDirectory, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	dir := cacheDirectory + "/Riot Games/Riot Client/Config/lockfile"

	_, err = os.Stat(dir)

	if errors.Is(err, fs.ErrNotExist) {
		return errors.New(VALORANT_NO_LOCAL_LOCKFILE_FOUND)
	}

	if err != nil {
		return err
	}

	LocalLockfileData, err := os.ReadFile(dir)
	if err != nil {
		return err
	}

	split := strings.Split(string(LocalLockfileData), ":")

	context.localLockfile = LocalLockfile{
		Name:      split[0],
		ProcessID: split[1],
		Port:      split[2],
		Password:  split[3],
		Protocol:  split[4],
	}

	return nil

}

/* [ - Local Region File - ] */
/* [ Returns Regional Data which is required for later API endpoints ] */

type LocalRegion struct {
	Region string
	Shard  string
}

func (context *ValorantAPIContext) setLocalRegionData() error {

	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	dir := userCacheDir + "/VALORANT/Saved/Logs/ShooterGame.log"

	file, err := os.ReadFile(dir)
	if err != nil {
		return err
	}

	shooterContents := (string(file))

	reg, err := regexp.Compile("https://glz-(.+?)-1.(.+?).a.pvp.net")
	if err != nil {
		return err
	}

	matches := reg.FindStringSubmatch(shooterContents)

	if len(matches) <= 2 {
		return fmt.Errorf(VALORANT_NO_REGION_SHARD)
	}

	context.localRegionData = LocalRegion{
		Region: matches[1],
		Shard:  matches[2],
	}

	return nil

}

/* [ - Return Context - ] */
/* [ Will be used for all later local API endpoint calls ] */

type ValorantAPIContext struct {
	localLockfile   LocalLockfile
	localRegionData LocalRegion
}

func GetLocalValorantAPIContext() (*ValorantAPIContext, error) {

	ReturnedContext := &ValorantAPIContext{}

	ReturnedContext.setLocalLockfile()
	ReturnedContext.setLocalRegionData()

	return ReturnedContext, nil

}

/* [ - Player Context - ] */
/* [ Will be required for player based requests ] */

type ValorantPlayerContext struct {
	UUID string
}

func GetValorantPlayer(uuid string) *ValorantPlayerContext {

	ReturnedContext := &ValorantPlayerContext{}

	ReturnedContext.UUID = uuid

	return ReturnedContext

}

/* [ Gets the local Player Context ] */

func (context ValorantAPIContext) GetLocalPlayerContext() (*ValorantPlayerContext, error) {

	LocalPlrInfo, err := context.GetLocalPlayerInfo()
	if err != nil {
		return nil, err
	}

	ReturnedContext := &ValorantPlayerContext{}

	ReturnedContext.UUID = LocalPlrInfo.LocalUUID

	return ReturnedContext, nil

}
