package checker

import(
  "db"
  "say"
  "qurl"
  "utils"
)

func CheckRepos(){
  say.L1("CheckRepos Daemon: started work")
  repos := db.GetRepos()
  for e, _ := range repos {
    repoinfo := db.GetRepoPretty(e)
    Req := "/v2/_catalog?n=&last="
    if body, _, ok := qurl.MakeQuery(Req, "GET", repoinfo, map[string]string{}); ok {
      dbcatalog := db.GetCatalog(e)
      arrint := body.(map[string]interface{})["repositories"].([]interface{})
      arrstr := make([]string, 0)
      for _, e := range arrint {
        Reqtag := "/v2/" + e.(string) + "/tags/list"
        if body, _, ok := qurl.MakeQuery(Reqtag, "GET", repoinfo, map[string]string{}); ok {
          if body.(map[string]interface{})["tags"] != nil {
            arrstr = append(arrstr, e.(string))
          }
        } else {
          say.L3("CheckRepos Daemon: cannot recieve response from registry, stopping work")
          break;
        }
      }
      if utils.IsSliceDifferent(dbcatalog, arrstr) {
        db.AddCatalog(e, arrstr)
      }
    } else {
      say.L3("CheckRepos Daemon: cannot recieve response from registry, stopping work")
    }
  }
  say.L1("CheckRepos Daemon: finished work")
}
