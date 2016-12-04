package handler

import(
  "db"
  "say"
  "utils"

  "net/url"
  "net/http"
)

func RepoGraph(w http.ResponseWriter, r *http.Request){
	if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
		say.L3(err.Error())
	} else {
		if v["reponame"] != nil {
			irepos := make(map[string]interface{})
			if v["reponame"][0] != "" {
				irepos["graphdata"] = db.GetSchemaFromPoint([]string{v["reponame"][0], "_namesgraph"})

				headerdata := make(map[string]string)
				headerdata["header"] = v["reponame"][0] + " : " + db.GetRepoPretty(v["reponame"][0])["host"]
				headerdata["currepo"] = v["reponame"][0]
				irepos["headerdata"] = headerdata

				irepos["repodata"] = make(map[string]interface{})
				irepos["repodata"].(map[string]interface{})["catalog"] = utils.Keys(db.GetRepos())
				irepos["repodata"].(map[string]interface{})["curname"] = v["reponame"][0]
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

			irepos["action"] = "graph"
			renderTemplate(w, "repograph", irepos)
		} else {
			say.L3("Name of repository not set in repograpHandler")
		}
	}
}
