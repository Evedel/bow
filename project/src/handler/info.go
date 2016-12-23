package handler

import(
  "db"
  "say"
  "utils"

  "sort"
  "strconv"
  "net/url"
  "net/http"
  "encoding/json"
)

func Info(w http.ResponseWriter, r *http.Request){
	irepos := make(map[string]interface{})
	headerdata := make(map[string]string)
	repos := utils.Keys(db.GetRepos())
  sort.Strings(repos)
  irepos["repos"] = repos

	if len(r.URL.Path) <= 6 {
		headerdata["header"] = ""
		headerdata["currepo"] = ""
		irepos["curname"] = ""
	} else {
		irepos["reponame"] = r.URL.Path[len("/info/"):]
		headerdata["header"] = irepos["reponame"].(string) + " : " + db.GetRepoPretty(irepos["reponame"].(string))["host"]
		headerdata["currepo"] = irepos["reponame"].(string)
		irepos["curname"] = irepos["reponame"].(string)
		if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
			say.L3(err.Error())
		} else {
			if len(v) != 0 {
				if v["curname"] != nil {
					irepos["curname"] = v["curname"][0]

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
						irepos["parent"] = db.GetAllPairsFromBucket(dbpath)
					}
				}
			}
		}
		names := db.GetCatalog(irepos["reponame"].(string))
    sort.Strings(names)
    irepos["catalog"] = names
	}

	irepos["headerdata"] = headerdata
	irepos["action"] = "repos"
	renderTemplate(w, "info", irepos)
}
