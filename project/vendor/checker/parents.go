package checker

import(
  "dt"
  "db"
  "say"

  "time"
  "utils"
  "strings"
  "encoding/json"
)

func checkParents(runchannel chan int){
  defer dt.Watch(time.Now(), "Check Parents Demon")

  runchannel <- 1
  say.L1("CheckParents Daemon: started work")
  repos := db.GetRepos()
  for key, value := range repos {
    if value == "" {
      names := db.GetAllPairsFromBucket([]string{key, "catalog"})
      for keyn, valuen := range names {
        if valuen == "" {
          tags := db.GetAllPairsFromBucket([]string{key, "catalog", keyn})
          for keyt, valuet := range tags {
            if (valuet == "") && (keyt[0:1] != "_"){
              say.L1("CheckParents Daemon: check [" + key + "->" + keyn + ":" + keyt + "]")
              if _, ok := db.GetAllPairsFromBucket([]string{key, "catalog", keyn, keyt})["history"]; !ok {
                db.PutBucketToBucket([]string{key, "catalog", keyn, keyt, "history"})
              }
              history := db.GetAllPairsFromBucket([]string{key, "catalog", keyn, keyt, "history"})
              histarr := []string{}
              var tmpstr string
              cmd := db.GetAllPairsFromBucket([]string{key, "_names", keyn + ":" + keyt})
              for _, valh := range history {
                var ch interface{}
                if err := json.Unmarshal([]byte(valh), &ch); err != nil {
                  say.L3(err.Error())
                } else {
                  tmpstr = ""
                  if ch.(map[string]interface{})["container_config"].(map[string]interface{})["Cmd"] != nil {
                    for valji, valj := range ch.(map[string]interface{})["container_config"].(map[string]interface{})["Cmd"].([]interface{}) {
                      if strings.Contains(valj.(string), " CMD ") ||
                         strings.Contains(valj.(string), " WORKDIR ") ||
                         strings.Contains(valj.(string), " ENTRYPOINT ") ||
                         strings.Contains(valj.(string), " VOLUME ") ||
                         strings.Contains(valj.(string), " EXPOSE "){
                         tmpstr = ""
                         break
                      } else {
                        tmpstr += valj.(string)
                        if (valji < len(ch.(map[string]interface{})["container_config"].(map[string]interface{})["Cmd"].([]interface{}))-1) {
                          tmpstr += " "
                        }
                      }
                    }
                  }
                  if tmpstr != "" {
                    histarr = append(histarr, tmpstr)
                  }
                }
              }
              var cmdslice []string
              cmdneedaddition := true
              for _, valcmd := range cmd {
                if err := json.Unmarshal([]byte(valcmd), &cmdslice); err != nil {
                  say.L3(err.Error())
                } else {
                  if ! utils.IsSliceDifferent(histarr, cmdslice) {
                    cmdneedaddition = false
                    break
                  }
                }
              }
              if cmdneedaddition {
                sizedt := time.Now().Local().Format("2006-01-02 15:04:05")
                fullcmd, _ := json.Marshal(histarr)
                db.PutSimplePairToBucket([]string{ key, "_names", keyn + ":" + keyt }, sizedt, string(fullcmd))
              }
              say.L1("Finding parent for [ " + keyn + ":" + keyt +  " ]")
              if pn, pt, pok := FindParent(histarr, key, keyn, keyt); pok {
                db.PutSimplePairToBucket([]string{ key, "catalog", keyn, keyt, "_parent" }, "name", pn)
                db.PutSimplePairToBucket([]string{ key, "catalog", keyn, keyt, "_parent" }, "tag",  pt)
              } else {
                db.PutSimplePairToBucket([]string{ key, "catalog", keyn, keyt, "_parent" }, "name", "")
                db.PutSimplePairToBucket([]string{ key, "catalog", keyn, keyt, "_parent" }, "tag",  "")
              }
            }
          }
        }
      }
    }
    db.DeleteBucket([]string{key, "_namesgraph"})
    BuildParentsGraph(key)
  }
  say.L1("CheckParents Daemon: finished work")
  <- runchannel
}

func FindParent(childcmd []string, repo string, namei string, tagi string) (name string, tag string, ok bool){
  say.L1("Searching for parent of [ " + namei + ":" + tagi + " ]")
  ok = false
  names := db.GetAllPairsFromBucket([]string{repo, "_names"})
  maxname := ""
  maxlayers := 0
  for kn, _ := range names {
    if strings.Split(kn, ":")[0] != namei {
      // cmd := map[datetime_of_push]commands_in_manifest
      cmd := db.GetAllPairsFromBucket([]string{repo, "_names", kn})
      for _, vc := range cmd {
        var parentcmd interface{}
        if err := json.Unmarshal([]byte(vc), &parentcmd); err != nil {
          say.L3(err.Error())
          // return -- don't shure why it was return
          // if one of cmd broken, we can skeap only this exact comand
          break
        } else {
          initParentLen := len(parentcmd.([]interface{}))
          if (maxlayers <= initParentLen) {
            // already has better parent
            if initParentLen != 0 {
              for _, eccmd := range childcmd {
                for ipcmd, epcmd := range parentcmd.([]interface{}) {
                  if epcmd == eccmd {
                    parentcmd = append(parentcmd.([]interface{})[:ipcmd], parentcmd.([]interface{})[ipcmd+1:]...)
                    break
                  }
                }
              }
              if len(parentcmd.([]interface{})) == 0 {
                if (maxlayers < initParentLen) && (len(childcmd) != initParentLen) {
                                               // else it has same length and layers => it is exact copy
                  maxlayers = initParentLen
                  maxname = kn
                }
              }
            }
          }
        }
      }
    }
  }
  if maxlayers == 0 {
    say.L1("Parent not found")
  } else {
    say.L1("Parent is [ "+ maxname +" ]")
    ok = true
    s := strings.Split(maxname, ":")
    name = s[0]
    tag = s[1]
  }
  return
}

func BuildParentsGraph(repo string){
  say.L1("Building parents tree for [ " + repo + " ]")
  fullnames := []string{}
  names := db.GetAllPairsFromBucket([]string{repo, "catalog"})
  for kn, _ := range names {
    tags := db.GetAllPairsFromBucket([]string{repo, "catalog", kn})
    for kt, _ := range tags {
      if kt[0:1] != "_" {
        fullnames = append(fullnames, kn + ":" + kt)
      }
    }
  }

  Depth := 0
  Base := [][]string{}
  Base = append(Base, []string{repo, "_namesgraph"})
  db.PutBucketToBucket(Base[0])

  for (len(fullnames) > 0) && (Depth < 100) {
    tmpBase := [][]string{}
    for i := len(fullnames)-1; i > -1; i-- {
      s := strings.Split(fullnames[i], ":")
      n := s[0]
      t := s[1]
      np := db.GetValueFromBucket([]string{ repo, "catalog", n, t, "_parent" }, "name")
      tp := db.GetValueFromBucket([]string{ repo, "catalog", n, t, "_parent" }, "tag")

      for j := 0; j < len(Base); j++ {
        if ( np + ":" + tp == Base[j][len(Base[j])-1] ) || (( Depth == 0 ) && ( np + ":" + tp == ":" )) {
          say.L1("Found parents [ " + np + ":" + tp + " => " + n + ":" + t + " ]")
          tmpPath := append(Base[j], n + ":" + t)
          cpPath := make([]string, len(tmpPath))
          copy(cpPath, tmpPath)
          tmpBase = append(tmpBase, cpPath)
          db.PutBucketToBucket(tmpPath)
          fullnames = append(fullnames[:i], fullnames[i+1:]...)
        }
      }
    }
    Base = tmpBase
    Depth++
  }
}
