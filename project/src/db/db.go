package db

import (
  "say"
  "conf"
  "time"
  "github.com/boltdb/bolt"
)

var DB *bolt.DB

func Init(){
  say.Info("DB: start init function")
  var err error
  DB, err = bolt.Open(conf.Env["dbpath"] + "/" + conf.Env["dbname"] + ".db", 0600,
    &bolt.Options{Timeout: 1 * time.Second})
  if err != nil {
    say.Error(err.Error())
  }
  err = DB.Update(func(tx *bolt.Tx) error {
    _, err := tx.CreateBucketIfNotExists([]byte("repositories"))
	if err != nil {
  	   return err
  	}
  	return nil
  })
  if err != nil{
    say.Error(err.Error())
  }
  say.Info("Done")
}
