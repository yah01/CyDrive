package store

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yah01/CyDrive/config"
	. "github.com/yah01/CyDrive/consts"
	"github.com/yah01/CyDrive/model"
	"io/ioutil"
	"os"
	"path/filepath"
)

type UserStore interface {
	GetUserByName(name string) *model.User
}

// Store users in memory
// load from a json file
type MemStore struct {
	userNameMap map[string]*model.User
}

func NewMemStore(userJson string) *MemStore {
	store := MemStore{userNameMap: make(map[string]*model.User)}

	data, _ := ioutil.ReadFile(userJson)

	userArray := make([]*model.User, 1)
	json.Unmarshal(data, &userArray)
	for _, user := range userArray {
		// Get the storage usage
		rootDirInfo,_ := os.Stat(user.RootDir)
		user.Usage = rootDirInfo.Size()

		store.userNameMap[user.Username] = user
	}

	return &store
}

func (store MemStore) GetUserByName(name string) *model.User {
	return store.userNameMap[name]
}

// Store users in a relational db
type RdbStore struct {
	db *gorm.DB
}

func NewRdbStore(config config.Config) *RdbStore {
	store := RdbStore{}
	store.db, _ = gorm.Open("mysql", config.PackDSN())
	return &store
}

func (store RdbStore) GetUserByName(name string) *model.User {
	var user model.User

	if store.db.First(user, "username = ?", name).RecordNotFound() {
		return nil
	}

	user.RootDir = filepath.Join(UserDataDir, fmt.Sprint(user.Id))
	user.WorkDir = user.RootDir
	return &user
}
