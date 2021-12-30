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

var sendEmailPagePath string = "/send_email"

var sendEmailFormRecipient string = "recipient"
var sendEmailFormEmailBody string = "email_body"

type SendEmailPageData struct {
  Username string
  SendEmailFormRecipient string
  SendEmailFormEmailBody string
  EmailsPagePath string
  SendEmailPagePath string
}

func GetSendEmailPageData(r *http.Request) (data SendEmailPageData) {
  data = SendEmailPageData {
    Username: authentication.GetSessionValue(r, "username"),
    EmailsPagePath: emailsPagePath,
    SendEmailFormRecipient: sendEmailFormRecipient,
    SendEmailFormEmailBody: sendEmailFormEmailBody,
    SendEmailPagePath: sendEmailPagePath,
  }
  return data
}

func SendEmailGet(w http.ResponseWriter, r *http.Request) {
  if !authentication.ValidateAuthentication(r) {
    http.Redirect(w, r, loginPagePath, 302)
    return
  }

  page := template.Must(template.ParseFiles("pages/html/send_email.html"))
  data := GetSendEmailPageData(r)
  page.Execute(w, data)
}

func SendEmailPost(w http.ResponseWriter, r *http.Request) {
  if !authentication.ValidateAuthentication(r) {
    http.Redirect(w, r, loginPagePath, 302)
    return
  }

  data := database.SendEmailData {
    Recipient: r.FormValue(sendEmailFormRecipient),
    EmailBody: r.FormValue(sendEmailFormEmailBody),
  }

  if !validator.ValidateStruct(data) {
    http.Redirect(w, r, errorPagePath, 302)
    return
  }

  userEmail := authentication.GetSessionValue(r, "email")
  publicKey, _ := database.GetKeyPairFromEmail(userEmail)
  encryptedEmail := pgp.EncryptMessage(publicKey, data.EmailBody)
  database.InsertSentEmail(userEmail, data.Recipient, encryptedEmail)

  http.Redirect(w, r, emailsPagePath, 302)
}

func ConfigSendEmail(r *mux.Router) {
  r.HandleFunc(sendEmailPagePath, middleware.Chain(SendEmailGet)).Methods("GET")
  r.HandleFunc(sendEmailPagePath, middleware.Chain(SendEmailPost)).Methods("POST")
}
