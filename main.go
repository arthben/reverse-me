package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	listen := flag.String("listen", "", "--listen=127.0.0.1:8080")
	target := flag.String("target", "", "--target=http://ip_target:port")
	flag.Parse()

	gapura := *target
	listener := *listen
	if len(gapura) == 0 || len(listener) == 0 {
		fmt.Println("How :\n reverse_me --listen=127.0.0.1:8080 --target=http://ip_target:port")
		os.Exit(1)
	}

	remote, err := url.Parse(gapura)
	if err != nil {
		panic(err)
	}

	// init log
	file, err := os.OpenFile("reverse-me.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalln("Unable to set logfile:", err.Error())
	}
	defer file.Close()

	log.SetOutput(file)

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc("/", handlerProxy(proxy, remote))

	fmt.Println(">> Proxy Listen on " + listener)

	err = http.ListenAndServe(listener, nil)
	if err != nil {
		fmt.Printf("Proxy Failed To Listen: %v\n", err)
		os.Exit(1)
	}
}

func handlerProxy(p *httputil.ReverseProxy, remote *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Scheme = remote.Scheme
		r.URL.Host = remote.Host

		// dump log request sebelum dikirimkan ke gateway tujuan
		reqDump, errDump := httputil.DumpRequestOut(r, true)
		if errDump != nil {
			log.Fatal(errDump)
		}
		log.Printf("\n\nREQUEST:\n%s", string(reqDump))

		p.ModifyResponse = modifyResponse()
		p.ServeHTTP(w, r)
	}
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("\n\nRESPONSE:\n%s\n\n", string(respDump))
		return nil
	}
}
