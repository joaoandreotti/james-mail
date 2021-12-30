package pages

import (
  "net/http"
  "james-mail/middleware"
  "james-mail/authentication"
  "github.com/gorilla/mux"
)

const logoutPagePath string = "/logout"

func LogoutGet(w http.ResponseWriter, r *http.Request) {
  if !authentication.ValidateAuthentication(r) {
    http.Redirect(w, r, errorPagePath, 302)
    return
  }

  authentication.Unauthenticate(w, r)
  http.Redirect(w, r, loginPagePath, 302)
}

func ConfigLogout(r *mux.Router) {
  r.HandleFunc(logoutPagePath, middleware.Chain(LogoutGet)).Methods("GET")
}
