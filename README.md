
#### URL Shortener

## Build and Run

In this project we use two containers, the first one is to do the build and the second one to expose the binary of the project.

```bash
make build
```

Run the URLShortener

```bash
make run-urlshortener
```

Using CURL

Generate shortener\
`curl -H "Content-Type: application/json" -X POST -d '{"url":"www.google.com"}' http://localhost:8080/encode/`

Response:
`{"success":true,"response":"http://localhost:8080/a"}`

Redirect:

`curl http://localhost:8080/a`

Telemetry
Information for related to a specific short URL can be obtained by 
`curl http://localhost:8080/urlinfo/a `
This provides the information on how many times this URL has been used, the first time the entry was generated
..

## Installation

You can install it using 'go get' or cloning the repository.

#### Use go get
```bash
go get github.com/sehgalswati/urlshortener
cd $GOPATH/src/github.com/sehgalswati/urlshortener
```
#### Cloning the repo
```bash
mkdir -p $GOPATH/src/github.com/sehgalswati
cd $GOPATH/src/github.com/sehgalswati
```
Clone repository 
```git clone https://github.com/sehgalswati/urlshortener.git```


Use GLIDE Package Management for Golang, for installation all packages 

https://github.com/Masterminds/glide
 
Run `glide install` in the folder.

Add your config for the method of persistence and other options in file config.json
```json
{
  "server": {
    "host": "0.0.0.0",
    "port": "8080"
  },
  "options": {
    "prefix": "localhost:8080/"
  },
  "postgres": {
    "host": "postgres_application",
    "port": "5432",
    "user": "urlshortener_db",
    "password": "postgrespassword",
    "db": "urlshortener_db"
  }
}
```

Telemetry
Information for related to a specific short URL can be obtained by
`curl http://localhost:8080/urlinfo/a `
This provides the information on how many times this URL has been used, the first time the entry was generated

