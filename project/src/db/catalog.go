package db

import (
  "say"
  "github.com/boltdb/bolt"
)

func GetCatalog(repo string) (catalog []string){
  say.Info("DB: select catalog for [" + repo + "]")
  if err := DB.View(func(tx *bolt.Tx) error {
    if b := tx.Bucket([]byte("repositories")); b != nil {
      if br := b.Bucket([]byte(repo)); br != nil {
        if brc := br.Bucket([]byte("catalog")); brc != nil {
          if err := brc.ForEach(func(k, v []byte) error {
            if brcn := brc.Bucket(k); brcn != nil {
              if _valid := string(brcn.Get([]byte("_valid"))); _valid == "1"{
                catalog = append(catalog, string(k))
              }
            }
            return nil
          }); err != nil {
            return err
          }
        }
      }
    }
    return nil
  }); err != nil {
    say.Error(err.Error())
  }
  say.Info("DB: Done")
  return
}
func AddCatalog(repo string, catalog []string) {
  say.Info("DB: insert catalog for [" + repo + "]")
  if err := DB.Update(func(tx *bolt.Tx) error {
    if b := tx.Bucket([]byte("repositories")); b != nil {
      if br := b.Bucket([]byte(repo)); br != nil {
        if brc, err := br.CreateBucketIfNotExists([]byte("catalog")); err == nil {
          if err := brc.ForEach(func(k, v []byte) error {
            if brcn := brc.Bucket(k); brcn != nil {
              brcn.Put([]byte("_valid"), []byte("0"))
            }
            return nil
          }); err != nil {
            return err
          }
          for _, e := range catalog{
            if brcn, err := brc.CreateBucketIfNotExists([]byte(e)); err == nil {
              brcn.Put([]byte("_valid"), []byte("1"))
            } else {
              return err
            }
          }
        } else {
          return err
        }
      }
    }
    return nil
  }); err != nil {
    say.Error(err.Error())
  }
  say.Info("DB: Done")
}
