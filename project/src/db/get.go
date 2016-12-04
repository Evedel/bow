package db

import (
  "say"
  "errors"
  "github.com/boltdb/bolt"
)

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
  if pathstr == "" {
    pathstr = "RootPoint"
  }
  say.L1("DB: GET BUCKET: open bucket for READ  [ " + pathstr + " ]")
  if err := DB.View(func(tx *bolt.Tx) error {
    b[0] = tx.Bucket([]byte("repositories"))
    for i, e := range path {
      if b[i] != nil {
        b[i+1] = b[i].Bucket([]byte(e))
      } else {
        return errors.New("DB: GET BUCKET: There is no such bucket [ " + path[i-1] + " ]")
      }
    }
    if b[len(path)] == nil {
      return errors.New("DB: GET BUCKET: There is no such bucket [ " + path[len(path)-1] + " ]")
    }
    if err := b[len(path)].ForEach(func(k, v []byte) error {
      pairs[string(k)] = string(v)
      return nil
    }); err != nil {
      return err
    }
    return nil
  }); err != nil {
    say.L3(err.Error())
  }
  say.L1("DB: GET BUCKET: Done")
  return
}
func GetValueFromBucket(path []string, key string) (value string){
  b := make([]*bolt.Bucket, len(path)+1)
  pathstr := ""
  if len(path) > 0 {
    pathstr = path[0]
  }
  for i := 1; i < len(path); i++ {
    pathstr = pathstr + "->" + path[i]
  }
  say.L1("DB: GET VALUE: open bucket for READ  [ " + pathstr + "=>" + key + " ]")
  if err := DB.View(func(tx *bolt.Tx) error {
    b[0] = tx.Bucket([]byte("repositories"))
    for i, e := range path {
      if b[i] != nil {
        b[i+1] = b[i].Bucket([]byte(e))
      } else {
        return errors.New("DB: GET VALUE: There is no such bucket [ " + path[i-1] + " ]")
      }
    }
    if b[len(path)] == nil {
      return errors.New("DB: GET VALUE: There is no such bucket [ " + path[len(path)-1] + " ]")
    }
    value =  string(b[len(path)].Get([]byte(key)))
    return nil
  }); err != nil {
    say.L3(err.Error())
  }
  say.L1("DB: GET VALUE: Done")
  return
}
