package checker

import (
  "db"
  "say"
  "time"
  "strconv"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

func DaemonManager() {
  for {
    say.Info("Manager Daemon: TicTac")
    go CheckRepos()
    go CheckTags()
    go CheckManifests()
    time.Sleep(60 * time.Second)
  }
}
func CatalogNeedUpdate(a []string, b []string) (bool) {
  if a == nil && b == nil {
    say.Info("Catalog are equally nil. Not need update")
    return false
  }
  if a == nil || b == nil {
    say.Info("One of catalog is empty. Need update")
    return true
  }
  if len(a) != len(b) {
    say.Info("Length of catalog are different. Need update")
    return true
  }
  numofequal := 0
  for _, bel := range b {
    for _, ael := range a {
      if bel == ael{
        numofequal++
      }
    }
  }
  if len(a) == numofequal {
    say.Info("Length of catalogs are equal to number of equal elements. Not need update")
    return false
  }
  say.Error("Cant evaluate any condition. Try update")
  say.Raw(a)
  say.Raw(b)
  return true
}
func CheckRepos(){
  say.Info("CheckRepos Daemon: started work")
  repos := db.GetRepos()
  for _, e := range repos {
    pretty := db.GetRepoPretty(e)
    Req := "https://" + pretty["repouser"] +
      ":" + pretty["repopass"] + "@" + pretty["repohost"] + "/v2/_catalog?n=&last="
    if Resp, err := http.Get(Req); err != nil {
      say.Error(err.Error())
      say.Error("CheckRepos Daemon: cannot recieve response from registry, stopping work")
      //TODO
      //[ Tue, 24 May 2016 14:48:36 UTC ]
      //[  ERROR  ]
      //Get https://xxx:xxx@123.45.67.89:5000
      ///v2/_catalog?n=&last=: x509: cannot validate
      //certificate for 123.45.67.89 because it doesn't contain any IP SANs
      //TODO
      //[ Tue, 24 May 2016 14:50:01 UTC ]
      //[  ERROR  ]
      //Get https://xxxxxx:xxxxx@some.wrong.host:5000
      ///v2/_catalog?n=&last=: dial tcp 123.45.67.89:5000: getsockopt: connection refused
    } else {
      if Body, err := ioutil.ReadAll(Resp.Body); err != nil {
        say.Error(err.Error())
      } else {
        //TODO
        //{"errors":[{"code":"UNAUTHORIZED","message":"authentication required",
        //"detail":[{"Type":"registry","Name":"catalog","Action":"*"}]}]}
        var c interface{}
        if err := json.Unmarshal(Body, &c); err != nil {
          say.Error(err.Error())
        } else {
          dbcatalog := db.GetCatalog(e)
          arrint := c.(map[string]interface{})["repositories"].([]interface{})
          arrstr := make([]string, len(arrint))
          for i, _ := range arrint {
            arrstr[i] = arrint[i].(string)
          }
          if CatalogNeedUpdate(dbcatalog, arrstr) {
            db.AddCatalog(e, arrstr)
          }
        }
      }
    }
  }
  say.Info("CheckRepos Daemon: finished work")
}
func CheckTags(){
  say.Info("CheckTags Daemon: started work")
  repos := db.GetRepos()
  for _, er := range repos {
    pretty := db.GetRepoPretty(er)
    catalog := db.GetCatalog(er)
    reponame := "https://" + pretty["repouser"] + ":" + pretty["repopass"] + "@" + pretty["repohost"]
    for _, en := range catalog {
      Reqt := reponame + "/v2/" + en + "/tags/list"
      if Resp, err := http.Get(Reqt); err != nil {
        say.Error(err.Error())
        say.Error("CheckTags Daemon: cannot recieve response from registry, stopping work")
      } else {
        if Body, err := ioutil.ReadAll(Resp.Body); err != nil {
          say.Error(err.Error())
        } else {
          var c interface{}
          if err := json.Unmarshal(Body, &c); err != nil {
            say.Error(err.Error())
          } else {
            if c.(map[string]interface{})["errors"] != nil {
              say.Error(pretty["repohost"] + "/" + en)
              say.Error(c.(map[string]interface{})["errors"].([]interface{})[0].(map[string]interface{})["message"].(string))
            } else {
              dbtags := db.GetTags(er, en)
              arrint := c.(map[string]interface{})["tags"].([]interface{})
              arrstr := make([]string, len(arrint))
              for i, _ := range arrint {
                arrstr[i] = arrint[i].(string)
              }
              if CatalogNeedUpdate(dbtags, arrstr) {
                db.AddTags(er, en, arrstr)
              }
            }
          }
        }
      }
    }
  }
  say.Info("CheckTags Daemon: finished work")
}
func CheckManifests(){
  say.Info("CheckManifests Daemon: started work")
  repos := db.GetRepos()
  for _, er := range repos {
    pretty := db.GetRepoPretty(er)
    catalog := db.GetCatalog(er)
    curlpath := "https://" + pretty["repouser"] + ":" + pretty["repopass"] + "@" + pretty["repohost"]
    for _, en := range catalog {
      dbtags := db.GetTags(er, en)
      for _, et := range dbtags {
        Reqt := curlpath + "/v2/" + en + "/manifests/" + et
        if Resp, err := http.Get(Reqt); err != nil {
          say.Error(err.Error())
          say.Error("CheckManifests Daemon: cannot recieve response from registry, stopping work")
        } else {
          if Body, err := ioutil.ReadAll(Resp.Body); err != nil {
            say.Error(err.Error())
          } else {
            var c interface{}
            if err := json.Unmarshal(Body, &c); err != nil {
              say.Error(err.Error())
            } else {
              dbdigest := db.GetTagDigest(er, en, et)
              curldigest := Resp.Header.Get("Docker-Content-Digest")
              if (dbdigest != curldigest){
                say.Raw(er + " / " + en + " / " + et)
                db.PutTagDigest(er, en, et, curldigest)
                var ch interface{}
                totalsize := 0
                fsshaarr := c.(map[string]interface{})["fsLayers"].([]interface{})
                historyarr := c.(map[string]interface{})["history"].([]interface{})

                db.DeleteTagSubBucket(er, en, et, "history")
                for i, _ := range fsshaarr {
                  fssha := fsshaarr[i].(map[string]interface{})["blobSum"].(string)
                  fssize := GetfsLayerSize(curlpath + "/v2/" + en + "/blobs/" + fssha)
                  history := historyarr[i].(map[string]interface{})["v1Compatibility"].(string)
                  historytrunc := history
                  if last := len(historytrunc) - 1; last >= 0 {
                      historytrunc = historytrunc[:last]
                  }
                  historynew := historytrunc + ",\"blobSum\":\"" + fssha + "\", \"blobSize\":\"" + fssize + "\"}"
                  if fsshanum, err := strconv.Atoi(fssize); err != nil {
                    say.Error(err.Error())
                  } else {
                    totalsize += fsshanum
                  }
                  if err := json.Unmarshal([]byte(history), &ch); err != nil {
                    say.Error(err.Error())
                  } else {
                    created := ch.(map[string]interface{})["created"].(string)
                    if i == 0 {
                      db.PutCatalogSubBucket(er, "_names", ch.(map[string]interface{})["parent"].(string), en + "/" + et)
                    }
                    db.PutTagSubBucket(er, en, et, "history", created, historynew)
                  }
                }
                db.PutTagSubBucket(er, en, et, "_totalsize", time.Now().Local().Format("2006-01-02 15:04:05"), strconv.Itoa(totalsize))
              } else {
                say.Info("CheckManifests Daemon: digests are the same, shouldnot update anything, stopping work")
              }
            }
          }
        }
      }
    }
  }
  say.Info("CheckManifests Daemon: finished work")
}
func GetfsLayerSize(link string ) (size string){
  if Resp, err := http.Head(link); err != nil {
    say.Error(err.Error())
    say.Error("CheckManifests Daemon: GetfsLayerSize cannot recieve response from registry, stopping work")
  } else {
    if _, err := ioutil.ReadAll(Resp.Body); err != nil {
      say.Error(err.Error())
    } else {
      size = Resp.Header.Get("Content-Length")
      return
    }
  }
  return ""
}
