package smtp

import (
  "strings"
  "james-mail/database"
  "james-mail/pgp"

  "github.com/emersion/go-sasl"
  "github.com/emersion/go-smtp"
)

const smtpServer string = "localhost:2525"

func CreatePlainSession(email, password string) sasl.Client {
  auth := sasl.NewPlainClient("", email, password)
  return auth
}

func SendAuthenticatedEmail(auth sasl.Client, sender, recipient, emailBody string) {
  publicKey, _ := database.GetKeyPairFromEmail(sender)
  encryptedEmail := pgp.EncryptMessage(publicKey, emailBody)
  database.InsertSentEmail(sender, recipient, encryptedEmail)
  emailBodyReader := strings.NewReader(emailBody)
  smtp.SendMail(smtpServer, auth, sender, []string{recipient}, emailBodyReader)
}
