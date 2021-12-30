package main

import (
  "net/http"
  "james-mail/pages"
  "james-mail/database"
  "github.com/gorilla/mux"
)

func main() {
  database.ConfigConnection()

  r := mux.NewRouter()
  pages.ConfigIndex(r)
  pages.ConfigRegister(r)
  pages.ConfigLogin(r)
  pages.ConfigError(r)
  pages.ConfigEmails(r)
  pages.ConfigSendEmail(r)
  pages.ConfigSentEmails(r)
  pages.ConfigLogout(r)
  pages.ConfigPgp(r)
  http.ListenAndServe(":8080", r)
}
