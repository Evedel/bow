package handler

import(
  "db"

  "net/http"

  "github.com/Evedel/glb/say"
)

func UpgradeDB(w http.ResponseWriter, r *http.Request){
	funcname := r.URL.Path[len("/upgrade/"):]
	say.L3("Starting upgrade for [ " + funcname + " ]", "","\n")
	if funcname == "totalsize" {
		db.UpgradeTotalSize()
	}
	if funcname == "falsenumnames" {
		db.UpgradeFalseNumericImage()
	}
	if funcname == "oldparentnames" {
		db.UpgradeOldParentNames()
	}
	http.Redirect(w, r, "/", 307)
}
