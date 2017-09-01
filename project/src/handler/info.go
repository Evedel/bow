package handler

import(
  "db"
  "dt"
  "say"
  "utils"

  "time"
  "sort"
  "strconv"
  "net/url"
  "net/http"
  "encoding/json"
)

func Info(w http.ResponseWriter, r *http.Request){
  defer dt.Watch(time.Now(), "Info Handler")

	irepos := make(map[string]interface{})
	headerdata := make(map[string]string)
	repos := utils.Keys(db.GetRepos())
  sort.Strings(repos)
  irepos["repos"] = repos

  if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
    say.L3(err.Error())
  } else {
    headerdata["header"] = ""
    headerdata["currepo"] = ""
    irepos["curname"] = ""

    if v["reponame"] != nil {
      if v["reponame"][0] != "" {
      	irepos["reponame"] = v["reponame"][0]

      	headerdata["header"] = irepos["reponame"].(string) + " : " +
          db.GetRepoPretty(irepos["reponame"].(string))["host"]
    		headerdata["currepo"] = irepos["reponame"].(string)
    		irepos["curname"] = irepos["reponame"].(string)

        namesStructure := db.GetCatalogStructure(irepos["reponame"].(string))
        nsarr := []string{}
        if (len(namesStructure) > 0) {
          sort.Slice(namesStructure, func(i, j int) bool {
            return namesStructure[i]["fullname"] < namesStructure[j]["fullname"]
          })
          nslast := namesStructure[0]["namespace"]
          nsarr = append(nsarr, nslast)
          for _, e := range namesStructure {
            if nslast != e["namespace"] {
              nsarr = append(nsarr, e["namespace"])
              nslast = e["namespace"]
            }
          }
        }
        irepos["namespaces"] = nsarr

        if v["curnamespace"] != nil {
          irepos["curnamespace"] = v["curnamespace"][0]

          nmarr := []string{}
          for _, e := range namesStructure {
            if irepos["curnamespace"] == e["namespace"] {
              nmarr = append(nmarr, e["nameshort"])
            }
          }
          irepos["shortnames"] = nmarr

          if v["curshortname"] != nil {
            irepos["curshortname"] = v["curshortname"][0]
            if irepos["curnamespace"].(string) == "_none" {
              irepos["curname"] = irepos["curshortname"].(string)
            } else {
              irepos["curname"] = irepos["curnamespace"].(string) + "/" + irepos["curshortname"].(string)
            }

            tags := db.GetTags(irepos["reponame"].(string), irepos["curname"].(string))
            sort.Strings(tags)
  					uploads := make(map[string]map[string]string)
  					totaluploads := make(map[string]int)
  					for _, e := range tags {
  						uploads[e] = make(map[string]string)
  						uploads[e] = db.GetAllPairsFromBucket([]string{
  							irepos["reponame"].(string),
  							"catalog",
  							irepos["curname"].(string),
  							e,
  							"_uploads" })
  						count := 0
  						for _, eu := range uploads[e] {
  							if num, err := strconv.Atoi(eu); err != nil {
  								say.L3(err.Error())
  							} else {
  								count += num
  							}
  						}
  						totaluploads[e] = count
  					}
  					irepos["tags"] = totaluploads
  					headerdata["header"] = headerdata["header"] + "/" + irepos["curname"].(string)
  					if v["curtag"] != nil {
  						irepos["curtag"] = v["curtag"][0]
  						irepos["uploads"] = uploads[irepos["curtag"].(string)]
  						headerdata["header"] = headerdata["header"] + ":" + irepos["curtag"].(string)
  						var dbpath = []string{
  							irepos["reponame"].(string),
  							"catalog",
  							irepos["curname"].(string),
  							irepos["curtag"].(string),
  							"history" }
  						strhist := db.GetAllPairsFromBucket(dbpath)
  						objhist := make(map[string]interface{})
  						lastkey := ""
  						layersnum := 0
  						for key, value := range  strhist {
  							var ch interface{}
  							_ = json.Unmarshal([]byte(value), &ch)
  							objhist[key] = ch
  							if lastkey < key {
  								lastkey = key
  							}
  							layersnum++
  						}
  						irepos["history"] = objhist
  						irepos["lastupdated"] = lastkey
  						irepos["layersnum"] = layersnum
  						dbpath[4] = "_totalsizehuman"
  						strsizehuman := db.GetAllPairsFromBucket(dbpath)
  						dbpath[4] = "_totalsizebytes"
  						strsizebytes := db.GetAllPairsFromBucket(dbpath)
  						lastkey = ""
  						for key, _ := range strsizehuman {
  							if lastkey < key {
  								lastkey = key
  							}
  						}
  						if strsizebytes != nil {
  							irepos["imagesizebytes"] = strsizebytes
  						}
  						if strsizehuman != nil {
  							irepos["imagesizehuman"] = strsizehuman
  						}
  						irepos["lastpushed"] = lastkey
  						dbpath[4] = "_parent"

              parentDescr := db.GetAllPairsFromBucket(dbpath)
              dbpath[2] = parentDescr["name"]
              dbpath[3] = "_namepair"
              npair := db.GetAllPairsFromBucket(dbpath[0:4])
              for e, k := range npair {
                parentDescr["namespace"] = e
                parentDescr["shortname"] = k
              }
  						irepos["parent"] = parentDescr
  					}
          }
        }
      }
    }
  }

	irepos["headerdata"] = headerdata
	irepos["action"] = "info"
	renderTemplate(w, "info", irepos)
}
