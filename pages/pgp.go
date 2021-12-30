package pages

import (
  "net/http"
  "html/template"
  "james-mail/middleware"
  "james-mail/authentication"
  "james-mail/database"
  "james-mail/validator"
  "james-mail/pgp"
  "github.com/gorilla/mux"
)

const pgpPagePath string = "/pgp"

const pgpFormPassword string = "password"

type PgpPageData struct {
  PublicKey string
  PrivateKey string
  Username string
  EmailsPagePath string
  PgpFormPassword string
  PgpPagePath string
}

func GetPgpPageData() (data PgpPageData) {
  data = PgpPageData {
    PgpFormPassword: pgpFormPassword,
    EmailsPagePath: emailsPagePath,
    PgpPagePath: pgpPagePath,
  }
  return data
}

func PgpGet(w http.ResponseWriter, r *http.Request) {
  if !authentication.ValidateAuthentication(r) {
    http.Redirect(w, r, loginPagePath, 302)
    return
  }

  page := template.Must(template.ParseFiles("pages/html/pgp.html"))
  email := authentication.GetSessionValue(r, "email")
  publicKey, privateKey := database.GetKeyPairFromEmail(email)
  username := authentication.GetSessionValue(r, "username")
  data := GetPgpPageData()
  data.Username = username
  data.PublicKey = publicKey
  data.PrivateKey = privateKey

  page.Execute(w, data)
}

func PgpGeneratePost(w http.ResponseWriter, r *http.Request) {
  if !authentication.ValidateAuthentication(r) {
    http.Redirect(w, r, loginPagePath, 302)
    return
  }

  data := database.PgpFormData {
    Password: r.FormValue(pgpFormPassword),
  }
  if !validator.ValidateStruct(data) {
    http.Redirect(w, r, errorPagePath, 302)
    return
  }

  email := authentication.GetSessionValue(r, "email")
  username := authentication.GetSessionValue(r, "username")
  publicKey, privateKey := pgp.GeneratePgpKeyPair(username, email, data.Password)
  database.UpdateKeyPairFromEmail(publicKey, privateKey, email)

  http.Redirect(w, r, pgpPagePath, 302)
}

func ConfigPgp(r *mux.Router) {
  r.HandleFunc(pgpPagePath, middleware.Chain(PgpGet)).Methods("GET")
  r.HandleFunc(pgpPagePath, middleware.Chain(PgpGeneratePost)).Methods("POST")
}
