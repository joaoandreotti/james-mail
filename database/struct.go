package database

type UserLoginData struct {
  Email string `validate:"required,email"`
  Password string `validate:"required,ascii,min=8,max=64"`
}

type UserRegistrationData struct {
  Username string `validate:"required,alphanum,min=2,max=32"`
  Email string `validate:"required,email"`
  Password string `validate:"required,ascii,min=8,max=128"`
}

type PgpFormData struct {
  Password string `validate:"required,ascii,min=8,max=128"`
}

type SendEmailData struct {
  Recipient string `validate:"required,email"`
  EmailBody string `validate:"required,ascii,min=1`
}

type StoredEmail struct {
  EmailAddress string
  EmailBody string
}
