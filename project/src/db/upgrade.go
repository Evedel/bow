package db

import(
  "say"
  "utils"
  "strconv"
)

func upto2(){
  say.L2("DB: INIT: DB Upgrade: Need upgrade")
  repos := GetRepos()
  for er, _ := range repos{
    repofull := GetSimplePairsFromBucket([]string{er})
    PutSimplePairToBucket([]string{er, "_info"}, "host", repofull["repohost"])
    PutSimplePairToBucket([]string{er, "_info"}, "pass", repofull["repopass"])
    PutSimplePairToBucket([]string{er, "_info"}, "user", repofull["repouser"])
    PutSimplePairToBucket([]string{er, "_info"}, "scheme", repofull["reposcheme"])
    PutSimplePairToBucket([]string{er, "_info"}, "name", er)
    PutSimplePairToBucket([]string{er, "_info"}, "secure", "true")
    DeleteKey([]string{er}, "repohost")
    DeleteKey([]string{er}, "repopass")
    DeleteKey([]string{er}, "repouser")
    DeleteKey([]string{er}, "reposcheme")
  }
  PutSimplePairToBucket([]string{"_info"}, "version", "2")
}

func upto1(){
  say.L2("DB: INIT: DB Upgrade: Version: 0")
  say.L2("DB: INIT: DB Upgrade: Need upgrade")
  PutSimplePairToBucket([]string{"_info"}, "version", "1")
}

func UpgradeOldParentNames(){
  say.L2("DB UPGRADE: make upgrade for old parent names")
  repos := GetRepos()
  for key, value := range repos {
    if value == "" {
      names := GetSimplePairsFromBucket([]string{key, "_names"})
      for keyn, valuen := range names {
        if valuen != "" {
          DeleteKey([]string{key, "_names" }, keyn)
        }
      }
    }
  }
}

func UpgradeFalseNumericImage(){
  say.L2("DB UPGRADE: make upgrade for false numeric image")
  repos := GetRepos()
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

func UpgradeTotalSize(){
  say.L2("DB UPGRADE: make upgrade for total size")
  repos := GetRepos()
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
                    strconv.Itoa(utils.FromHumanToByte(vals)))
                } else {
                  PutSimplePairToBucket([]string{key, "catalog", keyi, keyt, "_totalsizebytes"}, keys, vals)
                  num, _ := strconv.Atoi(vals)
                  PutSimplePairToBucket([]string{key, "catalog", keyi, keyt, "_totalsizehuman"}, keys,
                    utils.FromByteToHuman(num))
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
