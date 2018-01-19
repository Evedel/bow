package db

import (
  "errors"
  "strings"

  "github.com/boltdb/bolt"
  "github.com/Evedel/glb/say"
)

type Schema struct {
  Key string
  Children map[string]Schema
}

func GetSchemaFromPoint(path []string, filter string)(schema Schema){
  b := []*bolt.Bucket{}
  pathstr := ""
  if len(path) > 0 {
    pathstr = path[0]
  }
  for i := 1; i < len(path); i++ {
    pathstr = pathstr + "->" + path[i]
  }
  say.L3("DB: GET SCHEMA: open bucket for READ  [ ", pathstr, " ]\n")
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
    schema, _ = buildSchemaRecursive(b[len(path)], "root", filter, false)
    return nil
  }); err != nil {
    say.L1("", err, "\n")
  }
  say.L3("DB: GET SCHEMA: Done.", "", "\n")
  return
}

func buildSchemaRecursive(b *bolt.Bucket, s string, f string, b0 bool) (_sch Schema, bi bool) {
  bl := false
  if (f == "") || (b0) {
    bl = true
  } else {
    bl = strings.Contains(s, f)
    bi = bl
  }

  _psc := make(map[string]Schema)
  _ = b.ForEach(func(k, v []byte) error {
    bk := b.Bucket([]byte(k))
    if (bk != nil) {
      schtmp, bn := buildSchemaRecursive(bk, string(k), f, bl)
      if bl || bn {
        _psc[string(k)] = schtmp
        bi = true
      }
      return nil
    }
    return nil
  });
  _sch = Schema{s, _psc}
  return
}

func Schema2json(schema Schema) (json string) {
  cnum := len(schema.Children)
  if cnum == 0 {
    json = "{\"text\":{\"name\":\"" + schema.Key + "\"}}"
  } else {
    json = "{\"text\":{\"name\":\"" + schema.Key + "\"}, \"children\": ["
    iter := 0
    for k, _ := range schema.Children {
      iter++
      json += Schema2json(schema.Children[k])
      if iter < cnum {
        json += ","
      }
    }
    json += "]}"
  }
  return
}
