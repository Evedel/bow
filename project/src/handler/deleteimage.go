package handler

import(
  "db"
  "say"
  "qurl"
  "checker"

  "net/url"
  "net/http"
)

func DeleteImage(w http.ResponseWriter, r *http.Request){
	if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
		say.L3(err.Error())
	} else {
		if (v["reponame"] != nil) && (v["curname"] != nil) && (v["curtag"] != nil) {
			say.L1("Starting delete manifest [ " + v["reponame"][0] + "/" + v["curname"][0] + "/" + v["curtag"][0] + " ]")
			if qurl.DeleteTagFromRepo(v["reponame"][0], v["curname"][0], v["curtag"][0]) {
				db.PutSimplePairToBucket([]string{ v["reponame"][0], "catalog", v["curname"][0], v["curtag"][0]}, "_valid", "0")
				go checker.CheckTags()
				http.Redirect(w, r, "/info/" + v["reponame"][0] + "?curname=" + v["curname"][0], 307)
			}
		} else {
			say.L3("Something wrong with args in deleteHandler")
		}
	}
}
