  package db

import (
  "say"
  "github.com/boltdb/bolt"
)

func GetCatalog(repo string) (catalog []string){
  say.L1("DB: select catalog for [" + repo + "]")
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
    say.L3(err.Error())
  }
  say.L1("DB: Done")
  return
}
func AddCatalog(repo string, catalog []string) {
  say.L1("DB: insert catalog for [" + repo + "]")
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
    say.L3(err.Error())
  }
  say.L1("DB: Done")
}
func PutCatalogSubBucket(repo string, bucket string, key string, value string){
  say.L1("DB: insert in subbucket for catalog [ " + repo + "/" + bucket + " / " + key + " ]")
  if err := DB.Update(func(tx *bolt.Tx) error {
    if b := tx.Bucket([]byte("repositories")); b != nil {
      if br := b.Bucket([]byte(repo)); br != nil {
        if brs, err := br.CreateBucketIfNotExists([]byte(bucket)); err == nil {
          brs.Put([]byte(key), []byte(value))
        } else {
          return err
        }
      }
    }
    return nil
  }); err != nil {
    say.L3(err.Error())
  }
  say.L1("DB: Done ")
}
