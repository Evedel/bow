package db

import (
  "say"
  "github.com/boltdb/bolt"
)

func GetRepos() (repos []string){
  say.L1("DB: select list of repos")
  if err := DB.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte("repositories")); b != nil {
      err := b.ForEach(func(k, v []byte) error {
        repos = append(repos, string(k))
        return nil
      })
      return err
    }
    return nil
	}); err != nil {
    say.L3(err.Error())
	}
  say.L1("DB: Done")
  return
}
func CreateRepo(params map[string][]string) {
  say.L1("DB: insert repository info [" + params["reponame"][0] + "]")
  if err := DB.Update(func(tx *bolt.Tx) error {
    if b := tx.Bucket([]byte("repositories")); b != nil {
      if br, err := b.CreateBucketIfNotExists([]byte(params["reponame"][0])); err == nil {
        br.Put([]byte("repohost"), []byte(params["repohost"][0]))
        br.Put([]byte("repouser"), []byte(params["repouser"][0]))
        br.Put([]byte("repopass"), []byte(params["repopass"][0]))
      } else {
        return err
      }
    }
    return nil
  }); err != nil {
    say.L3(err.Error())
  }
  say.L1("DB: Done")
}
func GetRepoPretty(repo string) (pretty map[string]string){
  say.L1("DB: select pretty info for [" + repo + "]")
  pretty = make(map[string]string)
  if err := DB.View(func(tx *bolt.Tx) error {
    if b := tx.Bucket([]byte("repositories")); b != nil {
      if br := b.Bucket([]byte(repo)); br != nil {
        pretty["reponame"] = repo
        err := br.ForEach(func(k, v []byte) error {
          pretty[string(k)] = string(v)
          return nil
          });
        return err
      }
    }
    return nil
  }); err != nil {
    say.L3(err.Error())
  }
  say.L1("DB: Done")
  return
}
func DeleteRepo(repo string){
  say.L1("DB: Delete repository [" + repo + "]")
  if err := DB.Update(func(tx *bolt.Tx) error {
    if b := tx.Bucket([]byte("repositories")); b != nil {
      err := b.DeleteBucket([]byte(repo))
      return err
    }
    return nil
  }); err != nil {
    say.L3(err.Error())
  }
  say.L1("Done")
}
