package checker

import(
  "db"
  "say"
  "time"
  "utils"
  "strings"
  "strconv"
  "net/http"
  "encoding/json"
)

func CheckManifests(){
  say.L1("CheckManifests Daemon: started work")
  repos := db.GetRepos()
  for _, er := range repos {
    pretty := db.GetRepoPretty(er)
    catalog := db.GetCatalog(er)
    curlpath := pretty["reposcheme"] + "://" + pretty["repouser"] + ":" + pretty["repopass"] + "@" + pretty["repohost"]
    for _, en := range catalog {
      dbtags := db.GetTags(er, en)
      for _, et := range dbtags {
        Reqt := curlpath + "/v2/" + en + "/manifests/" + et
        if body, ok := utils.MakeQueryToRepo(Reqt); ok {
          client := &http.Client{}
          Reqtv2Digest, _ := http.NewRequest("GET", Reqt, nil)
          Reqtv2Digest.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
          if Respv2Digest, err := client.Do(Reqtv2Digest); err != nil {
            say.L3(err.Error())
            say.L3("CheckManifests Daemon: cannot recieve response from registry, stopping work")
          } else {
            defer Respv2Digest.Body.Close()
            dbdigest := db.GetTagDigest(er, en, et)
            curldigest := Respv2Digest.Header.Get("Docker-Content-Digest")
            if (dbdigest != curldigest){
              var ch interface{}
              totalsize := 0
              fsshaarr := body.(map[string]interface{})["fsLayers"].([]interface{})
              historyarr := body.(map[string]interface{})["history"].([]interface{})
              db.DeleteTagSubBucket(er, en, et, "history")
              for i, _ := range fsshaarr {
                fssha := fsshaarr[i].(map[string]interface{})["blobSum"].(string)
                fssize := utils.GetfsLayerSize(curlpath + "/v2/" + en + "/blobs/" + fssha)
                history := historyarr[i].(map[string]interface{})["v1Compatibility"].(string)
                historynew := history
                if fsshanum, err := strconv.Atoi(fssize); err != nil {
                  say.L3(err.Error())
                } else {
                  if last := len(historynew) - 1; last >= 0 {
                      historynew = historynew[:last]
                  }
                  historynew = historynew + ",\"blobSum\":\"" +
                    fssha + "\", \"blobSize\":\"" +
                    utils.FromByteToHuman(fsshanum) + "\"}"
                  totalsize += fsshanum
                }
                if err := json.Unmarshal([]byte(history), &ch); err != nil {
                  say.L3(err.Error())
                } else {
                  created := ch.(map[string]interface{})["created"].(string)
                  var indx int
                  if indx = strings.Index(created, "T"); indx > -1 {
                    created = created[:indx] + " " + created[indx+1:]
                    if indx = strings.Index(created, "."); indx > -1 {
                      created = created[:indx]
                    }
                  }
                  if indx > -1 {
                    db.PutSimplePairToBucket([]string{ er, "catalog", en, et, "history" }, created, historynew)
                  }
                }
              }
              sizedt := time.Now().Local().Format("2006-01-02 15:04:05")
              db.PutSimplePairToBucket([]string{ er, "catalog", en, et, "_totalsizehuman" }, sizedt, utils.FromByteToHuman(totalsize))
              db.PutSimplePairToBucket([]string{ er, "catalog", en, et, "_totalsizebytes" }, sizedt, strconv.Itoa(totalsize))
              db.PutTagDigest(er, en, et, curldigest)
            } else {
              say.L1("CheckManifests Daemon: digests are the same, shouldnot update anything, stopping work")
            }
          }
        } else {
          say.L3("CheckManifests Daemon: cannot recieve response from registry, stopping work")
        }
      }
    }
  }
  say.L1("CheckManifests Daemon: finished work")
}
