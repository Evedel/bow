package db

import (
  "errors"

  "github.com/boltdb/bolt"
  "github.com/Evedel/glb/say"
)

func DeleteBucket(path []string) {
  var err error
  b := make([]*bolt.Bucket, len(path)+1)
  pathstr := ""
  if len(path) > 0 {
    pathstr = path[0]
  }
  for i := 1; i < len(path); i++ {
    pathstr = pathstr + "->" + path[i]
  }
  say.L3("DB: DELETE BUCKET: open bucket for DELETE [ ",pathstr," ]\n")
  if err := DB.Update(func(tx *bolt.Tx) error {
    b[0] = tx.Bucket([]byte("repositories"))
    for i, e := range path {
      if b[i] != nil {
        b[i+1] = b[i].Bucket([]byte(e))
      } else {
        return errors.New("DB: DELETE BUCKET: There is no such bucket [ " + path[i-1] + " ]")
      }
    }
    if b[len(path)] == nil {
      return errors.New("DB: DELETE BUCKET: There is no such bucket [ " + path[len(path)-1] + " ]")
    }
    say.L3("DB: DELETE BUCKET: deleting bucket [ ",path[len(path)-1]," ]\n")
    if err = b[len(path)-1].DeleteBucket([]byte(path[len(path)-1])); err != nil {
      return err
    }
    return nil
  }); err != nil {
    say.L1("",err,"\n")
  }
  say.L3("DB: DELETE BUCKET: Done.", "","\n")
}

func DeleteKey(path []string, key string ) {
  b := make([]*bolt.Bucket, len(path)+1)
  pathstr := ""
  if len(path) > 0 {
    pathstr = path[0]
  }
  for i := 1; i < len(path); i++ {
    pathstr = pathstr + "->" + path[i]
  }
  say.L3("DB: DELETE KEY: open bucket for DELETE [ ",pathstr," ]\n")
  if err := DB.Update(func(tx *bolt.Tx) error {
    b[0] = tx.Bucket([]byte("repositories"))
    for i, e := range path {
      if b[i] != nil {
        b[i+1] = b[i].Bucket([]byte(e))
      } else {
        return errors.New("DB: DELETE KEY: There is no such bucket [ " + path[i-1] + " ]")
      }
    }
    if b[len(path)] == nil {
      return errors.New("DB: DELETE KEY: There is no such bucket [ " + path[len(path)-1] + " ]")
    }
    if b[len(path)].Get([]byte(key)) == nil {
      return errors.New("DB: DELETE KEY: There is no such key [ " + key + " ]")
    } else {
      say.L3("DB: DELETE KEY: deleting key [ ", key, " ]\n")
      if err := b[len(path)].Delete([]byte(key)); err != nil {
        return err
      }
      return nil
    }
  }); err != nil {
    say.L1("", err, "")
  }
  say.L3("DB: DELETE KEY: Done.", "","\n")
}
