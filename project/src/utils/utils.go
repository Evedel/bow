package utils

import(
  "say"
  "strings"
  "strconv"
)

func IsSliceDifferent(a []string, b []string) (bool) {
  al := len(a)
  bl := len(b)
  if a == nil && b == nil {
    say.L1("Slices are equally nill. Same.")
    return false
  }
  if a == nil || b == nil {
    say.L1("One of the slices is empty. Different.")
    return true
  }
  if al != bl {
    say.L1("Length of slices are different. Different.")
    return true
  }
  numofequal := 0
  for _, bel := range b {
    for _, ael := range a {
      if bel == ael{
        numofequal++
        break
      }
    }
  }
  if len(a) == numofequal {
    say.L1("Length of slices are same with number of equal elements. Same.")
    return false
  } else {
  say.L1("Length of slices are differ with number of equal elements. Different.")
    return true
  }
}

func FromByteToHuman(bytes int) (human string){
  human = strconv.Itoa(bytes) + " B"
  if bytes > 1024 {
    bytes = bytes / 1024
    human = strconv.Itoa(bytes) + " KB"
  }
  if bytes > 1024 {
    bytes = bytes / 1024
    human = strconv.Itoa(bytes) + " MB"
  }
  if bytes > 1024 {
    bytes = bytes / 1024
    human = strconv.Itoa(bytes) + " GB"
  }
  return
}
func FromHumanToByte(human string) (bytes int){
  space := strings.Index(human, " ")
  bytes = 0
  number := ""
  scale := ""
  if space != -1 {
    number = human[:space]
    scale = human[space+1:]
  }
  fscale := 0.0
  switch scale {
  case "B":
    bytes, _ = strconv.Atoi(number)
    return
  case "KB":
    fscale = 1024
  case "MB":
    fscale = 1024 * 1024
  case "GB":
    fscale = 1024 * 1024 * 1024
  }
  fnum, _ := strconv.ParseFloat(number, 64)
  bytes = int(fnum * fscale)
  return
}
