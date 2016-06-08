package db

import (
  "say"
  "github.com/boltdb/bolt"
)

func GetManifest(repo string, name string, tag string) (manifest []string){
  say.Info("DB: select manifest for [" + repo + "/" + name + "/" + tag + "]")
  if err := DB.View(func(tx *bolt.Tx) error {
    if b := tx.Bucket([]byte("repositories")); b != nil {
      if br := b.Bucket([]byte(repo)); br != nil {
        if brc := br.Bucket([]byte("catalog")); brc != nil {
          if brcn := brc.Bucket([]byte(name)); brcn != nil {
            if brcnt := brcn.Bucket([]byte(tag)); brcnt != nil {

            }
            // if err := brcn.ForEach(func(k, v []byte) error {
            //   if brcnt := brcn.Bucket(k); brcnt != nil {
            //     if _valid := string(brcnt.Get([]byte("_valid"))); _valid == "1"{
            //       tags = append(tags, string(k))
            //     }
            //   }
            //   return nil
            // }); err != nil {
            //   return err
            // }
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
