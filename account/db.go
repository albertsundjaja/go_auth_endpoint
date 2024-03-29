package account

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
)

// Datasource abstraction for testable DB
type Datasource interface {
	QueryUserPassword(username string) (password string, err error)
	InputNewUser(registerRequest RegisterRequest) (err error)
}

// implementation of the Datasource using levelDB
var levelDBPath = "levelDB"
type LevelDB struct {
	DB *leveldb.DB
}

// the accessible variable for our DB
var ProdDB LevelDB

func init() {
	// create new folder to store our leveldb if not exist
	_, err := os.Stat(levelDBPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(levelDBPath, os.ModePerm)
		if err != nil {
			log.Fatal("Error creating new folder for levelDB")
		}
	}

	ProdDB.DB, err = leveldb.OpenFile(levelDBPath, nil)
	if err != nil {
		log.Fatal("Error opening levelDB folder")
	}
}

// given the username will return the password
func (levelDB LevelDB) QueryUserPassword(username string) (password string, err error) {
	var func_name = "QueryUserPassword"

	pass, err := levelDB.DB.Get([]byte(username), nil)
	if err != nil {
		fmt.Println(func_name, err)
		return "", err
	}
	return string(pass), nil
}

// register new user with username as key and password as value
func (levelDB LevelDB) InputNewUser(registerReq RegisterRequest) (err error) {
	var func_name = "InputNewUser"

	err = levelDB.DB.Put([]byte(registerReq.Username), []byte(registerReq.Password), nil)
	if err != nil {
		fmt.Println(func_name, err)
		return err
	}

	return nil
}