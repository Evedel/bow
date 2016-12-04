package qurl

import(
  "say"
  "time"
  "flag"
  "utils"
  "strings"
  "net/http"
  "strconv"
  "testing"
  "encoding/json"
)

var TestAddress string
var TestInfo map[string]string
func init(){
  flag.StringVar(&TestAddress, "repo", "", "Repository to be tested")
  flag.Parse()
  tmpaddress := TestAddress
  TestInfo = make(map[string]string)
  TestInfo["scheme"]=tmpaddress[:strings.Index(tmpaddress, "://")]
  tmpaddress = tmpaddress[strings.Index(tmpaddress, "://")+3:]
  TestInfo["user"]=tmpaddress[:strings.Index(tmpaddress, ":")]
  tmpaddress = tmpaddress[strings.Index(tmpaddress, ":")+1:]
  TestInfo["pass"]=tmpaddress[:strings.Index(tmpaddress, "@")]
  tmpaddress = tmpaddress[strings.Index(tmpaddress, "@")+1:]
  TestInfo["host"]=tmpaddress
  TestInfo["secure"]="false"
}

func TestGetAPI(t *testing.T){
  say.L1(">> GET /v2/")
  if body_ping, ok := MakeSimpleQuery("/v2/", TestInfo); !ok {
    if body_ping == 404 {
      say.L3("API V2 is not supported by this registry")
    }
    t.Fail()
  } else {
    say.L1(">> GET /v2/_catalog")
    if body_catalog, ok := MakeSimpleQuery("/v2/_catalog", TestInfo); !ok {
      t.Fail()
    } else {
      var catalog interface{}
      if body_catalog.(map[string]interface{})["repositories"] == nil {
          say.L3("There is no 'repositories' field, check API specification")
          t.Fail()
      } else {
        catalog = body_catalog.(map[string]interface{})["repositories"]
        if len(catalog.([]interface{})) == 0 {
          say.L3("The catalog is empty, cannot make further checks")
          t.Fail()
        } else {
          for _, ei := range catalog.([]interface{}) {
            say.L1(">> GET /v2/" + ei.(string) + "/tags/list")
            if body_tags, ok := MakeSimpleQuery("/v2/" + ei.(string) + "/tags/list", TestInfo); ok {
              if body_tags.(map[string]interface{})["name"] != nil &&
                 body_tags.(map[string]interface{})["tags"] != nil {
                  for _, et := range body_tags.(map[string]interface{})["tags"].([]interface{}) {
                    say.L1(">> GET /v2/" + ei.(string) + "/manifests/" + et.(string))
                    if body_manifest, ok := MakeSimpleQuery("/v2/" + ei.(string) + "/manifests/" + et.(string), TestInfo); ok {
                      say.L1(">> CHECK MANIFEST FIELDS [" + ei.(string) + ":" + et.(string)+"]")
                      testManifestFields(t, body_manifest)
                      if !t.Failed() {
                        say.L1(">> Ok")
                      }
                      say.L1(">> CHECK CONTENT DIGEST [" + ei.(string) + ":" + et.(string)+"]")
                      testContentDigest(t, body_manifest)
                      if !t.Failed() {
                        say.L1(">> Ok")
                      }
                      say.L1(">> CHECK FSSHA & HISTORY [" + ei.(string) + ":" + et.(string)+"]")
                      testFSSHAandHistory(t, body_manifest)
                      if !t.Failed() {
                        say.L1(">> Ok")
                      }
                    }
                  }
              }
            }
          }
        }
      }
    }
  }
}

func testManifestFields(t *testing.T, manifest interface{}){
  if _, ok := manifest.(map[string]interface{})["name"].(string); !ok {
    say.L3("Something wrong with 'name' field.")
    say.L3("Probably API was changed.")
    say.L4(manifest.(map[string]interface{})["name"])
    t.Fail()
  }

  if _, ok := manifest.(map[string]interface{})["tag"].(string); !ok {
    say.L3("Something wrong with 'tag' field.")
    say.L3("Probably API was changed.")
    say.L4(manifest.(map[string]interface{})["tag"])
    t.Fail()
  }

  if _, ok := manifest.(map[string]interface{})["fsLayers"].([]interface{}); !ok {
    say.L3("Something wrong with 'fsLayers' field.")
    say.L3("Probably API was changed.")
    say.L4(manifest.(map[string]interface{})["fsLayers"])
    t.Fail()
  } else {
    if _, ok := manifest.(map[string]interface{})["fsLayers"].([]interface{})[0].(map[string]interface{})["blobSum"].(string); !ok {
      say.L3("Something wrong with 'blobSum' field.")
      say.L3("Probably API was changed.")
      say.L4(manifest.(map[string]interface{})["fsLayers"].([]interface{})[0].(map[string]interface{})["blobSum"])
      t.Fail()
    }
  }

  if e, ok := manifest.(map[string]interface{})["history"].([]interface{}); !ok {
    say.L3("Something wrong with 'history' field.")
    say.L3("Probably API was changed.")
    say.L4(manifest.(map[string]interface{})["history"])
    t.Fail()
  } else {
    if _, ok := e[0].(map[string]interface{})["v1Compatibility"].(string); !ok {
      say.L3("Something wrong with 'v1Compatibility' field.")
      say.L3("Probably API was changed.")
      say.L4(e[0].(map[string]interface{})["v1Compatibility"])
      t.Fail()
    }
  }
}

func testContentDigest(t *testing.T, Manifest interface{}){
  client := &http.Client{}
  Query := "/v2/" +
          Manifest.(map[string]interface{})["name"].(string) + "/manifests/" +
          Manifest.(map[string]interface{})["tag"].(string)
  say.L1(">> GET DIGEST " + Query)
  Reqtv2Digest, _ := http.NewRequest("GET", TestAddress + Query, nil)
  Reqtv2Digest.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
  if Respv2Digest, err := client.Do(Reqtv2Digest); err != nil {
    say.L3("Cannot recieve response from registry.")
    t.Fail()
  } else {
    defer Respv2Digest.Body.Close()
    contentdigest := Respv2Digest.Header.Get("Docker-Content-Digest")
    if contentdigest[0:7] != "sha256:"{
      say.L3("Digest was received in wrong format.")
      say.L3("Probably API was changed.")
      say.L4(contentdigest)
      say.L4(contentdigest[0:7])
      t.Fail()
    }
  }
}

func testFSSHAandHistory(t *testing.T, Manifest interface{}){
  fsshaarr := Manifest.(map[string]interface{})["fsLayers"].([]interface{})
  historyarr := Manifest.(map[string]interface{})["history"].([]interface{})
  for i, _ := range fsshaarr {
    fssha := fsshaarr[i].(map[string]interface{})["blobSum"].(string)
    fssize := GetfsLayerSize(TestAddress + "/v2/" +
      Manifest.(map[string]interface{})["name"].(string) + "/blobs/" + fssha)
    if fsshanum, err := strconv.Atoi(fssize); err != nil {
      say.L3(err.Error())
      say.L3("Cannot convert fssha size to int.")
      say.L3("Probably API was changed.")
      t.Fail()
    } else {
      if strconv.Itoa(utils.FromHumanToByte(utils.FromByteToHuman(fsshanum)))[0:2] != fssize[0:2] {
        say.L3("Connot convert size from Human to Byte and Back")
        t.Fail()
      }
    }
    var ch interface{}
    history := historyarr[i].(map[string]interface{})["v1Compatibility"].(string)
    if err := json.Unmarshal([]byte(history), &ch); err != nil {
      say.L3(err.Error())
      say.L3("Cannot convert v1Compatibility history to JSON.")
      say.L3("Probably API was changed.")
      t.Fail()
    } else {
      if created, ok := ch.(map[string]interface{})["created"].(string); !ok {
        say.L3("Something went wrong with 'history->v1Compatibility->created' field.")
        say.L3("Probably API was changed.")
        say.L4(ch.(map[string]interface{})["created"])
        t.Fail()
      } else {
        var indx int
        if indx = strings.Index(created, "T"); indx > -1 {
          created = created[:indx] + " " + created[indx+1:]
          if indx = strings.Index(created, "."); indx > -1 {
            created = created[:indx]
          }
        }
        if indx < 0 {
            say.L3("Timedate format was not in supposed form")
            t.Fail()
        }
        if dt, err := time.Parse("2006-01-02 15:04:05", created); err != nil {
          say.L3(err.Error())
          say.L3("Something went wrong with 'history->v1Compatibility->created' conversion.")
          say.L3("Probably API was changed.")
          say.L4(ch.(map[string]interface{})["created"])
          t.Fail()
        } else {
          if created != dt.Format("2006-01-02 15:04:05") {
            say.L3("Reverse time converson went wrong.")
            t.Fail()
          }
        }
      }
      if c_conf, ok := ch.(map[string]interface{})["container_config"]; !ok {
        say.L3("Something went wrong with 'manifest->history->container_config' fiels.")
        say.L3("Probably API was changed.")
        say.L4(ch)
        t.Fail()
      } else {
        if _, ok := c_conf.(map[string]interface{})["Cmd"].([]interface{}); !ok {
          say.L3("Something went wrong with 'manifest->history->container_config->Cmd' field.")
          say.L3("Probably API was changed.")
          say.L4(c_conf)
          t.Fail()
        }
      }
    }
  }
}
