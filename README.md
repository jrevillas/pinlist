Magnet
======

Magnet is a tiny self-hosted bookmarks management tool written in Go(lang). 

Works on Go 1.5 and Rethinkdb 2.1.

![alt text](https://github.com/mvader/magnet/raw/master/magnet.png "Magnet screenshot").

Requisites
-------
* [Rethinkdb](http://rethinkdb.com)
* [Golang](http://golang.org/doc/install)

Setup and run
-------

```bash
go get github.com/mvader/magnet magnet 
cd $GOPATH/src/github.com/mvader/magnet
go build
./magnet
```

The default values are:
```bash
RDB_PORT_28015_TCP_PORT = "28015"
RDB_PORT_28015_TCP_ADDR = "localhost"
MAGNET_SESSION_KEY = "Here be dragons"
MAGNET_PORT = ":3000"
MAGNET_SESSION_EXPIRE = "1296000"
```

For change this you can export variables like that.
```bash
export MAGNET_PORT=":8000"
./magnet
```

Docker
------

## Build
```bash
go get github.com/mvader/magnet magnet 
cd $GOPATH/src/github.com/mvader/magnet
GO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o magnet .
docker build -t magnet .
```
##Run with Docker

```bash
docker run --name magnet-rethinkdb -v "$PWD:/data" -d rethinkdb
docker run -it --link magnet-rethinkdb:rdb -p 3000:3000 magnet
```

Go dependencies 
-------
* [github.com/dancannon/gorethink](https://github.com/dancannon/gorethink)
* [github.com/gorilla/sessions](https://github.com/gorilla/sessions)
* [github.com/codegangsta/martini](https://github.com/codegangsta/martini)
* [github.com/hoisie/mustache](https://github.com/hoisie/mustache)
* [github.com/justinas/nosurf](https://github.com/justinas/nosurf)
