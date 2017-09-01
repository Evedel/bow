package handler

import(
  "dt"
  "db"
  "say"
  "utils"

  "time"
  "sort"
  "net/url"
  "net/http"
)

func RepoGraph(w http.ResponseWriter, r *http.Request){
  defer dt.Watch(time.Now(), "Graph Handler")

	if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
		say.L3(err.Error())
	} else {
    irepos := make(map[string]interface{})
    headerdata := make(map[string]string)
    repos := utils.Keys(db.GetRepos())
    sort.Strings(repos)
    irepos["repos"] = repos
    irepos["action"] = "graph"

    headerdata["header"] = ""
    headerdata["currepo"] = ""
    irepos["curname"] = ""

    irepos["headerdata"] = headerdata

    filter := ""

		if v["reponame"] == nil {
      say.L3("Name of repository not set in repograpHandler")
    } else {
			if v["reponame"][0] != "" {
        irepos["reponame"] = v["reponame"][0]

        headerdata["header"] = irepos["reponame"].(string) + " : " +
          db.GetRepoPretty(irepos["reponame"].(string))["host"]
        headerdata["currepo"] = irepos["reponame"].(string)
        irepos["curname"] = irepos["reponame"].(string)

        nsarr := []string{}
        namesStructure := db.GetCatalogStructure(irepos["reponame"].(string))
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
          if irepos["curnamespace"].(string) != "_none" {
            filter = irepos["curnamespace"].(string)
          }

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

            headerdata["header"] = headerdata["header"] + "/" + irepos["curname"].(string)
            filter = irepos["curname"].(string)

            tags := db.GetTags(irepos["reponame"].(string), irepos["curname"].(string))
            sort.Strings(tags)
            totaluploads := make(map[string]string)
            for _, e := range tags {
              totaluploads[e] = ""
            }
            irepos["tags"] = totaluploads
            headerdata["header"] = headerdata["header"] + "/" + irepos["curname"].(string)

            if v["curtag"] != nil {
              irepos["curtag"] = v["curtag"][0]
              headerdata["header"] = headerdata["header"] + ":" + irepos["curtag"].(string)

              filter = irepos["curname"].(string) + ":" + irepos["curtag"].(string)
            }
          }
        }

				irepos["graphdata"] = db.Schema2json(db.GetSchemaFromPoint([]string{irepos["reponame"].(string), "_namesgraph"}, filter))

				headerdata := make(map[string]string)
				headerdata["header"] = irepos["reponame"].(string) + " : " +
          db.GetRepoPretty(irepos["reponame"].(string))["host"]
				headerdata["currepo"] = irepos["reponame"].(string)
				irepos["headerdata"] = headerdata

				irepos["repodata"] = make(map[string]interface{})
        repos := utils.Keys(db.GetRepos())
        sort.Strings(repos)
				irepos["repodata"].(map[string]interface{})["catalog"] = repos
				irepos["repodata"].(map[string]interface{})["curname"] = irepos["reponame"].(string)
			} else {
				irepos["graphdata"] = ""

				headerdata := make(map[string]string)
				headerdata["header"] = ""
				headerdata["currepo"] = ""
				irepos["headerdata"] = headerdata

				irepos["repodata"] = make(map[string]interface{})
				irepos["repodata"].(map[string]interface{})["catalog"] = utils.Keys(db.GetRepos())
				irepos["repodata"].(map[string]interface{})["curname"] = ""
			}
		}
    renderTemplate(w, "repograph", irepos)
	}
}
