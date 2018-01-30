package handler

import(
  "say"

  "net/http"
  "html/template"
)

func renderTemplate(w http.ResponseWriter, tmpl string, c interface{}) {
	say.L1("Rendering template [ " + tmpl + " ]")
	templates := template.Must(template.ParseGlob("./templates/*"))
	err := templates.ExecuteTemplate(w, tmpl, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
