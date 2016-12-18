package checker

import(
  "db"
  "say"
  "qurl"
  "time"
  "utils"
  "strings"
  "strconv"
  "encoding/json"
)

func CheckManifests(){
  say.L1("CheckManifests Daemon: started work")
  repos := db.GetRepos()
  for er, _ := range repos {
    repoinfo := db.GetRepoPretty(er)
    catalog := db.GetCatalog(er)
    for _, en := range catalog {
      dbtags := db.GetTags(er, en)
      for _, et := range dbtags {
        Reqt := "/v2/" + en + "/manifests/" + et
        if body, _, ok := qurl.MakeQuery(Reqt, "GET", repoinfo, map[string]string{}); ok {
          ihdr := map[string]string{"Accept": "application/vnd.docker.distribution.manifest.v2+json"}
          if _, ohdr, ok := qurl.MakeQuery(Reqt, "GET", repoinfo, ihdr); !ok {
            say.L3("CheckManifests Daemon: cannot recieve digest header from registry, stopping work")
            return
          } else {
            olddidg := db.GetTagDigest(er, en, et)
            newdidg := ohdr["Docker-Content-Digest"][0]
            say.L3(olddidg)
            say.L3(newdidg)

            if (olddidg != newdidg){
              var ch interface{}
              totalsize := 0
              fsshaarr := body.(map[string]interface{})["fsLayers"].([]interface{})
              historyarr := body.(map[string]interface{})["history"].([]interface{})
              db.DeleteTagSubBucket(er, en, et, "history")
              for i, _ := range fsshaarr {
                fssha := fsshaarr[i].(map[string]interface{})["blobSum"].(string)
                var fssize string
                if _, fsshdr, okcl := qurl.MakeQuery("/v2/" + en + "/blobs/" + fssha, "HEAD", repoinfo, map[string]string{}); !okcl {
                  say.L3("CheckManifests Daemon: cannot recieve content length header from registry, stopping work")
                  fssize = "0"
                } else {
                  fssize = fsshdr["Content-Length"][0]
                }
                history := historyarr[i].(map[string]interface{})["v1Compatibility"].(string)
                historynew := history
                if fsshanum, err := strconv.Atoi(fssize); err != nil {
                  say.L3(err.Error())
                } else {
                  // Cut the carriage return
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
              db.PutTagDigest(er, en, et, newdidg)
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
