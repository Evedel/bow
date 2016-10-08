package checker

import(
  "db"
  "say"
  "utils"
)

func CheckRepos(){
  say.L1("CheckRepos Daemon: started work")
  repos := db.GetRepos()
  for _, e := range repos {
    pretty := db.GetRepoPretty(e)
    repopref := pretty["reposcheme"] + "://" + pretty["repouser"] + ":" + pretty["repopass"] + "@" + pretty["repohost"]
    Req := repopref + "/v2/_catalog?n=&last="
    if body, ok := utils.MakeQueryToRepo(Req); ok {
      dbcatalog := db.GetCatalog(e)
      arrint := body.(map[string]interface{})["repositories"].([]interface{})
      arrstr := make([]string, 0)
      for _, e := range arrint {
        Reqtag := repopref + "/v2/" + e.(string) + "/tags/list"
        if body, ok := utils.MakeQueryToRepo(Reqtag); ok {
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
