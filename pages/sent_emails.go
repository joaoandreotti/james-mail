package pages

import (
  "net/http"
  "html/template"
  "james-mail/middleware"
  "james-mail/authentication"
  "james-mail/validator"
  "james-mail/database"
  "github.com/gorilla/mux"
)

var sentEmailsPagePath string = "/sent_emails"

var sentEmailsFormPassword string = "password"

type SentEmailsPageData struct {
  Username string
  EmailsPagePath string
  SentEmailsPagePath string
  SentEmailsFormPassword string
  SentEmailsList template.HTML
}

func GetSentEmailsPageData(r *http.Request) (data SentEmailsPageData) {
  data = SentEmailsPageData {
    Username: authentication.GetSessionValue(r, "username"),
    EmailsPagePath: emailsPagePath,
    SentEmailsPagePath: sentEmailsPagePath,
    SentEmailsFormPassword: sentEmailsFormPassword,
    SentEmailsList: template.HTML(""),
  }
  return data
}

func SentEmailsGet(w http.ResponseWriter, r *http.Request) {
  if !authentication.ValidateAuthentication(r) {
    http.Redirect(w, r, loginPagePath, 302)
    return
  }

  userEmail := authentication.GetSessionValue(r, "email")
  encryptedEmailList := database.GetSentEmailsFromEmail(userEmail)
  data := GetSentEmailsPageData(r)
  data.SentEmailsList = template.HTML(GetEmailHtmlList(encryptedEmailList, "", ""))

  page := template.Must(template.ParseFiles("pages/html/sent_emails.html"))
  page.Execute(w, data)
}

func SentEmailsPost(w http.ResponseWriter, r *http.Request) {
  if !authentication.ValidateAuthentication(r) {
    http.Redirect(w, r, loginPagePath, 302)
    return
  }

  formData := database.PgpFormData {
    Password: r.FormValue(sentEmailsFormPassword),
  }
  if !validator.ValidateStruct(formData) {
    http.Redirect(w, r, errorPagePath, 302)
    return
  }

  userEmail := authentication.GetSessionValue(r, "email")
  encryptedEmailList := database.GetSentEmailsFromEmail(userEmail)
  _, privateKey := database.GetKeyPairFromEmail(userEmail)
  data := GetSentEmailsPageData(r)
  data.SentEmailsList = template.HTML(GetEmailHtmlList(encryptedEmailList, privateKey, formData.Password))

  page := template.Must(template.ParseFiles("pages/html/sent_emails.html"))
  page.Execute(w, data)
}

func ConfigSentEmails(r *mux.Router) {
  r.HandleFunc(sentEmailsPagePath, middleware.Chain(SentEmailsGet)).Methods("GET")
  r.HandleFunc(sentEmailsPagePath, middleware.Chain(SentEmailsPost)).Methods("POST")
}
