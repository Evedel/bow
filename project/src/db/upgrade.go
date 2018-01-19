package db

import(
  "utils"
  "strconv"
  "strings"

  "github.com/Evedel/glb/say"
)

func upto3(){
  say.L2("DB: INIT: DB Upgrade: Need upgrade.", "","\n")
  repos := GetRepos()
  for er, _ := range repos {
    names := GetAllPairsFromBucket([]string{er, "catalog"})
    for en, _ := range names {
      if en != "_valid" {
        PutBucketToBucket([]string{ er, "catalog", en, "_namepair"})
        idx := strings.Index(en, "/")
        if idx != -1 {
          PutSimplePairToBucket([]string{ er, "catalog", en, "_namepair"}, en[:idx], en[idx+1:])
        } else {
          PutSimplePairToBucket([]string{ er, "catalog", en, "_namepair"}, "_none", en)
        }
      }
    }
  }
  PutSimplePairToBucket([]string{"_info"}, "version", "3")
}

func upto2(){
  say.L2("DB: INIT: DB Upgrade: Need upgrade.", "","\n")
  repos := GetRepos()
  for er, _ := range repos{
    repofull := GetAllPairsFromBucket([]string{er})
    PutSimplePairToBucket([]string{er, "_info"}, "host", repofull["repohost"])
    PutSimplePairToBucket([]string{er, "_info"}, "pass", repofull["repopass"])
    PutSimplePairToBucket([]string{er, "_info"}, "user", repofull["repouser"])
    if scheme, ok := repofull["reposcheme"]; ok {
      PutSimplePairToBucket([]string{er, "_info"}, "scheme", scheme)
      DeleteKey([]string{er}, "reposcheme")
    } else {
      PutSimplePairToBucket([]string{er, "_info"}, "scheme", "http")
    }
    PutSimplePairToBucket([]string{er, "_info"}, "name", er)
    PutSimplePairToBucket([]string{er, "_info"}, "secure", "true")
    DeleteKey([]string{er}, "repohost")
    DeleteKey([]string{er}, "repopass")
    DeleteKey([]string{er}, "repouser")
  }
  PutSimplePairToBucket([]string{"_info"}, "version", "2")
}

func upto1(){
  say.L2("DB: INIT: DB Upgrade: Version: 0.", "","\n")
  say.L2("DB: INIT: DB Upgrade: Need upgrade.", "","\n")
  PutSimplePairToBucket([]string{"_info"}, "version", "1")
}

func UpgradeOldParentNames(){
  say.L2("DB UPGRADE: make upgrade for old parent names.", "","\n")
  repos := GetRepos()
  for key, value := range repos {
    if value == "" {
      names := GetAllPairsFromBucket([]string{key, "_names"})
      for keyn, valuen := range names {
        if valuen != "" {
          DeleteKey([]string{key, "_names" }, keyn)
        }
      }
    }
  }
}

func UpgradeFalseNumericImage(){
  say.L2("DB UPGRADE: make upgrade for false numeric image.", "","\n")
  repos := GetRepos()
  for key, value := range repos {
    if value == "" {
      imagenames := GetAllPairsFromBucket([]string{key, "catalog"})
      for keyi, _ := range imagenames {
        if _, err := strconv.Atoi(keyi); err == nil {
          DeleteBucket([]string{key, "catalog", keyi})
        }
      }
    }
  }
}

func UpgradeTotalSize(){
  say.L2("DB UPGRADE: make upgrade for total size.", "","\n")
  repos := GetRepos()
  for key, value := range repos {
    if value == "" {
      imagenames := GetAllPairsFromBucket([]string{key, "catalog"})
      for keyi, _ := range imagenames {
        tags := GetAllPairsFromBucket([]string{key, "catalog", keyi})
        for keyt, _ := range tags {
          if (keyt != "_uploads") && (keyt != "_valid") {
            fields := GetAllPairsFromBucket([]string{key, "catalog", keyi, keyt})
            if _, ok := fields["_totalsize"]; ok {
              totalsize := GetAllPairsFromBucket([]string{key, "catalog", keyi, keyt, "_totalsize"})
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
