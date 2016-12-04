package db

import (
  "say"
  "errors"
  "github.com/boltdb/bolt"
)

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
