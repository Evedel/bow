package dt

import(
  "say"
  "conf"

  "time"
)


func Watch(start time.Time, name string) {
  if conf.Env["timewatch"] == "yes" {
    elapsed := time.Since(start)
    say.L2("TimeWatch: [ " + name + " ] - " + elapsed.String())
  }
}
