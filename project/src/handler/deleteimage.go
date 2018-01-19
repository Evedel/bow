package handler

import(
  "db"
  "qurl"
  "checker"

  "net/url"
  "net/http"

  "github.com/Evedel/glb/say"
)

func DeleteImage(w http.ResponseWriter, r *http.Request){
	if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
		say.L1("", err, "\n")
	} else {
		if (v["reponame"] != nil) && (v["curname"] != nil) && (v["curtag"] != nil) {
			say.L3("Starting delete manifest [ " + v["reponame"][0] + "/" + v["curname"][0] + "/" + v["curtag"][0] + " ]", "","\n")
      query := "/v2/" + v["curname"][0] + "/manifests/" +
                db.GetValueFromBucket([]string{v["reponame"][0], "catalog", v["curname"][0], v["curtag"][0]}, "digest")
      inhdr := map[string]string{"Accept": "application/vnd.docker.distribution.manifest.v2+json"}
      if _, _, ok := qurl.MakeQuery(query, "DELETE", db.GetRepoPretty(v["reponame"][0]), inhdr); ok {
				db.PutSimplePairToBucket([]string{ v["reponame"][0], "catalog", v["curname"][0], v["curtag"][0]}, "_valid", "0")
				go checker.RunCheckTags()
				http.Redirect(w, r, "/info?reponame=" + v["reponame"][0] + "&curname=" + v["curname"][0], 307)
			}
		} else {
			say.L1("Something wrong with args in deleteHandler", "","\n")
		}
	}
}
