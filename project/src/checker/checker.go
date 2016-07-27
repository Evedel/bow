package checker

import (
  "db"
  "say"
  "conf"
  "time"
  "strconv"
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

func DaemonManager() {
  t, _ := strconv.Atoi(conf.Env["checker_time"])
  say.L2("DaemonManager: Sleep time is : " + conf.Env["checker_time"] + " seconds")
  for {
    say.L1("DaemonManager: TicTac")
    go CheckRepos()
    go CheckTags()
    go CheckManifests()
    go CheckParents()
    time.Sleep(time.Duration(t) * time.Second)
  }
}

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

func CheckRepos(){
  say.L1("CheckRepos Daemon: started work")
  repos := db.GetRepos()
  for _, e := range repos {
    pretty := db.GetRepoPretty(e)
    Req := pretty["reposcheme"] + "://" + pretty["repouser"] +
      ":" + pretty["repopass"] + "@" + pretty["repohost"] + "/v2/_catalog?n=&last="
    if body, ok := MakeQueryToRepo(Req); ok {
      dbcatalog := db.GetCatalog(e)
      arrint := body.(map[string]interface{})["repositories"].([]interface{})
      arrstr := make([]string, len(arrint))
      for i, _ := range arrint {
        arrstr[i] = arrint[i].(string)
      }
      if IsSliceDifferent(dbcatalog, arrstr) {
        db.AddCatalog(e, arrstr)
      }
    } else {
      say.L3("CheckRepos Daemon: cannot recieve response from registry, stopping work")
    }
  }
  say.L1("CheckRepos Daemon: finished work")
}

func CheckTags(){
  say.L1("CheckTags Daemon: started work")
  repos := db.GetRepos()
  for _, er := range repos {
    pretty := db.GetRepoPretty(er)
    catalog := db.GetCatalog(er)
    reponame := pretty["reposcheme"] + "://" + pretty["repouser"] + ":" + pretty["repopass"] + "@" + pretty["repohost"]
    for _, en := range catalog {
      Reqt := reponame + "/v2/" + en + "/tags/list"
      if body, ok := MakeQueryToRepo(Reqt); ok {
        dbtags := db.GetTags(er, en)
        arrint := body.(map[string]interface{})["tags"].([]interface{})
        arrstr := make([]string, len(arrint))
        for i, _ := range arrint {
          arrstr[i] = arrint[i].(string)
        }
        if IsSliceDifferent(dbtags, arrstr) {
          db.AddTags(er, en, arrstr)
        }
      } else {
        say.L3("CheckTags Daemon: cannot recieve response from registry, stopping work")
      }
    }
  }
  say.L1("CheckTags Daemon: finished work")
}

func CheckManifests(){
  say.L1("CheckManifests Daemon: started work")
  repos := db.GetRepos()
  for _, er := range repos {
    pretty := db.GetRepoPretty(er)
    catalog := db.GetCatalog(er)
    curlpath := pretty["reposcheme"] + "://" + pretty["repouser"] + ":" + pretty["repopass"] + "@" + pretty["repohost"]
    for _, en := range catalog {
      dbtags := db.GetTags(er, en)
      for _, et := range dbtags {
        Reqt := curlpath + "/v2/" + en + "/manifests/" + et
        if body, ok := MakeQueryToRepo(Reqt); ok {
          client := &http.Client{}
          Reqtv2Digest, _ := http.NewRequest("GET", Reqt, nil)
          Reqtv2Digest.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
          if Respv2Digest, err := client.Do(Reqtv2Digest); err != nil {
            say.L3(err.Error())
            say.L3("CheckManifests Daemon: cannot recieve response from registry, stopping work")
          } else {
            defer Respv2Digest.Body.Close()
            dbdigest := db.GetTagDigest(er, en, et)
            curldigest := Respv2Digest.Header.Get("Docker-Content-Digest")
            if (dbdigest != curldigest){
              var ch interface{}
              totalsize := 0
              fsshaarr := body.(map[string]interface{})["fsLayers"].([]interface{})
              historyarr := body.(map[string]interface{})["history"].([]interface{})
              db.DeleteTagSubBucket(er, en, et, "history")
              for i, _ := range fsshaarr {
                fssha := fsshaarr[i].(map[string]interface{})["blobSum"].(string)
                fssize := GetfsLayerSize(curlpath + "/v2/" + en + "/blobs/" + fssha)
                history := historyarr[i].(map[string]interface{})["v1Compatibility"].(string)
                historynew := history
                if fsshanum, err := strconv.Atoi(fssize); err != nil {
                  say.L3(err.Error())
                } else {
                  if last := len(historynew) - 1; last >= 0 {
                      historynew = historynew[:last]
                  }
                  historynew = historynew + ",\"blobSum\":\"" + fssha + "\", \"blobSize\":\"" + fromByteToHuman(fsshanum) + "\"}"
                  totalsize += fsshanum
                }
                if err := json.Unmarshal([]byte(history), &ch); err != nil {
                  say.L3(err.Error())
                } else {
                  created := ch.(map[string]interface{})["created"].(string)
                  created = created[0:10] + " " + created[11:len(created)-11]
                  db.PutSimplePairToBucket([]string{ er, "catalog", en, et, "history" }, created, historynew)
                }
              }
              sizedt := time.Now().Local().Format("2006-01-02 15:04:05")
              db.PutSimplePairToBucket([]string{ er, "catalog", en, et, "_totalsizehuman" }, sizedt, fromByteToHuman(totalsize))
              db.PutSimplePairToBucket([]string{ er, "catalog", en, et, "_totalsizebytes" }, sizedt, strconv.Itoa(totalsize))
              db.PutTagDigest(er, en, et, curldigest)
            } else {
              say.L1("CheckManifests Daemon: digests are the same, shouldnot update anything, stopping work")
            }
          }
        } else {
          say.L3("CheckManifests Daemon: cannot recieve response from registry, stopping work")
        }
      }
    }
  }
  say.L1("CheckManifests Daemon: finished work")
}

func CheckParents(){
  repos := db.GetSimplePairsFromBucket([]string{})
  for key, value := range repos {
    if value == "" {
      names := db.GetSimplePairsFromBucket([]string{key, "catalog"})
      for keyn, valuen := range names {
        if valuen == "" {
          tags := db.GetSimplePairsFromBucket([]string{key, "catalog", keyn})
          for keyt, valuet := range tags {
            if (valuet == "") && (keyt[0:1] != "_"){
              history := db.GetSimplePairsFromBucket([]string{key, "catalog", keyn, keyt, "history"})
              histarr := []string{}
              var tmpstr string
              cmd := db.GetSimplePairsFromBucket([]string{key, "_names", keyn + ":" + keyt})
              for _, valh := range history {
                var ch interface{}
                if err := json.Unmarshal([]byte(valh), &ch); err != nil {
                  say.L3(err.Error())
                } else {
                  tmpstr = ""
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
                  if ! IsSliceDifferent(histarr, cmdslice) {
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
}

func BuildParentsGraph(repo string){
  say.L1("Building parents tree for [ " + repo + " ]")
  fullnames := []string{}
  names := db.GetSimplePairsFromBucket([]string{repo, "catalog"})
  for kn, _ := range names {
    tags := db.GetSimplePairsFromBucket([]string{repo, "catalog", kn})
    for kt, _ := range tags {
      if kt[0:1] != "_" {
        fullnames = append(fullnames, kn + ":" + kt)
      }
    }
  }

  Depth := 0
  Base := [][]string{}
  L0 := []string{repo, "_namesgraph"}
  Base = append(Base, L0)
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

func GetfsLayerSize(link string ) (size string){
  if Resp, err := http.Head(link); err != nil {
    say.L3(err.Error())
    say.L3("CheckManifests Daemon: GetfsLayerSize cannot recieve response from registry, stopping work")
  } else {
    defer Resp.Body.Close()
    if _, err := ioutil.ReadAll(Resp.Body); err != nil {
      say.L3(err.Error())
    } else {
      size = Resp.Header.Get("Content-Length")
      return
    }
  }
  return ""
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

func DeleteTagFromRepo(repo string, name string, tag string) (ok bool){
  ok = false
  pretty := db.GetRepoPretty(repo)
  curlpath := pretty["reposcheme"] + "://" + pretty["repouser"] + ":" + pretty["repopass"] + "@" + pretty["repohost"]
  ReqtStr := curlpath + "/v2/" + name + "/manifests/" + db.GetValueFromBucket([]string{repo, "catalog", name, tag}, "digest")
  client := &http.Client{}
  Reqt, _ := http.NewRequest("DELETE", ReqtStr, nil)
  Reqt.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
  if Resp, err := client.Do(Reqt); err != nil {
    say.L3(err.Error())
    say.L3("Delete From Repository: cannot recieve response from registry, stopping work")
    return
  } else {
    defer Resp.Body.Close()
    if Resp.StatusCode == 202 {
      ok = true
    } else {
      say.L3(ReqtStr)
      say.L3(Resp.Status)
    }
  }
  return
}

func MakeQueryToRepo(query string) (body interface{}, ok bool){
  ok = false
  if response, err := http.Get(query); err != nil {
    say.L3(err.Error())
    return
  } else {
    defer response.Body.Close()
    if bodytmp, err := ioutil.ReadAll(response.Body); err != nil {
      say.L3(err.Error())
      return
    } else {
      var c interface{}
      if err := json.Unmarshal(bodytmp, &c); err != nil {
        say.L3(err.Error())
        return
      } else {
        if c.(map[string]interface{})["errors"] != nil {
          say.L3(query)
          say.L3(c.(map[string]interface{})["errors"].([]interface{})[0].(map[string]interface{})["message"].(string))
          return
        } else {
          body = c
          ok = true
        }
      }
    }
  }
  return
}

func FindParent(childcmd []string, repo string, namei string, tagi string) (name string, tag string, ok bool){
  say.L1("Searching for parent of [ " + namei + ":" + tagi + " ]")
  ok = true
  names := db.GetSimplePairsFromBucket([]string{repo, "_names"})
  maxname := ""
  maxlayers := 0
  for kn, _ := range names {
    if strings.Split(kn, ":")[0] != namei {
      cmd := db.GetSimplePairsFromBucket([]string{repo, "_names", kn})
      for _, vc := range cmd {
        var parentcmd interface{}
        if err := json.Unmarshal([]byte(vc), &parentcmd); err == nil {
          includecount := 0
          for _, childraw := range childcmd {
            cmdinparent := false
            for _, parentraw := range parentcmd.([]interface{}) {
              if parentraw == childraw {
                cmdinparent = true
                break
              }
            }
            if cmdinparent {
              includecount++
            }
          }
          if includecount == len(parentcmd.([]interface{})) {
            if len(parentcmd.([]interface{})) < len(childcmd) {
              if maxlayers < len(parentcmd.([]interface{})) {
                maxlayers = len(parentcmd.([]interface{}))
                maxname = kn
              }
            }
          }
        } else {
          say.L3(err.Error())
          ok = false
          return
        }
      }
    }
  }
  if maxlayers == 0 {
    ok = false
    say.L1("Parent not found")
  } else {
    say.L1("Parent is [ "+ maxname +" ]")
    s := strings.Split(maxname, ":")
    name = s[0]
    tag = s[1]
  }
  return
}
