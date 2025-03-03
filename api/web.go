package api

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ents-source/go-amember-api/amember"
)

var srv *http.Server
var amemberProApi *amember.Client

func Start(addr string, static string, paymentsApi *amember.Client) *sync.WaitGroup {
	wg := new(sync.WaitGroup)
	amemberProApi = paymentsApi

	//http.HandleFunc("/print", doPrint)
	http.Handle("/", http.FileServer(http.Dir(static)))

	go func() {
		srv = &http.Server{Addr: addr, Handler: http.DefaultServeMux}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}

		wg.Done()
	}()

	return wg
}

func Stop() {
	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}
}
