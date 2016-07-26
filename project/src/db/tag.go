package db

import (
  "strconv"
  "time"
  "say"
  "github.com/boltdb/bolt"
)

func GetTags(repo string, name string) (tags []string){
  say.L1("DB: select tags for [" + repo + "/" + name + "]")
  if err := DB.View(func(tx *bolt.Tx) error {
    if b := tx.Bucket([]byte("repositories")); b != nil {
      if br := b.Bucket([]byte(repo)); br != nil {
        if brc := br.Bucket([]byte("catalog")); brc != nil {
          if brcn := brc.Bucket([]byte(name)); brcn != nil {
            if err := brcn.ForEach(func(k, v []byte) error {
              if brcnt := brcn.Bucket(k); brcnt != nil {
                if _valid := string(brcnt.Get([]byte("_valid"))); _valid == "1"{
                  tags = append(tags, string(k))
                }
              }
              return nil
            }); err != nil {
              return err
            }
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
func AddTags(repo string, name string, tags []string){
  say.L1("DB: insert tags for [" + repo + "/" + name + "]")
  if err := DB.Update(func(tx *bolt.Tx) error {
    if b := tx.Bucket([]byte("repositories")); b != nil {
      if br := b.Bucket([]byte(repo)); br != nil {
        if brc := br.Bucket([]byte("catalog")); brc != nil {
          if brcn := brc.Bucket([]byte(name)); brcn != nil {
            if err := brcn.ForEach(func(k, v []byte) error {
              if brcnt := brcn.Bucket(k); brcnt != nil {
                brcnt.Put([]byte("_valid"), []byte("0"))
              }
              return nil
            }); err != nil {
              return err
            }
            for _, e := range tags{
              if brcnt, err := brcn.CreateBucketIfNotExists([]byte(e)); err == nil {
                brcnt.Put([]byte("_valid"), []byte("1"))
              }
            }
          }
        }
      }
    }
    return nil
  }); err != nil {
    say.L3(err.Error())
  }
  say.L1("DB: Done")
}
func GetTagDigest(repo string, name string, tag string) (digest string){
  say.L1("DB: select digest for [" + repo + "/" + name + "/" + tag + "]")
  if err := DB.Update(func(tx *bolt.Tx) error {
    if b := tx.Bucket([]byte("repositories")); b != nil {
      if br := b.Bucket([]byte(repo)); br != nil {
        if brc := br.Bucket([]byte("catalog")); brc != nil {
          if brcn := brc.Bucket([]byte(name)); brcn != nil {
            if brcnt := brcn.Bucket([]byte(tag)); brcnt != nil {
              digest = string(brcnt.Get([]byte("digest")));
            }
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
func PutTagDigest(repo string, name string, tag string, digest string){
  say.L1("DB: insert digest for [" + repo + "/" + name + "/" + tag + "]")
  if err := DB.Update(func(tx *bolt.Tx) error {
    if b := tx.Bucket([]byte("repositories")); b != nil {
      if br := b.Bucket([]byte(repo)); br != nil {
        if brc := br.Bucket([]byte("catalog")); brc != nil {
          if brcn := brc.Bucket([]byte(name)); brcn != nil {
            if brcnt := brcn.Bucket([]byte(tag)); brcnt != nil {
              brcnt.Put([]byte("digest"), []byte(digest))
              if brcntu, err := brcnt.CreateBucketIfNotExists([]byte("_uploads")); err == nil {
                say.L1("DB: incrimenting uploads for [" + repo + "/" + name + "/" + tag + "]")
                shortDate := time.Now().Local().Format("2006-01-02")
                if brcntud := brcntu.Get([]byte(shortDate)); brcntud != nil {
                  val, _ := strconv.Atoi(string(brcntud))
                  val++
                  brcntu.Put([]byte(shortDate), []byte(strconv.Itoa(val)))
                } else {
                  brcntu.Put([]byte(shortDate), []byte("1"))
                }
              } else {
                return err
              }
              if brcnu, err := brcn.CreateBucketIfNotExists([]byte("_uploads")); err == nil {
                say.L1("DB: incrimenting uploads for [" + repo + "/" + name + "]")
                shortDate := time.Now().Local().Format("2006-01-02")
                if brcnud := brcnu.Get([]byte(shortDate)); brcnud != nil {
                  val, _ := strconv.Atoi(string(brcnud))
                  val++
                  brcnu.Put([]byte(shortDate), []byte(strconv.Itoa(val)))
                } else {
                  brcnu.Put([]byte(shortDate), []byte("1"))
                }
              } else {
                return err
              }
            }
          }
        }
      }
    }
    return nil
  }); err != nil {
    say.L3(err.Error())
  }
  say.L1("DB: Done")
}
func GetTagSubbucket(repo string, name string, tag string, bucket string) (data map[string]string){
  data = make(map[string]string)
  say.L1("DB: select manifest for [" + repo + "/" + name + "/" + tag + "]")
  if err := DB.Update(func(tx *bolt.Tx) error {
    if b := tx.Bucket([]byte("repositories")); b != nil {
      if br := b.Bucket([]byte(repo)); br != nil {
        if brc := br.Bucket([]byte("catalog")); brc != nil {
          if brcn := brc.Bucket([]byte(name)); brcn != nil {
            if brcnt := brcn.Bucket([]byte(tag)); brcnt != nil {
              if brcnts := brcnt.Bucket([]byte(bucket)); brcnts != nil {
                if err := brcnts.ForEach(func(k, v []byte) error {
                  data[string(k)] = string(v)
                  return nil
                }); err != nil {
                  return err
                }
              }
            }
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
func DeleteTagSubBucket(repo string, name string, tag string, bucket string){
  say.L1("DB: delete subbucket for [ " + repo + "/" + name + "/" + tag + "/" + bucket + " ]")
  if err := DB.Update(func(tx *bolt.Tx) error {
    if b := tx.Bucket([]byte("repositories")); b != nil {
      if br := b.Bucket([]byte(repo)); br != nil {
        if brc := br.Bucket([]byte("catalog")); brc != nil {
          if brcn := brc.Bucket([]byte(name)); brcn != nil {
            if brcnt := brcn.Bucket([]byte(tag)); brcnt != nil {
              if brcnts := brcnt.Bucket([]byte(bucket)); brcnts != nil {
                if err := brcnt.DeleteBucket([]byte(bucket)); err != nil {
                  return err
                }
              }
            }
          }
        }
      }
    }
    return nil
  }); err != nil {
    say.L3(err.Error())
  }
  say.L1("DB: Done")
}
