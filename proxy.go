package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Config struct {
	Host string
}

func logResponse(res *http.Response) error {
	log.Printf("%s %s %s\n", res.Request.Method, res.Request.URL.Path, res.Status)
	return nil
}

func (c *Config) handleRequest(w http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(c.Host)

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = logResponse

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	proxy.ServeHTTP(w, req)
}

func main() {
	var c Config

	if len(os.Args) < 2 {
		fmt.Println("Missing upstream host")
		os.Exit(1)
	}

	c.Host = os.Args[1]

	port := ":1337"
	if len(os.Args) == 3 {
		port = ":" + os.Args[2]
	}

	http.HandleFunc("/", c.handleRequest)
	log.Printf("Listening on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
