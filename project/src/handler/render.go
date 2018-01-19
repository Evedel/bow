package handler

import(
  "net/http"
  "html/template"

  "github.com/Evedel/glb/say"
)

func renderTemplate(w http.ResponseWriter, tmpl string, c interface{}) {
	say.L3("Rendering template [ " + tmpl + " ]", "","\n")
	templates := template.Must(template.ParseGlob("./templates/*"))
	err := templates.ExecuteTemplate(w, tmpl, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
