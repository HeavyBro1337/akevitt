package main

import (
	"akevitt/core/database"
	"akevitt/core/database/utils"
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

const testDBName string = "test_akevitt.db"

func initDB(t *testing.T) *bolt.DB {
	db, err := bolt.Open(testDBName, 0600, nil)
	if err != nil {
		t.FailNow()
	}
	return db
}
func destroyDB(db *bolt.DB, t *testing.T) {
	db.Close()
	if err := os.Remove(testDBName); err != nil {
		t.FailNow()
	}
}
func Test_CreateWithEmptyCredentials(t *testing.T) {
	db := initDB(t)
	defer destroyDB(db, t)
	// Register
	err := database.CreateAccount(db, "", "")
	if err == nil { // No errors were made. But it should
		t.Fail()
	}

}
func Test_CreateDuplicateAccounts(t *testing.T) {
	db := initDB(t)
	defer destroyDB(db, t)
	// Register
	err := database.CreateAccount(db, "IamUserIwillDuplicateMyself", "000000")
	if err != nil {
		t.FailNow()
	}
	err = database.CreateAccount(db, "IamUserIwillDuplicateMyself", "000000")
	if err == nil {
		t.Fail()
	}
}

func Test_CreateAccountsWithEmptyPassword(t *testing.T) {
	db := initDB(t)
	defer destroyDB(db, t)
	err := database.CreateAccount(db, "Passwordless27", "")
	if err == nil {
		t.FailNow()
	}
}

func Test_HashString(t *testing.T) {
	originalText := "qwerty1234"
	hashedText := utils.HashString(originalText)
	t.Logf("Original text: %s\nHashed text: %s", originalText, hashedText)
	if originalText == hashedText {
		t.Fail()
	}
}
