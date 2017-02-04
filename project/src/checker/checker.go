package checker

import (
  "say"
  "conf"
  "time"
  "strconv"
)

var runchannel = make( map[string]chan int)


func DaemonManager() {
  t, _ := strconv.Atoi(conf.Env["checker_time"])
  runchannel["repo"] = make(chan int,1)
  runchannel["tags"] = make(chan int,1)
  runchannel["mnft"] = make(chan int,1)
  runchannel["prnt"] = make(chan int,1)

  say.L2("DaemonManager: Sleep time is : " + conf.Env["checker_time"] + " seconds")
  for {
    say.L1("DaemonManager: TicTac")
    go checkRepos(runchannel["repo"])
    go checkTags(runchannel["tags"])
    go checkManifests(runchannel["mnft"])
    go checkParents(runchannel["prnt"])
    time.Sleep(time.Duration(t) * time.Second)
  }
}

func StartManual() {
  say.L2("DaemonManager: Started all checkers manually")
  go checkRepos(runchannel["repo"])
  go checkTags(runchannel["tags"])
  go checkManifests(runchannel["mnft"])
  go checkParents(runchannel["prnt"])
}

func RunCheckTags(){
  go checkTags(runchannel["tags"])
}
