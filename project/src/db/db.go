package db

import (
  "say"
  "conf"
  "time"
  "errors"
  "github.com/boltdb/bolt"
)

var DB *bolt.DB

type Schema struct {
  Key string
  Children map[string]Schema
}

func Init(){
  say.L1("DB: INIT: start init function")
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
  say.L1("DB: INIT: Done")
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
  say.L1("DB: PUT PAIR: open bucket for WRITE [ " + pathstr + " ]")
  if err := DB.Update(func(tx *bolt.Tx) error {
    b[0] = tx.Bucket([]byte("repositories"))
    for i, e := range path {
      if b[i] != nil {
        b[i+1] = b[i].Bucket([]byte(e))
      } else {
        say.L1("DB: PUT PAIR: creating bucket [ " + path[i-1] + " ]")
        if b[i+1], err = b[i].CreateBucketIfNotExists([]byte(path[i-1])); err != nil {
          return err
        }
      }
    }
    if b[len(path)] == nil {
      say.L1("DB: PUT PAIR: creating bucket [ " + path[len(path)-1] + " ]")
      if b[len(path)], err = b[len(path)-1].CreateBucketIfNotExists([]byte(path[len(path)-1])); err != nil {
        return err
      }
    }
    say.L1("DB: PUT PAIR: putting key in bucket [ " + key + " ]")
    b[len(path)].Put([]byte(key), []byte(value))
    return nil
  }); err != nil {
    say.L3(err.Error())
  }
  say.L1("DB: PUT PAIR: Done")
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
  say.L1("DB: CREATE BUCKET: open bucket for WRITE [ " + pathstr + " ]")
  if err := DB.Update(func(tx *bolt.Tx) error {
    b[0] = tx.Bucket([]byte("repositories"))
    for i, e := range path {
      if b[i] != nil {
        b[i+1] = b[i].Bucket([]byte(e))
      } else {
        say.L1("DB: CREATE BUCKET: creating bucket [ " + path[i-1] + " ]")
        if b[i], err = b[i-1].CreateBucketIfNotExists([]byte(path[i-1])); err != nil {
          return err
        }
      }
    }
    say.L1("DB: CREATE BUCKET: creating bucket [ " + path[len(path)-1] + " ]")
    if b[len(path)], err = b[len(path)-1].CreateBucketIfNotExists([]byte(path[len(path)-1])); err != nil {
      return err
    }
    return nil
  }); err != nil {
    say.L3(err.Error())
  }
  say.L1("DB: CREATE BUCKET: Done")
}

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
  say.L1("DB: DELETE BUCKET: open bucket for DELETE [ " + pathstr + " ]")
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
    say.L1("DB: DELETE BUCKET: deleting bucket [ " + path[len(path)-1] + " ]")
    if err = b[len(path)-1].DeleteBucket([]byte(path[len(path)-1])); err != nil {
      return err
    }
    return nil
  }); err != nil {
    say.L3(err.Error())
  }
  say.L1("DB: DELETE BUCKET: Done")
}
func DeleteKeyFromDB(path []string, key string ) {
  b := make([]*bolt.Bucket, len(path)+1)
  pathstr := ""
  if len(path) > 0 {
    pathstr = path[0]
  }
  for i := 1; i < len(path); i++ {
    pathstr = pathstr + "->" + path[i]
  }
  say.L1("DB: DELETE KEY: open bucket for DELETE [ " + pathstr + " ]")
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
    say.L1("DB: DELETE KEY: deleting key [ " + key + " ]")
    if err := b[len(path)].Delete([]byte(key)); err != nil {
      return err
    }
    return nil
  }); err != nil {
    say.L3(err.Error())
  }
  say.L1("DB: DELETE KEY: Done")
}

func GetSchemaFromPoint(path []string)(schema string){
  b := []*bolt.Bucket{}
  pathstr := ""
  if len(path) > 0 {
    pathstr = path[0]
  }
  for i := 1; i < len(path); i++ {
    pathstr = pathstr + "->" + path[i]
  }
  say.L1("DB: GET SCHEMA: open bucket for READ  [ " + pathstr + " ]")
  if err := DB.View(func(tx *bolt.Tx) error {
    b = append(b, tx.Bucket([]byte("repositories")))
    for i, e := range path {
      if b[i] != nil {
        b = append(b, b[i].Bucket([]byte(e)))
      } else {
        return errors.New("DB: GET SCHEMA: There is no such bucket [ " + path[i-1] + " ]")
      }
    }
    if b[len(path)] == nil {
      return errors.New("DB: GET SCHEMA: There is no such bucket [ " + path[len(path)-1] + " ]")
    }
    schema = schema2json(build_schema_recursive(b[len(path)], "root"))
    return nil
  }); err != nil {
    say.L3(err.Error())
  }
  say.L1("DB: GET SCHEMA: Done")
  return
}
func build_schema_recursive(b *bolt.Bucket, s string) (_sch Schema) {
  _psc := make(map[string]Schema)
  _ = b.ForEach(func(k, v []byte) error {
    bk := b.Bucket([]byte(k))
    if (bk != nil) {
      _psc[string(k)] = build_schema_recursive(bk, string(k))
      return nil
    } else {
      _psc[string(k)] = Schema{ string(k), nil }
      return nil
    }
    return nil
  });
  _sch = Schema{s, _psc}
  return
}
func schema2json(schema Schema) (json string) {
  cnum := len(schema.Children)
  if cnum == 0 {
    json = "{\"text\":{\"name\":\"" + schema.Key + "\"}}"
  } else {
    json = "{\"text\":{\"name\":\"" + schema.Key + "\"}, \"children\": ["
    iter := 0
    for k, _ := range schema.Children {
      iter++
      json += schema2json(schema.Children[k])
      if iter < cnum {
        json += ","
      }
    }
    json += "]}"
  }
  return
}
