package main

import (
	"flag"
	"fmt"
	"github.com/emicklei/go-restful"
	"net/http"
	"tail-project/api"
)

// http service
func httpService(httpAddr *string) {
	http.Handle("/", http.FileServer(http.Dir(api.DataPath)))
	fmt.Println("Start http server listen build addr is ", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, restful.DefaultContainer); err != nil {
		panic(err)
	}
}

//
func main() {
	buildAddr := flag.String("build", ":6543", "server http buildAddr")
	dataPath := flag.String("dataPath", "./data", " project data path")
	hostNameAndPort := flag.String("hostNameAndPort", "", "host name and port ")
	flag.Parse()
	api.DataPath = *dataPath
	api.HostNameAndPort = *hostNameAndPort
	httpService(buildAddr)
}
