package utils

import(
  "db"
  "say"
  "strconv"
)

func UpgradeTotalSize(){
  say.L2("DB UPGRADE: make upgrade for total size")
  repos := db.GetSimplePairsFromBucket([]string{})
  for key, value := range repos {
    if value == "" {
      imagenames := db.GetSimplePairsFromBucket([]string{key, "catalog"})
      for keyi, _ := range imagenames {
        tags := db.GetSimplePairsFromBucket([]string{key, "catalog", keyi})
        for keyt, _ := range tags {
          if (keyt != "_uploads") && (keyt != "_valid") {
            fields := db.GetSimplePairsFromBucket([]string{key, "catalog", keyi, keyt})
            if _, ok := fields["_totalsize"]; ok {
              totalsize := db.GetSimplePairsFromBucket([]string{key, "catalog", keyi, keyt, "_totalsize"})
              for keys, vals := range totalsize {
                lastchar := vals[len(vals)-1:len(vals)]
                if lastchar == "B"{
                  db.PutSimplePairToBucket([]string{key, "catalog", keyi, keyt, "_totalsizehuman"}, keys, vals)
                  db.PutSimplePairToBucket([]string{key, "catalog", keyi, keyt, "_totalsizebytes"}, keys,
                    strconv.Itoa(FromHumanToByte(vals)))
                } else {
                  db.PutSimplePairToBucket([]string{key, "catalog", keyi, keyt, "_totalsizebytes"}, keys, vals)
                  num, _ := strconv.Atoi(vals)
                  db.PutSimplePairToBucket([]string{key, "catalog", keyi, keyt, "_totalsizehuman"}, keys, FromByteToHuman(num))
                }
              }
              db.DeleteBucket([]string{key, "catalog", keyi, keyt, "_totalsize"})
            }
          }
        }
      }
    }
  }
}
func UpgradeFalseNumericImage(){
  say.L2("DB UPGRADE: make upgrade for false numeric image")
  repos := db.GetSimplePairsFromBucket([]string{})
  for key, value := range repos {
    if value == "" {
      imagenames := db.GetSimplePairsFromBucket([]string{key, "catalog"})
      for keyi, _ := range imagenames {
        if _, err := strconv.Atoi(keyi); err == nil {
          db.DeleteBucket([]string{key, "catalog", keyi})
        }
      }
    }
  }
}
func UpgradeOldParentNames(){
  say.L2("DB UPGRADE: make upgrade for old parent names")
  repos := db.GetSimplePairsFromBucket([]string{})
  for key, value := range repos {
    if value == "" {
      names := db.GetSimplePairsFromBucket([]string{key, "_names"})
      for keyn, valuen := range names {
        if valuen != "" {
          db.DeleteKeyFromDB([]string{key, "_names" }, keyn)
        }
      }
    }
  }
}
