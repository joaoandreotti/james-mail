package main

import (
  "sync"
  "james-mail/pages"
  "james-mail/database"
  "james-mail/smtp"
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

  var wg sync.WaitGroup
  wg.Add(2)
  go smtp.StartSmtpServer(&wg)
  go pages.StartHttpServer(r, &wg)
  wg.Wait()
}
