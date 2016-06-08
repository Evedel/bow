package say

import(
  "conf"

  "github.com/fatih/color"

  "time"
  "fmt"
)
var mode string

func Info(str string){
  cyan := color.New(color.FgCyan).SprintFunc()
  if conf.Env["log_silent"] != "yes" && conf.Env["log_silent"] != "super" {
    fmt.Printf("[ " + time.Now().Format(time.RFC1123) + " ]" + cyan("[   info  ] ") + str + "\n")
  }
}

func Warn(str string){
  yellow := color.New(color.FgYellow).SprintFunc()
  if conf.Env["log_silent"] != "yes" && conf.Env["log_silent"] != "super" {
    fmt.Printf("[ " + time.Now().Format(time.RFC1123) + " ]" + yellow("[ warning ] ") + str + "\n")
  }
}

func Error(str string){
  if conf.Env["log_silent"] != "super" {
    finstr := "[ " + time.Now().Format(time.RFC1123) + " ][  ERROR  ] " + str
    color.Red(finstr)
  }
}

func Raw(str interface{}){
  fmt.Println(str)
}
