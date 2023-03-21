# Web Sample

## How to run


```bash
$ cd $GOPATH/src/github.com
$ mkdir pchchv
$ cd pchchv
$ git clone https://github.com/pchchv/web.git
$ cd web/cmd
$ go run *.go
```

Or if you have [Docker](https://www.docker.com/), open the terminal and:

```bash
$ git clone https://github.com/pchchv/web.git
$ cd web
$ docker run \
-p 8080:8080 \
-p 9595:9595 \
-v ${PWD}:/go/src/github.com/pchchv/web/ \
-w /go/src/github.com/pchchv/web/cmd \
--rm -ti golang:latest go run *.go
```

You can try the following API calls with the example application. It also uses all the features provided by the web

1. `http://localhost:8080/`
   - Loads an HTML page
2. `http://localhost:8080/matchall/`
   - Route with wildcard parameter configured
   - All URIs which begin with `/matchall` will be matched because it has a wildcard variable
   - e.g.
     - http://localhost:8080/matchall/hello
     - http://localhost:8080/matchall/hello/world
     - http://localhost:8080/matchall/hello/world/user
3. `http://localhost:8080/api/<param>`
   - Route with a named 'param' configured
   - It will match all requests which match `/api/<single parameter>`
   - e.g.
     - http://localhost:8080/api/hello
     - http://localhost:8080/api/world
4. `http://localhost:8080/error-setter`
   - Route which sets an error and sets response status 500
5. `http://localhost:8080/api/<param>`
   - Route with a named 'param' configured
   - It will match all requests which match `/api/<single parameter>`
   - e.g.
     - http://localhost:8080/api/hello
     - http://localhost:8080/api/world