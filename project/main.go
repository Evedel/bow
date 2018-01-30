package main

import (
	"db"
	"say"
	"conf"
	"handler"
	"checker"
	"net/http"
	_ "github.com/wader/disable_sendfile_vbox_linux"
)
func main() {
	conf.Init()
	db.Init()

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

	http.HandleFunc("/managerepos/", handler.ManageRepos)
	http.HandleFunc("/info", handler.Info)
	http.HandleFunc("/upgrade/", handler.UpgradeDB)
	http.HandleFunc("/delete", handler.DeleteImage)
	http.HandleFunc("/graph", handler.RepoGraph)
	http.HandleFunc("/update", handler.UpdateAll)
	http.HandleFunc("/", handler.Main)

	go checker.DaemonManager()

	say.L2("Server listening at [" + conf.Env["servadd"] + "]")
	if err := http.ListenAndServe(conf.Env["servadd"], nil); err != nil {
		say.L3(err.Error() + "\nListenAndServe()\nmain()\nmain.go\nmain")
	}
}
