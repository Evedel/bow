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
  var num float64
  num = float64(bytes)
  human = strconv.FormatFloat(num, 'f', 2, 64) + " B"
  human = strings.Replace(human, ".00 B", " B", 1)
  if num > 1024 {
    num = num / 1024
    human = strconv.FormatFloat(num, 'f', 2, 64) + " KB"
    human = strings.Replace(human, ".00 KB", " KB", 1)
  }
  if num > 1024 {
    num = num / 1024
    human = strconv.FormatFloat(num, 'f', 2, 64) + " MB"
    human = strings.Replace(human, ".00 MB", " MB", 1)
  }
  if num > 1024 {
    num = num / 1024
    human = strconv.FormatFloat(num, 'f', 2, 64) + " GB"
    human = strings.Replace(human, ".00 MB", " MB", 1)
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

func Keys(inmap map[string]string) (outkeys []string){
  outkeys = make([]string, 0, len(inmap))
  for k := range inmap {
      outkeys = append(outkeys, k)
  }
  return
}
