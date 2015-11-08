Magnet
======

Magnet is a tiny self-hosted bookmarks management tool written in Go(lang). 

Works on Go 1.5 and Rethinkdb 2.1.

![alt text](https://github.com/mvader/magnet/raw/master/magnet.png "Magnet screenshot").

Requisites
-------
* [Rethinkdb](http://rethinkdb.com)
* [Golang](http://golang.org/doc/install)

Setup
-------

```bash
go get github.com/mvader/magnet magnet 
cd $GOPATH/src/github.com/mvader/magnet
go build
cp config.sample.json config.json
# Edit config.json
nano config.json
./magnet
```

Go dependencies 
-------
* [github.com/dancannon/gorethink](https://github.com/dancannon/gorethink)
* [github.com/gorilla/sessions](https://github.com/gorilla/sessions)
* [github.com/codegangsta/martini](https://github.com/codegangsta/martini)
* [github.com/hoisie/mustache](https://github.com/hoisie/mustache)
* [github.com/justinas/nosurf](https://github.com/justinas/nosurf)
