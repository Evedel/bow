package handler

import(
  "db"
  "say"
  "utils"

  "net/url"
  "net/http"
)

func ManageRepos(w http.ResponseWriter, r *http.Request){
	urlc := r.URL.Path[len("/managerepos/"):]
	repos := db.GetRepos()
	var repopretty map[string]string
	if urlc == "add" {
		if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
			say.L3(err.Error())
		} else {
			if len(v) != 0 {
				db.CreateRepo(v)
				http.Redirect(w, r, "/managerepos/", 307)
			}
		}
	}
	if urlc == "edit" {
		if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
			say.L3(err.Error())
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
			say.L3(err.Error())
		} else {
			if len(v) == 1 {
				db.DeleteRepo(v["reponame"][0])
				http.Redirect(w, r, "/managerepos/", 307)
			}
		}
	}
	irepos := make(map[string]interface{}, len(urlc)+len(repos)+len(repopretty))

	irepos["path"] = urlc
	irepos["repos"] = utils.Keys(repos)
	irepos["chosen"] = repopretty
	irepos["action"] = "conf"

	renderTemplate(w, "managerepos", irepos)
}
