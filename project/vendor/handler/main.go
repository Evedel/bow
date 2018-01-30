package handler

import(
  "say"

  "net/http"
)

func Main(w http.ResponseWriter, r *http.Request){
	if r.URL.Path == "/" {
		Info(w, r)
	} else {
		if r.URL.Path != "/favicon.ico" {
			say.L3("Main Handler : Wrong query [" + r.URL.Path + "]")
			http.Redirect(w, r, "/", 307)
		}
	}
}
