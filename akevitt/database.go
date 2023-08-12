package akevitt

import "github.com/boltdb/bolt"

func createDatabase(engine *Akevitt) error {
	db, err := bolt.Open(engine.dbPath, 0600, nil)
	engine.db = db

	return err
}
