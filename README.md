## This is a playground. There are changes in different branches that include more idiomatic code.
Feel free to browse around.
When I get time I do a little more. Some changes have not been merged because tests have not been updated. 

# [Hash-O-Matic](https://github.com/dbyington/hash-o-matic) [![Build Status](https://travis-ci.org/dbyington/hash-o-matic.svg?branch=development)](https://travis-ci.org/dbyington/hash-o-matic)[![Coverage Status](https://coveralls.io/repos/github/dbyington/hash-o-matic/badge.svg?branch=development)](https://coveralls.io/github/dbyington/hash-o-matic?branch=development)
(note: coveralls.io has an [issue](https://github.com/lemurheavy/coveralls-public/issues/971) where the badge does not update. Supposedly fixed but since my badge isn't updating I'd say it's not fixed. If you want to see current test coverage click the badge to see the coveralls.io page. thx) 

A Go based hashing REST API server

This code is currently hosted on Heroku at [https://hash-o-matic.herokuapp.com](https://hash-o-matic.herokuapp.com). If you tried to go to this address and ended up here that's because a GET of / redirects here.

## What it does

The server accepts a form `POST` to `http://localhost:8080/hash` in the form of `password=stringToHa$h`. At this time the `%` character seems to be the only standard utf8 ASCII character not accepted, so avoid that one for now.
When the `POST` comes in the server responds with an id (numeric) that can be used to retrieve the calculated hash. The hash is retrieved by using the id in a `GET` to `http://localhost:8080/hash/{id}`.
The server responds to the `POST` immediately with the id, however system waits 5 seconds before actually generating the hash. Attempts to immediately retrieve the hash will not fail but will receive a 202 `Accepted` and a JSON body containing `{"ErrorMessage":"hash string not ready"}`. Meaning, "your request has been received. Further processing needs to be done to fulfill your request, come back later."

### Running the server

Running `hash-o-matic` will log connections to STDOUT, so you will see the queries come in. Logging control is planned for a future update.

```
$> hash-o-matic
2018/01/06 16:25:00 Server listening on: :8080
2018/01/06 16:34:31 [::1]:54591 GET / curl/7.49.0
2018/01/06 16:34:45 [::1]:54592 GET /hash curl/7.49.0
2018/01/06 16:35:19 [::1]:54597 POST /hash curl/7.49.0
2018/01/06 16:35:33 [::1]:54598 POST /hash curl/7.49.0
2018/01/06 16:35:42 [::1]:54599 GET /hash/1 curl/7.49.0
```

You can use any tool you wish to send the `POST` and `GET` requests to the server, including `curl`.
```
$> curl -XPOST http://localhost:8080/hash -d 'password=angryMonkey' -H 'Content-Type: application/x-www-form-urlencoded'
{"HashId":1}
$> curl http://localhost:8080/hash/1
{"HashString":"ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="}
```

### Stopping the server

Stopping the running `hash-o-matic` server can be done by simply entering `^C` or by sending a `PUT` request to `http://localhost:8080/shutdown`. Either shutdown method will wait for any current sessions to complete before closing the connection and exiting.


`^C` method
```
$> hash-o-matic
2018/01/06 17:24:50 Server listening on: :8080
^C2018/01/06 17:24:53 Received signal interrupt; shutting down
2018/01/06 17:24:53 Stopping server
2018/01/06 17:24:53 Shutdown complete.
$>
```
`PUT` method
Start server:
```
$> hash-o-matic
2018/01/06 17:25:40 Server listening on: :8080

```
Send `PUT`:
```
$> curl -XPUT http://localhost:8080/shutdown
shutting down...%
$>
```
Output by server:
```
2018/01/06 17:26:21 Received call to /shutdown, shutting down
2018/01/06 17:26:21 Stopping server
2018/01/06 17:26:21 Shutdown complete.
$>
```

## Getting Started

Clone this repo `git clone https://github.com/dbyington/hash-o-matic.git`
Then cd into `hash-o-matic` and run `go get`

### Prerequisites

You will need [Go](https://golang.org/) (Golang)


### Installing

Install Go
Clone this repo
`workdir $> git clone https://github.com/dbyington/hash-o-matic.git`
Change into the `hash-o-matic` directory and run `go install`

```
workdir $> cd hash-o-matic
workdir\hash-o-matic $> go install
```

Or you can run the server with the `go run` command.

```
workdir\hash-o-matic $> go run hash-o-matic.go
```

## Running the tests

To run the included tests run
```
hash-o-matic $> go test ./...
```

## Deployment

I doubt you want to actually deploy this for anything, and I recommend you don't.

## Built With

* [Pride](https://www.google.com/search?q=pride)
* [Go](https://golang.org)

## Contributing

Since this is more of a learning project for me I would appreciate any comments, updates, or fixes come to me as an issue with a description and a clue, rather than straight-up code. Thanks.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/dbyington/hash-o-matic/tags).

## Authors

* **Don Byington** - *Initial work* - [Don](https://github.com/dbyington)

See also the list of [contributors](https://github.com/dbyington/hash-o-matic/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* The Go documentation project [Go Doc](https://golang.org/doc/)
* [Go by Example](https://gobyexample.com/)
* [Golang Book](http://www.golang-book.com/)
```
