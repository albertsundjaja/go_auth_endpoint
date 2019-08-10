package account

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"testing"
)

var TestDB LevelDB

func TestMain(m *testing.M) {
	var levelDBPath = "TestLevelDB"
	err := os.Mkdir(levelDBPath, os.ModePerm)

	TestDB.DB, err = leveldb.OpenFile(levelDBPath, nil)
	if err != nil {
		log.Fatal("Error opening TestLevelDB folder")
	}

	code := m.Run()
	os.Exit(code)
}

func TestLevelDB_InputNewUser(t *testing.T) {
	var registerReg RegisterRequest
	registerReg.Username = "username_must_be_unique"
	registerReg.Password = "password"
	err := TestDB.InputNewUser(registerReg)
	if err != nil {
		t.Errorf("Expected error not to have occurred, got %s", err)
	}

	pass, err := TestDB.DB.Get([]byte(registerReg.Username), nil)
	if err != nil {
		t.Errorf("Expected error not to have occurred, got %s", err)
	}

	if string(pass) != registerReg.Password {
		t.Errorf("Registered user has differen password, expected %s, got %s", registerReg.Password, pass)
	}

	_= TestDB.DB.Delete([]byte(registerReg.Username), nil)
}

func TestLevelDB_QueryUserPassword(t *testing.T) {
	var username = "username_must_be_unique"
	var password = "password"

	_ = TestDB.DB.Put([]byte(username), []byte(password), nil)

	pass, err := TestDB.QueryUserPassword(username)
	if err != nil {
		t.Errorf("Expected error not to have occurred, got %s", err)
	}

	if pass != password {
		t.Errorf("Password returned does not match, expected %s, got %s", password, pass)
	}

	_= TestDB.DB.Delete([]byte(username), nil)
}

