package main

import (
	"db"
	"conf"
	"handler"
	"checker"

	"net/http"

	_ "github.com/wader/disable_sendfile_vbox_linux"
	"github.com/Evedel/glb/say"
)
func main() {
	conf.Init()
	db.Init()
	say.L0("", conf.Env, "\n")
	say.Init(conf.Env["log_level"])

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

	http.HandleFunc("/managerepos/", handler.ManageRepos)
	http.HandleFunc("/info", handler.Info)
	http.HandleFunc("/upgrade/", handler.UpgradeDB)
	http.HandleFunc("/delete", handler.DeleteImage)
	http.HandleFunc("/graph", handler.RepoGraph)
	http.HandleFunc("/update", handler.UpdateAll)
	http.HandleFunc("/", handler.Main)

	say.L0("", db.GetAllPairsFromBucket([]string{"basic", "_info"}), "\n")

	go checker.DaemonManager()

	say.L2("Main: Server listening at [" + conf.Env["servadd"] + "]", "","\n")
	if err := http.ListenAndServe(conf.Env["servadd"], nil); err != nil {
		say.L1("Main: Server start failed: ", err, "\n")
	}
}
