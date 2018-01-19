package checker

import (
  "conf"
  "time"
  "strconv"

  "github.com/Evedel/glb/say"
)

var runchannel = make( map[string]chan int)


func DaemonManager() {
  t, _ := strconv.Atoi(conf.Env["checker_time"])
  runchannel["repo"] = make(chan int,1)
  runchannel["tags"] = make(chan int,1)
  runchannel["mnft"] = make(chan int,1)
  runchannel["prnt"] = make(chan int,1)

  say.L2("DaemonManager: Sleep time is : ", conf.Env["checker_time"], " seconds.\n")
  for {
    say.L3("DaemonManager: TicTac.", "","\n")
    go checkRepos(runchannel["repo"])
    go checkTags(runchannel["tags"])
    go checkManifests(runchannel["mnft"])
    go checkParents(runchannel["prnt"])
    time.Sleep(time.Duration(t) * time.Second)
  }
}

func StartManual() {
  say.L3("DaemonManager: Started all checkers manually.", "","\n")
  go checkRepos(runchannel["repo"])
  go checkTags(runchannel["tags"])
  go checkManifests(runchannel["mnft"])
  go checkParents(runchannel["prnt"])
}

func RunCheckTags(){
  say.L3("DaemonManager: Run check tags only.", "","\n")
  go checkTags(runchannel["tags"])
}
