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

		if v["reponame"] == nil {
      say.L3("Name of repository not set in repograpHandler")
    } else {
			if v["reponame"][0] != "" {

        irepos["reponame"] = v["reponame"][0]

				irepos["graphdata"] = db.GetSchemaFromPoint([]string{irepos["reponame"].(string), "_namesgraph"})

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
