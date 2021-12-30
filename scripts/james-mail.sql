create keyspace james with replication = {'class' : 'SimpleStrategy', 'replication_factor': 1 };

create table james.received_emails(id UUID primary key, user_id UUID, sender text, email_body text);

create table james.sent_emails(id UUID primary key, user_id UUID, recipient text, email_body text);

create table james.userinfo(id UUID primary key, username text, email text, password text, private_key text, public_key text);
