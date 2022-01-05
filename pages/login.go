package pages

import (
  "net/http"
  "html/template"
  "james-mail/middleware"
  "james-mail/database"
  "james-mail/validator"
  "james-mail/authentication"
  "github.com/gorilla/mux"
)

const loginPagePath string = "/login"

const loginFormEmail string = "email"
const loginFormPassword string = "password"

type LoginPageData struct {
  LoginFormUsername string
  LoginFormEmail string
  LoginFormPassword string
  RegisterPagePath string
  LoginPagePath string
}

func GetLoginPageData() (data LoginPageData) {
  data = LoginPageData {
    LoginFormEmail: loginFormEmail,
    LoginFormPassword: loginFormPassword,
    RegisterPagePath: registerPagePath,
    LoginPagePath: loginPagePath,
  }
  return data
}

func LoginGet(w http.ResponseWriter, r *http.Request) {
  if authentication.ValidateAuthentication(r) {
    http.Redirect(w, r, emailsPagePath, 302)
    return
  }

  page := template.Must(template.ParseFiles("pages/html/login.html"))
  data := GetLoginPageData()

  page.Execute(w, data)
}

func LoginPost(w http.ResponseWriter, r *http.Request) {
  data := database.UserLoginData {
    Email: r.FormValue(loginFormEmail),
    Password: r.FormValue(loginFormPassword),
  }

  if !validator.ValidateStruct(data) {
    authentication.Unauthenticate(w, r)
    http.Redirect(w, r, errorPagePath, 302)
    return
  }

  if database.ValidateUserCreds(data.Email, data.Password) {
    username := database.GetUsernameFromEmail(data.Email)
    authentication.Authenticate(w, r, username, data.Email, data.Password)
    http.Redirect(w, r, emailsPagePath, 302)
    return
  }

  http.Redirect(w, r, errorPagePath, 302)
}

func ConfigLogin(r *mux.Router) {
  r.HandleFunc(loginPagePath, middleware.Chain(LoginGet)).Methods("GET")
  r.HandleFunc(loginPagePath, middleware.Chain(LoginPost)).Methods("POST")
}
