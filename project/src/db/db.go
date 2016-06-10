package db

import (
  "say"
  "conf"
  "time"
  "errors"
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
func GetSimplePairsFromBucket(path []string) (pairs map[string]string){
  b := make([]*bolt.Bucket, len(path)+1)
  pairs = make(map[string]string)
  say.Info("DB: open bucket [ " + path[0] + " ]")
  if err := DB.View(func(tx *bolt.Tx) error {
    b[0] = tx.Bucket([]byte("repositories"))
    for i, e := range path {
      if b[i] != nil {
        b[i+1] = b[i].Bucket([]byte(e))
      } else {
        return errors.New("There is no such bucket [ " + path[i-1] + " ]")
      }
    }
    if err := b[len(path)].ForEach(func(k, v []byte) error {
      pairs[string(k)] = string(v)
      return nil
    }); err != nil {
      return err
    }
    return nil
  }); err != nil {
    say.Error(err.Error())
  }
  say.Info("DB: Done")
  return
}
