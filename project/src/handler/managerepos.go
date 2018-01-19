package handler

import(
  "db"
  "utils"

  "sort"
  "net/url"
  "net/http"

  "github.com/Evedel/glb/say"
)

func ManageRepos(w http.ResponseWriter, r *http.Request){
	urlc := r.URL.Path[len("/managerepos/"):]
	repos := utils.Keys(db.GetRepos())
  sort.Strings(repos)
	var repopretty map[string]string
	if urlc == "add" {
		if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
      say.L1("", err, "\n")
		} else {
			if len(v) != 0 {
				db.CreateRepo(v)
				http.Redirect(w, r, "/managerepos/", 307)
			}
		}
	}
	if urlc == "edit" {
		if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
			say.L1("", err, "\n")
		} else {
			if len(v) == 1 {
				repopretty = db.GetRepoPretty(v["reponame"][0])
				repopretty["pass"] = ""
			}
			if len(v) > 1 {
				db.CreateRepo(v)
				http.Redirect(w, r, "/managerepos/", 307)
			}
		}
	}
	if urlc == "delete" {
		if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
			say.L1("", err, "\n")
		} else {
			if len(v) == 1 {
				db.DeleteRepo(v["reponame"][0])
				http.Redirect(w, r, "/managerepos/", 307)
			}
		}
	}
	irepos := make(map[string]interface{}, len(urlc)+len(repos)+len(repopretty))

	irepos["path"] = urlc
	irepos["repos"] = repos
	irepos["chosen"] = repopretty
	irepos["action"] = "conf"

	renderTemplate(w, "managerepos", irepos)
}
