package db

import (
  "conf"
  "time"

  "github.com/boltdb/bolt"
  "github.com/Evedel/glb/say"
)

var DB *bolt.DB

func Init(){
  say.L2("DB: INIT: Start", "","\n")
  var err error
  if DB, err = bolt.Open(conf.Env["dbpath"] + "/" + conf.Env["dbname"] + ".db", 0600,
    &bolt.Options{Timeout: 1 * time.Second}); err != nil {
      say.L1("DB: INIT: Cannot open file error: ", err.Error(), "\n")
  }

  for Upgrade() {}

  say.L2("DB: INIT: Done.", "","\n")
}

func Upgrade() (repeat bool){
  repeat = true
  say.L3("DB: INIT: DB Upgrade: Start.", "","\n")
  if _, ok := GetAllPairsFromBucket([]string{})["_info"]; !ok {
    PutBucketToBucket([]string{"_info"})
  }
  info := GetAllPairsFromBucket([]string{"_info"})
  if len(info) == 0 {
    upto1()
  } else {
    if version, ok := info["version"]; ok {
      say.L3("DB: INIT: DB Upgrade: Version: ",version,"\n")
      switch version{
      case "1":
        upto2()
      case "2":
        upto3()
      default:
        say.L3("DB: INIT: DB Upgrade: Actual version.", "","\n")
        repeat = false
      }
    } else {
      upto1()
    }
  }
  say.L3("DB: INIT: DB Upgrade: Finish", "","\n")
  return
}
