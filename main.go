package main

/*
go get github.com/valyala/fasthttp
go get github.com/hoisie/mustache
go get gopkg.in/gorethink/gorethink.v3
*/
import (
	"flag"

	"./routes"
	"github.com/valyala/fasthttp"
)

var (
	addr = flag.String("addr", ":18081", "TCP address to listen to")
)

func main() {
	routes.Prepare()

	s := &fasthttp.Server{
		Handler:            routes.Handle,
		MaxRequestBodySize: 1024 * 1024 * 10,
	}
	s.ListenAndServe(*addr)

}
