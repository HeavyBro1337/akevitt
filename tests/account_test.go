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
	_, err := database.CreateAccount(db, "", "")
	if err == nil { // No errors were made. But it should
		t.Fail()
	}

}
func Test_CreateDuplicateAccounts(t *testing.T) {
	db := initDB(t)
	defer destroyDB(db, t)
	// Register
	_, err := database.CreateAccount(db, "IamUserIwillDuplicateMyself", "000000")
	if err != nil {
		t.FailNow()
	}
	_, err = database.CreateAccount(db, "IamUserIwillDuplicateMyself", "000000")
	if err == nil {
		t.Fail()
	}
}

func Test_CreateAccountsWithEmptyPassword(t *testing.T) {
	db := initDB(t)
	defer destroyDB(db, t)
	_, err := database.CreateAccount(db, "Passwordless27", "")
	if err == nil {
		t.FailNow()
	}
}

func Test_RetrieveAccounts(t *testing.T) {
	db := initDB(t)
	defer destroyDB(db, t)

	id_heavybro, err := database.CreateAccount(db, "HeavyBro", "1337")
	if err != nil {
		t.FailNow()
	}
	id_hauser, err := database.CreateAccount(db, "Hauser", "999")
	if err != nil {
		t.FailNow()
	}
	retAcc, err := database.GetAccount(id_heavybro, db) // HeavyBro Account
	if err != nil {
		t.FailNow()
	}
	if retAcc.Username != "HeavyBro" || retAcc.Password != utils.HashString("1337") {
		t.Fail()
	}
	retAcc, err = database.GetAccount(id_hauser, db) // Hauser Account
	if err != nil {
		t.FailNow()
	}
	if retAcc.Username != "Hauser" || retAcc.Password != utils.HashString("999") {
		t.Fail()
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
