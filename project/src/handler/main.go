package handler

import(
  "net/http"

  "github.com/Evedel/glb/say"
)

func Main(w http.ResponseWriter, r *http.Request){
	if r.URL.Path == "/" {
		Info(w, r)
	} else {
		if r.URL.Path != "/favicon.ico" {
			say.L1("Main Handler : Wrong query [" + r.URL.Path + "]", "","\n")
			http.Redirect(w, r, "/", 307)
		}
	}
}
