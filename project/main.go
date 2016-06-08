package main

import (
	"say"
	"conf"
	"db"
	"net/url"
	"net/http"
	"html/template"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"checker"
// make it work in virtualbox
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
				irepos["tags"] = db.GetTags(irepos["reponame"].(string), v["curname"][0])
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
func calcImageSize(reponame string, name string, manifest interface{} ) int {
	var size int
	sliceMap := manifest.(map[string]interface{})["fsLayers"]
	for _, element := range sliceMap.([]interface{}) {
		reqt :=  reponame + "/v2/" + name + "/blobs/" + element.(map[string]interface {})["blobSum"].(string)
		if resp, err := http.Head(reqt); err != nil {
			say.Error(err.Error())
		} else {
			bytes := resp.Header.Get("Content-Length")
			if bn, err := strconv.Atoi(bytes); err != nil {
				say.Error(err.Error())
			} else {
				size += bn
			}
		}
	}
	return size
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
func getManifest(reponame string, name string, tag string, c *interface{})  {
	say.Info("Preparing data for manifest [ " + name + ":" + tag + " ]")
	reqt := reponame + "/v2/" + name + "/manifests/" + tag
	if resp, err := http.Get(reqt); err != nil {
		say.Error(err.Error())
	} else {
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			say.Error(err.Error())
		} else {
			if err := json.Unmarshal(body, c); err != nil {
				say.Error(err.Error())
			}
		}
	}
}
