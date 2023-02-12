package input

import (
	"errors"
	"os"
	"path"

	"git.fiblab.net/sim/simulet-go/utils/config"
	"git.fiblab.net/sim/simulet-go/utils/protoutil"
	"google.golang.org/protobuf/proto"
)

func preCheckCache(cacheDir string) bool {
	if cacheDir == "" {
		log.Info("disable input cache")
		return false
	} else {
		if stat, err := os.Stat(cacheDir); err == nil && stat.IsDir() {
			// 文件夹存在
			log.Infof("enable input cache at %s", cacheDir)
			return true
		} else {
			log.Errorf("disable input cache because invalid dir %s (not exist or file)", cacheDir)
			return false
		}
	}
}

func mustReadCacheOrDownloadAndSave[T any, PT interface {
	proto.Message
	*T
}](
	cacheDir string, mp config.MongoPath,
	download func() PT,
) PT {
	var pt PT = new(T)
	filePath := path.Join(cacheDir, mp.GetDb()+"."+mp.GetCol()+".pb")
	err := protoutil.UnmarshalFromFile(pt, filePath)
	if err == nil {
		log.Infof("cache: read %+v from %s", mp, filePath)
		return pt
	} else if errors.Is(err, os.ErrNotExist) {
		// download and save
		pt = download()
		return pt
	} else {
		log.Fatal(err)
		return nil
	}
}
