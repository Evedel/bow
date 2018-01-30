package checker

import(
  "dt"
  "db"
  "say"
  "qurl"
  "utils"

  "time"
)

func checkTags(runchannel chan int){
  defer dt.Watch(time.Now(), "Check Tags Demon")

  runchannel <- 1
  say.L1("CheckTags Daemon: started work")
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
            say.L2("CheckTags Daemon: Page with name [" + en + "] not found. Asuming it isn't valid in the moment")
            db.PutSimplePairToBucket([]string{ er, "catalog", en }, "_valid", "0")
          } else {
            say.L3("CheckTags Daemon: cannot recieve response from registry, stopping work")
          }
        }
      }
    }
  }
  say.L1("CheckTags Daemon: finished work")
  <- runchannel
}
