package database

import(
  "github.com/gocql/gocql"
  "fmt"
)

var cluster *gocql.ClusterConfig
var session *gocql.Session

var getUserInfoQuery string = `
  select *
  from userinfo
  where username=? and email=? ALLOW FILTERING`

var insertUserInfoQuery string = `
  insert into userinfo
  (id, email, password, username) values
  (?, ?, ?, ?)`

var getUserPasswordQuery string = `
  select password
  from userinfo
  where email=? ALLOW FILTERING`

var getUsernameFromEmailQuery string = `
  select username
  from userinfo
  where email=? ALLOW FILTERING`

var getUserIdFromEmailQuery string = `
  select id
  from userinfo
  where email=? ALLOW FILTERING`

var getKeyPairFromEmailQuery string = `
  select public_key, private_key
  from userinfo
  where email=? ALLOW FILTERING`

var updateKeyPairFromIdQuery string = `
  update userinfo
  set public_key=?, private_key=?
  where id=?`

var insertSentEmailQuery string = `
  insert into sent_emails
  (id, user_id, recipient, email_body) values
  (?, ?, ?, ?)`

var insertReceivedEmailQuery string = `
  insert into received_emails
  (id, user_id, sender, email_body) values
  (?, ?, ?, ?)`

var getSentEmailsQuery string = `
  select recipient, email_body
  from sent_emails
  where user_id=? ALLOW FILTERING`

var getReceivedEmailsQuery string = `
  select sender, email_body
  from received_emails
  where user_id=? ALLOW FILTERING`

func InsertNewUser(username, email, password string) {
  hashedPassword := GenerateHashedPassword(password)
  total := session.Query(getUserInfoQuery, username, email).Iter().NumRows()
  if total == 0 {
    _ = session.Query(insertUserInfoQuery,
    gocql.TimeUUID(), email, hashedPassword, username).Exec()
  }
}

func ValidateUserCreds(email, password string) bool {
  var storedPassword string
  err := session.Query(getUserPasswordQuery, email).Scan(&storedPassword)
  if err != nil {
    return false
  }
  return ValidateHashedPassword(password, storedPassword)
}

func GetUsernameFromEmail(email string) string {
  var username string
  err := session.Query(getUsernameFromEmailQuery, email).Scan(&username)
  if err != nil {
    return "Not found"
  }
  return username
}

func GetKeyPairFromEmail(email string) (publicKey, privateKey string) {
  _ = session.Query(getKeyPairFromEmailQuery, email).Scan(&publicKey, &privateKey)
  return publicKey, privateKey
}

func UpdateKeyPairFromEmail(publicKey, privateKey, email string) {
  var id gocql.UUID
  _ = session.Query(getUserIdFromEmailQuery, email).Scan(&id)
  _ = session.Query(updateKeyPairFromIdQuery, publicKey, privateKey, id).Exec()
}

func InsertSentEmail(userEmail, recipient, emailBody string) {
  var id gocql.UUID
  _ = session.Query(getUserIdFromEmailQuery, userEmail).Scan(&id)
  _ = session.Query(insertSentEmailQuery, gocql.TimeUUID(), id, recipient, emailBody).Exec()
}

func InsertReceivedEmail(userEmail, sender, emailBody string) {
  var id gocql.UUID
  _ = session.Query(getUserIdFromEmailQuery, userEmail).Scan(&id)
  _ = session.Query(insertReceivedEmailQuery, gocql.TimeUUID(), id, sender, emailBody).Exec()
}

func GetSentEmailsFromEmail(email string) (emails []StoredEmail) {
  var id gocql.UUID
  _ = session.Query(getUserIdFromEmailQuery, email).Scan(&id)
  iter := session.Query(getSentEmailsQuery, id).Iter()
  for {
    row := make(map[string]interface{})
    if !iter.MapScan(row) {
      break
    }
    emailAddress, _ := row["recipient"]
    emailBody, _ := row["email_body"]
    emails = append(emails,
      StoredEmail {
        EmailAddress: fmt.Sprint(emailAddress),
        EmailBody: fmt.Sprint(emailBody),
      })
  }
  return emails
}

func GetReceivedEmailsFromEmail(email string) (emails []StoredEmail) {
  var id gocql.UUID
  _ = session.Query(getUserIdFromEmailQuery, email).Scan(&id)
  iter := session.Query(getReceivedEmailsQuery, id).Iter()
  for {
    row := make(map[string]interface{})
    if !iter.MapScan(row) {
      break
    }
    emailAddress, _ := row["sender"]
    emailBody, _ := row["email_body"]
    emails = append(emails,
      StoredEmail {
        EmailAddress: fmt.Sprint(emailAddress),
        EmailBody: fmt.Sprint(emailBody),
      })
  }
  return emails
}

func ConfigConnection() {
  cluster = gocql.NewCluster("172.17.0.2")
  cluster.Keyspace = "james"
  cluster.Consistency = gocql.Quorum
  session, _ = cluster.CreateSession()
}
