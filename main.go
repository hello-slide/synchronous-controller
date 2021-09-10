package main

import (
	"fmt"
	"net/http"

	networkutil "github.com/hello-slide/network-util"
	"github.com/hello-slide/synchronous-controller/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Roothandler)
	mux.HandleFunc("/sync/host", handler.HostHandler)
	mux.HandleFunc("/sync/visitor", handler.VisitorHandler)

	networkHandler := networkutil.CorsConfig.Handler(mux)

	if err := http.ListenAndServe(":3000", networkHandler); err != nil {
		fmt.Println(err)
	}
}
