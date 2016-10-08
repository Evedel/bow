package checker

import (
  "say"
  "conf"
  "time"
  "strconv"
)

func DaemonManager() {
  t, _ := strconv.Atoi(conf.Env["checker_time"])
  say.L2("DaemonManager: Sleep time is : " + conf.Env["checker_time"] + " seconds")
  for {
    say.L1("DaemonManager: TicTac")
    go CheckRepos()
    go CheckTags()
    go CheckManifests()
    go CheckParents()
    time.Sleep(time.Duration(t) * time.Second)
  }
}
