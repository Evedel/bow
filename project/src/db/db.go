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
  pathstr := ""
  if len(path) > 0 {
    pathstr = path[0]
  }
  for i := 1; i < len(path); i++ {
    pathstr = pathstr + "->" + path[i]
  }
  say.Info("DB: open bucket for READ  [ " + pathstr + " ]")
  if err := DB.View(func(tx *bolt.Tx) error {
    b[0] = tx.Bucket([]byte("repositories"))
    for i, e := range path {
      if b[i] != nil {
        b[i+1] = b[i].Bucket([]byte(e))
      } else {
        return errors.New("There is no such bucket [ " + path[i-1] + " ]")
      }
    }
    if b[len(path)] == nil {
      return errors.New("There is no such bucket [ " + path[len(path)-1] + " ]")
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
func PutSimplePairToBucket(path []string, key string, value string){
  var err error
  b := make([]*bolt.Bucket, len(path)+1)
  pathstr := ""
  if len(path) > 0 {
    pathstr = path[0]
  }
  for i := 1; i < len(path); i++ {
    pathstr = pathstr + "->" + path[i]
  }
  say.Info("DB: open bucket for WRITE [ " + pathstr + " ]")
  if err := DB.Update(func(tx *bolt.Tx) error {
    b[0] = tx.Bucket([]byte("repositories"))
    for i, e := range path {
      if b[i] != nil {
        b[i+1] = b[i].Bucket([]byte(e))
      } else {
        say.Info("DB : creating bucket [ " + path[i-1] + " ]")
        if b[i+1], err = b[i].CreateBucketIfNotExists([]byte(path[i-1])); err != nil {
          return err
        }
      }
    }
    if b[len(path)] == nil {
      say.Info("DB: creating bucket [ " + path[len(path)-1] + " ]")
      if b[len(path)], err = b[len(path)-1].CreateBucketIfNotExists([]byte(path[len(path)-1])); err != nil {
        return err
      }
    }
    say.Info("DB: putting key in bucket [ " + key + " ]")
    b[len(path)].Put([]byte(key), []byte(value))
    return nil
  }); err != nil {
    say.Error(err.Error())
  }
  say.Info("DB: Done")
}
func DeleteBucketFromDB(path []string) {
  var err error
  b := make([]*bolt.Bucket, len(path)+1)
  pathstr := ""
  if len(path) > 0 {
    pathstr = path[0]
  }
  for i := 1; i < len(path); i++ {
    pathstr = pathstr + "->" + path[i]
  }
  say.Info("DB: open bucket for DELETE [ " + pathstr + " ]")
  if err := DB.Update(func(tx *bolt.Tx) error {
    b[0] = tx.Bucket([]byte("repositories"))
    for i, e := range path {
      if b[i] != nil {
        b[i+1] = b[i].Bucket([]byte(e))
      } else {
        return errors.New("There is no such bucket [ " + path[i-1] + " ]")
      }
    }
    if b[len(path)] == nil {
      return errors.New("There is no such bucket [ " + path[len(path)-1] + " ]")
    }
    say.Info("DB: deleting bucket [ " + path[len(path)-1] + " ]")
    if err = b[len(path)-1].DeleteBucket([]byte(path[len(path)-1])); err != nil {
      return err
    }
    return nil
  }); err != nil {
    say.Error(err.Error())
  }
  say.Info("DB: Done")
}
