package checker

import(
  "dt"
  "db"
  "qurl"
  "utils"

  "time"

  "github.com/Evedel/glb/say"
)

func checkTags(runchannel chan int){
  defer dt.Watch(time.Now(), "Check Tags Demon")

  runchannel <- 1
  say.L3("CheckTags Daemon: started work", "","\n")
  repos := db.GetRepos()
  for er, _ := range repos {
    repoinfo := db.GetRepoPretty(er)
    catalog := db.GetCatalog(er)
    for _, en := range catalog {
      Reqt := "/v2/" + en + "/tags/list"
      if body, _, ok := qurl.MakeQuery(Reqt, "GET", repoinfo, map[string]string{}); ok {
        dbtags := db.GetTags(er, en)
        arrint := make([]interface{}, 0)
        if body.(map[string]interface{})["tags"] == nil {
          db.PutSimplePairToBucket([]string{ er, "catalog", en }, "_valid", "0")
        } else {
          db.PutSimplePairToBucket([]string{ er, "catalog", en }, "_valid", "1")
          arrint = body.(map[string]interface{})["tags"].([]interface{})
        }
        arrstr := make([]string, len(arrint))
        for i, _ := range arrint {
          arrstr[i] = arrint[i].(string)
        }
        if utils.IsSliceDifferent(dbtags, arrstr) {
          db.AddTags(er, en, arrstr)
        }
      } else {
        if body != nil {
          if body.(int) == 404 {
            say.L2("CheckTags Daemon: Page with name [" + en + "] not found. Asuming it isn't valid in the moment", "","\n")
            db.PutSimplePairToBucket([]string{ er, "catalog", en }, "_valid", "0")
          } else {
            say.L1("CheckTags Daemon: cannot recieve response from registry, stopping work", "","\n")
          }
        }
      }
    }
  }
  say.L3("CheckTags Daemon: finished work", "","\n")
  <- runchannel
}
