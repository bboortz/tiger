# tiger

tiger is a tiny webserver for testing purposes written in go.


# Features

* Provides different endpoints for different use-cases
* Zero Dependencies
* Available as docker build


# Endpoints

* [Index](http://localhost:8080/)
* [Static Page](http://localhost:8080/static)
* [Show Version](http://localhost:8080/version)
* [Show Headers](http://localhost:8080/headers)
* [Show Cookies](http://localhost:8080/cookies)
* [Show Environment Variables](http://localhost:8080/environ)
* [Show Memory Profiling](http://localhost:8080/memory)


# How to use

* build `make build`
* run `make run`
* dev & test `make dev`
