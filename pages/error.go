package pages

import (
  "net/http"
  "html/template"
  "james-mail/middleware"
  "github.com/gorilla/mux"
)

var errorPagePath string = "/error"

func ErrorGet(w http.ResponseWriter, r *http.Request) {
  page := template.Must(template.ParseFiles("pages/html/error.html"))
  page.Execute(w, "")
}

func ConfigError(r *mux.Router) {
  r.HandleFunc(errorPagePath, middleware.Chain(ErrorGet)).Methods("GET")
}
