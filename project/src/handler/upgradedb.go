package handler

import(
  "db"
  "say"

  "net/http"
)

func UpgradeDB(w http.ResponseWriter, r *http.Request){
	funcname := r.URL.Path[len("/upgrade/"):]
	say.L1("Starting upgrade for [ " + funcname + " ]")
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
