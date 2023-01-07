package main_test

import (
	akevitt "akevitt/cmd"
	"fmt"
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
	acc := akevitt.Account{Username: "", Password: ""}
	_, err := akevitt.CreateAccount(db, acc)
	if err == nil { // No errors were made. But it should
		fmt.Printf("Test failed: %s", err.Error())
		t.Fail()
	}

}
func Test_CreateDuplicateAccounts(t *testing.T) {
	db := initDB(t)
	defer destroyDB(db, t)
	// Register
	dupAcc := akevitt.Account{Username: "IamUserIwillDuplicateMyself", Password: "000000"}
	_, err := akevitt.CreateAccount(db, dupAcc)
	if err != nil {
		t.FailNow()
	}
	_, err = akevitt.CreateAccount(db, dupAcc)
	if err == nil {
		t.Fail()
	}
}

func Test_CreateAccountsWithEmptyPassword(t *testing.T) {
	db := initDB(t)
	defer destroyDB(db, t)
	pwdlessAcc := akevitt.Account{Username: "Passwordless27", Password: ""}
	_, err := akevitt.CreateAccount(db, pwdlessAcc)
	if err == nil {
		t.FailNow()
	}
}

func Test_RetrieveAccounts(t *testing.T) {
	db := initDB(t)
	defer destroyDB(db, t)
	hbAcc := akevitt.Account{Username: "HeavyBro", Password: "1337"}
	hsrAcc := akevitt.Account{Username: "Hauser", Password: "999"}

	id_heavybro, err := akevitt.CreateAccount(db, hbAcc)
	if err != nil {
		t.FailNow()
	}
	id_hauser, err := akevitt.CreateAccount(db, hsrAcc)
	if err != nil {
		t.FailNow()
	}
	retAcc, err := akevitt.GetAccount(id_heavybro, db) // HeavyBro Account
	if err != nil {
		t.FailNow()
	}
	if retAcc.Username != "HeavyBro" || retAcc.Password != "1337" {
		t.Fail()
	}
	retAcc, err = akevitt.GetAccount(id_hauser, db) // Hauser Account
	if err != nil {
		t.FailNow()
	}
	if retAcc.Username != "Hauser" || retAcc.Password != "999" {
		t.Fail()
	}
}
