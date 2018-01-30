package handler

import(
  "checker"

  "net/http"
)

func UpdateAll(w http.ResponseWriter, r *http.Request){
  checker.StartManual()
  http.Redirect(w, r, "/info", 307)
}
