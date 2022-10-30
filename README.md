# mail-hub

Send email service

# Tech stack

- Golang: cobra, mux, negroni middleware...
- Message queue: Kafka
- Deploy: cmake, local registry, docker, docker compose, nginx

# Code structure

- cmd: cobra command for starting up application
- models: ORM models
- pkg: self-implemented packages for re-use
- util: shared constants and functions
- deploy: docker compose
- api: implement api server and handler
- consumer: implement queue subscriber
- scheduler: cronjob

# How to run

- Required: golang, cmake, docker and docker compose

step 1: Start local registry:
```
cd deploy
./start-local-registry
```

step 2: Build api
```
cd api
make dev
```

step 3: Build consumer
```
cd consumer
make dev
```

step 4: Update service config
```
- Update service config send-mail-consumer from file deploy/docker-compose.yml

environment:
       - KAFKABROKERS=172.17.0.1:9092
       - MAILSERVICE=sendgrid
       - SENDGRIDAPIKEY=sendgrid_api_key
       - MAILGUNDOMAIN=mailgun_domain
       - MAILGUNAPIKEY=mailgun_api_key

*NOTE* MAILSERVICE value: sendgrid, mailgun, mixed
With MAILSERVICE=mixed :
- Consumer use either Sendgrid and Mailgun which contribute to consumming queue message
```

step 5: docker-compose up

step 6: Update /etc/hosts: add this line: `127.0.0.1 mail-api.local`

# How to test:

- Use Postman with curl:

```
curl --location --request POST 'http://mail-api.local/api/mail/send' \
--header 'Authorization: Basic bWFpbC1odWI6bWFpbC1odWJAMTIz' \
--header 'Content-Type: application/json' \
--data-raw '{
    "from": {
        "name": "your name",
        "email": "your_verified_email@abc.com"
    },
    "to": {
        "name": "your recipent",
        "email": "your_receipent@gmail.com"
    },
    "subject": "Mail test",
    "html": "<p>Ok checked</p>"
}'
```

**Note**: Sender email should be verified email account from email service. 

# TODO:

- I implemented simple authenticated middleware with basic auth verification.
Furthermore, improve it by using better auth mechanism.
- Mail log is necessary for re-sending and handle webhook from mail service is needed.
- Smart failover mechanism should be applied. Ex: monitoring 3rd mail service health.
- In this version, only support switch by configuration manually.
- Because of time limitation, code may be not clean. If you have any suggestion, let you contribute to repo.


