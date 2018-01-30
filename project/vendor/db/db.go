package db

import (
  "say"
  "conf"
  "time"
  "github.com/boltdb/bolt"
)

var DB *bolt.DB

func Init(){
  say.L2("DB: INIT: Start")
  var err error
  if DB, err = bolt.Open(conf.Env["dbpath"] + "/" + conf.Env["dbname"] + ".db", 0600,
    &bolt.Options{Timeout: 1 * time.Second}); err != nil {
      say.L3("DB: INIT: OPEN FILE: " + err.Error())
  }
  if err = DB.Update(func(tx *bolt.Tx) error {
    if _, err := tx.CreateBucketIfNotExists([]byte("repositories")); err != nil {
      return err
    } else {
      return nil
    }
  }); err != nil {
    say.L3("DB: INIT: CREATE ROOT POINT: " + err.Error())
  }

  for Upgrade() {}

  say.L2("DB: INIT: Done")
}

func Upgrade() (repeat bool){
  repeat = true
  say.L1("DB: INIT: DB Upgrade: Start")
  if _, ok := GetAllPairsFromBucket([]string{})["_info"]; !ok {
    PutBucketToBucket([]string{"_info"})
  }
  info := GetAllPairsFromBucket([]string{"_info"})
  if len(info) == 0 {
    upto1()
  } else {
    if version, ok := info["version"]; ok {
      say.L2("DB: INIT: DB Upgrade: Version: " + version)
      switch version{
      case "1":
        upto2()
      case "2":
        upto3()
      default:
        say.L2("DB: INIT: DB Upgrade: Actual version")
        repeat = false
      }
    } else {
      upto1()
    }
  }
  say.L1("DB: INIT: DB Upgrade: Finish")
  return
}
