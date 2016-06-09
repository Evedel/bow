package main

import (
	"say"
	"conf"
	"db"
	"strconv"
	"net/url"
	"net/http"
	"html/template"
	"checker"
	_ "github.com/wader/disable_sendfile_vbox_linux"
)
func main() {
	conf.Init()
	db.Init()

	http.HandleFunc("/managerepos/", mrepoHandler)
	http.HandleFunc("/info/", infoHandler)
	http.HandleFunc("/", welcomeHandler)

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

	go checker.DaemonManager()

	say.Info("Server listening at [" + conf.Env["servadd"] + "]")
	if err := http.ListenAndServe(conf.Env["servadd"], nil); err != nil {
		say.Error(err.Error() + "\nListenAndServe()\nmain()\nmain.go\nmain")
	}
}
func welcomeHandler(w http.ResponseWriter, r *http.Request){
	repos := db.GetRepos()
	irepos := make(map[string]interface{}, len(repos))
	irepos["repos"] = repos
	renderTemplate(w, "welcome", irepos)
}
func mrepoHandler(w http.ResponseWriter, r *http.Request){
	urlc := r.URL.Path[len("/managerepos/"):]
	repos := db.GetRepos()
	var repopretty map[string]string
	if urlc == "add" {
		if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
			say.Raw(err)
		} else {
			if len(v) != 0 {
				db.CreateRepo(v)
				http.Redirect(w, r, "/managerepos/", 307)
			}
		}
	}
	if urlc == "edit" {
		if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
			say.Raw(err)
		} else {
			if len(v) == 1 {
				repopretty = db.GetRepoPretty(v["reponame"][0])
			}
			if len(v) > 1 {
				db.CreateRepo(v)
				http.Redirect(w, r, "/managerepos/", 307)
			}
		}
	}
	if urlc == "delete" {
		if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
			say.Raw(err)
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

	renderTemplate(w, "managerepos", irepos)
}
func infoHandler(w http.ResponseWriter, r *http.Request){
	irepos := make(map[string]interface{})
	irepos["reponame"] = r.URL.Path[len("/info/"):]
	repo := db.GetRepoPretty(irepos["reponame"].(string))
	irepos["header"] = irepos["reponame"].(string) + " : " + repo["repohost"]

	if v, err := url.ParseQuery(r.URL.RawQuery); err != nil {
		say.Raw(err)
	} else {
		if len(v) != 0 {
			if v["curname"] != nil {
				tags := db.GetTags(irepos["reponame"].(string), v["curname"][0])
				uploads := make(map[string]int)
				for _, e := range tags {
					uploads[e] = totalUploads(irepos["reponame"].(string), v["curname"][0], e)
				}
				irepos["tags"] = uploads
				irepos["curname"] = v["curname"][0]
				irepos["header"] = irepos["header"].(string) + "/" + irepos["curname"].(string)
				if v["curtag"] != nil {
					irepos["curtag"] = v["curtag"][0]
					irepos["header"] = irepos["header"].(string) + "/" + irepos["curtag"].(string)
				}
			}
		}
	}

	irepos["catalog"] = db.GetCatalog(irepos["reponame"].(string))
	renderTemplate(w, "info", irepos)
}
func totalUploads(repo string, name string, tag string) (count int){
	uploads := db.GetTagSubbucket(repo, name, tag, "_uploads")
	for _, e := range uploads {
		if num, err := strconv.Atoi(e); err != nil {
			say.Error(err.Error())
		} else {
			count += num
		}
	}
	return
}
func renderTemplate(w http.ResponseWriter, tmpl string, c interface{}) {
	say.Info("Rendering template [ " + tmpl + " ]")
	templates := template.Must(template.ParseGlob("./templates/*"))
	err := templates.ExecuteTemplate(w, tmpl, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
