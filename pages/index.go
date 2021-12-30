package pages

import (
  "net/http"
  "html/template"
  "james-mail/middleware"
  "github.com/gorilla/mux"
)

func IndexGet(w http.ResponseWriter, r *http.Request) {
  page := template.Must(template.ParseFiles("pages/html/index.html"))
  page.Execute(w, "")
}

func ConfigIndex(r *mux.Router) {
  r.HandleFunc("/", middleware.Chain(IndexGet)).Methods("GET")
}
