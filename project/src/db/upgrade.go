package db

import(
  "say"
  "strings"
  "strconv"
)

func UpgradeTotalSize(){
  say.L2("DB UPGRADE: make upgrade for total size")
  repos := GetSimplePairsFromBucket([]string{})
  for key, value := range repos {
    if value == "" {
      imagenames := GetSimplePairsFromBucket([]string{key, "catalog"})
      for keyi, _ := range imagenames {
        tags := GetSimplePairsFromBucket([]string{key, "catalog", keyi})
        for keyt, _ := range tags {
          if (keyt != "_uploads") && (keyt != "_valid") {
            fields := GetSimplePairsFromBucket([]string{key, "catalog", keyi, keyt})
            if _, ok := fields["_totalsize"]; ok {
              totalsize := GetSimplePairsFromBucket([]string{key, "catalog", keyi, keyt, "_totalsize"})
              for keys, vals := range totalsize {
                lastchar := vals[len(vals)-1:len(vals)]
                if lastchar == "B"{
                  PutSimplePairToBucket([]string{key, "catalog", keyi, keyt, "_totalsizehuman"}, keys, vals)
                  PutSimplePairToBucket([]string{key, "catalog", keyi, keyt, "_totalsizebytes"}, keys,
                    strconv.Itoa(fromHumanToByte(vals)))
                } else {
                  PutSimplePairToBucket([]string{key, "catalog", keyi, keyt, "_totalsizebytes"}, keys, vals)
                  num, _ := strconv.Atoi(vals)
                  PutSimplePairToBucket([]string{key, "catalog", keyi, keyt, "_totalsizehuman"}, keys, fromByteToHuman(num))
                }
              }
              DeleteBucket([]string{key, "catalog", keyi, keyt, "_totalsize"})
            }
          }
        }
      }
    }
  }
}
func UpgradeFalseNumericImage(){
  say.L2("DB UPGRADE: make upgrade for false numeric image")
  repos := GetSimplePairsFromBucket([]string{})
  for key, value := range repos {
    if value == "" {
      imagenames := GetSimplePairsFromBucket([]string{key, "catalog"})
      for keyi, _ := range imagenames {
        if _, err := strconv.Atoi(keyi); err == nil {
          DeleteBucket([]string{key, "catalog", keyi})
        }
      }
    }
  }
}
func UpgradeOldParentNames(){
  say.L2("DB UPGRADE: make upgrade for old parent names")
  repos := GetSimplePairsFromBucket([]string{})
  for key, value := range repos {
    if value == "" {
      names := GetSimplePairsFromBucket([]string{key, "_names"})
      for keyn, valuen := range names {
        if valuen != "" {
          DeleteKeyFromDB([]string{key, "_names" }, keyn)
        }
      }
    }
  }
}
func fromByteToHuman(bytes int) (human string){
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
func fromHumanToByte(human string) (bytes int){
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
