package pages

import (
  "net/http"
  "html/template"
  "james-mail/middleware"
  "james-mail/authentication"
  "james-mail/database"
  "james-mail/validator"
  "james-mail/smtp"
  "github.com/gorilla/mux"
)

var emailsPagePath string = "/emails"

var emailsReceivedFormPassword string = "password"

type EmailsPageData struct {
  Username string
  EmailsReceived template.HTML
  PgpPagePath string
  LogoutPagePath string
  SendEmailPagePath string
  SentEmailsPagePath string
  EmailsReceivedFormPassword string
  EmailsPagePath string
}

func GetEmailsPageData(r *http.Request) EmailsPageData{
  data := EmailsPageData {
    Username: authentication.GetSessionValue(r, "username"),
    PgpPagePath: pgpPagePath,
    LogoutPagePath: logoutPagePath,
    SendEmailPagePath: sendEmailPagePath,
    SentEmailsPagePath: sentEmailsPagePath,
    EmailsReceived: template.HTML(""),
    EmailsReceivedFormPassword: emailsReceivedFormPassword,
    EmailsPagePath: emailsPagePath,
  }
  return data
}

func EmailsGet(w http.ResponseWriter, r *http.Request) {
  if !authentication.ValidateAuthentication(r) {
    http.Redirect(w, r, loginPagePath, 302)
    return
  }

  userEmail := authentication.GetSessionValue(r, "email")
  emailList := smtp.GetReceivedEmailsList(userEmail, "")
  data := GetEmailsPageData(r)
  data.EmailsReceived = template.HTML(GetEmailHtmlList(emailList))

  page := template.Must(template.ParseFiles("pages/html/emails.html"))
  page.Execute(w, data)
}

func EmailsReceivedPost(w http.ResponseWriter, r *http.Request) {
  if !authentication.ValidateAuthentication(r) {
    http.Redirect(w, r, loginPagePath, 302)
    return
  }

  formData := database.PgpFormData {
    Password: r.FormValue(emailsReceivedFormPassword),
  }
  if !validator.ValidateStruct(formData) {
    http.Redirect(w, r, errorPagePath, 302)
    return
  }

  userEmail := authentication.GetSessionValue(r, "email")
  emailList := smtp.GetReceivedEmailsList(userEmail, formData.Password)
  data := GetEmailsPageData(r)
  data.EmailsReceived = template.HTML(GetEmailHtmlList(emailList))

  page := template.Must(template.ParseFiles("pages/html/emails.html"))
  page.Execute(w, data)
}

func ConfigEmails(r *mux.Router) {
  r.HandleFunc(emailsPagePath, middleware.Chain(EmailsGet)).Methods("GET")
  r.HandleFunc(emailsPagePath, middleware.Chain(EmailsReceivedPost)).Methods("POST")
}
