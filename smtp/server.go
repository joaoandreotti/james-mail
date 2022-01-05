package smtp

import (
  "errors"
  "io"
  "log"
  "time"
  "sync"
  "james-mail/database"
  "james-mail/validator"
  "james-mail/pgp"
  "github.com/emersion/go-smtp"
)

type Backend struct{}

func (bkd *Backend) NewSession(_ smtp.ConnectionState, _ string) (smtp.Session, error) {
  return &Session{}, nil
}

type Session struct {
  Sender string
  Recipient string
  EmailBody string
}

func (s *Session) AuthPlain(email, password string) error {
  data := database.UserLoginData {
    Email: email,
    Password: password,
  }
  if !validator.ValidateStruct(data) {
    return errors.New("Invalid credentials")
  }
  if database.ValidateUserCreds(email, password) {
    return errors.New("Invalid email or password")
  }
  return nil
}

func (bkd *Backend) AnonymousLogin(_ *smtp.ConnectionState) (smtp.Session, error) {
  return &Session{}, nil
}

func (bkd *Backend) Login(_ *smtp.ConnectionState, username string, password string) (smtp.Session, error) {
  return &Session{}, nil
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
  s.Sender = from
  return nil
}

func (s *Session) Rcpt(to string) error {
  s.Recipient = to
  return nil
}

func (s *Session) Data(r io.Reader) error {
  if b, err := io.ReadAll(r); err != nil {
    return err
  } else {
    s.EmailBody = string(b)
  }
  return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
  publicKey, _ := database.GetKeyPairFromEmail(s.Recipient)
  encryptedEmail := pgp.EncryptMessage(publicKey, s.EmailBody)
  database.InsertReceivedEmail(s.Sender, s.Recipient, encryptedEmail)
  return nil
}

func StartSmtpServer (wg *sync.WaitGroup) {
  defer wg.Done()
  be := &Backend{}

  s := smtp.NewServer(be)

  s.Addr = ":2525"
  s.Domain = "localhost"
  s.ReadTimeout = 10 * time.Second
  s.WriteTimeout = 10 * time.Second
  s.MaxMessageBytes = 1024 * 1024
  s.MaxRecipients = 50
  s.AllowInsecureAuth = true

  log.Println("Starting smtp server at", s.Addr)
  s.ListenAndServe()
}
