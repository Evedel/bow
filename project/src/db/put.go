package db

import (
  "github.com/boltdb/bolt"
  "github.com/Evedel/glb/say"
)

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
  say.L3("DB: PUT PAIR: open bucket for WRITE [ ", pathstr, " ]\n")
  if err := DB.Update(func(tx *bolt.Tx) error {
    b[0] = tx.Bucket([]byte("repositories"))
    for i, e := range path {
      if b[i] != nil {
        b[i+1] = b[i].Bucket([]byte(e))
      } else {
        say.L3("DB: PUT PAIR: creating bucket [ ", path[i-1], " ]\n")
        if b[i], err = b[i-1].CreateBucketIfNotExists([]byte(path[i-1])); err != nil {
          return err
        }
      }
    }
    if b[len(path)] == nil {
      say.L3("DB: PUT PAIR: creating bucket [ ", path[len(path)-1], " ]\n")
      if b[len(path)], err = b[len(path)-1].CreateBucketIfNotExists([]byte(path[len(path)-1])); err != nil {
        return err
      }
    }
    say.L3("DB: PUT PAIR: putting key in bucket [ ", key, " ]\n")
    b[len(path)].Put([]byte(key), []byte(value))
    return nil
  }); err != nil {
    say.L1("", err, "\n")
  }
  say.L3("DB: PUT PAIR: Done.", "", "\n")
}

func PutBucketToBucket(path []string){
  var err error
  b := make([]*bolt.Bucket, len(path)+1)
  pathstr := ""
  if len(path) > 0 {
    pathstr = path[0]
  }
  for i := 1; i < len(path); i++ {
    pathstr = pathstr + "->" + path[i]
  }
  say.L3("DB: CREATE BUCKET: open bucket for WRITE [ ", pathstr, " ]\n")
  if err := DB.Update(func(tx *bolt.Tx) error {
    b[0] = tx.Bucket([]byte("repositories"))
    for i, e := range path {
      if b[i] != nil {
        b[i+1] = b[i].Bucket([]byte(e))
      } else {
        say.L3("DB: CREATE BUCKET: creating bucket [ ", path[i-1], " ]\n")
        if b[i], err = b[i-1].CreateBucketIfNotExists([]byte(path[i-1])); err != nil {
          return err
        }
      }
    }
    if b[len(path)-1].Bucket([]byte(path[len(path)-1])) == nil {
      say.L3("DB: CREATE BUCKET: creating bucket [ ", path[len(path)-1], " ]\n")
      if b[len(path)], err = b[len(path)-1].CreateBucketIfNotExists([]byte(path[len(path)-1])); err != nil {
        return err
      }
    } else {
      say.L3("DB: CREATE BUCKET: bucket already exist [ ", pathstr, " ]\n")
    }
    return nil
  }); err != nil {
    say.L1("", err, "\n")
  }
  say.L3("DB: CREATE BUCKET: Done.", "", "\n")
}
