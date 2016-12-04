package qurl

import(
  "db"
  "say"
  "strings"
  "strconv"
  "net/http"
  "io/ioutil"
  "crypto/tls"
  "encoding/json"
)

func makequery(rqst *http.Request, secure bool) (body interface{}, header map[string][]string, status int){
  var client *http.Client
  if secure {
    client = &http.Client{}
  } else {
    tr := &http.Transport{
       TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client = &http.Client{Transport: tr}
  }
  if resp, err := client.Do(rqst); err != nil {
    status = -1
    say.L3(err.Error())
    say.L3("Probably something wrong with network configuration or registry state")
  } else {
    defer resp.Body.Close()
    status = resp.StatusCode
    header = resp.Header
    if status == 200 || status == 401 {
      if bodytmp, err := ioutil.ReadAll(resp.Body); err != nil {
        status = -1
        say.L3(err.Error())
      } else {
        var c interface{}
        if err := json.Unmarshal(bodytmp, &c); err != nil {
          status = -1
          say.L3(err.Error())
        } else {
          body = c
          if c.(map[string]interface{})["errors"] != nil  && status != 401 {
            say.L3(c.(map[string]interface{})["errors"].([]interface{})[0].(map[string]interface{})["message"].(string))
          }
        }
      }
    }
  }
  return
}

func getbearertoken(wwwauth string, user string, pass string, secure bool) (token string, ok bool){
  ok = false
  query := ""
  splitted := strings.Split(wwwauth, ",")
  idaddrs := strings.Index(splitted[0], "\"")
  idservc := strings.Index(splitted[1], "\"")
  idscope := strings.Index(splitted[2], "\"")
  if (splitted[0][:idaddrs-1]=="Bearer realm") &&
     (splitted[1][:idservc-1]=="service") &&
     (splitted[2][:idscope-1]=="scope"){

     addrs := splitted[0][idaddrs+1:len(splitted[0])-1]
     servc := splitted[1][idservc+1:len(splitted[1])-1]
     spaceinservice := strings.Index(servc, " ")
     servc = servc[:spaceinservice] + "+" + servc[spaceinservice+1:]
     scope := splitted[2][idscope+1:len(splitted[2])-1]
     endofschemeinadress := strings.Index(addrs, "://")
     addrs = addrs[:endofschemeinadress+3] + user + ":" + pass + "@" + addrs[endofschemeinadress+3:]
     query = addrs + "?account=" + user + "&service=" + servc + "&scope=" + scope
     if reqst, err := http.NewRequest("GET", query, nil); err != nil {
       say.L3(err.Error())
     } else {
       body, _, c := makequery(reqst, secure)
       if c == 200 {
         token = body.(map[string]interface{})["token"].(string)
         ok = true
       }
     }
   } else {
     say.L3("GetBearerToken: Registry sent wrong Www-Authenticate header.")
     say.L3(wwwauth)
   }
   return
}

func MakeSimpleQuery(query string, info map[string]string) (body interface{}, ok bool){
  ok = false
  var c int
  var h map[string][]string
  secure := true
  if info["secure"] == "false" { secure = false}
  tquery := info["scheme"] + "://" + info["user"] + ":" + info["pass"] + "@" + info["host"] + query
  if reqst, err := http.NewRequest("GET", tquery, nil); err != nil {
    say.L3(err.Error())
  } else {
    body, h, c = makequery(reqst, secure)
    if c == 200 {
      ok = true
    } else {
      switch c {
      case 401:
        if h["Www-Authenticate"][0][0:5] == "Basic" {
          say.L3("MakeSimpleQuery: SCode [401] : Unauthorized response is returned (credentials problem, check user/pass pair)")
        } else if h["Www-Authenticate"][0][0:6] == "Bearer" {
          say.L1("MakeSimpleQuery: SCode [401] : Bearer auth. Trying to get auth token.")
          if token, oktok := getbearertoken(h["Www-Authenticate"][0], info["user"], info["pass"], secure); oktok {
            say.L1("MakeSimpleQuery: Token recieved. Retriying query.")
            tquery = info["scheme"] + "://" + info["host"] + query
            if reqst, err := http.NewRequest("GET", tquery, nil); err != nil {
              say.L3(err.Error())
            } else {
              reqst.Header.Add("Authorization", "Bearer "+token)
              body, h, c = makequery(reqst, secure)
              if c == 200 {
                ok = true
              } else {
                switch c {
                case 401:
                  say.L3("MakeSimpleQuery: Token code [401] : Unauthorized response is returned (credentials problem, check user/pass pair or communication between registry and auth server)")
                  say.L3(body.(map[string]interface{})["errors"].([]interface{})[0].(map[string]interface{})["message"].(string))
                case -1:
                }
              }
            }
          }
        }
      case -1:
      default:  say.L3("MakeSimpleQuery: Cannot diagnose problem. SCode \n[ " + strconv.Itoa(c) + " ] ")
      }
    }
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
      if Resp.StatusCode == 405 {
        say.L3(Resp.Status)
        say.L3("You need to add '-e REGISTRY_STORAGE_DELETE_ENABLED=true'")
        say.L3("Follow instructions here: https://github.com/Evedel/bow#image-deletion")
      } else {
        say.L3("Delete manifest: " + Resp.Status)
      }
      say.L3(ReqtStr)
    }
  }
  return
}

func GetfsLayerSize(link string ) (size string){
  if Resp, err := http.Head(link); err != nil {
    say.L3(err.Error())
    say.L3("GetfsLayerSize: Cannot recieve response from registry, stopping work")
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
