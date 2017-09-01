package checker

import(
  "dt"
  "db"
  "say"
  "qurl"
  "utils"

  "time"
  "strings"
)

func checkRepos(runchannel chan int){
  defer dt.Watch(time.Now(), "Check Repositories Demon")

  runchannel <- 1
  say.L1("CheckRepos Daemon: started work")
  repos := db.GetRepos()
  for e, _ := range repos {
    repoinfo := db.GetRepoPretty(e)
    Req := "/v2/_catalog?n=&last="

    arrstr := make([]string, 0)
    repeat := true
    for repeat {
      repeat = false
      if body, headers, ok := qurl.MakeQuery(Req, "GET", repoinfo, map[string]string{}); ok {
        arrint := body.(map[string]interface{})["repositories"].([]interface{})
        for _, e := range arrint {
          // after tag was deleted by 'delete' button in UI
          // if it was the last tag in namespace/name
          // the slice will be empty, and manifest will return 404
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
        if link, okh := headers["Link"]; okh {
          // until
          // Link:[</v2/_catalog?last=0099-evedel%2Fbow&n=100>; rel="next"]
          //        |                  Req            |
          from  := strings.Index(link[0], "<")
          to    := strings.Index(link[0], "&")
          if ((from != -1) && (to != -1)) {
            Req   = link[0][from+1 : to]
            repeat = true
          }
        }
      } else {
        say.L3("CheckRepos Daemon: cannot recieve response from registry, stopping work")
      }
    }

    dbcatalog := db.GetCatalog(e)
    if utils.IsSliceDifferent(dbcatalog, arrstr) {
      db.AddCatalog(e, arrstr)
    }

  }
  say.L1("CheckRepos Daemon: finished work")
  <-runchannel
}
