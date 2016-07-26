package say

import(
  "conf"

  "github.com/fatih/color"

  "time"
  "fmt"
)
var mode string

func L1(str string){
  cyan := color.New(color.FgCyan).SprintFunc()
  if conf.Env["log_silent"] != "yes" && conf.Env["log_silent"] != "super" {
    fmt.Printf("[ " + time.Now().Format(time.RFC1123) + " ]" + cyan("[ L1 ] ") + str + "\n")
  }
}

func L2(str string){
  yellow := color.New(color.FgYellow).SprintFunc()
  if conf.Env["log_silent"] != "super" {
    fmt.Printf("[ " + time.Now().Format(time.RFC1123) + " ]" + yellow("[ L2 ] ") + str + "\n")
  }
}

func L3(str string){
  if conf.Env["log_silent"] != "super" {
    finstr := "[ " + time.Now().Format(time.RFC1123) + " ][ L3 ] " + str
    color.Red(finstr)
  }
}

func L4(str interface{}){
  fmt.Println(str)
}
